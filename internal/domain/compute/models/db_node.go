package models

import (
	"encoding/json"
	"fmt"
	"time"

	"fuse/internal/domain/compute"
	"fuse/internal/infrastructure/db"

	"github.com/google/uuid"
)

type DBNode struct {
	db.Model

	OwnerID        string          `gorm:"not null;size:36;uniqueIndex:idx_compute_owner_installation"`
	InstallationID string          `gorm:"not null;size:36;uniqueIndex:idx_compute_owner_installation"`
	Name           string          `gorm:"not null;size:255"`
	Hostname       string          `gorm:"not null;size:255"`
	AgentVersion   string          `gorm:"not null;size:64"`
	Status         string          `gorm:"not null;size:32;default:'registered'"`
	Capabilities   json.RawMessage `gorm:"type:jsonb;not null"`
	RegisteredAt   time.Time       `gorm:"not null"`
}

type DBCLICredential struct {
	db.Model

	OwnerID    string    `gorm:"not null;size:36;index"`
	Name       string    `gorm:"not null;size:255"`
	TokenHash  string    `gorm:"not null;size:64;uniqueIndex"`
	ExpiresAt  time.Time `gorm:"not null;index"`
	LastUsedAt *time.Time
	RevokedAt  *time.Time `gorm:"index"`
}

func (DBCLICredential) TableName() string { return "compute_cli_credentials" }

func FromDomainCredential(credential *compute.CLICredential) *DBCLICredential {
	if credential == nil {
		return nil
	}
	return &DBCLICredential{
		Model:   db.Model{ID: credential.ID, CreatedAt: credential.CreatedAt, UpdatedAt: credential.UpdatedAt},
		OwnerID: credential.OwnerID.String(), Name: credential.Name, TokenHash: credential.TokenHash,
		ExpiresAt: credential.ExpiresAt, LastUsedAt: credential.LastUsedAt, RevokedAt: credential.RevokedAt,
	}
}

func (d *DBCLICredential) ToDomain() (*compute.CLICredential, error) {
	ownerID, err := uuid.Parse(d.OwnerID)
	if err != nil {
		return nil, fmt.Errorf("parse credential owner ID: %w", err)
	}
	return &compute.CLICredential{
		ID: d.ID, OwnerID: ownerID, Name: d.Name, TokenHash: d.TokenHash,
		ExpiresAt: d.ExpiresAt, LastUsedAt: d.LastUsedAt, RevokedAt: d.RevokedAt,
		CreatedAt: d.CreatedAt, UpdatedAt: d.UpdatedAt,
	}, nil
}

func (DBNode) TableName() string { return "compute_nodes" }

func FromDomainNode(node *compute.Node) (*DBNode, error) {
	if node == nil {
		return nil, compute.ErrInvalidNode
	}

	capabilities, err := json.Marshal(node.Capabilities)
	if err != nil {
		return nil, fmt.Errorf("marshal compute capabilities: %w", err)
	}

	return &DBNode{
		Model:   db.Model{ID: node.ID, CreatedAt: node.CreatedAt, UpdatedAt: node.UpdatedAt},
		OwnerID: node.OwnerID.String(), InstallationID: node.InstallationID.String(),
		Name: node.Name, Hostname: node.Hostname, AgentVersion: node.AgentVersion,
		Status: string(node.Status), Capabilities: capabilities, RegisteredAt: node.RegisteredAt,
	}, nil
}

func (d *DBNode) ToDomain() (*compute.Node, error) {
	if d == nil {
		return nil, compute.ErrInvalidNode
	}

	ownerID, err := uuid.Parse(d.OwnerID)
	if err != nil {
		return nil, fmt.Errorf("parse compute owner ID: %w", err)
	}
	installationID, err := uuid.Parse(d.InstallationID)
	if err != nil {
		return nil, fmt.Errorf("parse compute installation ID: %w", err)
	}

	var capabilities compute.Capabilities
	if err := json.Unmarshal(d.Capabilities, &capabilities); err != nil {
		return nil, fmt.Errorf("unmarshal compute capabilities: %w", err)
	}

	return &compute.Node{
		ID: d.ID, OwnerID: ownerID, InstallationID: installationID,
		Name: d.Name, Hostname: d.Hostname, AgentVersion: d.AgentVersion,
		Status: compute.Status(d.Status), Capabilities: capabilities,
		RegisteredAt: d.RegisteredAt, CreatedAt: d.CreatedAt, UpdatedAt: d.UpdatedAt,
	}, nil
}
