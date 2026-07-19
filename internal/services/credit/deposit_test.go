package credit

import (
	"context"
	stdErrors "errors"
	"sort"
	"testing"

	domain "fuse/internal/domain/credit"

	"github.com/google/uuid"
)

var errForcedTransactionFailure = stdErrors.New(
	"forced transaction repository failure",
)

func TestService_Deposit_CreatesAccountAndTransaction(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ownerID := uuid.New()

	unitOfWork := newFakeUnitOfWork()
	service := NewService(unitOfWork, nil)

	amount := mustAmount(t, 500)

	input := DepositInput{
		OwnerID:           ownerID,
		Amount:            amount,
		ReferenceID:       "payment-123",
		ExternalReference: "cs_test_123",
		IdempotencyKey:    "stripe:event:evt_123",
	}

	err := service.Deposit(ctx, input)
	if err != nil {
		t.Fatalf("Deposit() returned unexpected error: %v", err)
	}

	account, err := unitOfWork.state.findAccountByOwnerID(ownerID)
	if err != nil {
		t.Fatalf("expected credit account to exist: %v", err)
	}

	if account.Balance != amount {
		t.Fatalf(
			"expected balance %d, got %d",
			amount.Value(),
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

	if transaction.Type != domain.TransactionTypeDeposit {
		t.Errorf(
			"expected transaction type %q, got %q",
			domain.TransactionTypeDeposit,
			transaction.Type,
		)
	}

	if transaction.Source != domain.TransactionSourcePurchase {
		t.Errorf(
			"expected transaction source %q, got %q",
			domain.TransactionSourcePurchase,
			transaction.Source,
		)
	}

	if transaction.Amount != amount {
		t.Errorf(
			"expected transaction amount %d, got %d",
			amount.Value(),
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

	if transaction.ExternalReference != input.ExternalReference {
		t.Errorf(
			"expected external reference %q, got %q",
			input.ExternalReference,
			transaction.ExternalReference,
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

func TestService_Deposit_DuplicateIdempotencyKeyDoesNotDepositTwice(
	t *testing.T,
) {
	t.Parallel()

	ctx := context.Background()
	ownerID := uuid.New()

	unitOfWork := newFakeUnitOfWork()
	service := NewService(unitOfWork, nil)

	input := DepositInput{
		OwnerID:           ownerID,
		Amount:            mustAmount(t, 500),
		ReferenceID:       "payment-123",
		ExternalReference: "cs_test_123",
		IdempotencyKey:    "stripe:event:evt_123",
	}

	if err := service.Deposit(ctx, input); err != nil {
		t.Fatalf("first Deposit() returned unexpected error: %v", err)
	}

	if err := service.Deposit(ctx, input); err != nil {
		t.Fatalf("duplicate Deposit() returned unexpected error: %v", err)
	}

	account, err := unitOfWork.state.findAccountByOwnerID(ownerID)
	if err != nil {
		t.Fatalf("expected credit account to exist: %v", err)
	}

	if account.Balance.Value() != 500 {
		t.Fatalf(
			"expected duplicate deposit to preserve balance 500, got %d",
			account.Balance.Value(),
		)
	}

	transactions := unitOfWork.state.transactionsForAccount(account.ID)
	if len(transactions) != 1 {
		t.Fatalf(
			"expected exactly 1 transaction after duplicate deposit, got %d",
			len(transactions),
		)
	}
}

func TestService_Deposit_TransactionFailureRollsBackAccountAndBalance(
	t *testing.T,
) {
	t.Parallel()

	ctx := context.Background()
	ownerID := uuid.New()

	unitOfWork := newFakeUnitOfWork()
	unitOfWork.transactionCreateErr = errForcedTransactionFailure

	service := NewService(unitOfWork, nil)

	err := service.Deposit(ctx, DepositInput{
		OwnerID:        ownerID,
		Amount:         mustAmount(t, 500),
		ReferenceID:    "payment-123",
		IdempotencyKey: "stripe:event:evt_123",
	})
	if !stdErrors.Is(err, errForcedTransactionFailure) {
		t.Fatalf(
			"expected transaction failure, got %v",
			err,
		)
	}

	_, err = unitOfWork.state.findAccountByOwnerID(ownerID)
	if !stdErrors.Is(err, domain.ErrAccountNotFound) {
		t.Fatalf(
			"expected account creation to be rolled back, got %v",
			err,
		)
	}

	if len(unitOfWork.state.transactions) != 0 {
		t.Fatalf(
			"expected no transactions after rollback, got %d",
			len(unitOfWork.state.transactions),
		)
	}
}

func TestService_Deposit_TransactionFailureRollsBackExistingBalance(
	t *testing.T,
) {
	t.Parallel()

	ctx := context.Background()
	ownerID := uuid.New()

	unitOfWork := newFakeUnitOfWork()
	unitOfWork.seedAccount(t, ownerID, 100)
	unitOfWork.transactionCreateErr = errForcedTransactionFailure

	service := NewService(unitOfWork, nil)

	err := service.Deposit(ctx, DepositInput{
		OwnerID:        ownerID,
		Amount:         mustAmount(t, 50),
		ReferenceID:    "payment-456",
		IdempotencyKey: "stripe:event:evt_456",
	})
	if !stdErrors.Is(err, errForcedTransactionFailure) {
		t.Fatalf(
			"expected transaction failure, got %v",
			err,
		)
	}

	account, err := unitOfWork.state.findAccountByOwnerID(ownerID)
	if err != nil {
		t.Fatalf("expected existing account to remain: %v", err)
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

func TestService_Deposit_ValidatesInput(t *testing.T) {
	t.Parallel()

	validOwnerID := uuid.New()
	validAmount := mustAmount(t, 100)

	tests := []struct {
		name     string
		input    DepositInput
		expected error
	}{
		{
			name: "missing owner ID",
			input: DepositInput{
				Amount:         validAmount,
				IdempotencyKey: "deposit-1",
			},
			expected: domain.ErrOwnerIDRequired,
		},
		{
			name: "zero amount",
			input: DepositInput{
				OwnerID:        validOwnerID,
				Amount:         domain.Amount(0),
				IdempotencyKey: "deposit-1",
			},
			expected: domain.ErrAmountMustBePositive,
		},
		{
			name: "missing idempotency key",
			input: DepositInput{
				OwnerID: validOwnerID,
				Amount:  validAmount,
			},
			expected: domain.ErrIdempotencyKeyRequired,
		},
		{
			name: "blank idempotency key",
			input: DepositInput{
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

			err := service.Deposit(context.Background(), test.input)
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

type fakeUnitOfWork struct {
	state                *fakeCreditState
	transactionCreateErr error
}

func newFakeUnitOfWork() *fakeUnitOfWork {
	return &fakeUnitOfWork{
		state: newFakeCreditState(),
	}
}

func (uow *fakeUnitOfWork) WithinTransaction(
	ctx context.Context,
	operation func(
		accounts domain.AccountRepository,
		transactions domain.TransactionRepository,
	) error,
) error {
	transactionState := uow.state.clone()

	accountRepository := &fakeAccountRepository{
		state: transactionState,
	}

	transactionRepository := &fakeTransactionRepository{
		state:     transactionState,
		createErr: uow.transactionCreateErr,
	}

	if err := operation(
		accountRepository,
		transactionRepository,
	); err != nil {
		return err
	}

	uow.state = transactionState
	return nil
}

func (uow *fakeUnitOfWork) seedAccount(
	t *testing.T,
	ownerID uuid.UUID,
	balanceValue int64,
) {
	t.Helper()

	balance := mustAmount(t, balanceValue)

	account, err := domain.NewAccount(ownerID, balance)
	if err != nil {
		t.Fatalf("failed to create seed account: %v", err)
	}

	uow.state.accounts[ownerID] = cloneAccount(account)
}

type fakeCreditState struct {
	accounts           map[uuid.UUID]*domain.Account
	transactions       map[uuid.UUID]*domain.Transaction
	transactionIDByKey map[string]uuid.UUID
}

func newFakeCreditState() *fakeCreditState {
	return &fakeCreditState{
		accounts:           make(map[uuid.UUID]*domain.Account),
		transactions:       make(map[uuid.UUID]*domain.Transaction),
		transactionIDByKey: make(map[string]uuid.UUID),
	}
}

func (state *fakeCreditState) clone() *fakeCreditState {
	cloned := newFakeCreditState()

	for ownerID, account := range state.accounts {
		cloned.accounts[ownerID] = cloneAccount(account)
	}

	for transactionID, transaction := range state.transactions {
		cloned.transactions[transactionID] = cloneTransaction(transaction)
	}

	for idempotencyKey, transactionID := range state.transactionIDByKey {
		cloned.transactionIDByKey[idempotencyKey] = transactionID
	}

	return cloned
}

func (state *fakeCreditState) findAccountByOwnerID(
	ownerID uuid.UUID,
) (*domain.Account, error) {
	account, exists := state.accounts[ownerID]
	if !exists {
		return nil, domain.ErrAccountNotFound
	}

	return cloneAccount(account), nil
}

func (state *fakeCreditState) transactionsForAccount(
	accountID uuid.UUID,
) []*domain.Transaction {
	transactions := make([]*domain.Transaction, 0)

	for _, transaction := range state.transactions {
		if transaction.AccountID == accountID {
			transactions = append(
				transactions,
				cloneTransaction(transaction),
			)
		}
	}

	sort.Slice(
		transactions,
		func(left, right int) bool {
			return transactions[left].CreatedAt.Before(
				transactions[right].CreatedAt,
			)
		},
	)

	return transactions
}

type fakeAccountRepository struct {
	state *fakeCreditState
}

func (repository *fakeAccountRepository) Create(
	_ context.Context,
	account *domain.Account,
) error {
	if account == nil {
		return domain.ErrInvalidAccount
	}

	if _, exists := repository.state.accounts[account.OwnerID]; exists {
		return domain.ErrInvalidAccount
	}

	repository.state.accounts[account.OwnerID] = cloneAccount(account)
	return nil
}

func (repository *fakeAccountRepository) FindByOwnerID(
	_ context.Context,
	ownerID uuid.UUID,
) (*domain.Account, error) {
	return repository.state.findAccountByOwnerID(ownerID)
}

func (repository *fakeAccountRepository) Update(
	_ context.Context,
	account *domain.Account,
) error {
	if account == nil {
		return domain.ErrInvalidAccount
	}

	if _, exists := repository.state.accounts[account.OwnerID]; !exists {
		return domain.ErrAccountNotFound
	}

	repository.state.accounts[account.OwnerID] = cloneAccount(account)
	return nil
}

type fakeTransactionRepository struct {
	state     *fakeCreditState
	createErr error
}

func (repository *fakeTransactionRepository) Create(
	_ context.Context,
	transaction *domain.Transaction,
) error {
	if transaction == nil {
		return domain.ErrInvalidTransaction
	}

	if repository.createErr != nil {
		return repository.createErr
	}

	if _, exists := repository.state.transactionIDByKey[transaction.IdempotencyKey]; exists {
		return domain.ErrTransactionAlreadyExists
	}

	repository.state.transactions[transaction.ID] = cloneTransaction(
		transaction,
	)
	repository.state.transactionIDByKey[transaction.IdempotencyKey] = transaction.ID

	return nil
}

func (repository *fakeTransactionRepository) FindByID(
	_ context.Context,
	id uuid.UUID,
) (*domain.Transaction, error) {
	transaction, exists := repository.state.transactions[id]
	if !exists {
		return nil, domain.ErrInvalidTransaction
	}

	return cloneTransaction(transaction), nil
}

func (repository *fakeTransactionRepository) FindByIdempotencyKey(
	_ context.Context,
	idempotencyKey string,
) (*domain.Transaction, error) {
	transactionID, exists := repository.state.transactionIDByKey[idempotencyKey]
	if !exists {
		return nil, domain.ErrInvalidTransaction
	}

	transaction, exists := repository.state.transactions[transactionID]
	if !exists {
		return nil, domain.ErrInvalidTransaction
	}

	return cloneTransaction(transaction), nil
}

func (repository *fakeTransactionRepository) ListByAccountID(
	_ context.Context,
	accountID uuid.UUID,
	limit int,
	offset int,
) ([]*domain.Transaction, error) {
	transactions := repository.state.transactionsForAccount(accountID)

	if offset < 0 {
		offset = 0
	}

	if offset >= len(transactions) {
		return []*domain.Transaction{}, nil
	}

	transactions = transactions[offset:]

	if limit > 0 && limit < len(transactions) {
		transactions = transactions[:limit]
	}

	return transactions, nil
}

func mustAmount(t *testing.T, value int64) domain.Amount {
	t.Helper()

	amount, err := domain.NewAmount(value)
	if err != nil {
		t.Fatalf("failed to create amount %d: %v", value, err)
	}

	return amount
}

func cloneAccount(account *domain.Account) *domain.Account {
	if account == nil {
		return nil
	}

	cloned := *account
	return &cloned
}

func cloneTransaction(
	transaction *domain.Transaction,
) *domain.Transaction {
	if transaction == nil {
		return nil
	}

	cloned := *transaction
	return &cloned
}
