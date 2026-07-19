package postgres

import (
	"context"
	stdErrors "errors"
	"fuse/internal/domain/credit"
	"fuse/internal/domain/credit/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreditAccountRepository struct {
	db *gorm.DB
}

func NewCreditAccountRepository(db *gorm.DB) credit.AccountRepository {
	return &CreditAccountRepository{db: db}
}

func (r *CreditAccountRepository) Create(ctx context.Context, account *credit.Account) error {
	if account == nil {
		return credit.ErrAccountNotFound
	}

	dbAccount, err := models.FromDomainAccount(account)
	if err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Create(dbAccount).Error; err != nil {
		return err
	}

	converted, err := dbAccount.ToDomainAccount()
	if err != nil {
		return err
	}

	*account = *converted

	return nil
}

func (r *CreditAccountRepository) FindByOwnerID(ctx context.Context, ownerID uuid.UUID) (*credit.Account, error) {
	if ownerID == uuid.Nil {
		return nil, credit.ErrOwnerIDRequired
	}

	var dbAccount models.DBAccount
	err := r.db.WithContext(ctx).
		Where("owner_id = ?", ownerID.String()).
		First(&dbAccount).Error

	if stdErrors.Is(err, gorm.ErrRecordNotFound) {
		return nil, credit.ErrAccountNotFound
	}

	if err != nil {
		return nil, err
	}

	account, err := dbAccount.ToDomainAccount()
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (r *CreditAccountRepository) Update(ctx context.Context, account *credit.Account) error {
	if account == nil {
		return credit.ErrAccountNotFound
	}

	if account.ID == uuid.Nil {
		return credit.ErrInvalidAccount
	}

	dbAccount, err := models.FromDomainAccount(account)
	if err != nil {
		return err
	}

	result := r.db.WithContext(ctx).
		Model(&models.DBAccount{}).
		Where("id = ?", account.ID).
		Updates(map[string]any{
			"balance":    dbAccount.Balance,
			"updated_at": dbAccount.UpdatedAt,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return credit.ErrAccountNotFound
	}

	return nil
}
