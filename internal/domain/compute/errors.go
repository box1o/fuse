package compute

import "fuse/pkg/errors"

var (
	ErrInvalidNode         = errors.New("INVALID_COMPUTE_NODE", "the compute node is invalid")
	ErrNodeNotFound        = errors.New("COMPUTE_NODE_NOT_FOUND", "compute node not found")
	ErrOwnerIDEmpty        = errors.New("INVALID_COMPUTE_OWNER_ID", "compute node owner ID cannot be empty")
	ErrInstallationIDEmpty = errors.New("INVALID_COMPUTE_INSTALLATION_ID", "installation ID cannot be empty")
	ErrNodeNameEmpty       = errors.New("INVALID_COMPUTE_NODE_NAME", "compute node name cannot be empty")
	ErrNodeNameInvalid     = errors.New("INVALID_COMPUTE_NODE_NAME", "compute node name is invalid")
	ErrHostnameEmpty       = errors.New("INVALID_COMPUTE_HOSTNAME", "hostname cannot be empty")
	ErrHostnameInvalid     = errors.New("INVALID_COMPUTE_HOSTNAME", "hostname is invalid")
	ErrAgentVersionEmpty   = errors.New("INVALID_COMPUTE_AGENT_VERSION", "agent version cannot be empty")
	ErrAgentVersionInvalid = errors.New("INVALID_COMPUTE_AGENT_VERSION", "agent version is invalid")
	ErrCapabilitiesInvalid = errors.New("INVALID_COMPUTE_CAPABILITIES", "compute capabilities are invalid")
	ErrCreateNodeFailed    = errors.New("CREATE_COMPUTE_NODE_FAILED", "failed to create compute node")
	ErrUpdateNodeFailed    = errors.New("UPDATE_COMPUTE_NODE_FAILED", "failed to update compute node")
	ErrDatabaseOperation   = errors.New("COMPUTE_DATABASE_OPERATION_FAILED", "compute database operation failed")
	ErrInvalidCredential   = errors.New("INVALID_COMPUTE_CREDENTIAL", "compute CLI credential is invalid")
	ErrCredentialNotFound  = errors.New("COMPUTE_CREDENTIAL_NOT_FOUND", "compute CLI credential not found")
	ErrCredentialExpired   = errors.New("COMPUTE_CREDENTIAL_EXPIRED", "compute CLI credential has expired")
	ErrCredentialRevoked   = errors.New("COMPUTE_CREDENTIAL_REVOKED", "compute CLI credential has been revoked")
	ErrCreateCredential    = errors.New("CREATE_COMPUTE_CREDENTIAL_FAILED", "failed to create compute CLI credential")
	ErrUpdateCredential    = errors.New("UPDATE_COMPUTE_CREDENTIAL_FAILED", "failed to update compute CLI credential")
)
