package deviceauth

import (
	"context"
	"strings"
	"time"

	domainCLI "fuse/internal/domain/cli"
	"fuse/internal/domain/user"
)

func (service *Service) ValidateToken(ctx context.Context, accessToken string) (*domainCLI.Credential, *user.User, error) {
	accessToken = strings.TrimSpace(accessToken)
	if accessToken == "" {
		return nil, nil, domainCLI.ErrCredentialNotFound
	}

	credential, err := service.credentials.FindByTokenHash(ctx, hash(accessToken))
	if err != nil {
		return nil, nil, err
	}
	if err := credential.Validate(time.Now()); err != nil {
		return nil, nil, err
	}

	owner, err := service.users.FindByID(ctx, credential.OwnerID)
	if err != nil {
		return nil, nil, err
	}

	credential.MarkUsed(time.Now())
	if err := service.credentials.Update(ctx, credential); err != nil {
		return nil, nil, err
	}

	return credential, owner, nil
}

func (service *Service) RevokeToken(ctx context.Context, accessToken string) error {
	accessToken = strings.TrimSpace(accessToken)
	if accessToken == "" {
		return domainCLI.ErrCredentialNotFound
	}

	credential, err := service.credentials.FindByTokenHash(ctx, hash(accessToken))
	if err != nil {
		return err
	}

	credential.Revoke(time.Now())
	return service.credentials.Update(ctx, credential)
}

func (service *Service) revokeCredential(ctx context.Context, credential *domainCLI.Credential) {
	credential.Revoke(time.Now())
	_ = service.credentials.Update(ctx, credential)
}
