package compute

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, node *Node) error

	FindByID(ctx context.Context, id uuid.UUID) (*Node, error)
	ListByWorkspaceID( ctx context.Context, workspaceID uuid.UUID,) ([]*Node, error)
	ListByOwnerID( ctx context.Context, ownerID uuid.UUID,) ([]*Node, error)

	Update(ctx context.Context, node *Node) error
	Delete(ctx context.Context, id uuid.UUID) error
}
