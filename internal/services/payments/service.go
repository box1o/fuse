package payments

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	domain "fuse/internal/domain/payments"
	"fuse/pkg/config"

	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v83"
	checkoutsession "github.com/stripe/stripe-go/v83/checkout/session"
	"github.com/stripe/stripe-go/v83/customer"
	"github.com/stripe/stripe-go/v83/subscription"
	"github.com/stripe/stripe-go/v83/webhook"
)

type Service struct {
	cfg  *config.Config
	repo domain.Repository
}

type CheckoutSessionResult struct {
	SessionID string `json:"session_id"`
	URL       string `json:"url"`
}

func NewService(cfg *config.Config, repo domain.Repository) *Service {
	stripe.Key = cfg.Stripe.SecretKey

	return &Service{
		cfg:  cfg,
		repo: repo,
	}
}

func (s *Service) CreateCheckoutSession(
	ctx context.Context,
	userID uuid.UUID,
	successURL string,
	cancelURL string,
	priceID string,
) (*CheckoutSessionResult, error) {
	if userID == uuid.Nil {
		return nil, domain.ErrUserIDRequired
	}

	successURL = strings.TrimSpace(successURL)
	if successURL == "" {
		return nil, fmt.Errorf("success URL is required")
	}

	cancelURL = strings.TrimSpace(cancelURL)
	if cancelURL == "" {
		return nil, fmt.Errorf("cancel URL is required")
	}

	priceID = strings.TrimSpace(priceID)
	if priceID == "" {
		return nil, fmt.Errorf("price ID is required")
	}

	account, err := s.repo.FindBillingAccountByUserID(ctx, userID)
	if err != nil && !domain.ErrBillingAccountNotFound.Is(err) {
		return nil, err
	}

	customerID := ""
	if account != nil {
		customerID = account.StripeCustomerID
	}

	if customerID == "" {
		cust, err := customer.New(&stripe.CustomerParams{
			Metadata: map[string]string{
				"user_id": userID.String(),
			},
		})
		if err != nil {
			return nil, fmt.Errorf("create stripe customer: %w", err)
		}

		account, err = domain.NewBillingAccount(userID, cust.ID)
		if err != nil {
			return nil, err
		}
		if err := s.repo.CreateBillingAccount(ctx, account); err != nil && !domain.ErrBillingAccountAlreadyExists.Is(err) {
			return nil, err
		}
		customerID = cust.ID
	}

	params := &stripe.CheckoutSessionParams{
		Mode:       stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		SuccessURL: stripe.String(successURL),
		CancelURL:  stripe.String(cancelURL),
		Customer:   stripe.String(customerID),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(priceID),
				Quantity: stripe.Int64(1),
			},
		},
		Metadata: map[string]string{
			"user_id": userID.String(),
		},
	}

	session, err := checkoutsession.New(params)
	if err != nil {
		return nil, fmt.Errorf("create checkout session: %w", err)
	}

	return &CheckoutSessionResult{
		SessionID: session.ID,
		URL:       session.URL,
	}, nil
}

func (s *Service) RecordUsage(
	ctx context.Context,
	userID uuid.UUID,
	resourceType domain.ResourceType,
	quantity int64,
	occurredAt time.Time,
	idempotencyKey string,
) (*domain.UsageRecord, error) {
	record, err := domain.NewUsageRecord(userID, resourceType, quantity, occurredAt, idempotencyKey)
	if err != nil {
		return nil, err
	}
	if err := s.repo.CreateUsageRecord(ctx, record); err != nil {
		return nil, err
	}
	return record, nil
}

func (s *Service) CancelSubscription(ctx context.Context, userID uuid.UUID) error {
	sub, err := s.repo.FindSubscriptionByUserID(ctx, userID)
	if err != nil {
		return err
	}

	params := &stripe.SubscriptionParams{
		CancelAtPeriodEnd: stripe.Bool(true),
	}
	params.AddMetadata("user_id", userID.String())

	if _, err := subscription.Update(sub.StripeSubscriptionID, params); err != nil {
		return fmt.Errorf("cancel subscription: %w", err)
	}

	sub.CancelAtPeriodEnd = true
	sub.UpdatedAt = time.Now().UTC()
	return s.repo.UpsertSubscription(ctx, sub)
}

func (s *Service) HandleWebhook(ctx context.Context, payload []byte, signature string) error {
	signature = strings.TrimSpace(signature)
	if signature == "" {
		return fmt.Errorf("missing webhook signature")
	}

	event, err := webhook.ConstructEvent(payload, signature, s.cfg.Stripe.WebhookSecret)
	if err != nil {
		return fmt.Errorf("construct webhook event: %w", err)
	}

	processed, err := s.repo.WebhookEventExists(ctx, event.ID)
	if err != nil {
		return err
	}
	if processed {
		return nil
	}

	if err := s.repo.CreateWebhookEvent(ctx, &domain.WebhookEvent{
		StripeEventID: event.ID,
		EventType:     string(event.Type),
		ProcessedAt:   time.Now().UTC(),
	}); err != nil && !domain.ErrWebhookEventAlreadyProcessed.Is(err) {
		return err
	}

	switch event.Type {
	case stripe.EventTypeCustomerSubscriptionCreated,
		stripe.EventTypeCustomerSubscriptionUpdated,
		stripe.EventTypeCustomerSubscriptionDeleted:
		var sub stripe.Subscription
		if err := json.Unmarshal(event.Data.Raw, &sub); err != nil {
			return fmt.Errorf("decode subscription event: %w", err)
		}
		return s.upsertSubscriptionFromStripe(ctx, &sub)
	case stripe.EventTypeInvoicePaid:
		return nil
	case stripe.EventTypeInvoicePaymentFailed:
		return nil
	default:
		return nil
	}
}

func (s *Service) upsertSubscriptionFromStripe(ctx context.Context, sub *stripe.Subscription) error {
	if sub == nil {
		return domain.ErrInvalidSubscription
	}

	userIDStr := sub.Metadata["user_id"]
	userID, err := uuid.Parse(userIDStr)
	if err != nil || userID == uuid.Nil {
		return fmt.Errorf("subscription missing user metadata")
	}

	status := domain.SubscriptionStatus(sub.Status)

	account, err := s.repo.FindBillingAccountByUserID(ctx, userID)
	if err != nil && !domain.ErrBillingAccountNotFound.Is(err) {
		return err
	}
	if account == nil && sub.Customer != nil {
		account, err = s.repo.FindBillingAccountByStripeCustomerID(ctx, sub.Customer.ID)
		if err != nil {
			return err
		}
	}
	if account == nil {
		return domain.ErrBillingAccountNotFound
	}

	if len(sub.Items.Data) == 0 {
		return fmt.Errorf("subscription has no items")
	}

	currentPeriodStart := time.Unix(sub.Items.Data[0].CurrentPeriodStart, 0).UTC()
	currentPeriodEnd := time.Unix(sub.Items.Data[0].CurrentPeriodEnd, 0).UTC()

	agg, err := domain.NewSubscription(
		account.UserID,
		sub.ID,
		status,
		currentPeriodStart,
		currentPeriodEnd,
		sub.CancelAtPeriodEnd,
	)
	if err != nil {
		return err
	}

	return s.repo.UpsertSubscription(ctx, agg)
}
