package deviceauthapi

import "time"

type CreateCodeRequest struct {
	ClientName string `json:"client_name"`
}

type CodeResponse struct {
	DeviceCode              string `json:"device_code"`
	UserCode                string `json:"user_code"`
	VerificationURI         string `json:"verification_uri"`
	VerificationURIComplete string `json:"verification_uri_complete"`
	ExpiresIn               int    `json:"expires_in"`
	Interval                int    `json:"interval"`
}

type TokenRequest struct {
	DeviceCode string `json:"device_code"`
}

type TokenResponse struct {
	AccessToken string    `json:"access_token"`
	TokenType   string    `json:"token_type"`
	ExpiresAt   time.Time `json:"expires_at"`
	OwnerID     string    `json:"owner_id"`
	OwnerName   string    `json:"owner_name"`
	OwnerEmail  string    `json:"owner_email"`
}

type StatusResponse struct {
	Authenticated  bool      `json:"authenticated"`
	OwnerID        string    `json:"owner_id"`
	CredentialName string    `json:"name"`
	OwnerName      string    `json:"owner_name"`
	OwnerEmail     string    `json:"owner_email"`
	ExpiresAt      time.Time `json:"expires_at"`
}
