package credit

import (
	"context"

	domain "fuse/internal/domain/credit"

	"github.com/google/uuid"
)

// ListActivePacks returns all credit packs currently available for purchase.
func (s *Service) ListActivePacks(ctx context.Context) ([]*domain.Pack, error) {
	packs, err := s.packs.ListActive(ctx)
	if err != nil {
		return nil, err
	}

	return packs, nil
}

// GetActivePack returns a purchasable credit pack.
func (s *Service) GetActivePack(ctx context.Context, packID uuid.UUID) (*domain.Pack, error) {
	if packID == uuid.Nil {
		return nil, domain.ErrPackNotFound
	}

	pack, err := s.packs.FindByID(ctx, packID)
	if err != nil {
		return nil, err
	}

	if !pack.Active {
		return nil, domain.ErrPackInactive
	}

	return pack, nil
}
