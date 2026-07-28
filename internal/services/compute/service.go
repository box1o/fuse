package compute

import (
	"context"

	domain "fuse/internal/domain/compute"

	"github.com/google/uuid"
)

type Service struct {
	nodes domain.Repository
}

type CreateNodeInput struct {
	OwnerID      uuid.UUID
	WorkspaceID  uuid.UUID
	Name         string
	Capabilities domain.Capabilities
}

func NewService(nodes domain.Repository) *Service {
	return &Service{
		nodes: nodes,
	}
}

func (service *Service) CreateNode(ctx context.Context, input CreateNodeInput) (*domain.Node, error) {
	node, err := domain.NewNode(
		input.OwnerID,
		input.WorkspaceID,
		input.Name,
		input.Capabilities,
	)
	if err != nil {
		return nil, err
	}

	if err := service.nodes.Create(ctx, node); err != nil {
		return nil, err
	}

	return node, nil
}

func (service *Service) GetNode(ctx context.Context, nodeID uuid.UUID) (*domain.Node, error) {
	if nodeID == uuid.Nil {
		return nil, domain.ErrNodeIDRequired
	}

	return service.nodes.FindByID(ctx, nodeID)
}

func (service *Service) ListNodesByWorkspaceID(ctx context.Context, workspaceID uuid.UUID) ([]*domain.Node, error) {
	if workspaceID == uuid.Nil {
		return nil, domain.ErrWorkspaceIDRequired
	}

	return service.nodes.ListByWorkspaceID(ctx, workspaceID)
}

func (service *Service) ListNodesByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]*domain.Node, error) {
	if ownerID == uuid.Nil {
		return nil, domain.ErrOwnerIDRequired
	}

	return service.nodes.ListByOwnerID(ctx, ownerID)
}

func (service *Service) RenameNode(ctx context.Context, nodeID uuid.UUID, name string) (*domain.Node, error) {
	node, err := service.getNodeForUpdate(ctx, nodeID)
	if err != nil {
		return nil, err
	}

	if err := node.Rename(name); err != nil {
		return nil, err
	}

	if err := service.nodes.Update(ctx, node); err != nil {
		return nil, err
	}

	return node, nil
}

func (service *Service) UpdateNodeStatus(ctx context.Context, nodeID uuid.UUID, status domain.NodeStatus) (*domain.Node, error) {
	node, err := service.getNodeForUpdate(ctx, nodeID)
	if err != nil {
		return nil, err
	}

	if err := node.UpdateStatus(status); err != nil {
		return nil, err
	}

	if err := service.nodes.Update(ctx, node); err != nil {
		return nil, err
	}

	return node, nil
}

func (service *Service) UpdateNodeCapabilities(ctx context.Context, nodeID uuid.UUID, capabilities domain.Capabilities) (*domain.Node, error) {
	node, err := service.getNodeForUpdate(ctx, nodeID)
	if err != nil {
		return nil, err
	}

	if err := node.UpdateCapabilities(capabilities); err != nil {
		return nil, err
	}

	if err := service.nodes.Update(ctx, node); err != nil {
		return nil, err
	}

	return node, nil
}

func (service *Service) DeleteNode(ctx context.Context, nodeID uuid.UUID) error {
	if nodeID == uuid.Nil {
		return domain.ErrNodeIDRequired
	}

	return service.nodes.Delete(ctx, nodeID)
}

func (service *Service) getNodeForUpdate(ctx context.Context, nodeID uuid.UUID) (*domain.Node, error) {
	if nodeID == uuid.Nil {
		return nil, domain.ErrNodeIDRequired
	}

	return service.nodes.FindByID(ctx, nodeID)
}
