package application

import (
	"context"
	"fmt"
	"time"

	"fuse/cli/internal/ports"
	"fuse/pkg/deviceauthapi"
)

type AuthService struct {
	gateway     ports.DeviceAuthGateway
	credentials ports.CredentialStore
	browser     ports.Browser
	waiter      ports.Waiter
	presenter   ports.LoginPresenter
}

func NewAuthService(gateway ports.DeviceAuthGateway, credentials ports.CredentialStore, browser ports.Browser, waiter ports.Waiter, presenter ports.LoginPresenter) *AuthService {
	return &AuthService{
		gateway:     gateway,
		credentials: credentials,
		browser:     browser,
		waiter:      waiter,
		presenter:   presenter,
	}
}

func (s *AuthService) Login(ctx context.Context) (deviceauthapi.TokenResponse, error) {
	challenge, err := s.gateway.CreateCode(ctx, "Fuse CLI")
	if err != nil {
		return deviceauthapi.TokenResponse{}, fmt.Errorf("start device login: %w", err)
	}
	s.presenter.ShowDeviceCode(challenge.UserCode, challenge.VerificationURI)
	if err := s.browser.Open(challenge.VerificationURIComplete); err != nil {
		s.presenter.ShowDeviceCode(challenge.UserCode, challenge.VerificationURIComplete)
	}

	interval := time.Duration(challenge.Interval) * time.Second
	if interval <= 0 {
		interval = 5 * time.Second
	}
	loginCtx, cancel := context.WithTimeout(ctx, time.Duration(challenge.ExpiresIn)*time.Second)
	defer cancel()

	for {
		if err := s.waiter.Wait(loginCtx, interval); err != nil {
			return deviceauthapi.TokenResponse{}, fmt.Errorf("device authorization expired: %w", err)
		}
		token, pending, err := s.gateway.ExchangeCode(loginCtx, challenge.DeviceCode)
		if pending {
			continue
		}
		if err != nil {
			return deviceauthapi.TokenResponse{}, fmt.Errorf("complete device login: %w", err)
		}
		if err := s.credentials.Save(token.AccessToken); err != nil {
			return deviceauthapi.TokenResponse{}, fmt.Errorf("store CLI credential: %w", err)
		}
		s.presenter.ShowAuthenticated(token.OwnerName, token.OwnerEmail, token.ExpiresAt)
		return token, nil
	}
}

func (s *AuthService) Status(ctx context.Context) (deviceauthapi.StatusResponse, error) {
	token, err := s.credentials.Load()
	if err != nil {
		return deviceauthapi.StatusResponse{}, fmt.Errorf("not logged in")
	}
	return s.gateway.Status(ctx, token)
}

func (s *AuthService) Logout(ctx context.Context) error {
	token, err := s.credentials.Load()
	if err == nil {
		_ = s.gateway.Logout(ctx, token)
	}
	if err := s.credentials.Delete(); err != nil {
		return fmt.Errorf("delete local credential: %w", err)
	}
	return nil
}
