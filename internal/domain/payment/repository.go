package payment

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, payment *Payment) error
	FindByID(ctx context.Context, id uuid.UUID) (*Payment, error)
	FindByProviderSessionID(ctx context.Context, provider Provider, providerSessionID string) (*Payment, error)
	Update(ctx context.Context, payment *Payment) error
}
