package cli

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, credential *Credential) error
	FindByID(ctx context.Context, id uuid.UUID) (*Credential, error)
	FindByTokenHash(ctx context.Context, tokenHash string) (*Credential, error)
	Update(ctx context.Context, credential *Credential) error
}
