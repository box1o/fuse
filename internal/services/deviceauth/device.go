package deviceauth

import (
	"context"
	"strings"
	"time"

	domainCLI "fuse/internal/domain/cli"
	"fuse/pkg/deviceauthapi"

	"github.com/google/uuid"
)

type deviceState struct {
	UserCode     string    `json:"user_code"`
	ClientName   string    `json:"client_name"`
	Status       string    `json:"status"`
	OwnerID      uuid.UUID `json:"owner_id,omitempty"`
	CredentialID uuid.UUID `json:"credential_id,omitempty"`
	AccessToken  string    `json:"access_token,omitempty"`
	ExpiresAt    time.Time `json:"expires_at"`
}

func (service *Service) Create(ctx context.Context, clientName string) (*deviceauthapi.CodeResponse, error) {
	clientName = strings.TrimSpace(clientName)
	if clientName == "" {
		clientName = "Fuse CLI"
	}
	if len(clientName) > 100 {
		return nil, ErrInvalidCode.WithDetail("client_name must not exceed 100 characters")
	}

	deviceCode, err := randomToken()
	if err != nil {
		return nil, err
	}

	userCode, err := service.uniqueUserCode(ctx)
	if err != nil {
		return nil, err
	}

	deviceHash := hash(deviceCode)
	state := deviceState{
		UserCode:   userCode,
		ClientName: clientName,
		Status:     "pending",
		ExpiresAt:  time.Now().UTC().Add(deviceTTL),
	}

	if err := service.saveState(ctx, deviceHash, state, deviceTTL); err != nil {
		return nil, err
	}
	if err := service.redis.Set(ctx, userKey(userCode), deviceHash, deviceTTL); err != nil {
		_ = service.redis.Delete(ctx, deviceKey(deviceHash))
		return nil, err
	}

	verificationURI := service.frontendURL + "/device"
	return &deviceauthapi.CodeResponse{
		DeviceCode:              deviceCode,
		UserCode:                userCode,
		VerificationURI:         verificationURI,
		VerificationURIComplete: verificationURI + "?code=" + userCode,
		ExpiresIn:               int(deviceTTL.Seconds()),
		Interval:                pollInterval,
	}, nil
}

func (service *Service) Inspect(ctx context.Context, userCode string) (*deviceauthapi.DeviceRequest, error) {
	deviceHash, err := service.redis.Get(ctx, userKey(userCode))
	if err != nil {
		return nil, ErrAuthorizationExpired
	}

	state, err := service.loadState(ctx, deviceHash)
	if err != nil {
		return nil, err
	}
	if state.Status != "pending" {
		return nil, ErrAlreadyHandled
	}

	return &deviceauthapi.DeviceRequest{
		UserCode:   state.UserCode,
		ClientName: state.ClientName,
		Status:     state.Status,
		ExpiresAt:  state.ExpiresAt,
	}, nil
}

func (service *Service) Approve(ctx context.Context, userCode string, ownerID uuid.UUID) error {
	if ownerID == uuid.Nil {
		return ErrInvalidCode.WithDetail("authenticated owner is required")
	}

	release, err := service.acquireDecisionLock(ctx, userCode)
	if err != nil {
		return err
	}
	defer release()

	deviceHash, err := service.redis.Get(ctx, userKey(userCode))
	if err != nil {
		return ErrAuthorizationExpired
	}

	state, err := service.loadState(ctx, deviceHash)
	if err != nil {
		return err
	}
	if state.Status == "approved" && state.OwnerID == ownerID {
		return nil
	}
	if state.Status != "pending" {
		return ErrAlreadyHandled
	}

	accessToken, err := randomToken()
	if err != nil {
		return err
	}

	credential, err := domainCLI.NewCredential(ownerID, state.ClientName, hash(accessToken))
	if err != nil {
		return err
	}
	if err := service.credentials.Create(ctx, credential); err != nil {
		return err
	}

	state.Status = "approved"
	state.OwnerID = ownerID
	state.CredentialID = credential.ID
	state.AccessToken = accessToken

	remaining := time.Until(state.ExpiresAt)
	if remaining <= 0 {
		service.revokeCredential(ctx, credential)
		return ErrAuthorizationExpired
	}
	if err := service.saveState(ctx, deviceHash, *state, remaining); err != nil {
		service.revokeCredential(ctx, credential)
		return err
	}

	return nil
}

func (service *Service) Deny(ctx context.Context, userCode string, ownerID uuid.UUID) error {
	if ownerID == uuid.Nil {
		return ErrInvalidCode.WithDetail("authenticated owner is required")
	}

	release, err := service.acquireDecisionLock(ctx, userCode)
	if err != nil {
		return err
	}
	defer release()

	deviceHash, err := service.redis.Get(ctx, userKey(userCode))
	if err != nil {
		return ErrAuthorizationExpired
	}

	state, err := service.loadState(ctx, deviceHash)
	if err != nil {
		return err
	}
	if state.Status != "pending" {
		return ErrAlreadyHandled
	}

	state.Status = "denied"
	state.OwnerID = ownerID

	remaining := time.Until(state.ExpiresAt)
	if remaining <= 0 {
		return ErrAuthorizationExpired
	}

	return service.saveState(ctx, deviceHash, *state, remaining)
}

func (service *Service) Exchange(ctx context.Context, deviceCode string) (*deviceauthapi.TokenResponse, error) {
	deviceCode = strings.TrimSpace(deviceCode)
	if deviceCode == "" {
		return nil, ErrInvalidCode
	}

	state, err := service.loadState(ctx, hash(deviceCode))
	if err != nil {
		return nil, err
	}

	switch state.Status {
	case "pending":
		return nil, ErrAuthorizationPending
	case "denied":
		return nil, ErrAuthorizationDenied
	case "approved":
	default:
		return nil, ErrInvalidCode
	}

	if state.AccessToken == "" || state.CredentialID == uuid.Nil || state.OwnerID == uuid.Nil {
		return nil, ErrInvalidCode
	}

	credential, err := service.credentials.FindByID(ctx, state.CredentialID)
	if err != nil {
		return nil, err
	}
	if err := credential.Validate(time.Now()); err != nil {
		return nil, err
	}

	owner, err := service.users.FindByID(ctx, state.OwnerID)
	if err != nil {
		return nil, err
	}

	return &deviceauthapi.TokenResponse{
		AccessToken: state.AccessToken,
		TokenType:   "Bearer",
		ExpiresAt:   credential.ExpiresAt,
		OwnerID:     owner.ID.String(),
		OwnerName:   owner.Name,
		OwnerEmail:  owner.Email,
	}, nil
}
