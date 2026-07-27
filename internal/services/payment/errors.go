package payment

import "fuse/pkg/errors"

var (
	ErrProviderUnavailable           = errors.New("PAYMENT_PROVIDER_UNAVAILABLE", "payment provider is unavailable")
	ErrCheckoutSessionCreationFailed = errors.New("CHECKOUT_SESSION_CREATION_FAILED", "checkout session could not be created")
	ErrPriceReferenceRequired        = errors.New("PAYMENT_PRICE_REFERENCE_REQUIRED", "payment price reference is required")
	ErrSuccessURLRequired            = errors.New("PAYMENT_SUCCESS_URL_REQUIRED", "payment success URL is required")
	ErrCancelURLRequired             = errors.New("PAYMENT_CANCEL_URL_REQUIRED", "payment cancel URL is required")
	ErrPriceNotFound                 = errors.New("PAYMENT_PRICE_NOT_FOUND", "payment price was not found")
	ErrInvalidPrice                  = errors.New("PAYMENT_INVALID_PRICE", "payment price configuration is invalid")

	ErrCheckoutAmountMismatch   = errors.New("PAYMENT_CHECKOUT_AMOUNT_MISMATCH", "checkout amount does not match the configured payment amount")
	ErrCheckoutCurrencyMismatch = errors.New("PAYMENT_CHECKOUT_CURRENCY_MISMATCH", "checkout currency does not match the configured payment currency")

	ErrWebhookSignatureInvalid = errors.New("PAYMENT_WEBHOOK_SIGNATURE_INVALID", "payment webhook signature is invalid")
	ErrWebhookSecretRequired   = errors.New("PAYMENT_WEBHOOK_SECRET_REQUIRED", "payment webhook secret is required")
	ErrInvalidWebhookEvent     = errors.New("PAYMENT_INVALID_WEBHOOK_EVENT", "payment webhook event is invalid")
)
