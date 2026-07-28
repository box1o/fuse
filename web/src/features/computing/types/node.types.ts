type NodeStatus = | "pending" | "online" | "offline" | "disabled";

interface Capabilities {
    cpu_cores: number;
    memory_mb: number;
    gpu_count: number;
    npu_supported: boolean;
}

interface ComputeNode {
    id: string;
    owner_id: string;
    workspace_id: string;
    name: string;
    status: NodeStatus;
    capabilities: Capabilities;
    created_at: string;
    updated_at: string;
}

interface CreateNodeRequest {
    workspace_id: string;
    name: string;
    capabilities: Capabilities;
}

interface RenameNodeRequest {
    name: string;
}

interface UpdateNodeStatusRequest {
    status: NodeStatus;
}

interface UpdateNodeCapabilitiesRequest {
    capabilities: Capabilities;
}

export type {
    Capabilities,
    ComputeNode,
    CreateNodeRequest,
    NodeStatus,
    RenameNodeRequest,
    UpdateNodeCapabilitiesRequest,
    UpdateNodeStatusRequest,
};
