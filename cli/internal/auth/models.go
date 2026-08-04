package auth

import "time"

type DeviceCode struct {
	DeviceCode              string        `json:"-"`
	UserCode                string        `json:"user_code"`
	VerificationURI         string        `json:"verification_uri"`
	VerificationURIComplete string        `json:"verification_uri_complete,omitempty"`
	ExpiresIn               time.Duration `json:"expires_in"`
	PollInterval            time.Duration `json:"poll_interval"`
}

type Credential struct {
	AccessToken string    `json:"access_token"`
	TokenType   string    `json:"token_type"`
	ExpiresAt   time.Time `json:"expires_at"`
}

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Session struct {
	User       User       `json:"user"`
	Credential Credential `json:"credential"`
}

type Status struct {
	Authenticated bool      `json:"authenticated"`
	User          User      `json:"user"`
	ExpiresAt     time.Time `json:"expires_at"`
}
