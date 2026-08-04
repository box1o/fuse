package auth

import appErrors "fuse/pkg/errors"

var (
	ErrNotAuthenticated      = appErrors.New("CLI_NOT_AUTHENTICATED", "not authenticated")
	ErrAuthorizationPending  = appErrors.New("CLI_AUTHORIZATION_PENDING", "authorization pending")
	ErrAuthorizationSlowDown = appErrors.New("CLI_AUTHORIZATION_SLOW_DOWN", "authorization polling slowed down")
	ErrAuthorizationDenied   = appErrors.New("CLI_AUTHORIZATION_DENIED", "authorization denied")
	ErrAuthorizationExpired  = appErrors.New("CLI_AUTHORIZATION_EXPIRED", "authorization expired")
	ErrInvalidResponse       = appErrors.New("CLI_INVALID_AUTH_RESPONSE", "invalid authentication response")
)
