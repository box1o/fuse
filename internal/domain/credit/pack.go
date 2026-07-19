package credit

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Pack struct {
	ID      uuid.UUID `json:"id"`
	Code    string    `json:"code"`
	Name    string    `json:"name"`
	Credits Amount    `json:"credits"`
	Active  bool      `json:"active"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewPack(code, name string, credits Amount) (*Pack, error) {
	code = strings.TrimSpace(code)
	name = strings.TrimSpace(name)

	if code == "" {
		return nil, ErrPackCodeRequired
	}

	if name == "" {
		return nil, ErrPackNameRequired
	}

	if !credits.IsPositive() {
		return nil, ErrAmountMustBePositive
	}

	now := time.Now().UTC()

	return &Pack{
		ID:        uuid.New(),
		Code:      code,
		Name:      name,
		Credits:   credits,
		Active:    true,
		CreatedAt: now,
		UpdatedAt: now,
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
