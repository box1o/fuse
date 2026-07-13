package ports

import (
	"context"
	"time"

	"fuse/pkg/computeapi"
	"fuse/pkg/deviceauthapi"
)

type Settings struct {
	APIURL         string
	AppURL         string
	InstallationID string
}

type APIEndpoint interface {
	SetBaseURL(baseURL string) error
}

type InstallationStore interface{ EnsureInstallationID() (string, error) }

type CredentialStore interface {
	Save(token string) error
	Load() (string, error)
	Delete() error
}

type DeviceAuthGateway interface {
	CreateCode(ctx context.Context, clientName string) (deviceauthapi.CodeResponse, error)
	ExchangeCode(ctx context.Context, deviceCode string) (deviceauthapi.TokenResponse, bool, error)
	Status(ctx context.Context, token string) (deviceauthapi.StatusResponse, error)
	Logout(ctx context.Context, token string) error
}

type NodeGateway interface {
	Register(ctx context.Context, token string, request computeapi.RegisterNodeRequest) (computeapi.Node, bool, error)
	List(ctx context.Context, token string) ([]computeapi.Node, error)
	Get(ctx context.Context, token, nodeID string) (computeapi.Node, error)
	Update(ctx context.Context, token, nodeID string, request computeapi.UpdateNodeRequest) (computeapi.Node, error)
	Delete(ctx context.Context, token, nodeID string) error
}

type CapabilityDetector interface {
	Detect(ctx context.Context) (computeapi.Capabilities, error)
}

type Browser interface{ Open(url string) error }

type Waiter interface {
	Wait(ctx context.Context, duration time.Duration) error
}

type LoginPresenter interface {
	ShowDeviceCode(userCode, verificationURI string)
	ShowAuthenticated(ownerName, ownerEmail string, expiresAt time.Time)
}

type Output interface {
	Printf(format string, args ...any)
	Println(values ...any)
	JSON(value any) error
	Capabilities(value computeapi.Capabilities)
}

type Prompter interface{ Confirm(prompt string) bool }
