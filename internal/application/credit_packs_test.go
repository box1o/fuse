package application

import (
	"context"
	"testing"

	domainCredit "fuse/internal/domain/credit"
	"fuse/pkg/config"

	"github.com/google/uuid"
)

func TestSyncConfiguredCreditPacksCreatesSelectablePacks(t *testing.T) {
	repository := &packRepositoryStub{packs: map[string]*domainCredit.Pack{}}
	prices := map[string]config.StripePriceConfig{
		"credits_240":  {PriceID: "price_240", Amount: 500, Currency: "usd"},
		"credits_2400": {PriceID: "price_2400", Amount: 3650, Currency: "usd"},
		"credits_5000": {PriceID: "price_5000", Amount: 8550, Currency: "usd"},
	}

	if err := syncConfiguredCreditPacks(context.Background(), repository, prices); err != nil {
		t.Fatalf("sync configured credit packs: %v", err)
	}

	if len(repository.packs) != 3 {
		t.Fatalf("expected 3 credit packs, got %d", len(repository.packs))
	}

	pack := repository.packs["credits_2400"]
	if pack == nil {
		t.Fatal("expected credits_2400 pack")
	}
	if pack.Credits != 2400 || pack.PriceAmount != 3650 || pack.Currency != "USD" || !pack.Active {
		t.Fatalf("unexpected pack: %#v", pack)
	}
}

func TestSyncConfiguredCreditPacksRefreshesExistingPack(t *testing.T) {
	repository := &packRepositoryStub{packs: map[string]*domainCredit.Pack{
		"credits_240": {
			ID: uuid.New(), Code: "credits_240", Name: "Old", Credits: 1,
			StripePriceID: "old", PriceAmount: 1, Currency: "EUR", Active: false,
		},
	}}

	err := syncConfiguredCreditPacks(context.Background(), repository, map[string]config.StripePriceConfig{
		"credits_240": {PriceID: "price_240", Amount: 500, Currency: "usd"},
	})
	if err != nil {
		t.Fatalf("sync configured credit packs: %v", err)
	}

	pack := repository.packs["credits_240"]
	if pack.Credits != 240 || pack.StripePriceID != "price_240" || pack.PriceAmount != 500 || pack.Currency != "USD" || !pack.Active {
		t.Fatalf("unexpected refreshed pack: %#v", pack)
	}
}

type packRepositoryStub struct {
	packs map[string]*domainCredit.Pack
}

func (repository *packRepositoryStub) Create(_ context.Context, pack *domainCredit.Pack) error {
	repository.packs[pack.Code] = pack
	return nil
}

func (repository *packRepositoryStub) FindByID(_ context.Context, id uuid.UUID) (*domainCredit.Pack, error) {
	for _, pack := range repository.packs {
		if pack.ID == id {
			return pack, nil
		}
	}
	return nil, domainCredit.ErrPackNotFound
}

func (repository *packRepositoryStub) FindByCode(_ context.Context, code string) (*domainCredit.Pack, error) {
	pack, ok := repository.packs[code]
	if !ok {
		return nil, domainCredit.ErrPackNotFound
	}
	return pack, nil
}

func (repository *packRepositoryStub) ListActive(_ context.Context) ([]*domainCredit.Pack, error) {
	return nil, nil
}

func (repository *packRepositoryStub) Update(_ context.Context, pack *domainCredit.Pack) error {
	if _, ok := repository.packs[pack.Code]; !ok {
		return domainCredit.ErrPackNotFound
	}
	repository.packs[pack.Code] = pack
	return nil
}

var _ domainCredit.PackRepository = (*packRepositoryStub)(nil)
