package compute

import "fuse/pkg/errors"

var (
	ErrInvalidNode         = errors.New("INVALID_COMPUTE_NODE", "compute node is invalid")
	ErrOwnerIDRequired     = errors.New("OWNER_ID_REQUIRED", "owner ID is required")
	ErrWorkspaceIDRequired = errors.New("WORKSPACE_ID_REQUIRED", "workspace ID is required")
	ErrNodeNameRequired    = errors.New("NODE_NAME_REQUIRED", "compute node name is required")
	ErrInvalidNodeStatus   = errors.New("INVALID_NODE_STATUS", "compute node status is invalid")
	ErrInvalidCPUCores     = errors.New("INVALID_CPU_CORES", "CPU cores must be greater than zero")
	ErrInvalidMemory       = errors.New("INVALID_MEMORY", "memory must be greater than zero")
	ErrInvalidGPUCount     = errors.New("INVALID_GPU_COUNT", "GPU count cannot be negative")

	ErrNodeNotFound		   = errors.New("COMPUTE_NODE_NOT_FOUND", "compute node was not found")
	ErrNodeIDRequired      = errors.New("COMPUTE_NODE_ID_REQUIRED", "compute node ID is required")
	ErrCreateNodeFailed    = errors.New("CREATE_COMPUTE_NODE_FAILED", "failed to create compute node")
	ErrUpdateNodeFailed    = errors.New("UPDATE_COMPUTE_NODE_FAILED", "failed to update compute node")
	ErrDeleteNodeFailed    = errors.New("DELETE_COMPUTE_NODE_FAILED", "failed to delete compute node")
	ErrDatabaseOperation   = errors.New("DATABASE_OPERATION_FAILED", "database operation failed")
)
