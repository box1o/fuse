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
	GetSubscription(ctx context.Context, userID uuid.UUID) (*Subscription, error)
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

	// Credit grants
	GetCreditBalance(ctx context.Context, userID uuid.UUID) (int64, error)
	CreateCreditGrant(ctx context.Context, userID uuid.UUID, amount int64) error // Used after Stripe confirms a credit-pack payment.
	GetGrantedCredits(ctx context.Context, userID uuid.UUID) (int64, error)      // Used to get the total granted credits for a user.
	GetUsedCredits(ctx context.Context, userID uuid.UUID) (int64, error)         //
}
