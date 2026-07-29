package auth

import (
	"context"
	"errors"
	"fmt"
	"time"
)

func (s *Service) Login(ctx context.Context, options LoginOptions) (Session, error) {
	code, err := s.gateway.CreateDeviceCode(ctx, options.clientName())
	if err != nil {
		return Session{}, fmt.Errorf("create device authorization: %w", err)
	}

	if code.DeviceCode == "" ||
		code.UserCode == "" ||
		code.VerificationURI == "" ||
		code.ExpiresIn <= 0 ||
		code.PollInterval <= 0 {
		return Session{}, ErrInvalidResponse
	}

	if options.Observer != nil {
		options.Observer.DeviceCodeReady(code)
	}

	if options.OpenBrowser {
		url := code.VerificationURIComplete
		if url == "" {
			url = code.VerificationURI
		}

		_ = s.browser.Open(url)
	}

	loginCtx, cancel := context.WithTimeout(ctx, code.ExpiresIn)
	defer cancel()

	if options.Observer != nil {
		options.Observer.WaitingForAuthorization()
	}

	pollInterval := code.PollInterval
	for {
		if err := s.waiter.Wait(loginCtx, pollInterval); err != nil {
			if ctx.Err() != nil {
				return Session{}, ctx.Err()
			}

			return Session{}, fmt.Errorf("%w: %v", ErrAuthorizationExpired, err)
		}

		session, err := s.gateway.ExchangeDeviceCode(loginCtx, code.DeviceCode)
		if errors.Is(err, ErrAuthorizationPending) {
			continue
		}
		if errors.Is(err, ErrAuthorizationSlowDown) {
			pollInterval += 5 * time.Second
			continue
		}
		if err != nil {
			return Session{}, fmt.Errorf("exchange device authorization: %w", err)
		}

		if session.Credential.AccessToken == "" {
			return Session{}, ErrInvalidResponse
		}

		if err := s.credentials.Save(ctx, session.Credential); err != nil {
			return Session{}, fmt.Errorf("save credential: %w", err)
		}

		return session, nil
	}
}
