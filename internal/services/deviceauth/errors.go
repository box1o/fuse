package deviceauth

import "fuse/pkg/errors"

var (
	ErrInvalidCode          = errors.New("INVALID_DEVICE_CODE", "device code is invalid")
	ErrAuthorizationPending = errors.New("DEVICE_AUTHORIZATION_PENDING", "authorization is pending")
	ErrAuthorizationDenied  = errors.New("DEVICE_AUTHORIZATION_DENIED", "authorization was denied")
	ErrAuthorizationExpired = errors.New("DEVICE_AUTHORIZATION_EXPIRED", "device authorization expired")
	ErrAlreadyHandled       = errors.New("DEVICE_AUTHORIZATION_CONFLICT", "device authorization was already handled")
)
