package credit

import (
	"context"
	stdErrors "errors"
	"testing"

	domain "fuse/internal/domain/credit"

	"github.com/google/uuid"
)

func TestService_Spend_DecreasesBalanceAndCreatesTransaction(
	t *testing.T,
) {
	t.Parallel()

	ctx := context.Background()
	ownerID := uuid.New()

	unitOfWork := newFakeUnitOfWork()
	unitOfWork.seedAccount(t, ownerID, 100)

	service := NewService(unitOfWork, nil)

	input := SpendInput{
		OwnerID:        ownerID,
		Amount:         mustAmount(t, 30),
		ReferenceID:    "compute-job-123",
		IdempotencyKey: "compute-job:123:credit-charge",
	}

	err := service.Spend(ctx, input)
	if err != nil {
		t.Fatalf("Spend() returned unexpected error: %v", err)
	}

	account, err := unitOfWork.state.findAccountByOwnerID(ownerID)
	if err != nil {
		t.Fatalf("expected credit account to exist: %v", err)
	}

	if account.Balance.Value() != 70 {
		t.Fatalf(
			"expected balance 70, got %d",
			account.Balance.Value(),
		)
	}

	transactions := unitOfWork.state.transactionsForAccount(account.ID)
	if len(transactions) != 1 {
		t.Fatalf(
			"expected 1 transaction, got %d",
			len(transactions),
		)
	}

	transaction := transactions[0]

	if transaction.Type != domain.TransactionTypeSpend {
		t.Errorf(
			"expected transaction type %q, got %q",
			domain.TransactionTypeSpend,
			transaction.Type,
		)
	}

	if transaction.Source != domain.TransactionSourceComputeJob {
		t.Errorf(
			"expected transaction source %q, got %q",
			domain.TransactionSourceComputeJob,
			transaction.Source,
		)
	}

	if transaction.Amount.Value() != 30 {
		t.Errorf(
			"expected transaction amount 30, got %d",
			transaction.Amount.Value(),
		)
	}

	if transaction.ReferenceID != input.ReferenceID {
		t.Errorf(
			"expected reference ID %q, got %q",
			input.ReferenceID,
			transaction.ReferenceID,
		)
	}

	if transaction.IdempotencyKey != input.IdempotencyKey {
		t.Errorf(
			"expected idempotency key %q, got %q",
			input.IdempotencyKey,
			transaction.IdempotencyKey,
		)
	}
}

func TestService_Spend_DuplicateIdempotencyKeyDoesNotSpendTwice(
	t *testing.T,
) {
	t.Parallel()

	ctx := context.Background()
	ownerID := uuid.New()

	unitOfWork := newFakeUnitOfWork()
	unitOfWork.seedAccount(t, ownerID, 100)

	service := NewService(unitOfWork, nil)

	input := SpendInput{
		OwnerID:        ownerID,
		Amount:         mustAmount(t, 30),
		ReferenceID:    "compute-job-123",
		IdempotencyKey: "compute-job:123:credit-charge",
	}

	if err := service.Spend(ctx, input); err != nil {
		t.Fatalf("first Spend() returned unexpected error: %v", err)
	}

	if err := service.Spend(ctx, input); err != nil {
		t.Fatalf("duplicate Spend() returned unexpected error: %v", err)
	}

	account, err := unitOfWork.state.findAccountByOwnerID(ownerID)
	if err != nil {
		t.Fatalf("expected credit account to exist: %v", err)
	}

	if account.Balance.Value() != 70 {
		t.Fatalf(
			"expected duplicate spend to preserve balance 70, got %d",
			account.Balance.Value(),
		)
	}

	transactions := unitOfWork.state.transactionsForAccount(account.ID)
	if len(transactions) != 1 {
		t.Fatalf(
			"expected exactly 1 transaction after duplicate spend, got %d",
			len(transactions),
		)
	}
}

