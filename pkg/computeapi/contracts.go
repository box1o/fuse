package computeapi

import "time"

const CapabilitySchemaVersion = "1.0"

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

type RegisterNodeRequest struct {
	InstallationID string       `json:"installation_id"`
	Name           string       `json:"name"`
	Hostname       string       `json:"hostname"`
	AgentVersion   string       `json:"agent_version"`
	Capabilities   Capabilities `json:"capabilities"`
}

type UpdateNodeRequest struct {
	Name     *string `json:"name,omitempty"`
	Disabled *bool   `json:"disabled,omitempty"`
}

type Node struct {
	ID             string       `json:"id"`
	OwnerID        string       `json:"owner_id"`
	InstallationID string       `json:"installation_id"`
	Name           string       `json:"name"`
	Hostname       string       `json:"hostname"`
	AgentVersion   string       `json:"agent_version"`
	Status         string       `json:"status"`
	Capabilities   Capabilities `json:"capabilities"`
	RegisteredAt   time.Time    `json:"registered_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
	CreatedAt      time.Time    `json:"created_at"`
}
