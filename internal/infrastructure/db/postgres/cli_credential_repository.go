package postgres

import (
	"context"
	"errors"
	"strings"

	domainCLI "fuse/internal/domain/cli"
	cliModels "fuse/internal/domain/cli/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CLICredentialRepository struct {
	db *gorm.DB
}

var _ domainCLI.Repository = (*CLICredentialRepository)(nil)

func NewCLICredentialRepository(db *gorm.DB) domainCLI.Repository {
	return &CLICredentialRepository{db: db}
}

func (repository *CLICredentialRepository) Create(ctx context.Context, credential *domainCLI.Credential) error {
	if credential == nil || credential.ID == uuid.Nil {
		return domainCLI.ErrCredentialNotFound
	}

	if err := repository.db.WithContext(ctx).Create(cliModels.FromDomain(credential)).Error; err != nil {
		if isUniqueConstraintError(err) {
			return domainCLI.ErrCredentialConflict
		}

		return err
	}

	return nil
}

func (repository *CLICredentialRepository) FindByID(ctx context.Context, id uuid.UUID) (*domainCLI.Credential, error) {
	if id == uuid.Nil {
		return nil, domainCLI.ErrCredentialNotFound
	}

	var credential cliModels.DBCredential
	if err := repository.db.WithContext(ctx).First(&credential, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainCLI.ErrCredentialNotFound
		}

		return nil, err
	}

	return credential.ToDomain(), nil
}

func (repository *CLICredentialRepository) FindByTokenHash(ctx context.Context, tokenHash string) (*domainCLI.Credential, error) {
	tokenHash = strings.TrimSpace(tokenHash)
	if tokenHash == "" {
		return nil, domainCLI.ErrTokenHashRequired
	}

	var credential cliModels.DBCredential
	if err := repository.db.WithContext(ctx).First(&credential, "token_hash = ?", tokenHash).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainCLI.ErrCredentialNotFound
		}

		return nil, err
	}

	return credential.ToDomain(), nil
}

func (repository *CLICredentialRepository) Update(ctx context.Context, credential *domainCLI.Credential) error {
	if credential == nil || credential.ID == uuid.Nil {
		return domainCLI.ErrCredentialNotFound
	}

	result := repository.db.WithContext(ctx).
		Model(&cliModels.DBCredential{}).
		Where("id = ?", credential.ID).
		Updates(map[string]any{
			"name":         credential.Name,
			"last_used_at": credential.LastUsedAt,
			"revoked_at":   credential.RevokedAt,
			"updated_at":   credential.UpdatedAt,
		})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domainCLI.ErrCredentialNotFound
	}

	return nil
}
