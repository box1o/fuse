package models

import (
	"time"

	domainCLI "fuse/internal/domain/cli"
	"fuse/internal/infrastructure/db"

	"github.com/google/uuid"
)

type DBCredential struct {
	db.Model
	OwnerID    uuid.UUID `gorm:"type:uuid;not null;index"`
	Name       string    `gorm:"not null;size:100"`
	TokenHash  string    `gorm:"not null;uniqueIndex;size:64"`
	ExpiresAt  time.Time `gorm:"not null;index"`
	LastUsedAt *time.Time
	RevokedAt  *time.Time `gorm:"index"`
}

func (DBCredential) TableName() string {
	return "cli_credentials"
}

func FromDomain(credential *domainCLI.Credential) *DBCredential {
	return &DBCredential{
		Model: db.Model{
			ID:        credential.ID,
			CreatedAt: credential.CreatedAt,
			UpdatedAt: credential.UpdatedAt,
		},
		OwnerID:    credential.OwnerID,
		Name:       credential.Name,
		TokenHash:  credential.TokenHash,
		ExpiresAt:  credential.ExpiresAt,
		LastUsedAt: credential.LastUsedAt,
		RevokedAt:  credential.RevokedAt,
	}
}

func (credential *DBCredential) ToDomain() *domainCLI.Credential {
	return &domainCLI.Credential{
		ID:         credential.ID,
		OwnerID:    credential.OwnerID,
		Name:       credential.Name,
		TokenHash:  credential.TokenHash,
		ExpiresAt:  credential.ExpiresAt,
		LastUsedAt: credential.LastUsedAt,
		RevokedAt:  credential.RevokedAt,
		CreatedAt:  credential.CreatedAt,
		UpdatedAt:  credential.UpdatedAt,
	}
}
