package application

import (
	"context"
	stdErrors "errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	domainCredit "fuse/internal/domain/credit"
	"fuse/pkg/config"
)

func syncConfiguredCreditPacks(
	ctx context.Context,
	repository domainCredit.PackRepository,
	prices map[string]config.StripePriceConfig,
) error {
	codes := make([]string, 0, len(prices))
	for code := range prices {
		codes = append(codes, code)
	}
	sort.Strings(codes)

	for _, code := range codes {
		price := prices[code]
		credits, err := creditsFromPackCode(code)
		if err != nil {
			return err
		}

		pack, err := repository.FindByCode(ctx, code)
		if err != nil && !stdErrors.Is(err, domainCredit.ErrPackNotFound) {
			return fmt.Errorf("find configured credit pack %q: %w", code, err)
		}

		if stdErrors.Is(err, domainCredit.ErrPackNotFound) {
			pack, err = domainCredit.NewPack(domainCredit.Pack{
				Code:          code,
				Name:          fmt.Sprintf("%d Credits", credits),
				Credits:       domainCredit.Amount(credits),
				StripePriceID: price.PriceID,
				PriceAmount:   price.Amount,
				Currency:      price.Currency,
			})
			if err != nil {
				return fmt.Errorf("build configured credit pack %q: %w", code, err)
			}

			if err := repository.Create(ctx, pack); err != nil {
				return fmt.Errorf("create configured credit pack %q: %w", code, err)
			}
			continue
		}

		pack.Name = fmt.Sprintf("%d Credits", credits)
		pack.Credits = domainCredit.Amount(credits)
		pack.StripePriceID = strings.TrimSpace(price.PriceID)
		pack.PriceAmount = price.Amount
		pack.Currency = strings.ToUpper(strings.TrimSpace(price.Currency))
		pack.Active = true

		if err := repository.Update(ctx, pack); err != nil {
			return fmt.Errorf("update configured credit pack %q: %w", code, err)
		}
	}

	return nil
}

func creditsFromPackCode(code string) (int64, error) {
	value := strings.TrimPrefix(strings.TrimSpace(code), "credits_")
	credits, err := strconv.ParseInt(value, 10, 64)
	if err != nil || credits <= 0 {
		return 0, fmt.Errorf("invalid credit pack code %q", code)
	}
	return credits, nil
}
