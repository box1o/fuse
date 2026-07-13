package compute

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	FindByOwnerAndInstallationID(ctx context.Context, ownerID, installationID uuid.UUID) (*Node, error)
	FindByID(ctx context.Context, ownerID, nodeID uuid.UUID) (*Node, error)
	ListByOwner(ctx context.Context, ownerID uuid.UUID) ([]*Node, error)
	Create(ctx context.Context, node *Node) error
	Update(ctx context.Context, node *Node) error
	Delete(ctx context.Context, ownerID, nodeID uuid.UUID) error

	CreateCredential(ctx context.Context, credential *CLICredential) error
	FindCredentialByHash(ctx context.Context, tokenHash string) (*CLICredential, error)
	UpdateCredential(ctx context.Context, credential *CLICredential) error
}
