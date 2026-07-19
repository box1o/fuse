package credit

import "fuse/pkg/errors"

var (
	// Account
	ErrNegativeAmount       = errors.New("NEGATIVE_AMOUNT", "amount cannot be negative")
	ErrInsufficientBalance  = errors.New("INSUFFICIENT_BALANCE", "insufficient balance")
	ErrAmountMustBePositive = errors.New("AMOUNT_MUST_BE_POSITIVE", "amount must be greater than zero")
	ErrInvalidAccount       = errors.New("INVALID_ACCOUNT", "account is invalid")
	ErrOwnerIDRequired      = errors.New("OWNER_ID_REQUIRED", "owner ID is required")
	ErrBalanceOverflow      = errors.New("BALANCE_OVERFLOW", "credit balance exceeds the supported maximum")
	ErrAccountNotFound      = errors.New("CREDIT_ACCOUNT_NOT_FOUND", "Credit account was not found")

	// Pack
	ErrPackCodeRequired = errors.New("PACK_CODE_REQUIRED", "credit pack code is required")
	ErrPackNameRequired = errors.New("PACK_NAME_REQUIRED", "credit pack name is required")
	ErrPackNotFound     = errors.New("PACK_NOT_FOUND", "credit pack was not found")

	// Transcation
	ErrAccountIDRequired             = errors.New("ACCOUNT_ID_REQUIRED", "account id is required")
	ErrInvalidTransaction            = errors.New("INVALID_TRANSACTION", "transaction is invalid")
	ErrInvalidTransactionType        = errors.New("INVALID_TRANSACTION_TYPE", "transaction type is invalid")
	ErrInvalidTransactionSource      = errors.New("INVALID_TRANSACTION_SOURCE", "transaction source is invalid")
	ErrIdempotencyKeyRequired        = errors.New("IDEMPOTENCY_KEY_REQUIRED", "idempotency key required")
	ErrInvalidTransactionCombination = errors.New("INVALID_TRANSACTION_COMBINATION", "transaction type and source combination is invalid")
	ErrTransactionAlreadyExists      = errors.New("CREDIT_TRANSACTION_ALREADY_EXISTS", "Credit transaction has already been processed")
)
