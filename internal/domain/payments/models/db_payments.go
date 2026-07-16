package models

import (
	"time"

	"fuse/internal/domain/payments"
	"fuse/internal/infrastructure/db"

	"github.com/google/uuid"
)

type DBBillingAccount struct {
	db.Model

	UserID           uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`
	StripeCustomerID string    `gorm:"not null;size:255;uniqueIndex"`
}

func (DBBillingAccount) TableName() string {
	return "billing_accounts"
}

func BillingAccountFromDomain(
	account *payments.BillingAccount,
) *DBBillingAccount {
	if account == nil {
		return nil
	}

	return &DBBillingAccount{
		Model: db.Model{
			ID:        account.ID,
			CreatedAt: account.CreatedAt,
			UpdatedAt: account.UpdatedAt,
		},
		UserID:           account.UserID,
		StripeCustomerID: account.StripeCustomerID,
	}
}

func (m *DBBillingAccount) ToDomain() *payments.BillingAccount {
	if m == nil {
		return nil
	}

	return &payments.BillingAccount{
		ID:               m.ID,
		UserID:           m.UserID,
		StripeCustomerID: m.StripeCustomerID,
		CreatedAt:        m.CreatedAt,
		UpdatedAt:        m.UpdatedAt,
	}
}

type DBSubscription struct {
	db.Model

	UserID               uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`
	StripeSubscriptionID string    `gorm:"not null;size:255;uniqueIndex"`
	Status               string    `gorm:"not null;size:50"`
	CurrentPeriodStart   time.Time `gorm:"not null"`
	CurrentPeriodEnd     time.Time `gorm:"not null"`
	CancelAtPeriodEnd    bool      `gorm:"not null;default:false"`
}

func (DBSubscription) TableName() string {
	return "billing_subscriptions"
}

func SubscriptionFromDomain(
	subscription *payments.Subscription,
) *DBSubscription {
	if subscription == nil {
		return nil
	}

	return &DBSubscription{
		Model: db.Model{
			ID:        subscription.ID,
			CreatedAt: subscription.CreatedAt,
			UpdatedAt: subscription.UpdatedAt,
		},
		UserID:               subscription.UserID,
		StripeSubscriptionID: subscription.StripeSubscriptionID,
		Status:               string(subscription.Status),
		CurrentPeriodStart:   subscription.CurrentPeriodStart,
		CurrentPeriodEnd:     subscription.CurrentPeriodEnd,
		CancelAtPeriodEnd:    subscription.CancelAtPeriodEnd,
	}
}

func (m *DBSubscription) ToDomain() *payments.Subscription {
	if m == nil {
		return nil
	}

	return &payments.Subscription{
		ID:                   m.ID,
		UserID:               m.UserID,
		StripeSubscriptionID: m.StripeSubscriptionID,
		Status:               payments.SubscriptionStatus(m.Status),
		CurrentPeriodStart:   m.CurrentPeriodStart,
		CurrentPeriodEnd:     m.CurrentPeriodEnd,
		CancelAtPeriodEnd:    m.CancelAtPeriodEnd,
		CreatedAt:            m.CreatedAt,
		UpdatedAt:            m.UpdatedAt,
	}
}

type DBUsageRecord struct {
	db.Model

	UserID         uuid.UUID `gorm:"type:uuid;not null;index"`
	ResourceType   string    `gorm:"not null;size:20;index"`
	Quantity       int64     `gorm:"not null"`
	OccurredAt     time.Time `gorm:"not null;index"`
	IdempotencyKey string    `gorm:"not null;size:255;uniqueIndex"`
	StripeEventID  string    `gorm:"size:255"`
	Status         string    `gorm:"not null;size:20;index"`
}

func (DBUsageRecord) TableName() string {
	return "billing_usage_records"
}

func UsageRecordFromDomain(record *payments.UsageRecord) *DBUsageRecord {
	if record == nil {
		return nil
	}

	return &DBUsageRecord{
		Model: db.Model{
			ID:        record.ID,
			CreatedAt: record.CreatedAt,
			UpdatedAt: record.UpdatedAt,
		},
		UserID:         record.UserID,
		ResourceType:   string(record.ResourceType),
		Quantity:       record.Quantity,
		OccurredAt:     record.OccurredAt,
		IdempotencyKey: record.IdempotencyKey,
		StripeEventID:  record.StripeEventID,
		Status:         string(record.Status),
	}
}

func (m *DBUsageRecord) ToDomain() *payments.UsageRecord {
	if m == nil {
		return nil
	}

	return &payments.UsageRecord{
		ID:             m.ID,
		UserID:         m.UserID,
		ResourceType:   payments.ResourceType(m.ResourceType),
		Quantity:       m.Quantity,
		OccurredAt:     m.OccurredAt,
		IdempotencyKey: m.IdempotencyKey,
		StripeEventID:  m.StripeEventID,
		Status:         payments.UsageStatus(m.Status),
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}

type DBWebhookEvent struct {
	StripeEventID string    `gorm:"primaryKey;size:255"`
	EventType     string    `gorm:"not null;size:255;index"`
	ProcessedAt   time.Time `gorm:"not null"`
}

func (DBWebhookEvent) TableName() string {
	return "stripe_webhook_events"
}

func WebhookEventFromDomain(
	event *payments.WebhookEvent,
) *DBWebhookEvent {
	if event == nil {
		return nil
	}

	return &DBWebhookEvent{
		StripeEventID: event.StripeEventID,
		EventType:     event.EventType,
		ProcessedAt:   event.ProcessedAt,
	}
}

func (m *DBWebhookEvent) ToDomain() *payments.WebhookEvent {
	if m == nil {
		return nil
	}

	return &payments.WebhookEvent{
		StripeEventID: m.StripeEventID,
		EventType:     m.EventType,
		ProcessedAt:   m.ProcessedAt,
	}
}
