package payments

import (
	"context"
	"testing"
	"time"

	domain "fuse/internal/domain/payments"
	"fuse/pkg/config"

	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v83"
)

type fakePaymentsRepo struct {
	billingAccountsByUserID map[uuid.UUID]*domain.BillingAccount
	billingAccountsByStripe map[string]*domain.BillingAccount
	subscriptionsByUserID   map[uuid.UUID]*domain.Subscription
	createdSubscriptions    []*domain.Subscription
}

func newFakePaymentsRepo() *fakePaymentsRepo {
	return &fakePaymentsRepo{
		billingAccountsByUserID: map[uuid.UUID]*domain.BillingAccount{},
		billingAccountsByStripe: map[string]*domain.BillingAccount{},
		subscriptionsByUserID:    map[uuid.UUID]*domain.Subscription{},
	}
}

func (r *fakePaymentsRepo) CreateBillingAccount(_ context.Context, account *domain.BillingAccount) error {
	r.billingAccountsByUserID[account.UserID] = account
	r.billingAccountsByStripe[account.StripeCustomerID] = account
	return nil
}

func (r *fakePaymentsRepo) FindBillingAccountByUserID(_ context.Context, userID uuid.UUID) (*domain.BillingAccount, error) {
	if account, ok := r.billingAccountsByUserID[userID]; ok {
		return account, nil
	}
	return nil, domain.ErrBillingAccountNotFound
}

func (r *fakePaymentsRepo) FindBillingAccountByStripeCustomerID(_ context.Context, stripeCustomerID string) (*domain.BillingAccount, error) {
	if account, ok := r.billingAccountsByStripe[stripeCustomerID]; ok {
		return account, nil
	}
	return nil, domain.ErrBillingAccountNotFound
}

func (r *fakePaymentsRepo) UpsertSubscription(_ context.Context, subscription *domain.Subscription) error {
	copy := *subscription
	r.subscriptionsByUserID[subscription.UserID] = &copy
	r.createdSubscriptions = append(r.createdSubscriptions, &copy)
	return nil
}

func (r *fakePaymentsRepo) FindSubscriptionByUserID(_ context.Context, userID uuid.UUID) (*domain.Subscription, error) {
	if sub, ok := r.subscriptionsByUserID[userID]; ok {
		return sub, nil
	}
	return nil, domain.ErrSubscriptionNotFound
}

func (r *fakePaymentsRepo) FindSubscriptionByStripeID(_ context.Context, stripeSubscriptionID string) (*domain.Subscription, error) {
	for _, sub := range r.subscriptionsByUserID {
		if sub.StripeSubscriptionID == stripeSubscriptionID {
			return sub, nil
		}
	}
	return nil, domain.ErrSubscriptionNotFound
}

func (r *fakePaymentsRepo) CreateUsageRecord(context.Context, *domain.UsageRecord) error { return nil }
func (r *fakePaymentsRepo) FindUsageRecordByID(context.Context, uuid.UUID) (*domain.UsageRecord, error) {
	return nil, domain.ErrUsageRecordNotFound
}
func (r *fakePaymentsRepo) ListPendingUsage(context.Context, int) ([]*domain.UsageRecord, error) {
	return nil, nil
}
func (r *fakePaymentsRepo) UpdateUsageRecord(context.Context, *domain.UsageRecord) error { return nil }
func (r *fakePaymentsRepo) WebhookEventExists(context.Context, string) (bool, error) { return false, nil }
func (r *fakePaymentsRepo) CreateWebhookEvent(context.Context, *domain.WebhookEvent) error { return nil }

