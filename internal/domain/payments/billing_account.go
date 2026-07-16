package payments

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// BillingAccount associates a user with their Stripe customer.
//
// StripeCustomerID is stored as a string so the domain does not depend
// on Stripe SDK types.
type BillingAccount struct {
	ID               uuid.UUID `json:"id"`
	UserID           uuid.UUID `json:"user_id"`
	StripeCustomerID string    `json:"stripe_customer_id"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func NewBillingAccount(
	userID uuid.UUID,
	stripeCustomerID string,
) (*BillingAccount, error) {
	if userID == uuid.Nil {
		return nil, ErrUserIDRequired
	}

	stripeCustomerID = strings.TrimSpace(stripeCustomerID)
	if stripeCustomerID == "" {
		return nil, ErrStripeCustomerIDRequired
	}

	now := time.Now().UTC()

	return &BillingAccount{
		ID:               uuid.New(),
		UserID:           userID,
		StripeCustomerID: stripeCustomerID,
		CreatedAt:        now,
		UpdatedAt:        now,
	}, nil
}
