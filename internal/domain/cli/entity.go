package cli

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

const CredentialDuration = 90 * 24 * time.Hour

type Credential struct {
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

func NewCredential(ownerID uuid.UUID, name string, tokenHash string) (*Credential, error) {
	if ownerID == uuid.Nil {
		return nil, ErrOwnerRequired
	}

	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrNameRequired
	}

	tokenHash = strings.TrimSpace(tokenHash)
	if tokenHash == "" {
		return nil, ErrTokenHashRequired
	}

	now := time.Now().UTC()
	return &Credential{
		ID:        uuid.New(),
		OwnerID:   ownerID,
		Name:      name,
		TokenHash: tokenHash,
		ExpiresAt: now.Add(CredentialDuration),
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (credential *Credential) Validate(now time.Time) error {
	if credential == nil || credential.ID == uuid.Nil {
		return ErrCredentialNotFound
	}

	if credential.RevokedAt != nil {
		return ErrCredentialRevoked
	}

	if !now.UTC().Before(credential.ExpiresAt) {
		return ErrCredentialExpired
	}

	return nil
}

func (credential *Credential) MarkUsed(now time.Time) {
	usedAt := now.UTC()
	credential.LastUsedAt = &usedAt
	credential.UpdatedAt = usedAt
}

func (credential *Credential) Revoke(now time.Time) {
	revokedAt := now.UTC()
	credential.RevokedAt = &revokedAt
	credential.UpdatedAt = revokedAt
}