func TestCancelThenRebuy_UpsertsNewSubscriptionAndKeepsLatestEntitlement(t *testing.T) {
	ctx := context.Background()
	repo := newFakePaymentsRepo()
	svc := NewService(&config.Config{Stripe: config.StripeConfig{WebhookSecret: "whsec_test"}}, repo)

	userID := uuid.New()
	account, err := domain.NewBillingAccount(userID, "cus_123")
	if err != nil {
		t.Fatalf("expected billing account creation to succeed, got error: %v", err)
	}
	repo.billingAccountsByUserID[userID] = account
	repo.billingAccountsByStripe[account.StripeCustomerID] = account

	start := time.Date(2026, time.July, 1, 12, 0, 0, 0, time.UTC)
	end := time.Date(2026, time.July, 31, 12, 0, 0, 0, time.UTC)
	canceledSub, err := domain.NewSubscription(userID, "sub_old", domain.SubscriptionStatusActive, start, end, true)
	if err != nil {
		t.Fatalf("expected subscription creation to succeed, got error: %v", err)
	}
	if err := repo.UpsertSubscription(ctx, canceledSub); err != nil {
		t.Fatalf("expected initial subscription save to succeed, got error: %v", err)
	}

	if !canceledSub.IsProAt(time.Date(2026, time.July, 20, 12, 0, 0, 0, time.UTC)) {
		t.Fatalf("expected canceled subscription to remain pro until its end date")
	}
	if canceledSub.IsProAt(time.Date(2026, time.August, 1, 0, 0, 0, 0, time.UTC)) {
		t.Fatalf("expected canceled subscription to stop being pro after its end date")
	}

	rebuy := stripe.Subscription{
		ID:     "sub_new",
		Status: stripe.SubscriptionStatusActive,
		Metadata: map[string]string{
			"user_id": userID.String(),
		},
		Customer: &stripe.Customer{ID: account.StripeCustomerID},
		Items: &stripe.SubscriptionItemList{
			Data: []*stripe.SubscriptionItem{
				{
					CurrentPeriodStart: start.Unix(),
					CurrentPeriodEnd:   end.AddDate(0, 1, 0).Unix(),
				},
			},
		},
	}

	if err := svc.upsertSubscriptionFromStripe(ctx, &rebuy); err != nil {
		t.Fatalf("expected rebuy webhook to succeed, got error: %v", err)
	}

	latest, err := repo.FindSubscriptionByUserID(ctx, userID)
	if err != nil {
		t.Fatalf("expected latest subscription lookup to succeed, got error: %v", err)
	}

	if latest.StripeSubscriptionID != "sub_new" {
		t.Fatalf("expected latest subscription to be the new one, got %s", latest.StripeSubscriptionID)
	}

	if !latest.IsProAt(time.Date(2026, time.August, 15, 0, 0, 0, 0, time.UTC)) {
		t.Fatalf("expected new subscription to keep pro access during its new billing period")
	}
}

func TestCreateCheckoutSession_BlocksWhileProSubscriptionIsStillActive(t *testing.T) {
	ctx := context.Background()
	repo := newFakePaymentsRepo()
	svc := NewService(&config.Config{
		Stripe: config.StripeConfig{
			SecretKey:  "sk_test_123",
			ProPriceID: "price_pro",
		},
	}, repo)

	userID := uuid.New()
	account, err := domain.NewBillingAccount(userID, "cus_123")
	if err != nil {
		t.Fatalf("expected billing account creation to succeed, got error: %v", err)
	}
	repo.billingAccountsByUserID[userID] = account
	repo.billingAccountsByStripe[account.StripeCustomerID] = account

	now := time.Now().UTC()
	sub, err := domain.NewSubscription(
		userID,
		"sub_old",
		domain.SubscriptionStatusActive,
		now.Add(-7*24*time.Hour),
		now.Add(7*24*time.Hour),
		true,
	)
	if err != nil {
		t.Fatalf("expected subscription creation to succeed, got error: %v", err)
	}
	if err := repo.UpsertSubscription(ctx, sub); err != nil {
		t.Fatalf("expected subscription save to succeed, got error: %v", err)
	}

	result, err := svc.CreateCheckoutSession(ctx, userID, "https://example.com/success", "https://example.com/cancel", "price_pro")
	if err != domain.ErrProSubscriptionStillActive {
		t.Fatalf("expected ErrProSubscriptionStillActive, got result=%v err=%v", result, err)
	}
	if result != nil {
		t.Fatalf("expected no checkout session result, got %+v", result)
	}
}
