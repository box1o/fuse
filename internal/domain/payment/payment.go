package payment

import (
	"fuse/internal/infrastructure/db"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusPending   Status = "pending"
	StatusCompleted Status = "completed"
	StatusFailed    Status = "failed"
	StatusCanceled  Status = "canceled"
)

func (status Status) IsValid() bool {
	switch status {
	case StatusPending,
		StatusCompleted,
		StatusFailed,
		StatusCanceled:
		return true
	default:
		return false
	}
}

func (status Status) IsFinal() bool {
	switch status {
	case StatusCompleted,
		StatusFailed,
		StatusCanceled:
		return true
	default:
		return false
	}
}

type Provider string

const (
	ProviderStripe Provider = "stripe"
)

func (provider Provider) IsValid() bool {
	return provider == ProviderStripe
}

type Payment struct {
	db.Model
	ID                uuid.UUID  `json:"id"`
	OwnerID           uuid.UUID  `json:"owner_id"`
	CreditPackID      uuid.UUID  `json:"credit_pack_id"`
	Credits           int64      `json:"credits"`
	Amount            int64      `json:"amount"`
	Currency          string     `json:"currency"`
	Status            Status     `json:"status"`
	Provider          Provider   `json:"provider"`
	ProviderSessionID string     `json:"provider_session_id,omitempty"`
	ProviderPaymentID string     `json:"provider_payment_id,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	CompletedAt       *time.Time `json:"completed_at,omitempty"`
}

type NewPaymentInput struct {
	OwnerID      uuid.UUID
	CreditPackID uuid.UUID
	Credits      int64
	Amount       int64
	Currency     string
	Provider     Provider
}

func NewPayment(input NewPaymentInput) (*Payment, error) {
	if input.OwnerID == uuid.Nil {
		return nil, ErrOwnerIDRequired
	}

	if input.CreditPackID == uuid.Nil {
		return nil, ErrCreditPackIDRequired
	}

	if input.Credits <= 0 {
		return nil, ErrCreditsMustBePositive
	}

	if input.Amount <= 0 {
		return nil, ErrAmountMustBePositive
	}

	currency, err := normalizeCurrency(input.Currency)
	if err != nil {
		return nil, err
	}

	if !input.Provider.IsValid() {
		return nil, ErrInvalidProvider
	}

	now := time.Now().UTC()

	return &Payment{
		ID:           uuid.New(),
		OwnerID:      input.OwnerID,
		CreditPackID: input.CreditPackID,
		Credits:      input.Credits,
		Amount:       input.Amount,
		Currency:     currency,
		Status:       StatusPending,
		Provider:     input.Provider,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

type RestorePaymentInput struct {
	ID                uuid.UUID
	OwnerID           uuid.UUID
	CreditPackID      uuid.UUID
	Credits           int64
	Amount            int64
	Currency          string
	Status            Status
	Provider          Provider
	ProviderSessionID string
	ProviderPaymentID string
	CreatedAt         time.Time
	UpdatedAt         time.Time
	CompletedAt       *time.Time
}

func RestorePayment(input RestorePaymentInput) (*Payment, error) {
	if input.ID == uuid.Nil {
		return nil, ErrPaymentNotFound
	}

	if input.OwnerID == uuid.Nil {
		return nil, ErrOwnerIDRequired
	}

	if input.CreditPackID == uuid.Nil {
		return nil, ErrCreditPackIDRequired
	}

	if input.Credits <= 0 {
		return nil, ErrCreditsMustBePositive
	}

	if input.Amount <= 0 {
		return nil, ErrAmountMustBePositive
	}

	currency, err := normalizeCurrency(input.Currency)
	if err != nil {
		return nil, err
	}

	if !input.Status.IsValid() {
		return nil, ErrInvalidStatus
	}

	if !input.Provider.IsValid() {
		return nil, ErrInvalidProvider
	}

	if input.CreatedAt.IsZero() || input.UpdatedAt.IsZero() {
		return nil, ErrPaymentNotFound
	}

	return &Payment{
		ID:                input.ID,
		OwnerID:           input.OwnerID,
		CreditPackID:      input.CreditPackID,
		Credits:           input.Credits,
		Amount:            input.Amount,
		Currency:          currency,
		Status:            input.Status,
		Provider:          input.Provider,
		ProviderSessionID: strings.TrimSpace(input.ProviderSessionID),
		ProviderPaymentID: strings.TrimSpace(input.ProviderPaymentID),
		CreatedAt:         input.CreatedAt,
		UpdatedAt:         input.UpdatedAt,
		CompletedAt:       input.CompletedAt,
	}, nil
}

func (payment *Payment) AttachProviderSession(providerSessionID string) error {
	if payment == nil {
		return ErrPaymentNotFound
	}

	if payment.Status != StatusPending {
		return ErrPaymentNotPending
	}

	providerSessionID = strings.TrimSpace(providerSessionID)
	if providerSessionID == "" {
		return ErrProviderSessionIDRequired
	}

	if payment.ProviderSessionID != "" {
		return ErrProviderSessionAlreadyAttached
	}

	payment.ProviderSessionID = providerSessionID
	payment.UpdatedAt = time.Now().UTC()

	return nil
}

func (payment *Payment) Complete(providerPaymentID string) error {
	if payment == nil {
		return ErrPaymentNotFound
	}

	if payment.Status != StatusPending {
		return ErrPaymentNotPending
	}

	now := time.Now().UTC()

	payment.Status = StatusCompleted
	payment.ProviderPaymentID = strings.TrimSpace(providerPaymentID)
	payment.CompletedAt = &now
	payment.UpdatedAt = now

	return nil
}

func (payment *Payment) Fail() error {
	if payment == nil {
		return ErrPaymentNotFound
	}

	if payment.Status != StatusPending {
		return ErrPaymentNotPending
	}

	payment.Status = StatusFailed
	payment.UpdatedAt = time.Now().UTC()

	return nil
}

func (payment *Payment) Cancel() error {
	if payment == nil {
		return ErrPaymentNotFound
	}

	if payment.Status != StatusPending {
		return ErrPaymentNotPending
	}

	payment.Status = StatusCanceled
	payment.UpdatedAt = time.Now().UTC()

	return nil
}

func normalizeCurrency(currency string) (string, error) {
	currency = strings.ToUpper(strings.TrimSpace(currency))

	if len(currency) != 3 {
		return "", ErrInvalidCurrency
	}

	for _, character := range currency {
		if character < 'A' || character > 'Z' {
			return "", ErrInvalidCurrency
		}
	}

	return currency, nil
}
