package credit

import (
	"context"
	stdErrors "errors"
	"strings"

	domain "fuse/internal/domain/credit"

	"github.com/google/uuid"
)

type SpendInput struct {
	OwnerID        uuid.UUID
	Amount         domain.Amount
	ReferenceID    string
	IdempotencyKey string
}

func (s *Service) Spend(ctx context.Context, input SpendInput) error {
	if input.OwnerID == uuid.Nil {
		return domain.ErrOwnerIDRequired
	}

	if !input.Amount.IsPositive() {
		return domain.ErrAmountMustBePositive
	}

	if strings.TrimSpace(input.IdempotencyKey) == "" {
		return domain.ErrIdempotencyKeyRequired
	}

	err := s.unitOfWork.WithinTransaction(
		ctx,
		func(accounts domain.AccountRepository, transactions domain.TransactionRepository) error {
			account, err := accounts.FindByOwnerID(ctx, input.OwnerID)
			if err != nil {
				return err
			}

			if err := account.Spend(input.Amount); err != nil {
				return err
			}

			transaction, err := domain.NewTransaction(
				domain.NewTransactionInput{
					AccountID:      account.ID,
					Type:           domain.TransactionTypeSpend,
					Source:         domain.TransactionSourceComputeJob,
					Amount:         input.Amount,
					ReferenceID:    input.ReferenceID,
					IdempotencyKey: input.IdempotencyKey,
				},
			)
			if err != nil {
				return err
			}

			if err := accounts.Update(ctx, account); err != nil {
				return err
			}

			if err := transactions.Create(ctx, transaction); err != nil {
				return err
			}

			return nil
		},
	)

	if stdErrors.Is(err, domain.ErrTransactionAlreadyExists) {
		// Retrying the same compute-job charge must not spend twice.
		return nil
	}

	return err
}
