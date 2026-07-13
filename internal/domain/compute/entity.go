package compute

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

const CapabilitySchemaVersion = "1.0"

type Status string

const (
	StatusRegistered Status = "registered"
	StatusDisabled   Status = "disabled"
)

type OperatingSystem struct {
	Name         string `json:"name"`
	Version      string `json:"version,omitempty"`
	Kernel       string `json:"kernel,omitempty"`
	Architecture string `json:"architecture"`
}

type CPU struct {
	Vendor        string `json:"vendor,omitempty"`
	Model         string `json:"model"`
	PhysicalCores int    `json:"physical_cores"`
	LogicalCores  int    `json:"logical_cores"`
}

type Memory struct {
	TotalBytes uint64 `json:"total_bytes"`
}

type Storage struct {
	TotalBytes uint64 `json:"total_bytes"`
}

type ContainerRuntime struct {
	Name      string `json:"name"`
	Version   string `json:"version,omitempty"`
	Available bool   `json:"available"`
}

type Accelerator struct {
	Kind             string `json:"kind"`
	Vendor           string `json:"vendor"`
	Model            string `json:"model"`
	DeviceID         string `json:"device_id,omitempty"`
	MemoryBytes      uint64 `json:"memory_bytes,omitempty"`
	DriverVersion    string `json:"driver_version,omitempty"`
	RuntimeAvailable bool   `json:"runtime_available"`
}

type Capabilities struct {
	SchemaVersion    string           `json:"schema_version"`
	OS               OperatingSystem  `json:"os"`
	CPU              CPU              `json:"cpu"`
	Memory           Memory           `json:"memory"`
	Storage          Storage          `json:"storage"`
	ContainerRuntime ContainerRuntime `json:"container_runtime"`
	Accelerators     []Accelerator    `json:"accelerators"`
}

type Node struct {
	ID             uuid.UUID    `json:"id"`
	OwnerID        uuid.UUID    `json:"owner_id"`
	InstallationID uuid.UUID    `json:"installation_id"`
	Name           string       `json:"name"`
	Hostname       string       `json:"hostname"`
	AgentVersion   string       `json:"agent_version"`
	Status         Status       `json:"status"`
	Capabilities   Capabilities `json:"capabilities"`
	RegisteredAt   time.Time    `json:"registered_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
	CreatedAt      time.Time    `json:"created_at"`
}

type RegisterNodeInput struct {
	InstallationID uuid.UUID
	Name           string
	Hostname       string
	AgentVersion   string
	Capabilities   Capabilities
}

func NewNode(ownerID uuid.UUID, input RegisterNodeInput) (*Node, error) {
	if err := ValidateRegistration(ownerID, input); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	return &Node{
		ID:             uuid.New(),
		OwnerID:        ownerID,
		InstallationID: input.InstallationID,
		Name:           strings.TrimSpace(input.Name),
		Hostname:       strings.TrimSpace(input.Hostname),
		AgentVersion:   strings.TrimSpace(input.AgentVersion),
		Status:         StatusRegistered,
		Capabilities:   normalizeCapabilities(input.Capabilities),
		RegisteredAt:   now,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

func (n *Node) ApplyRegistration(input RegisterNodeInput) error {
	if n == nil {
		return ErrInvalidNode
	}
	if err := ValidateRegistration(n.OwnerID, input); err != nil {
		return err
	}

	n.Name = strings.TrimSpace(input.Name)
	n.Hostname = strings.TrimSpace(input.Hostname)
	n.AgentVersion = strings.TrimSpace(input.AgentVersion)
	n.Capabilities = normalizeCapabilities(input.Capabilities)
	if n.Status != StatusDisabled {
		n.Status = StatusRegistered
	}
	n.UpdatedAt = time.Now().UTC()
	return nil
}

func (n *Node) Rename(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return ErrNodeNameEmpty
	}
	if len(name) > 255 {
		return ErrNodeNameInvalid.WithDetail("node name cannot exceed 255 characters")
	}
	n.Name = name
	n.UpdatedAt = time.Now().UTC()
	return nil
}

func (n *Node) SetDisabled(disabled bool) {
	if disabled {
		n.Status = StatusDisabled
	} else {
		n.Status = StatusRegistered
	}
	n.UpdatedAt = time.Now().UTC()
}

func ValidateRegistration(ownerID uuid.UUID, input RegisterNodeInput) error {
	switch {
	case ownerID == uuid.Nil:
		return ErrOwnerIDEmpty
	case input.InstallationID == uuid.Nil:
		return ErrInstallationIDEmpty
	case strings.TrimSpace(input.Name) == "":
		return ErrNodeNameEmpty
	case len(strings.TrimSpace(input.Name)) > 255:
		return ErrNodeNameInvalid.WithDetail("node name cannot exceed 255 characters")
	case strings.TrimSpace(input.Hostname) == "":
		return ErrHostnameEmpty
	case len(strings.TrimSpace(input.Hostname)) > 255:
		return ErrHostnameInvalid.WithDetail("hostname cannot exceed 255 characters")
	case strings.TrimSpace(input.AgentVersion) == "":
		return ErrAgentVersionEmpty
	case len(strings.TrimSpace(input.AgentVersion)) > 64:
		return ErrAgentVersionInvalid.WithDetail("agent version cannot exceed 64 characters")
	case strings.TrimSpace(input.Capabilities.SchemaVersion) != "" && input.Capabilities.SchemaVersion != CapabilitySchemaVersion:
		return ErrCapabilitiesInvalid.WithDetail("unsupported capability schema version")
	case strings.TrimSpace(input.Capabilities.OS.Name) == "":
		return ErrCapabilitiesInvalid.WithDetail("operating system name is required")
	case strings.TrimSpace(input.Capabilities.OS.Architecture) == "":
		return ErrCapabilitiesInvalid.WithDetail("operating system architecture is required")
	case strings.TrimSpace(input.Capabilities.CPU.Model) == "":
		return ErrCapabilitiesInvalid.WithDetail("CPU model is required")
	case input.Capabilities.CPU.LogicalCores < 1:
		return ErrCapabilitiesInvalid.WithDetail("CPU logical cores must be greater than zero")
	case input.Capabilities.CPU.PhysicalCores < 0:
		return ErrCapabilitiesInvalid.WithDetail("CPU physical cores cannot be negative")
	case input.Capabilities.Memory.TotalBytes == 0:
		return ErrCapabilitiesInvalid.WithDetail("total memory must be greater than zero")
	}

	for _, accelerator := range input.Capabilities.Accelerators {
		if strings.TrimSpace(accelerator.Kind) == "" || strings.TrimSpace(accelerator.Vendor) == "" || strings.TrimSpace(accelerator.Model) == "" {
			return ErrCapabilitiesInvalid.WithDetail("accelerator kind, vendor, and model are required")
		}
	}

	return nil
}

func normalizeCapabilities(capabilities Capabilities) Capabilities {
	if strings.TrimSpace(capabilities.SchemaVersion) == "" {
		capabilities.SchemaVersion = CapabilitySchemaVersion
	}
	if capabilities.Accelerators == nil {
		capabilities.Accelerators = make([]Accelerator, 0)
	}
	return capabilities
}
