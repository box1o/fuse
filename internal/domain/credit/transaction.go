package credit

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type TransactionType string

const (
	TransactionTypeDeposit TransactionType = "deposit"
	TransactionTypeSpend   TransactionType = "spend"
	TransactionTypeRefund  TransactionType = "refund"
)

type TransactionSource string

const (
	TransactionSourcePurchase TransactionSource = "purchase"
	TransactionSourceUsage    TransactionSource = "usage"
)

type Transaction struct {
	ID        uuid.UUID         `json:"id"`
	AccountID uuid.UUID         `json:"account_id"`
	Type      TransactionType   `json:"type"`
	Source    TransactionSource `json:"source"`
	Amount    Amount            `json:"amount"`

	// ReferenceID points to the internal entity that caused the transaction.
	// Examples: payment ID, usage ID, refund ID.
	ReferenceID string `json:"reference_id,omitempty"`

	// ExternalReference stores a provider identifier when applicable.
	// Example: Stripe Checkout Session ID.
	ExternalReference string `json:"external_reference,omitempty"`

	// IdempotencyKey prevents the same operation from being processed twice.
	IdempotencyKey string `json:"idempotency_key"`

	CreatedAt time.Time `json:"created_at"`
}

type NewTransactionInput struct {
	AccountID         uuid.UUID
	Type              TransactionType
	Source            TransactionSource
	Amount            Amount
	ReferenceID       string
	ExternalReference string
	IdempotencyKey    string
}

func NewTransaction(input NewTransactionInput) (*Transaction, error) {
	if input.AccountID == uuid.Nil {
		return nil, ErrAccountIDRequired
	}

	if !input.Amount.IsPositive() {
		return nil, ErrAmountMustBePositive
	}

	if !input.Type.IsValid() {
		return nil, ErrInvalidTransactionType
	}

	if !input.Source.IsValid() {
		return nil, ErrInvalidTransactionSource
	}

	idempotencyKey := strings.TrimSpace(input.IdempotencyKey)
	if idempotencyKey == "" {
		return nil, ErrIdempotencyKeyRequired
	}

	if err := validateTransactionCombination(input.Type, input.Source); err != nil {
		return nil, err
	}

	return &Transaction{
		ID:                uuid.New(),
		AccountID:         input.AccountID,
		Type:              input.Type,
		Source:            input.Source,
		Amount:            input.Amount,
		ReferenceID:       strings.TrimSpace(input.ReferenceID),
		ExternalReference: strings.TrimSpace(input.ExternalReference),
		IdempotencyKey:    idempotencyKey,
		CreatedAt:         time.Now().UTC(),
	}, nil
}

func (transactionType TransactionType) IsValid() bool {
	switch transactionType {
	case TransactionTypeDeposit,
		TransactionTypeSpend,
		TransactionTypeRefund:
		return true
	default:
		return false
	}
}

func (transactionSource TransactionSource) IsValid() bool {
	switch transactionSource {
	case TransactionSourcePurchase,
		TransactionSourceUsage:
		return true
	default:
		return false
	}
}

func validateTransactionCombination(
	transactionType TransactionType,
	transactionSource TransactionSource,
) error {
	switch {
	case transactionType == TransactionTypeDeposit &&
		transactionSource == TransactionSourcePurchase:
		return nil

	case transactionType == TransactionTypeSpend &&
		transactionSource == TransactionSourceUsage:
		return nil

	case transactionType == TransactionTypeRefund &&
		transactionSource == TransactionSourceUsage:
		return nil

	default:
		return ErrInvalidTransactionCombination
	}
}

func (transaction *Transaction) IsCredit() bool {
	if transaction == nil {
		return false
	}

	switch transaction.Type {
	case TransactionTypeDeposit,
		TransactionTypeRefund:
		return true
	default:
		return false
	}
}

func (transaction *Transaction) IsDebit() bool {
	if transaction == nil {
		return false
	}

	return transaction.Type == TransactionTypeSpend
}
