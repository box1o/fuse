package models

import (
	"fmt"
	"fuse/internal/domain/credit"
	"fuse/internal/infrastructure/db"

	"github.com/google/uuid"
)

type DBTransaction struct {
	db.Model

	AccountID string `gorm:"type:uuid;not null;index" json:"account_id"`

	Type   string `gorm:"not null;size:50;index" json:"type"`
	Source string `gorm:"not null;size:50;index" json:"source"`
	Amount int64  `gorm:"not null;check:amount > 0" json:"amount"`

	ReferenceID       string `gorm:"size:255;index" json:"reference_id,omitempty"`
	ExternalReference string `gorm:"size:255;index" json:"external_reference,omitempty"`

	IdempotencyKey string `gorm:"not null;size:255;uniqueIndex" json:"idempotency_key"`
}

func (DBTransaction) TableName() string {
	return "credit_transactions"
}

func FromDomainTransaction(transaction *credit.Transaction) (*DBTransaction, error) {
	if transaction == nil {
		return nil, credit.ErrInvalidTransaction
	}

	return &DBTransaction{Model: db.Model{ID: transaction.ID, CreatedAt: transaction.CreatedAt},
		AccountID:         transaction.AccountID.String(),
		Type:              string(transaction.Type),
		Source:            string(transaction.Source),
		Amount:            transaction.Amount.Value(),
		ReferenceID:       transaction.ReferenceID,
		ExternalReference: transaction.ExternalReference,
		IdempotencyKey:    transaction.IdempotencyKey,
	}, nil
}

func (transaction *DBTransaction) ToDomainTransaction() (*credit.Transaction, error) {
	if transaction == nil {
		return nil, credit.ErrInvalidTransaction
	}

	accountID, err := uuid.Parse(transaction.AccountID)
	if err != nil {
		return nil, fmt.Errorf(
			"parse credit transaction account ID: %w",
			err,
		)
	}

	amount, err := credit.NewAmount(transaction.Amount)
	if err != nil {
		return nil, fmt.Errorf(
			"create credit transaction amount: %w",
			err,
		)
	}

	transactionType := credit.TransactionType(transaction.Type)
	if !transactionType.IsValid() {
		return nil, credit.ErrInvalidTransactionType
	}

	source := credit.TransactionSource(transaction.Source)
	if !source.IsValid() {
		return nil, credit.ErrInvalidTransactionSource
	}

	return &credit.Transaction{
		ID:                transaction.ID,
		AccountID:         accountID,
		Type:              transactionType,
		Source:            source,
		Amount:            amount,
		ReferenceID:       transaction.ReferenceID,
		ExternalReference: transaction.ExternalReference,
		IdempotencyKey:    transaction.IdempotencyKey,
		CreatedAt:         transaction.CreatedAt,
	}, nil
}
