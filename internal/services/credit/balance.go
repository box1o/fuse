package credit

import (
	"context"
	stdErrors "errors"

	domain "fuse/internal/domain/credit"

	"github.com/google/uuid"
)

func (s *Service) GetBalance(
	ctx context.Context,
	ownerID uuid.UUID,
) (domain.Amount, error) {
	zeroBalance, err := domain.NewAmount(0)
	if err != nil {
		return zeroBalance, err
	}

	if ownerID == uuid.Nil {
		return zeroBalance, domain.ErrOwnerIDRequired
	}

	account, err := s.accounts.FindByOwnerID(ctx, ownerID)
	if stdErrors.Is(err, domain.ErrAccountNotFound) {
		return zeroBalance, nil
	}

	if err != nil {
		return zeroBalance, err
	}

	return account.Balance, nil
}
