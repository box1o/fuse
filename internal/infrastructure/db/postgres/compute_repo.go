package postgres

import (
	"context"
	stdErrors "errors"

	"fuse/internal/domain/compute"
	"fuse/internal/domain/compute/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ComputeRepository struct {
	db *gorm.DB
}

var _ compute.Repository = (*ComputeRepository)(nil)

func NewComputeRepository(db *gorm.DB) compute.Repository {
	return &ComputeRepository{
		db: db,
	}
}

func (repository *ComputeRepository) Create(ctx context.Context, node *compute.Node) error {
	if node == nil {
		return compute.ErrInvalidNode
	}

	dbNode, err := models.FromDomain(node)
	if err != nil {
		return err
	}

	if err := repository.db.WithContext(ctx).Create(dbNode).Error; err != nil {
		return compute.ErrCreateNodeFailed.WithErr(err)
	}

	convertedNode, err := dbNode.ToDomain()
	if err != nil {
		return err
	}

	*node = *convertedNode
	return nil
}

func (repository *ComputeRepository) FindByID(ctx context.Context, id uuid.UUID) (*compute.Node, error) {
	if id == uuid.Nil {
		return nil, compute.ErrNodeIDRequired
	}

	var dbNode models.DBNode

	err := repository.db.WithContext(ctx).First(&dbNode, "id = ?", id).Error
	if stdErrors.Is(err, gorm.ErrRecordNotFound) {
		return nil, compute.ErrNodeNotFound
	}

	if err != nil {
		return nil, compute.ErrDatabaseOperation.WithErr(err)
	}

	node, err := dbNode.ToDomain()
	if err != nil {
		return nil, err
	}

	return node, nil
}

func (repository *ComputeRepository) ListByWorkspaceID(ctx context.Context, workspaceID uuid.UUID) ([]*compute.Node, error) {
	if workspaceID == uuid.Nil {
		return nil, compute.ErrWorkspaceIDRequired
	}

	var dbNodes []models.DBNode

	err := repository.db.WithContext(ctx).Where("workspace_id = ?", workspaceID.String()).Order("created_at DESC").Find(&dbNodes).Error
	if err != nil {
		return nil, compute.ErrDatabaseOperation.WithErr(err)
	}

	return convertToComputeNodes(dbNodes)
}

func (repository *ComputeRepository) ListByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]*compute.Node, error) {
	if ownerID == uuid.Nil {
		return nil, compute.ErrOwnerIDRequired
	}

	var dbNodes []models.DBNode

	err := repository.db.WithContext(ctx).Where("owner_id = ?", ownerID.String()).Order("created_at DESC").Find(&dbNodes).Error
	if err != nil {
		return nil, compute.ErrDatabaseOperation.WithErr(err)
	}

	return convertToComputeNodes(dbNodes)
}

func (repository *ComputeRepository) Update(ctx context.Context, node *compute.Node) error {
	if node == nil {
		return compute.ErrInvalidNode
	}

	if node.ID == uuid.Nil {
		return compute.ErrNodeIDRequired
	}

	dbNode, err := models.FromDomain(node)
	if err != nil {
		return err
	}

	result := repository.db.
		WithContext(ctx).
		Model(&models.DBNode{}).
		Where("id = ?", node.ID).
		Updates(map[string]any{
			"owner_id":      dbNode.OwnerID,
			"workspace_id":  dbNode.WorkspaceID,
			"name":          dbNode.Name,
			"status":        dbNode.Status,
			"cpu_cores":     dbNode.CPUCores,
			"memory_mb":     dbNode.MemoryMB,
			"gpu_count":     dbNode.GPUCount,
			"npu_supported": dbNode.NPUSupported,
			"updated_at":    dbNode.UpdatedAt,
		})

	if result.Error != nil {
		return compute.ErrUpdateNodeFailed.WithErr(result.Error)
	}

	if result.RowsAffected == 0 {
		return compute.ErrNodeNotFound
	}

	return nil
}

func (repository *ComputeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return compute.ErrNodeIDRequired
	}

	result := repository.db.WithContext(ctx).Delete(&models.DBNode{}, "id = ?", id)
	if result.Error != nil {
		return compute.ErrDeleteNodeFailed.WithErr(result.Error)
	}

	if result.RowsAffected == 0 {
		return compute.ErrNodeNotFound
	}

	return nil
}

func convertToComputeNodes(dbNodes []models.DBNode) ([]*compute.Node, error) {
	nodes := make([]*compute.Node, 0, len(dbNodes))

	for index := range dbNodes {
		node, err := dbNodes[index].ToDomain()
		if err != nil {
			return nil, err
		}

		nodes = append(nodes, node)
	}

	return nodes, nil
}
