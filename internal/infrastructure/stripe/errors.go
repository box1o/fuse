package stripe

import "fuse/pkg/errors"

var (
	ErrSecretKeyRequired             = errors.New("SECRET_KEY_REQUIRED", "Stripe secret key is required")
	ErrCheckoutSessionCreationFailed = errors.New("CHECKOUT_SESSION_CREATION_FAILED", "Stripe checkout session creation failed")
)
