package models

import (
	"fmt"

	"fuse/internal/domain/credit"
	"fuse/internal/infrastructure/db"

	"github.com/google/uuid"
)

type DBAccount struct {
	db.Model

	OwnerID string `gorm:"type:uuid;not null;uniqueIndex" json:"owner_id"`
	Balance int64  `gorm:"not null;default:0;check:balance >= 0" json:"balance"`
}

func (DBAccount) TableName() string {
	return "credit_accounts"
}

func FromDomainAccount(account *credit.Account) (*DBAccount, error) {
	if account == nil {
		return nil, credit.ErrInvalidAccount
	}

	return &DBAccount{
		Model:   db.Model{ID: account.ID, CreatedAt: account.CreatedAt, UpdatedAt: account.UpdatedAt},
		OwnerID: account.OwnerID.String(),
		Balance: account.Balance.Value(),
	}, nil
}

func (account *DBAccount) ToDomainAccount() (*credit.Account, error) {
	if account == nil {
		return nil, credit.ErrInvalidAccount
	}

	ownerID, err := uuid.Parse(account.OwnerID)
	if err != nil {
		return nil, fmt.Errorf("parse account owner ID: %w", err)
	}

	balance, err := credit.NewAmount(account.Balance)
	if err != nil {
		return nil, fmt.Errorf("create account balance: %w", err)
	}

	return &credit.Account{
		ID:        account.ID,
		OwnerID:   ownerID,
		Balance:   balance,
		CreatedAt: account.CreatedAt,
		UpdatedAt: account.UpdatedAt,
	}, nil
}
