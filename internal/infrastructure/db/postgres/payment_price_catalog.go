package postgres

import (
	"context"
	"strings"

	domainCredit "fuse/internal/domain/credit"
	paymentService "fuse/internal/services/payment"
)

type PaymentPriceCatalog struct {
	packs domainCredit.PackRepository
}

var _ paymentService.PriceCatalog = (*PaymentPriceCatalog)(nil)

func NewPaymentPriceCatalog(packs domainCredit.PackRepository) *PaymentPriceCatalog {
	return &PaymentPriceCatalog{
		packs: packs,
	}
}

func (catalog *PaymentPriceCatalog) FindByPackCode(ctx context.Context, packCode string) (*paymentService.Price, error) {
	normalizedPackCode := strings.TrimSpace(packCode)
	if normalizedPackCode == "" {
		return nil, paymentService.ErrPriceNotFound
	}

	pack, err := catalog.packs.FindByCode(ctx, normalizedPackCode)
	if err != nil {
		return nil, err
	}

	if !pack.Active {
		return nil, paymentService.ErrPriceNotFound
	}

	price := &paymentService.Price{
		Reference: strings.TrimSpace(pack.StripePriceID),
		Amount:    pack.PriceAmount,
		Currency: strings.ToUpper(
			strings.TrimSpace(pack.Currency),
		),
	}

	if price.Reference == "" {
		return nil, paymentService.ErrPriceReferenceRequired
	}

	if price.Amount <= 0 || len(price.Currency) != 3 {
		return nil, paymentService.ErrInvalidPrice
	}

	return price, nil
}
