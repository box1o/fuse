package compute

import (
	domain "fuse/internal/domain/compute"
	"fuse/pkg/computeapi"
)

func capabilitiesToDomain(value computeapi.Capabilities) domain.Capabilities {
	accelerators := make([]domain.Accelerator, 0, len(value.Accelerators))
	for _, accelerator := range value.Accelerators {
		accelerators = append(accelerators, domain.Accelerator{
			Kind: accelerator.Kind, Vendor: accelerator.Vendor, Model: accelerator.Model,
			DeviceID: accelerator.DeviceID, MemoryBytes: accelerator.MemoryBytes,
			DriverVersion: accelerator.DriverVersion, RuntimeAvailable: accelerator.RuntimeAvailable,
		})
	}
	return domain.Capabilities{
		SchemaVersion:    value.SchemaVersion,
		OS:               domain.OperatingSystem{Name: value.OS.Name, Version: value.OS.Version, Kernel: value.OS.Kernel, Architecture: value.OS.Architecture},
		CPU:              domain.CPU{Vendor: value.CPU.Vendor, Model: value.CPU.Model, PhysicalCores: value.CPU.PhysicalCores, LogicalCores: value.CPU.LogicalCores},
		Memory:           domain.Memory{TotalBytes: value.Memory.TotalBytes},
		Storage:          domain.Storage{TotalBytes: value.Storage.TotalBytes},
		ContainerRuntime: domain.ContainerRuntime{Name: value.ContainerRuntime.Name, Version: value.ContainerRuntime.Version, Available: value.ContainerRuntime.Available},
		Accelerators:     accelerators,
	}
}

func nodeToAPI(node *domain.Node) computeapi.Node {
	capabilities := computeapi.Capabilities{
		SchemaVersion:    node.Capabilities.SchemaVersion,
		OS:               computeapi.OperatingSystem{Name: node.Capabilities.OS.Name, Version: node.Capabilities.OS.Version, Kernel: node.Capabilities.OS.Kernel, Architecture: node.Capabilities.OS.Architecture},
		CPU:              computeapi.CPU{Vendor: node.Capabilities.CPU.Vendor, Model: node.Capabilities.CPU.Model, PhysicalCores: node.Capabilities.CPU.PhysicalCores, LogicalCores: node.Capabilities.CPU.LogicalCores},
		Memory:           computeapi.Memory{TotalBytes: node.Capabilities.Memory.TotalBytes},
		Storage:          computeapi.Storage{TotalBytes: node.Capabilities.Storage.TotalBytes},
		ContainerRuntime: computeapi.ContainerRuntime{Name: node.Capabilities.ContainerRuntime.Name, Version: node.Capabilities.ContainerRuntime.Version, Available: node.Capabilities.ContainerRuntime.Available},
		Accelerators:     make([]computeapi.Accelerator, 0, len(node.Capabilities.Accelerators)),
	}
	for _, accelerator := range node.Capabilities.Accelerators {
		capabilities.Accelerators = append(capabilities.Accelerators, computeapi.Accelerator{
			Kind: accelerator.Kind, Vendor: accelerator.Vendor, Model: accelerator.Model,
			DeviceID: accelerator.DeviceID, MemoryBytes: accelerator.MemoryBytes,
			DriverVersion: accelerator.DriverVersion, RuntimeAvailable: accelerator.RuntimeAvailable,
		})
	}
	return computeapi.Node{
		ID: node.ID.String(), OwnerID: node.OwnerID.String(), InstallationID: node.InstallationID.String(),
		Name: node.Name, Hostname: node.Hostname, AgentVersion: node.AgentVersion, Status: string(node.Status),
		Capabilities: capabilities, RegisteredAt: node.RegisteredAt, UpdatedAt: node.UpdatedAt, CreatedAt: node.CreatedAt,
	}
}
