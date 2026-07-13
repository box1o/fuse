package compute

import (
	"time"

	"github.com/google/uuid"
)

type CLICredential struct {
	ID         uuid.UUID  `json:"id"`
	OwnerID    uuid.UUID  `json:"owner_id"`
	Name       string     `json:"name"`
	TokenHash  string     `json:"-"`
	ExpiresAt  time.Time  `json:"expires_at"`
	LastUsedAt *time.Time `json:"last_used_at,omitempty"`
	RevokedAt  *time.Time `json:"revoked_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

func NewCLICredential(ownerID uuid.UUID, name, tokenHash string, expiresAt time.Time) (*CLICredential, error) {
	if ownerID == uuid.Nil || tokenHash == "" || !expiresAt.After(time.Now().UTC()) {
		return nil, ErrInvalidCredential
	}
	now := time.Now().UTC()
	return &CLICredential{
		ID: uuid.New(), OwnerID: ownerID, Name: name, TokenHash: tokenHash,
		ExpiresAt: expiresAt.UTC(), CreatedAt: now, UpdatedAt: now,
	}, nil
}

func (c *CLICredential) IsActive(now time.Time) bool {
	return c != nil && c.RevokedAt == nil && now.Before(c.ExpiresAt)
}
