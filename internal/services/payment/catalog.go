package payment

import (
	"context"
	domainCredit "fuse/internal/domain/credit"

	"github.com/google/uuid"
)

type CreditPackReader interface {
	GetActivePack(ctx context.Context, packID uuid.UUID) (*domainCredit.Pack, error)
}

type PriceCatalog interface {
	FindByPackCode(ctx context.Context, packCode string) (*Price, error)
}

type Price struct {
	Reference string
	Amount    int64
	Currency  string
}
