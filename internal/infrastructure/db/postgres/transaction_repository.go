package postgres

import (
	"context"
	stdErrors "errors"
	"fuse/internal/domain/credit"
	"fuse/internal/domain/credit/models"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreditTransactionRepository struct {
	db *gorm.DB
}

const (
	defaultTransactionLimit = 50
	maxTransactionLimit     = 200
)

var _ credit.TransactionRepository = (*CreditTransactionRepository)(nil)

func NewCreditTransactionRepository(db *gorm.DB) credit.TransactionRepository {
	return &CreditTransactionRepository{
		db: db,
	}
}

func (r *CreditTransactionRepository) Create(ctx context.Context, transaction *credit.Transaction) error {
	if transaction == nil {
		return credit.ErrInvalidTransaction
	}

	dbTransaction, err := models.FromDomainTransaction(transaction)
	if err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Create(dbTransaction).Error; err != nil {
		if isUniqueConstraintError(err) {
			return credit.ErrTransactionAlreadyExists
		}

		return err
	}

	converted, err := dbTransaction.ToDomainTransaction()
	if err != nil {
		return err
	}

	*transaction = *converted

	return nil
}

func (r *CreditTransactionRepository) FindByID(ctx context.Context, id uuid.UUID) (*credit.Transaction, error) {
	if id == uuid.Nil {
		return nil, credit.ErrInvalidTransaction
	}

	var dbTransaction models.DBTransaction

	err := r.db.
		WithContext(ctx).
		First(&dbTransaction, "id = ?", id).
		Error

	if stdErrors.Is(err, gorm.ErrRecordNotFound) {
		return nil, credit.ErrInvalidTransaction
	}

	if err != nil {
		return nil, err
	}

	transaction, err := dbTransaction.ToDomainTransaction()
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (r *CreditTransactionRepository) FindByIdempotencyKey(ctx context.Context, idempotencyKey string) (*credit.Transaction, error) {
	idempotencyKey = strings.TrimSpace(idempotencyKey)
	if idempotencyKey == "" {
		return nil, credit.ErrIdempotencyKeyRequired
	}

	var dbTransaction models.DBTransaction

	err := r.db.WithContext(ctx).
		Where("idempotency_key = ?", idempotencyKey).
		First(&dbTransaction).Error

	if stdErrors.Is(err, gorm.ErrRecordNotFound) {
		return nil, credit.ErrInvalidTransaction
	}

	if err != nil {
		return nil, err
	}

	transaction, err := dbTransaction.ToDomainTransaction()
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (r *CreditTransactionRepository) ListByAccountID(ctx context.Context, accountID uuid.UUID, limit int, offset int) ([]*credit.Transaction, error) {
	if accountID == uuid.Nil {
		return nil, credit.ErrAccountIDRequired
	}

	if limit <= 0 {
		limit = defaultTransactionLimit
	}

	if limit > maxTransactionLimit {
		limit = maxTransactionLimit
	}

	if offset < 0 {
		offset = 0
	}

	var dbTransactions []models.DBTransaction

	err := r.db.WithContext(ctx).
		Where("account_id = ?", accountID.String()).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&dbTransactions).Error

	if err != nil {
		return nil, err
	}

	transactions := make(
		[]*credit.Transaction,
		0,
		len(dbTransactions),
	)

	for i := range dbTransactions {
		transaction, err := dbTransactions[i].ToDomainTransaction()
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func isUniqueConstraintError(err error) bool {
	if err == nil {
		return false
	}

	message := strings.ToLower(err.Error())

	return strings.Contains(message, "duplicate key") ||
		strings.Contains(message, "unique constraint")
}
