package credit

import (
	"math"
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID        uuid.UUID `json:"id"`
	OwnerID   uuid.UUID `json:"owner_id"`
	Balance   Amount    `json:"balance"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

func NewAccount(ownerID uuid.UUID, initialBalance Amount) (*Account, error) {
	if ownerID == uuid.Nil {
		return nil, ErrOwnerIDRequired
	}

	if initialBalance.IsNegative() {
		return nil, ErrNegativeAmount
	}

	now := time.Now().UTC()
	return &Account{
		ID:        uuid.New(),
		OwnerID:   ownerID,
		Balance:   initialBalance,
		UpdatedAt: now,
		CreatedAt: now,
	}, nil
}

func (ac *Account) Spend(amount Amount) error {
	if ac == nil {
		return ErrInvalidAccount
	}

	if !amount.IsPositive() {
		return ErrAmountMustBePositive
	}

	if ac.Balance < amount {
		return ErrInsufficientBalance
	}

	ac.Balance -= amount
	ac.UpdatedAt = time.Now().UTC()

	return nil
}

func (ac *Account) Deposit(amount Amount) error {
	if ac == nil {
		return ErrInvalidAccount
	}
	if !amount.IsPositive() {
		return ErrAmountMustBePositive
	}

	if amount.Value() > math.MaxInt64-ac.Balance.Value() {
		return ErrBalanceOverflow
	}

	ac.Balance += amount
	ac.UpdatedAt = time.Now().UTC()

	return nil
}
