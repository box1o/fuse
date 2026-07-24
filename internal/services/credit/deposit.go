package credit

import (
	"context"
	stdErrors "errors"
	"strings"

	domain "fuse/internal/domain/credit"

	"github.com/google/uuid"
)

type DepositInput struct {
	OwnerID           uuid.UUID
	Amount            domain.Amount
	ReferenceID       string
	ExternalReference string
	IdempotencyKey    string
}

func (s *Service) Deposit(ctx context.Context, input DepositInput) error {
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
				if !stdErrors.Is(err, domain.ErrAccountNotFound) {
					return err
				}

				zeroBalance, amountErr := domain.NewAmount(0)
				if amountErr != nil {
					return amountErr
				}

				account, err = domain.NewAccount(input.OwnerID, zeroBalance)
				if err != nil {
					return err
				}

				if err := accounts.Create(ctx, account); err != nil {
					return err
				}
			}

			if err := account.Deposit(input.Amount); err != nil {
				return err
			}

			transaction, err := domain.NewTransaction(
				domain.NewTransactionInput{
					AccountID:         account.ID,
					Type:              domain.TransactionTypeDeposit,
					Source:            domain.TransactionSourcePurchase,
					Amount:            input.Amount,
					ReferenceID:       input.ReferenceID,
					ExternalReference: input.ExternalReference,
					IdempotencyKey:    input.IdempotencyKey,
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
		// The previous operation already completed successfully.
		// The database transaction rolls back the duplicate attempt.
		return nil
	}

	return err
}
