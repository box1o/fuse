package compute

import (
	"time"

	"github.com/google/uuid"
)

type NodeType string

const (
	GPUNode NodeType = "gpu"
	CPUNode NodeType = "cpu"
	NPUNode NodeType = "npu"
)

func (t NodeType) IsValid() bool {
	switch t {
	case GPUNode, CPUNode, NPUNode:
		return true
	default:
		return false
	}
}

type Status string

const (
	StatusActive   Status = "active"
	StatusPending  Status = "pending"
	StatusInactive Status = "inactive"
)

func (s Status) IsValid() bool {
	switch s {
	case StatusActive, StatusPending, StatusInactive:
		return true
	default:
		return false
	}
}

type Node struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	OwnerID uuid.UUID `json:"owner_id"`

	NodeType NodeType `json:"node_type"`
	Status   Status   `json:"status"`

	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

type ComputeStack struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	OwnerID uuid.UUID `json:"owner_id"`

	Nodes []Node `json:"nodes"`

	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

func NewNode(
	name string,
	ownerID uuid.UUID,
	nodeType NodeType,
) *Node {
	now := time.Now().UTC()

	return &Node{
		ID:        uuid.New(),
		Name:      name,
		OwnerID:   ownerID,
		NodeType:  nodeType,
		Status:    StatusPending,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (n *Node) UpdateStatus(status Status) {
	n.Status = status
	n.UpdatedAt = time.Now().UTC()
}

func NewComputeStack(
	name string,
	ownerID uuid.UUID,
	nodes []Node,
) *ComputeStack {
	now := time.Now().UTC()

	if nodes == nil {
		nodes = make([]Node, 0)
	}

	return &ComputeStack{
		ID:        uuid.New(),
		Name:      name,
		OwnerID:   ownerID,
		Nodes:     nodes,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
