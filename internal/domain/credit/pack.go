package credit

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Pack struct {
	ID            uuid.UUID `json:"id"`
	Code          string    `json:"code"`
	Name          string    `json:"name"`
	Credits       Amount    `json:"credits"`
	Active        bool      `json:"active"`
	StripePriceID string    `json:"-"`
	PriceAmount   int64     `json:"price_amount"`
	Currency      string    `json:"currency"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewPack(input Pack) (*Pack, error) {
	code := strings.TrimSpace(input.Code)
	name := strings.TrimSpace(input.Name)
	stripePriceID := strings.TrimSpace(input.StripePriceID)
	currency := strings.ToUpper(strings.TrimSpace(input.Currency))

	if code == "" {
		return nil, ErrPackCodeRequired
	}

	if name == "" {
		return nil, ErrPackNameRequired
	}

	if stripePriceID == "" {
		return nil, ErrPackStripePriceIDRequired
	}

	if input.PriceAmount <= 0 {
		return nil, ErrPackPriceAmountInvalid
	}

	if len(currency) != 3 {
		return nil, ErrPackCurrencyInvalid
	}

	if !input.Credits.IsPositive() {
		return nil, ErrAmountMustBePositive
	}

	now := time.Now().UTC()

	return &Pack{
		ID:            uuid.New(),
		Code:          code,
		Name:          name,
		Credits:       input.Credits,
		StripePriceID: stripePriceID,
		PriceAmount:   input.PriceAmount,
		Currency:      currency,
		Active:        true,
		CreatedAt:     now,
		UpdatedAt:     now,
	}, nil
}

func (p *Pack) Activate() error {
	if p == nil {
		return ErrPackNotFound
	}

	p.Active = true
	p.UpdatedAt = time.Now().UTC()

	return nil
}

func (p *Pack) Deactivate() error {
	if p == nil {
		return ErrPackNotFound
	}

	p.Active = false
	p.UpdatedAt = time.Now().UTC()

	return nil
}

func (p *Pack) ChangeCreditAmount(amount Amount) error {
	if p == nil {
		return ErrPackNotFound
	}

	if !amount.IsPositive() {
		return ErrAmountMustBePositive
	}

	p.Credits = amount
	p.UpdatedAt = time.Now().UTC()

	return nil
}
