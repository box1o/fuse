package auth

import (
	"context"
	"time"
)

type Gateway interface {
	CreateDeviceCode(ctx context.Context, clientName string) (DeviceCode, error)
	ExchangeDeviceCode(ctx context.Context, deviceCode string) (Session, error)
	Status(ctx context.Context, accessToken string) (Status, error)
	Logout(ctx context.Context, accessToken string) error
}

type CredentialStore interface {
	Save(ctx context.Context, credential Credential) error
	Load(ctx context.Context) (Credential, error)
	Delete(ctx context.Context) error
}

type Browser interface {
	Open(url string) error
}

type Waiter interface {
	Wait(ctx context.Context, duration time.Duration) error
}

type LoginObserver interface {
	DeviceCodeReady(code DeviceCode)
	WaitingForAuthorization()
}
