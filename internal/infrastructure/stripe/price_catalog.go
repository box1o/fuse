package stripe

import (
	"context"
	"strings"

	paymentService "fuse/internal/services/payment"
	"fuse/pkg/config"
)

type ConfigPriceCatalog struct {
	prices map[string]paymentService.Price
}

var _ paymentService.PriceCatalog = (*ConfigPriceCatalog)(nil)

func NewConfigPriceCatalog(prices map[string]config.StripePriceConfig) (*ConfigPriceCatalog, error) {
	if len(prices) == 0 {
		return nil, paymentService.ErrPriceNotFound
	}

	catalogPrices := make(
		map[string]paymentService.Price,
		len(prices),
	)

	for packCode, configuredPrice := range prices {
		normalizedPackCode := strings.TrimSpace(packCode)
		if normalizedPackCode == "" {
			return nil, paymentService.ErrInvalidPrice
		}

		price := paymentService.Price{
			Reference: strings.TrimSpace(configuredPrice.PriceID),
			Amount:    configuredPrice.Amount,
			Currency: strings.ToUpper(
				strings.TrimSpace(configuredPrice.Currency),
			),
		}

		if err := validateConfiguredPrice(price); err != nil {
			return nil, err
		}

		catalogPrices[normalizedPackCode] = price
	}

	return &ConfigPriceCatalog{
		prices: catalogPrices,
	}, nil
}

func (catalog *ConfigPriceCatalog) FindByPackCode(_ context.Context, packCode string) (*paymentService.Price, error) {
	packCode = strings.TrimSpace(packCode)
	if packCode == "" {
		return nil, paymentService.ErrPriceNotFound
	}

	price, found := catalog.prices[packCode]
	if !found {
		return nil, paymentService.ErrPriceNotFound
	}

	// Return a copy so callers cannot mutate the catalog configuration.
	return &paymentService.Price{
		Reference: price.Reference,
		Amount:    price.Amount,
		Currency:  price.Currency,
	}, nil
}

func validateConfiguredPrice(price paymentService.Price) error {
	if strings.TrimSpace(price.Reference) == "" {
		return paymentService.ErrPriceReferenceRequired
	}

	if price.Amount <= 0 {
		return paymentService.ErrInvalidPrice
	}

	if !isISOStyleCurrency(price.Currency) {
		return paymentService.ErrInvalidPrice
	}

	return nil
}

func isISOStyleCurrency(currency string) bool {
	if len(currency) != 3 {
		return false
	}

	for _, character := range currency {
		if character < 'A' || character > 'Z' {
			return false
		}
	}

	return true
}
