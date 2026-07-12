package models

import (
	"fmt"

	"fuse/internal/domain/compute"
	"fuse/internal/infrastructure/db"

	"github.com/google/uuid"
)

// DBNode is the database representation of a compute node.
type DBNode struct {
	db.Model

	Name     string `gorm:"not null;size:255" json:"name"`
	OwnerID  string `gorm:"not null;size:36;index" json:"owner_id"`
	NodeType string `gorm:"not null;size:20" json:"node_type"`
	Status   string `gorm:"not null;default:'pending';size:20" json:"status"`

	ComputeStacks []DBComputeStack `gorm:"many2many:compute_stack_nodes;" json:"compute_stacks,omitempty"`
}

func (DBNode) TableName() string {
	return "nodes"
}

func FromDomainNode(node *compute.Node) *DBNode {
	if node == nil {
		return nil
	}

	return &DBNode{
		Model: db.Model{
			ID:        node.ID,
			CreatedAt: node.CreatedAt,
			UpdatedAt: node.UpdatedAt,
		},
		Name:     node.Name,
		OwnerID:  node.OwnerID.String(),
		NodeType: string(node.NodeType),
		Status:   string(node.Status),
	}
}

// ToDomain converts a database node into a domain node.
// It returns an error instead of panicking when OwnerID is invalid.
func (d *DBNode) ToDomain() (*compute.Node, error) {
	if d == nil {
		return nil, nil
	}

	ownerID, err := uuid.Parse(d.OwnerID)
	if err != nil {
		return nil, fmt.Errorf("parse node owner ID %q: %w", d.OwnerID, err)
	}

	return &compute.Node{
		ID:        d.ID,
		Name:      d.Name,
		OwnerID:   ownerID,
		NodeType:  compute.NodeType(d.NodeType),
		Status:    compute.Status(d.Status),
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}, nil
}

// DBComputeStack is the database representation of a compute stack.
type DBComputeStack struct {
	db.Model

	Name    string `gorm:"not null;size:255" json:"name"`
	OwnerID string `gorm:"not null;size:36;index" json:"owner_id"`

	Nodes []DBNode `gorm:"many2many:compute_stack_nodes;" json:"nodes"`
}

func (DBComputeStack) TableName() string {
	return "compute_stacks"
}

func FromDomainComputeStack(stack *compute.ComputeStack) *DBComputeStack {
	if stack == nil {
		return nil
	}

	nodes := make([]DBNode, 0, len(stack.Nodes))

	for i := range stack.Nodes {
		dbNode := FromDomainNode(&stack.Nodes[i])
		if dbNode != nil {
			nodes = append(nodes, *dbNode)
		}
	}

	return &DBComputeStack{
		Model: db.Model{
			ID:        stack.ID,
			CreatedAt: stack.CreatedAt,
			UpdatedAt: stack.UpdatedAt,
		},
		Name:    stack.Name,
		OwnerID: stack.OwnerID.String(),
		Nodes:   nodes,
	}
}

func (d *DBComputeStack) ToDomain() (*compute.ComputeStack, error) {
	if d == nil {
		return nil, nil
	}

	ownerID, err := uuid.Parse(d.OwnerID)
	if err != nil {
		return nil, fmt.Errorf(
			"parse compute stack owner ID %q: %w",
			d.OwnerID,
			err,
		)
	}

	nodes := make([]compute.Node, 0, len(d.Nodes))

	for i := range d.Nodes {
		node, err := d.Nodes[i].ToDomain()
		if err != nil {
			return nil, fmt.Errorf(
				"convert node %s to domain: %w",
				d.Nodes[i].ID,
				err,
			)
		}

		if node != nil {
			nodes = append(nodes, *node)
		}
	}

	return &compute.ComputeStack{
		ID:        d.ID,
		Name:      d.Name,
		OwnerID:   ownerID,
		Nodes:     nodes,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}, nil
}
