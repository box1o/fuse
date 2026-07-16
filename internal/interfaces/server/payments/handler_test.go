package payments

import (
	"testing"

	"fuse/pkg/config"
)

func TestResolvePriceID_UsesProPrice(t *testing.T) {
	h := &Handler{
		cfg: &config.Config{
			Stripe: config.StripeConfig{
				CPUPriceID: "price_cpu",
				GPUPriceID: "price_gpu",
				NPUPriceID: "price_npu",
				ProPriceID: "price_pro",
			},
		},
	}

	priceID, err := h.resolvePriceID("cpu")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if priceID != "price_pro" {
		t.Fatalf("expected pro subscription price, got %s", priceID)
	}
}

func TestResolvePriceID_RejectsPlaceholderProPrice(t *testing.T) {
	h := &Handler{
		cfg: &config.Config{
			Stripe: config.StripeConfig{
				ProPriceID: "price_pro_placeholder",
			},
		},
	}

	priceID, err := h.resolvePriceID("pro")
	if err == nil {
		t.Fatalf("expected error, got price id %s", priceID)
	}
}
