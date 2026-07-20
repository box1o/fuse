package stripe

import (
	"context"
	stdErrors "errors"
	"testing"

	paymentService "fuse/internal/services/payment"
	"fuse/pkg/config"
)

func TestConfigPriceCatalog_FindByPackCode(t *testing.T) {
	t.Parallel()

	catalog, err := NewConfigPriceCatalog(
		map[string]config.StripePriceConfig{
			"credits_500": {
				PriceID:  " price_test_500 ",
				Amount:   999,
				Currency: " usd ",
			},
		},
	)
	if err != nil {
		t.Fatalf(
			"NewConfigPriceCatalog() returned unexpected error: %v",
			err,
		)
	}

	price, err := catalog.FindByPackCode(
		context.Background(),
		" credits_500 ",
	)
	if err != nil {
		t.Fatalf(
			"FindByPackCode() returned unexpected error: %v",
			err,
		)
	}

	if price.Reference != "price_test_500" {
		t.Errorf(
			"expected reference %q, got %q",
			"price_test_500",
			price.Reference,
		)
	}

	if price.Amount != 999 {
		t.Errorf("expected amount 999, got %d", price.Amount)
	}

	if price.Currency != "USD" {
		t.Errorf(
			"expected currency USD, got %q",
			price.Currency,
		)
	}
}

func TestConfigPriceCatalog_ReturnsIndependentPriceCopies(
	t *testing.T,
) {
	t.Parallel()

	catalog, err := NewConfigPriceCatalog(
		map[string]config.StripePriceConfig{
			"credits_500": {
				PriceID:  "price_test_500",
				Amount:   999,
				Currency: "USD",
			},
		},
	)
	if err != nil {
		t.Fatalf("create catalog: %v", err)
	}

	firstPrice, err := catalog.FindByPackCode(
		context.Background(),
		"credits_500",
	)
	if err != nil {
		t.Fatalf("find first price: %v", err)
	}

	firstPrice.Amount = 1
	firstPrice.Currency = "EUR"

	secondPrice, err := catalog.FindByPackCode(
		context.Background(),
		"credits_500",
	)
	if err != nil {
		t.Fatalf("find second price: %v", err)
	}

	if secondPrice.Amount != 999 {
		t.Errorf(
			"expected immutable amount 999, got %d",
			secondPrice.Amount,
		)
	}

	if secondPrice.Currency != "USD" {
		t.Errorf(
			"expected immutable currency USD, got %q",
			secondPrice.Currency,
		)
	}
}

func TestConfigPriceCatalog_RejectsInvalidConfiguration(
	t *testing.T,
) {
	t.Parallel()

	tests := []struct {
		name   string
		prices map[string]config.StripePriceConfig
	}{
		{
			name:   "empty catalog",
			prices: nil,
		},
		{
			name: "empty pack code",
			prices: map[string]config.StripePriceConfig{
				" ": {
					PriceID:  "price_test_500",
					Amount:   999,
					Currency: "USD",
				},
			},
		},
		{
			name: "missing price ID",
			prices: map[string]config.StripePriceConfig{
				"credits_500": {
					Amount:   999,
					Currency: "USD",
				},
			},
		},
		{
			name: "invalid amount",
			prices: map[string]config.StripePriceConfig{
				"credits_500": {
					PriceID:  "price_test_500",
					Amount:   0,
					Currency: "USD",
				},
			},
		},
		{
			name: "invalid currency",
			prices: map[string]config.StripePriceConfig{
				"credits_500": {
					PriceID:  "price_test_500",
					Amount:   999,
					Currency: "US1",
				},
			},
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewConfigPriceCatalog(test.prices)
			if err == nil {
				t.Fatal("expected catalog construction error")
			}
		})
	}
}

func TestConfigPriceCatalog_ReturnsNotFoundForUnknownPack(
	t *testing.T,
) {
	t.Parallel()

	catalog, err := NewConfigPriceCatalog(
		map[string]config.StripePriceConfig{
			"credits_500": {
				PriceID:  "price_test_500",
				Amount:   999,
				Currency: "USD",
			},
		},
	)
	if err != nil {
		t.Fatalf("create catalog: %v", err)
	}

	_, err = catalog.FindByPackCode(
		context.Background(),
		"credits_unknown",
	)

	if !stdErrors.Is(err, paymentService.ErrPriceNotFound) {
		t.Fatalf(
			"expected price not found error, got %v",
			err,
		)
	}
}