func TestService_Spend_InsufficientBalanceRollsBackOperation(
	t *testing.T,
) {
	t.Parallel()

	ctx := context.Background()
	ownerID := uuid.New()

	unitOfWork := newFakeUnitOfWork()
	unitOfWork.seedAccount(t, ownerID, 20)

	service := NewService(unitOfWork, nil)

	err := service.Spend(ctx, SpendInput{
		OwnerID:        ownerID,
		Amount:         mustAmount(t, 30),
		ReferenceID:    "compute-job-123",
		IdempotencyKey: "compute-job:123:credit-charge",
	})
	if !stdErrors.Is(err, domain.ErrInsufficientBalance) {
		t.Fatalf(
			"expected insufficient balance error, got %v",
			err,
		)
	}

	account, err := unitOfWork.state.findAccountByOwnerID(ownerID)
	if err != nil {
		t.Fatalf("expected credit account to exist: %v", err)
	}

	if account.Balance.Value() != 20 {
		t.Fatalf(
			"expected balance to remain 20, got %d",
			account.Balance.Value(),
		)
	}

	if len(unitOfWork.state.transactions) != 0 {
		t.Fatalf(
			"expected no transactions, got %d",
			len(unitOfWork.state.transactions),
		)
	}
}

func TestService_Spend_TransactionFailureRollsBackBalance(
	t *testing.T,
) {
	t.Parallel()

	ctx := context.Background()
	ownerID := uuid.New()

	unitOfWork := newFakeUnitOfWork()
	unitOfWork.seedAccount(t, ownerID, 100)
	unitOfWork.transactionCreateErr = errForcedTransactionFailure

	service := NewService(unitOfWork, nil)

	err := service.Spend(ctx, SpendInput{
		OwnerID:        ownerID,
		Amount:         mustAmount(t, 30),
		ReferenceID:    "compute-job-123",
		IdempotencyKey: "compute-job:123:credit-charge",
	})
	if !stdErrors.Is(err, errForcedTransactionFailure) {
		t.Fatalf(
			"expected transaction failure, got %v",
			err,
		)
	}

	account, err := unitOfWork.state.findAccountByOwnerID(ownerID)
	if err != nil {
		t.Fatalf("expected credit account to exist: %v", err)
	}

	if account.Balance.Value() != 100 {
		t.Fatalf(
			"expected balance to roll back to 100, got %d",
			account.Balance.Value(),
		)
	}

	if len(unitOfWork.state.transactions) != 0 {
		t.Fatalf(
			"expected no transactions after rollback, got %d",
			len(unitOfWork.state.transactions),
		)
	}
}

func TestService_Spend_AccountNotFound(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ownerID := uuid.New()

	unitOfWork := newFakeUnitOfWork()
	service := NewService(unitOfWork, nil)

	err := service.Spend(ctx, SpendInput{
		OwnerID:        ownerID,
		Amount:         mustAmount(t, 30),
		ReferenceID:    "compute-job-123",
		IdempotencyKey: "compute-job:123:credit-charge",
	})
	if !stdErrors.Is(err, domain.ErrAccountNotFound) {
		t.Fatalf(
			"expected account-not-found error, got %v",
			err,
		)
	}
}

func TestService_Spend_ValidatesInput(t *testing.T) {
	t.Parallel()

	validOwnerID := uuid.New()
	validAmount := mustAmount(t, 30)

	tests := []struct {
		name     string
		input    SpendInput
		expected error
	}{
		{
			name: "missing owner ID",
			input: SpendInput{
				Amount:         validAmount,
				IdempotencyKey: "spend-1",
			},
			expected: domain.ErrOwnerIDRequired,
		},
		{
			name: "zero amount",
			input: SpendInput{
				OwnerID:        validOwnerID,
				Amount:         domain.Amount(0),
				IdempotencyKey: "spend-1",
			},
			expected: domain.ErrAmountMustBePositive,
		},
		{
			name: "missing idempotency key",
			input: SpendInput{
				OwnerID: validOwnerID,
				Amount:  validAmount,
			},
			expected: domain.ErrIdempotencyKeyRequired,
		},
		{
			name: "blank idempotency key",
			input: SpendInput{
				OwnerID:        validOwnerID,
				Amount:         validAmount,
				IdempotencyKey: "   ",
			},
			expected: domain.ErrIdempotencyKeyRequired,
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			service := NewService(newFakeUnitOfWork(), nil)

			err := service.Spend(context.Background(), test.input)
			if !stdErrors.Is(err, test.expected) {
				t.Fatalf(
					"expected error %v, got %v",
					test.expected,
					err,
				)
			}
		})
	}
}
