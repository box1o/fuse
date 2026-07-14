package payments

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type SubscriptionStatus string

const (
	SubscriptionStatusIncomplete        SubscriptionStatus = "incomplete"
	SubscriptionStatusIncompleteExpired SubscriptionStatus = "incomplete_expired"
	SubscriptionStatusTrialing          SubscriptionStatus = "trialing"
	SubscriptionStatusActive            SubscriptionStatus = "active"
	SubscriptionStatusPastDue           SubscriptionStatus = "past_due"
	SubscriptionStatusCanceled          SubscriptionStatus = "canceled"
	SubscriptionStatusUnpaid            SubscriptionStatus = "unpaid"
	SubscriptionStatusPaused            SubscriptionStatus = "paused"
)

func (s SubscriptionStatus) IsValid() bool {
	switch s {
	case SubscriptionStatusIncomplete,
		SubscriptionStatusIncompleteExpired,
		SubscriptionStatusTrialing,
		SubscriptionStatusActive,
		SubscriptionStatusPastDue,
		SubscriptionStatusCanceled,
		SubscriptionStatusUnpaid,
		SubscriptionStatusPaused:
		return true
	default:
		return false
	}
}

type Subscription struct {
	ID                   uuid.UUID          `json:"id"`
	WorkspaceID          uuid.UUID          `json:"workspace_id"`
	StripeSubscriptionID string             `json:"stripe_subscription_id"`
	Status               SubscriptionStatus `json:"status"`
	CurrentPeriodStart   time.Time          `json:"current_period_start"`
	CurrentPeriodEnd     time.Time          `json:"current_period_end"`
	CancelAtPeriodEnd    bool               `json:"cancel_at_period_end"`
	CreatedAt            time.Time          `json:"created_at"`
	UpdatedAt            time.Time          `json:"updated_at"`
}

func NewSubscription(
	workspaceID uuid.UUID,
	stripeSubscriptionID string,
	status SubscriptionStatus,
	currentPeriodStart time.Time,
	currentPeriodEnd time.Time,
	cancelAtPeriodEnd bool,
) (*Subscription, error) {
	if workspaceID == uuid.Nil {
		return nil, ErrWorkspaceIDRequired
	}

	stripeSubscriptionID = strings.TrimSpace(stripeSubscriptionID)
	if stripeSubscriptionID == "" {
		return nil, ErrStripeSubscriptionIDRequired
	}

	if !status.IsValid() {
		return nil, ErrInvalidSubscriptionStatus
	}

	if currentPeriodStart.IsZero() {
		return nil, ErrCurrentPeriodStartRequired
	}

	if currentPeriodEnd.IsZero() {
		return nil, ErrCurrentPeriodEndRequired
	}

	if !currentPeriodEnd.After(currentPeriodStart) {
		return nil, ErrInvalidSubscriptionPeriod
	}

	now := time.Now().UTC()

	return &Subscription{
		ID:                   uuid.New(),
		WorkspaceID:          workspaceID,
		StripeSubscriptionID: stripeSubscriptionID,
		Status:               status,
		CurrentPeriodStart:   currentPeriodStart.UTC(),
		CurrentPeriodEnd:     currentPeriodEnd.UTC(),
		CancelAtPeriodEnd:    cancelAtPeriodEnd,
		CreatedAt:            now,
		UpdatedAt:            now,
	}, nil
}

func (s *Subscription) IsActive() bool {
	if s == nil {
		return false
	}

	return s.Status == SubscriptionStatusActive ||
		s.Status == SubscriptionStatusTrialing
}

func (s *Subscription) Update(
	status SubscriptionStatus,
	currentPeriodStart time.Time,
	currentPeriodEnd time.Time,
	cancelAtPeriodEnd bool,
) error {
	if !status.IsValid() {
		return ErrInvalidSubscriptionStatus
	}

	if currentPeriodStart.IsZero() {
		return ErrCurrentPeriodStartRequired
	}

	if currentPeriodEnd.IsZero() {
		return ErrCurrentPeriodEndRequired
	}

	if !currentPeriodEnd.After(currentPeriodStart) {
		return ErrInvalidSubscriptionPeriod
	}

	s.Status = status
	s.CurrentPeriodStart = currentPeriodStart.UTC()
	s.CurrentPeriodEnd = currentPeriodEnd.UTC()
	s.CancelAtPeriodEnd = cancelAtPeriodEnd
	s.UpdatedAt = time.Now().UTC()

	return nil
}
