package payments

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	// Billing accounts
	CreateBillingAccount(ctx context.Context, account *BillingAccount) error
	FindBillingAccountByUserID(ctx context.Context, userID uuid.UUID) (*BillingAccount, error)
	FindBillingAccountByStripeCustomerID(ctx context.Context, stripeCustomerID string) (*BillingAccount, error)

	// Subscriptions
	UpsertSubscription(ctx context.Context, subscription *Subscription) error
	FindSubscriptionByUserID(ctx context.Context, userID uuid.UUID) (*Subscription, error)
	FindSubscriptionByStripeID(ctx context.Context, stripeSubscriptionID string) (*Subscription, error)

	// Usage records
	CreateUsageRecord(ctx context.Context, record *UsageRecord) error
	FindUsageRecordByID(ctx context.Context, id uuid.UUID) (*UsageRecord, error)
	ListPendingUsage(ctx context.Context, limit int) ([]*UsageRecord, error)
	UpdateUsageRecord(ctx context.Context, record *UsageRecord) error

	// Webhook idempotency
	WebhookEventExists(ctx context.Context, stripeEventID string) (bool, error)
	CreateWebhookEvent(ctx context.Context, event *WebhookEvent) error
}
