package compute

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type NodeStatus string

const (
	NodeStatusPending  NodeStatus = "pending"
	NodeStatusOnline   NodeStatus = "online"
	NodeStatusOffline  NodeStatus = "offline"
	NodeStatusDisabled NodeStatus = "disabled"
)

func (status NodeStatus) IsValid() bool {
	switch status {
	case NodeStatusPending,
		NodeStatusOnline,
		NodeStatusOffline,
		NodeStatusDisabled:
		return true
	default:
		return false
	}
}

type Capabilities struct {
	CPUCores int `json:"cpu_cores" example:"8" minimum:"1"`
	MemoryMB int `json:"memory_mb" example:"16384" minimum:"1"`
	GPUCount int `json:"gpu_count" example:"1" minimum:"0"`
	NPUSupported bool `json:"npu_supported" example:"false"`
}

func (capabilities Capabilities) Validate() error {
	if capabilities.CPUCores <= 0 {
		return ErrInvalidCPUCores
	}

	if capabilities.MemoryMB <= 0 {
		return ErrInvalidMemory
	}

	if capabilities.GPUCount < 0 {
		return ErrInvalidGPUCount
	}

	return nil
}

type Node struct {
	ID           uuid.UUID    `json:"id"`
	OwnerID      uuid.UUID    `json:"owner_id"`
	WorkspaceID  uuid.UUID    `json:"workspace_id"`
	Name         string       `json:"name"`
	Status       NodeStatus   `json:"status"`
	Capabilities Capabilities `json:"capabilities"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

func NewNode( ownerID uuid.UUID, workspaceID uuid.UUID, name string, capabilities Capabilities,) (*Node, error) {
	if ownerID == uuid.Nil {
		return nil, ErrOwnerIDRequired
	}

	if workspaceID == uuid.Nil {
		return nil, ErrWorkspaceIDRequired
	}

	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrNodeNameRequired
	}

	if err := capabilities.Validate(); err != nil {
		return nil, err
	}

	now := time.Now().UTC()

	return &Node{
		ID:           uuid.New(),
		OwnerID:      ownerID,
		WorkspaceID:  workspaceID,
		Name:         name,
		Status:       NodeStatusPending,
		Capabilities: capabilities,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

func (node *Node) UpdateStatus(status NodeStatus) error {
	if node == nil {
		return ErrInvalidNode
	}

	if !status.IsValid() {
		return ErrInvalidNodeStatus
	}

	node.Status = status
	node.UpdatedAt = time.Now().UTC()

	return nil
}

func (node *Node) UpdateCapabilities(capabilities Capabilities) error {
	if node == nil {
		return ErrInvalidNode
	}

	if err := capabilities.Validate(); err != nil {
		return err
	}

	node.Capabilities = capabilities
	node.UpdatedAt = time.Now().UTC()

	return nil
}

func (node *Node) Rename(name string) error {
	if node == nil {
		return ErrInvalidNode
	}

	name = strings.TrimSpace(name)
	if name == "" {
		return ErrNodeNameRequired
	}

	node.Name = name
	node.UpdatedAt = time.Now().UTC()

	return nil
}
