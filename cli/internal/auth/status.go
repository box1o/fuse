package auth

import (
	"context"
	"fmt"
)

func (s *Service) Status(ctx context.Context) (Status, error) {
	credential, err := s.credentials.Load(ctx)
	if err != nil {
		return Status{}, fmt.Errorf("load credential: %w", err)
	}
	if credential.AccessToken == "" {
		return Status{}, ErrNotAuthenticated
	}

	status, err := s.gateway.Status(ctx, credential.AccessToken)
	if err != nil {
		return Status{}, fmt.Errorf("validate credential: %w", err)
	}
	if !status.Authenticated {
		return Status{}, ErrNotAuthenticated
	}

	return status, nil
}
