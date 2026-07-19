package credit

import (
	"context"

	domain "fuse/internal/domain/credit"
)

type UnitOfWork interface {
	WithinTransaction(
		ctx context.Context,
		operation func(
			accounts domain.AccountRepository,
			transactions domain.TransactionRepository,
		) error,
	) error
}

type Service struct {
	unitOfWork UnitOfWork
	packs      domain.PackRepository
}

func NewService(unitOfWork UnitOfWork, packs domain.PackRepository) *Service {
	return &Service{
		unitOfWork: unitOfWork,
		packs:      packs,
	}
}
