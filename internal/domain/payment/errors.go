package payment

import "fuse/pkg/errors"

var (
	ErrPaymentNotFound                = errors.New("PAYMENT_NOT_FOUND", "payment was not found")
	ErrPaymentAlreadyExists           = errors.New("PAYMENT_ALREADY_EXISTS", "payment has already been created")
	ErrOwnerIDRequired                = errors.New("PAYMENT_OWNER_ID_REQUIRED", "payment owner ID is required")
	ErrCreditPackIDRequired           = errors.New("PAYMENT_CREDIT_PACK_ID_REQUIRED", "credit pack ID is required")
	ErrCreditsMustBePositive          = errors.New("PAYMENT_CREDITS_MUST_BE_POSITIVE", "payment credits must be greater than zero")
	ErrAmountMustBePositive           = errors.New("PAYMENT_AMOUNT_MUST_BE_POSITIVE", "payment amount must be greater than zero")
	ErrCurrencyRequired               = errors.New("PAYMENT_CURRENCY_REQUIRED", "payment currency is required")
	ErrInvalidCurrency                = errors.New("PAYMENT_INVALID_CURRENCY", "payment currency is invalid")
	ErrInvalidProvider                = errors.New("PAYMENT_INVALID_PROVIDER", "payment provider is invalid")
	ErrInvalidStatus                  = errors.New("PAYMENT_INVALID_STATUS", "payment status is invalid")
	ErrProviderSessionIDRequired      = errors.New("PAYMENT_PROVIDER_SESSION_ID_REQUIRED", "payment provider session ID is required")
	ErrProviderSessionAlreadyAttached = errors.New("PAYMENT_PROVIDER_SESSION_ALREADY_ATTACHED", "payment provider session has already been attached")
	ErrPaymentNotPending              = errors.New("PAYMENT_NOT_PENDING", "only a pending payment can be changed")
)
