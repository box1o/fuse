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

func NewComputeRepository(db *gorm.DB) compute.Repository {
	return &ComputeRepository{db: db}
}

func (r *ComputeRepository) FindByOwnerAndInstallationID(ctx context.Context, ownerID, installationID uuid.UUID) (*compute.Node, error) {
	if ownerID == uuid.Nil {
		return nil, compute.ErrOwnerIDEmpty
	}
	if installationID == uuid.Nil {
		return nil, compute.ErrInstallationIDEmpty
	}

	var dbNode models.DBNode
	err := r.db.WithContext(ctx).
		Where("owner_id = ? AND installation_id = ?", ownerID.String(), installationID.String()).
		First(&dbNode).Error
	if stdErrors.Is(err, gorm.ErrRecordNotFound) {
		return nil, compute.ErrNodeNotFound
	}
	if err != nil {
		return nil, compute.ErrDatabaseOperation.WithErr(err)
	}

	node, err := dbNode.ToDomain()
	if err != nil {
		return nil, compute.ErrDatabaseOperation.WithErr(err)
	}
	return node, nil
}

func (r *ComputeRepository) FindByID(ctx context.Context, ownerID, nodeID uuid.UUID) (*compute.Node, error) {
	var dbNode models.DBNode
	err := r.db.WithContext(ctx).Where("owner_id = ? AND id = ?", ownerID.String(), nodeID).First(&dbNode).Error
	if stdErrors.Is(err, gorm.ErrRecordNotFound) {
		return nil, compute.ErrNodeNotFound
	}
	if err != nil {
		return nil, compute.ErrDatabaseOperation.WithErr(err)
	}
	return dbNode.ToDomain()
}

func (r *ComputeRepository) ListByOwner(ctx context.Context, ownerID uuid.UUID) ([]*compute.Node, error) {
	var dbNodes []models.DBNode
	if err := r.db.WithContext(ctx).Where("owner_id = ?", ownerID.String()).Order("created_at DESC").Find(&dbNodes).Error; err != nil {
		return nil, compute.ErrDatabaseOperation.WithErr(err)
	}
	nodes := make([]*compute.Node, 0, len(dbNodes))
	for i := range dbNodes {
		node, err := dbNodes[i].ToDomain()
		if err != nil {
			return nil, compute.ErrDatabaseOperation.WithErr(err)
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}

func (r *ComputeRepository) Create(ctx context.Context, node *compute.Node) error {
	dbNode, err := models.FromDomainNode(node)
	if err != nil {
		return err
	}
	if err := r.db.WithContext(ctx).Create(dbNode).Error; err != nil {
		return compute.ErrCreateNodeFailed.WithErr(err)
	}
	converted, err := dbNode.ToDomain()
	if err != nil {
		return compute.ErrDatabaseOperation.WithErr(err)
	}
	*node = *converted
	return nil
}

func (r *ComputeRepository) Update(ctx context.Context, node *compute.Node) error {
	dbNode, err := models.FromDomainNode(node)
	if err != nil {
		return err
	}
	result := r.db.WithContext(ctx).Save(dbNode)
	if result.Error != nil {
		return compute.ErrUpdateNodeFailed.WithErr(result.Error)
	}
	if result.RowsAffected == 0 {
		return compute.ErrNodeNotFound
	}
	return nil
}

func (r *ComputeRepository) Delete(ctx context.Context, ownerID, nodeID uuid.UUID) error {
	result := r.db.WithContext(ctx).Where("owner_id = ? AND id = ?", ownerID.String(), nodeID).Delete(&models.DBNode{})
	if result.Error != nil {
		return compute.ErrDatabaseOperation.WithErr(result.Error)
	}
	if result.RowsAffected == 0 {
		return compute.ErrNodeNotFound
	}
	return nil
}

func (r *ComputeRepository) CreateCredential(ctx context.Context, credential *compute.CLICredential) error {
	dbCredential := models.FromDomainCredential(credential)
	if dbCredential == nil {
		return compute.ErrInvalidCredential
	}
	if err := r.db.WithContext(ctx).Create(dbCredential).Error; err != nil {
		return compute.ErrCreateCredential.WithErr(err)
	}
	return nil
}

func (r *ComputeRepository) FindCredentialByHash(ctx context.Context, tokenHash string) (*compute.CLICredential, error) {
	var dbCredential models.DBCLICredential
	err := r.db.WithContext(ctx).Where("token_hash = ?", tokenHash).First(&dbCredential).Error
	if stdErrors.Is(err, gorm.ErrRecordNotFound) {
		return nil, compute.ErrCredentialNotFound
	}
	if err != nil {
		return nil, compute.ErrDatabaseOperation.WithErr(err)
	}
	return dbCredential.ToDomain()
}

func (r *ComputeRepository) UpdateCredential(ctx context.Context, credential *compute.CLICredential) error {
	dbCredential := models.FromDomainCredential(credential)
	if dbCredential == nil {
		return compute.ErrInvalidCredential
	}
	if err := r.db.WithContext(ctx).Save(dbCredential).Error; err != nil {
		return compute.ErrUpdateCredential.WithErr(err)
	}
	return nil
}
