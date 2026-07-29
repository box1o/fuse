package api

import (
	"context"
	"errors"
	"net/http"
	"time"

	"fuse/cli/internal/auth"
)

type AuthGateway struct {
	client *Client
}

func NewAuthGateway(client *Client) *AuthGateway {
	return &AuthGateway{client: client}
}

func (g *AuthGateway) CreateDeviceCode(ctx context.Context, clientName string) (auth.DeviceCode, error) {
	request := struct {
		ClientName string `json:"client_name"`
	}{ClientName: clientName}
	var response struct {
		DeviceCode              string `json:"device_code"`
		UserCode                string `json:"user_code"`
		VerificationURI         string `json:"verification_uri"`
		VerificationURIComplete string `json:"verification_uri_complete"`
		ExpiresIn               int    `json:"expires_in"`
		Interval                int    `json:"interval"`
	}

	if _, err := g.client.Do(ctx, http.MethodPost, "/auth/device/code", "", request, &response); err != nil {
		return auth.DeviceCode{}, err
	}
	if response.Interval == 0 {
		response.Interval = 5
	}

	return auth.DeviceCode{
		DeviceCode:              response.DeviceCode,
		UserCode:                response.UserCode,
		VerificationURI:         response.VerificationURI,
		VerificationURIComplete: response.VerificationURIComplete,
		ExpiresIn:               time.Duration(response.ExpiresIn) * time.Second,
		PollInterval:            time.Duration(response.Interval) * time.Second,
	}, nil
}

func (g *AuthGateway) ExchangeDeviceCode(ctx context.Context, deviceCode string) (auth.Session, error) {
	request := struct {
		DeviceCode string `json:"device_code"`
	}{DeviceCode: deviceCode}
	var response tokenResponse

	status, err := g.client.Do(ctx, http.MethodPost, "/auth/device/token", "", request, &response)
	if status == http.StatusAccepted {
		return auth.Session{}, auth.ErrAuthorizationPending
	}
	if err != nil {
		var responseErr *ResponseError
		if errors.As(err, &responseErr) {
			switch responseErr.Code {
			case "authorization_pending":
				return auth.Session{}, auth.ErrAuthorizationPending
			case "slow_down":
				return auth.Session{}, auth.ErrAuthorizationSlowDown
			case "access_denied", "authorization_declined":
				return auth.Session{}, auth.ErrAuthorizationDenied
			case "expired_token":
				return auth.Session{}, auth.ErrAuthorizationExpired
			}
		}
		switch status {
		case http.StatusForbidden:
			return auth.Session{}, auth.ErrAuthorizationDenied
		case http.StatusGone:
			return auth.Session{}, auth.ErrAuthorizationExpired
		default:
			return auth.Session{}, err
		}
	}

	return response.session(), nil
}

func (g *AuthGateway) Status(ctx context.Context, accessToken string) (auth.Status, error) {
	var response struct {
		Authenticated bool      `json:"authenticated"`
		OwnerID       string    `json:"owner_id"`
		OwnerName     string    `json:"owner_name"`
		OwnerEmail    string    `json:"owner_email"`
		ExpiresAt     time.Time `json:"expires_at"`
	}

	_, err := g.client.Do(ctx, http.MethodGet, "/auth/cli/status", accessToken, nil, &response)
	if err != nil {
		var responseErr *ResponseError
		if errors.As(err, &responseErr) && responseErr.StatusCode == http.StatusUnauthorized {
			return auth.Status{}, auth.ErrNotAuthenticated
		}
		return auth.Status{}, err
	}

	return auth.Status{
		Authenticated: response.Authenticated,
		User:          auth.User{ID: response.OwnerID, Name: response.OwnerName, Email: response.OwnerEmail},
		ExpiresAt:     response.ExpiresAt,
	}, nil
}

func (g *AuthGateway) Logout(ctx context.Context, accessToken string) error {
	_, err := g.client.Do(ctx, http.MethodPost, "/auth/cli/logout", accessToken, nil, nil)
	return err
}

type tokenResponse struct {
	AccessToken string    `json:"access_token"`
	TokenType   string    `json:"token_type"`
	ExpiresAt   time.Time `json:"expires_at"`
	OwnerID     string    `json:"owner_id"`
	OwnerName   string    `json:"owner_name"`
	OwnerEmail  string    `json:"owner_email"`
}

func (r tokenResponse) session() auth.Session {
	return auth.Session{
		User: auth.User{ID: r.OwnerID, Name: r.OwnerName, Email: r.OwnerEmail},
		Credential: auth.Credential{
			AccessToken: r.AccessToken,
			TokenType:   r.TokenType,
			ExpiresAt:   r.ExpiresAt,
		},
	}
}
