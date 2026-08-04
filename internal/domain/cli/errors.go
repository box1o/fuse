package cli

import "fuse/pkg/errors"

var (
	ErrCredentialNotFound = errors.New("CLI_CREDENTIAL_NOT_FOUND", "CLI credential was not found")
	ErrCredentialExpired  = errors.New("CLI_CREDENTIAL_EXPIRED", "CLI credential has expired")
	ErrCredentialRevoked  = errors.New("CLI_CREDENTIAL_REVOKED", "CLI credential has been revoked")
	ErrOwnerRequired      = errors.New("CLI_OWNER_REQUIRED", "CLI credential owner is required")
	ErrNameRequired       = errors.New("CLI_NAME_REQUIRED", "CLI credential name is required")
	ErrTokenHashRequired  = errors.New("CLI_TOKEN_HASH_REQUIRED", "CLI credential token hash is required")
	ErrCredentialConflict = errors.New("CLI_CREDENTIAL_CONFLICT", "CLI credential already exists")
)
