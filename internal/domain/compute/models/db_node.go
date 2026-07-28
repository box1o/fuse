package models

import (
	"fmt"

	"fuse/internal/domain/compute"
	"fuse/internal/infrastructure/db"

	"github.com/google/uuid"
)

type DBNode struct {
	db.Model

	OwnerID     string `gorm:"type:uuid;not null;index" json:"owner_id"`
	WorkspaceID string `gorm:"type:uuid;not null;index" json:"workspace_id"`

	Name   string `gorm:"not null;size:255" json:"name"`
	Status string `gorm:"not null;size:50;index" json:"status"`

	CPUCores     int  `gorm:"not null;check:cpu_cores > 0" json:"cpu_cores"`
	MemoryMB     int  `gorm:"not null;check:memory_mb > 0" json:"memory_mb"`
	GPUCount     int  `gorm:"not null;default:0;check:gpu_count >= 0" json:"gpu_count"`
	NPUSupported bool `gorm:"not null;default:false" json:"npu_supported"`
}

func (DBNode) TableName() string {
	return "compute_nodes"
}

func FromDomain(node *compute.Node) (*DBNode, error) {
	if node == nil || node.ID == uuid.Nil {
		return nil, compute.ErrInvalidNode
	}

	return &DBNode{
		Model: db.Model{
			ID:        node.ID,
			CreatedAt: node.CreatedAt,
			UpdatedAt: node.UpdatedAt,
		},
		OwnerID:     node.OwnerID.String(),
		WorkspaceID: node.WorkspaceID.String(),
		Name:        node.Name,
		Status:      string(node.Status),
		CPUCores:    node.Capabilities.CPUCores,
		MemoryMB:    node.Capabilities.MemoryMB,
		GPUCount:    node.Capabilities.GPUCount,
		NPUSupported: node.Capabilities.NPUSupported,
	}, nil
}

func (node *DBNode) ToDomain() (*compute.Node, error) {
	if node == nil {
		return nil, compute.ErrInvalidNode
	}

	ownerID, err := uuid.Parse(node.OwnerID)
	if err != nil {
		return nil, fmt.Errorf("parse compute node owner ID: %w", err)
	}

	workspaceID, err := uuid.Parse(node.WorkspaceID)
	if err != nil {
		return nil, fmt.Errorf("parse compute node workspace ID: %w", err)
	}

	status := compute.NodeStatus(node.Status)
	if !status.IsValid() {
		return nil, compute.ErrInvalidNodeStatus
	}

	capabilities := compute.Capabilities{
		CPUCores:     node.CPUCores,
		MemoryMB:     node.MemoryMB,
		GPUCount:     node.GPUCount,
		NPUSupported: node.NPUSupported,
	}

	if err := capabilities.Validate(); err != nil {
		return nil, err
	}

	return &compute.Node{
		ID:           node.ID,
		OwnerID:      ownerID,
		WorkspaceID:  workspaceID,
		Name:         node.Name,
		Status:       status,
		Capabilities: capabilities,
		CreatedAt:    node.CreatedAt,
		UpdatedAt:    node.UpdatedAt,
	}, nil
}
