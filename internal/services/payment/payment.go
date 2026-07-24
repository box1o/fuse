package payment

import (
	"context"

	domainPayment "fuse/internal/domain/payment"

	"github.com/google/uuid"
)

type Provider interface {
	CreateCheckoutSession(
		ctx context.Context,
		input CreateCheckoutSessionInput,
	) (*CheckoutSession, error)
}

type CreateCheckoutSessionInput struct {
	PaymentID      uuid.UUID
	OwnerID        uuid.UUID
	CreditPackID   uuid.UUID
	PriceReference string
	SuccessURL     string
	CancelURL      string
}

type CheckoutSession struct {
	Provider          domainPayment.Provider
	SessionID         string
	URL               string
	ProviderPaymentID string
	Amount            int64
	Currency          string
}
