package auth

import (
	"context"
	"errors"
	"fmt"
)

func (s *Service) Logout(ctx context.Context) error {
	credential, err := s.credentials.Load(ctx)
	if err != nil {
		loadErr := fmt.Errorf("load credential: %w", err)
		if deleteErr := s.credentials.Delete(ctx); deleteErr != nil {
			return errors.Join(loadErr, fmt.Errorf("delete local credential: %w", deleteErr))
		}
		return loadErr
	}
	if credential.AccessToken == "" {
		return ErrNotAuthenticated
	}

	var revokeErr error
	if err := s.gateway.Logout(ctx, credential.AccessToken); err != nil {
		revokeErr = fmt.Errorf("revoke credential: %w", err)
	}

	var deleteErr error
	if err := s.credentials.Delete(ctx); err != nil {
		deleteErr = fmt.Errorf("delete local credential: %w", err)
	}

	return errors.Join(revokeErr, deleteErr)
}
