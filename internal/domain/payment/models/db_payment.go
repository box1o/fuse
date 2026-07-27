package models

import (
	"fmt"

	"fuse/internal/domain/payment"
	"fuse/internal/infrastructure/db"
	"time"

	"github.com/google/uuid"
)

type DBPayment struct {
	db.Model

	OwnerID           string     `gorm:"type:uuid;not null;index" json:"owner_id"`
	CreditPackID      string     `gorm:"type:uuid;not null;index" json:"credit_pack_id"`
	Credits           int64      `gorm:"not null;check:credits > 0" json:"credits"`
	Amount            int64      `gorm:"not null;check:amount > 0" json:"amount"`
	Currency          string     `gorm:"not null;size:3" json:"currency"`
	Status            string     `gorm:"not null;size:32;index" json:"status"`
	Provider          string     `gorm:"not null;size:32;index;uniqueIndex:idx_payments_provider_session,priority:1;uniqueIndex:idx_payments_provider_payment,priority:1" json:"provider"`
	ProviderSessionID *string    `gorm:"size:255;uniqueIndex:idx_payments_provider_session,priority:2" json:"provider_session_id,omitempty"`
	ProviderPaymentID *string    `gorm:"size:255;uniqueIndex:idx_payments_provider_payment,priority:2" json:"provider_payment_id,omitempty"`
	CompletedAt       *time.Time `json:"completed_at,omitempty"`
}

func (DBPayment) TableName() string {
	return "payments"
}

func FromDomainPayment(domainPayment *payment.Payment) (*DBPayment, error) {
	if domainPayment == nil || domainPayment.ID == uuid.Nil {
		return nil, payment.ErrPaymentNotFound
	}

	return &DBPayment{
		Model: db.Model{
			ID:        domainPayment.ID,
			CreatedAt: domainPayment.CreatedAt,
			UpdatedAt: domainPayment.UpdatedAt,
		},
		OwnerID:      domainPayment.OwnerID.String(),
		CreditPackID: domainPayment.CreditPackID.String(),
		Credits:      domainPayment.Credits,
		Amount:       domainPayment.Amount,
		Currency:     domainPayment.Currency,
		Status:       string(domainPayment.Status),
		Provider:     string(domainPayment.Provider),
		ProviderSessionID: optionalString(
			domainPayment.ProviderSessionID,
		),
		ProviderPaymentID: optionalString(
			domainPayment.ProviderPaymentID,
		),
		CompletedAt: domainPayment.CompletedAt,
	}, nil
}

func (dbPayment *DBPayment) ToDomainPayment() (*payment.Payment, error) {
	if dbPayment == nil {
		return nil, payment.ErrPaymentNotFound
	}

	ownerID, err := uuid.Parse(dbPayment.OwnerID)
	if err != nil {
		return nil, fmt.Errorf("parse payment owner ID: %w", err)
	}

	creditPackID, err := uuid.Parse(dbPayment.CreditPackID)
	if err != nil {
		return nil, fmt.Errorf("parse payment credit pack ID: %w", err)
	}

	return payment.RestorePayment(payment.RestorePaymentInput{
		ID:                dbPayment.ID,
		OwnerID:           ownerID,
		CreditPackID:      creditPackID,
		Credits:           dbPayment.Credits,
		Amount:            dbPayment.Amount,
		Currency:          dbPayment.Currency,
		Status:            payment.Status(dbPayment.Status),
		Provider:          payment.Provider(dbPayment.Provider),
		ProviderSessionID: stringValue(dbPayment.ProviderSessionID),
		ProviderPaymentID: stringValue(dbPayment.ProviderPaymentID),
		CreatedAt:         dbPayment.CreatedAt,
		UpdatedAt:         dbPayment.UpdatedAt,
		CompletedAt:       dbPayment.CompletedAt,
	})
}

func optionalString(value string) *string {
	if value == "" {
		return nil
	}

	return &value
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}

	return *value
}
