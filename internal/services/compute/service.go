package compute

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	stdErrors "errors"
	"strings"
	"time"

	domain "fuse/internal/domain/compute"

	"github.com/google/uuid"
)

type Service struct {
	repository domain.Repository
}

type UpdateNodeInput struct {
	Name     *string
	Disabled *bool
}

func NewService(repository domain.Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) RegisterNode(ctx context.Context, ownerID uuid.UUID, input domain.RegisterNodeInput) (*domain.Node, bool, error) {
	if err := domain.ValidateRegistration(ownerID, input); err != nil {
		return nil, false, err
	}

	existing, err := s.repository.FindByOwnerAndInstallationID(ctx, ownerID, input.InstallationID)
	if err == nil {
		if err := existing.ApplyRegistration(input); err != nil {
			return nil, false, err
		}
		if err := s.repository.Update(ctx, existing); err != nil {
			return nil, false, err
		}
		return existing, false, nil
	}
	if !stdErrors.Is(err, domain.ErrNodeNotFound) {
		return nil, false, err
	}

	node, err := domain.NewNode(ownerID, input)
	if err != nil {
		return nil, false, err
	}
	if err := s.repository.Create(ctx, node); err != nil {
		// Resolve a concurrent registration of the same owner/installation pair
		// as an idempotent update rather than returning a conflict.
		existing, findErr := s.repository.FindByOwnerAndInstallationID(ctx, ownerID, input.InstallationID)
		if findErr != nil {
			return nil, false, err
		}
		if applyErr := existing.ApplyRegistration(input); applyErr != nil {
			return nil, false, applyErr
		}
		if updateErr := s.repository.Update(ctx, existing); updateErr != nil {
			return nil, false, updateErr
		}
		return existing, false, nil
	}
	return node, true, nil
}

func (s *Service) ListNodes(ctx context.Context, ownerID uuid.UUID) ([]*domain.Node, error) {
	if ownerID == uuid.Nil {
		return nil, domain.ErrOwnerIDEmpty
	}
	return s.repository.ListByOwner(ctx, ownerID)
}

func (s *Service) GetNode(ctx context.Context, ownerID, nodeID uuid.UUID) (*domain.Node, error) {
	if ownerID == uuid.Nil || nodeID == uuid.Nil {
		return nil, domain.ErrNodeNotFound
	}
	return s.repository.FindByID(ctx, ownerID, nodeID)
}

func (s *Service) UpdateNode(ctx context.Context, ownerID, nodeID uuid.UUID, input UpdateNodeInput) (*domain.Node, error) {
	node, err := s.GetNode(ctx, ownerID, nodeID)
	if err != nil {
		return nil, err
	}
	if input.Name != nil {
		if err := node.Rename(*input.Name); err != nil {
			return nil, err
		}
	}
	if input.Disabled != nil {
		node.SetDisabled(*input.Disabled)
	}
	if err := s.repository.Update(ctx, node); err != nil {
		return nil, err
	}
	return node, nil
}

func (s *Service) DeleteNode(ctx context.Context, ownerID, nodeID uuid.UUID) error {
	if ownerID == uuid.Nil || nodeID == uuid.Nil {
		return domain.ErrNodeNotFound
	}
	return s.repository.Delete(ctx, ownerID, nodeID)
}

func (s *Service) IssueCLICredential(ctx context.Context, ownerID uuid.UUID, name string) (string, *domain.CLICredential, error) {
	secret := make([]byte, 32)
	if _, err := rand.Read(secret); err != nil {
		return "", nil, err
	}
	token := "fuse_" + base64.RawURLEncoding.EncodeToString(secret)
	credential, err := domain.NewCLICredential(ownerID, strings.TrimSpace(name), hashToken(token), time.Now().UTC().Add(90*24*time.Hour))
	if err != nil {
		return "", nil, err
	}
	if err := s.repository.CreateCredential(ctx, credential); err != nil {
		return "", nil, err
	}
	return token, credential, nil
}

func (s *Service) AuthenticateCLIToken(ctx context.Context, token string) (*domain.CLICredential, error) {
	if !strings.HasPrefix(token, "fuse_") {
		return nil, domain.ErrCredentialNotFound
	}
	credential, err := s.repository.FindCredentialByHash(ctx, hashToken(token))
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	if credential.RevokedAt != nil {
		return nil, domain.ErrCredentialRevoked
	}
	if !credential.IsActive(now) {
		return nil, domain.ErrCredentialExpired
	}
	credential.LastUsedAt = &now
	credential.UpdatedAt = now
	if err := s.repository.UpdateCredential(ctx, credential); err != nil {
		return nil, err
	}
	return credential, nil
}

func (s *Service) RevokeCLICredential(ctx context.Context, credential *domain.CLICredential) error {
	if credential == nil {
		return domain.ErrInvalidCredential
	}
	now := time.Now().UTC()
	credential.RevokedAt = &now
	credential.UpdatedAt = now
	return s.repository.UpdateCredential(ctx, credential)
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
