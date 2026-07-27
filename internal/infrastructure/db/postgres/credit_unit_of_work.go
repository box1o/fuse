package postgres

import (
	"context"

	domain "fuse/internal/domain/credit"

	"gorm.io/gorm"
)

type CreditUnitOfWork struct {
	db *gorm.DB
}

func NewCreditUnitOfWork(db *gorm.DB) *CreditUnitOfWork {
	return &CreditUnitOfWork{
		db: db,
	}
}

func (uow *CreditUnitOfWork) WithinTransaction(
	ctx context.Context,
	operation func(
		accounts domain.AccountRepository,
		transactions domain.TransactionRepository,
	) error,
) error {
	return uow.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		accounts := NewCreditAccountRepository(tx)
		transactions := NewCreditTransactionRepository(tx)

		return operation(accounts, transactions)
	})
}
