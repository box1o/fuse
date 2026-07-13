export type ComputeNodeStatus = "registered" | "disabled";

export interface OperatingSystem {
    name: string;
    version?: string;
    kernel?: string;
    architecture: string;
}

export interface CPUCapabilities {
    vendor?: string;
    model: string;
    physical_cores: number;
    logical_cores: number;
}

export interface MemoryCapabilities {
    total_bytes: number;
}

export interface StorageCapabilities {
    total_bytes: number;
}

export interface ContainerRuntime {
    name: string;
    version?: string;
    available: boolean;
}

export interface Accelerator {
    kind: string;
    vendor: string;
    model: string;
    device_id?: string;
    memory_bytes?: number;
    driver_version?: string;
    runtime_available: boolean;
}

export interface ComputeCapabilities {
    schema_version: string;
    os: OperatingSystem;
    cpu: CPUCapabilities;
    memory: MemoryCapabilities;
    storage: StorageCapabilities;
    container_runtime: ContainerRuntime;
    accelerators: Accelerator[];
}

export interface ComputeNode {
    id: string;
    owner_id: string;
    installation_id: string;
    name: string;
    hostname: string;
    agent_version: string;
    status: ComputeNodeStatus;
    capabilities: ComputeCapabilities;
    registered_at: string;
    updated_at: string;
    created_at: string;
}

export interface RegisterComputeNodeRequest {
    installation_id: string;
    name: string;
    hostname: string;
    agent_version: string;
    capabilities: ComputeCapabilities;
}

export interface UpdateComputeNodeRequest {
    name?: string;
    disabled?: boolean;
}

export type ComputeNodeAvailability = "online" | "offline";
export type ComputeNodeFilter = "all" | ComputeNodeAvailability;
