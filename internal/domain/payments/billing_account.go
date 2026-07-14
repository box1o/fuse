package payments

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// BillingAccount associates a workspace with its Stripe customer.
//
// StripeCustomerID is stored as a string so the domain does not depend
// on Stripe SDK types.
type BillingAccount struct {
	ID               uuid.UUID `json:"id"`
	WorkspaceID      uuid.UUID `json:"workspace_id"`
	StripeCustomerID string    `json:"stripe_customer_id"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func NewBillingAccount(
	workspaceID uuid.UUID,
	stripeCustomerID string,
) (*BillingAccount, error) {
	if workspaceID == uuid.Nil {
		return nil, ErrWorkspaceIDRequired
	}

	stripeCustomerID = strings.TrimSpace(stripeCustomerID)
	if stripeCustomerID == "" {
		return nil, ErrStripeCustomerIDRequired
	}

	now := time.Now().UTC()

	return &BillingAccount{
		ID:               uuid.New(),
		WorkspaceID:      workspaceID,
		StripeCustomerID: stripeCustomerID,
		CreatedAt:        now,
		UpdatedAt:        now,
	}, nil
}
