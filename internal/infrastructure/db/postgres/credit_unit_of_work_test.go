package postgres_test

import (
	"context"
	stdErrors "errors"
	"os"
	"testing"
	"time"

	domain "fuse/internal/domain/credit"
	creditModels "fuse/internal/domain/credit/models"
	"fuse/internal/infrastructure/db/postgres"

	"github.com/google/uuid"
	gormPostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestCreditUnitOfWork_RollsBackAccountUpdateWhenTransactionCreationFails(
	t *testing.T,
) {
	dsn := os.Getenv("FUSE_TEST_DATABASE_DSN")
	if dsn == "" {
		t.Skip("FUSE_TEST_DATABASE_DSN is not configured")
	}

	db, err := gorm.Open(
		gormPostgres.Open(dsn),
		&gorm.Config{
			NowFunc: func() time.Time {
				return time.Now().UTC()
			},
		},
	)
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}

	if err := db.AutoMigrate(
		&creditModels.DBAccount{},
		&creditModels.DBTransaction{},
	); err != nil {
		t.Fatalf("migrate credit test tables: %v", err)
	}

	ctx := context.Background()
	ownerID := uuid.New()
	idempotencyKey := "test:duplicate:" + uuid.NewString()

	t.Cleanup(func() {
		cleanupCreditTestData(
			t,
			db,
			ownerID,
			idempotencyKey,
		)
	})

	accountRepository := postgres.NewCreditAccountRepository(db)
	transactionRepository := postgres.NewCreditTransactionRepository(db)

	initialBalance := mustCreateAmount(t, 100)

	account, err := domain.NewAccount(ownerID, initialBalance)
	if err != nil {
		t.Fatalf("create domain account: %v", err)
	}

	if err := accountRepository.Create(ctx, account); err != nil {
		t.Fatalf("persist initial account: %v", err)
	}

	existingTransaction, err := domain.NewTransaction(
		domain.NewTransactionInput{
			AccountID:      account.ID,
			Type:           domain.TransactionTypeSpend,
			Source:         domain.TransactionSourceComputeJob,
			Amount:         mustCreateAmount(t, 1),
			ReferenceID:    "existing-compute-job",
			IdempotencyKey: idempotencyKey,
		},
	)
	if err != nil {
		t.Fatalf("create existing transaction: %v", err)
	}

	if err := transactionRepository.Create(
		ctx,
		existingTransaction,
	); err != nil {
		t.Fatalf("persist existing transaction: %v", err)
	}

	unitOfWork := postgres.NewCreditUnitOfWork(db)

	err = unitOfWork.WithinTransaction(
		ctx,
		func(
			accounts domain.AccountRepository,
			transactions domain.TransactionRepository,
		) error {
			transactionalAccount, err := accounts.FindByOwnerID(
				ctx,
				ownerID,
			)
			if err != nil {
				return err
			}

			if err := transactionalAccount.Spend(
				mustCreateAmount(t, 30),
			); err != nil {
				return err
			}

			if err := accounts.Update(
				ctx,
				transactionalAccount,
			); err != nil {
				return err
			}

			duplicateTransaction, err := domain.NewTransaction(
				domain.NewTransactionInput{
					AccountID:      transactionalAccount.ID,
					Type:           domain.TransactionTypeSpend,
					Source:         domain.TransactionSourceComputeJob,
					Amount:         mustCreateAmount(t, 30),
					ReferenceID:    "new-compute-job",
					IdempotencyKey: idempotencyKey,
				},
			)
			if err != nil {
				return err
			}

			return transactions.Create(
				ctx,
				duplicateTransaction,
			)
		},
	)

	if !stdErrors.Is(err, domain.ErrTransactionAlreadyExists) {
		t.Fatalf(
			"expected ErrTransactionAlreadyExists, got %v",
			err,
		)
	}

	persistedAccount, err := accountRepository.FindByOwnerID(
		ctx,
		ownerID,
	)
	if err != nil {
		t.Fatalf("reload account after rollback: %v", err)
	}

	if persistedAccount.Balance.Value() != 100 {
		t.Fatalf(
			"expected balance to roll back to 100, got %d",
			persistedAccount.Balance.Value(),
		)
	}

	var transactionCount int64

	if err := db.WithContext(ctx).
		Model(&creditModels.DBTransaction{}).
		Where("account_id = ?", account.ID.String()).
		Count(&transactionCount).
		Error; err != nil {
		t.Fatalf("count persisted transactions: %v", err)
	}

	if transactionCount != 1 {
		t.Fatalf(
			"expected only the original transaction, got %d",
			transactionCount,
		)
	}
}

func mustCreateAmount(
	t *testing.T,
	value int64,
) domain.Amount {
	t.Helper()

	amount, err := domain.NewAmount(value)
	if err != nil {
		t.Fatalf("create amount %d: %v", value, err)
	}

	return amount
}

func cleanupCreditTestData(
	t *testing.T,
	db *gorm.DB,
	ownerID uuid.UUID,
	idempotencyKey string,
) {
	t.Helper()

	var account creditModels.DBAccount

	result := db.
		Where("owner_id = ?", ownerID.String()).
		First(&account)

	if result.Error != nil &&
		!stdErrors.Is(result.Error, gorm.ErrRecordNotFound) {
		t.Errorf("find test account during cleanup: %v", result.Error)
		return
	}

	if result.Error == nil {
		if err := db.
			Where("account_id = ?", account.ID.String()).
			Delete(&creditModels.DBTransaction{}).
			Error; err != nil {
			t.Errorf("delete test transactions: %v", err)
		}
	}

	if err := db.
		Where("idempotency_key = ?", idempotencyKey).
		Delete(&creditModels.DBTransaction{}).
		Error; err != nil {
		t.Errorf("delete transaction by idempotency key: %v", err)
	}

	if err := db.
		Where("owner_id = ?", ownerID.String()).
		Delete(&creditModels.DBAccount{}).
		Error; err != nil {
		t.Errorf("delete test account: %v", err)
	}
}
