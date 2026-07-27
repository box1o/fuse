package credit

import (
	"context"

	"github.com/google/uuid"
)

type AccountRepository interface {
	Create(ctx context.Context, account *Account) error
	FindByOwnerID(ctx context.Context, ownerID uuid.UUID) (*Account, error)
	Update(ctx context.Context, account *Account) error
}

type PackRepository interface {
	Create(ctx context.Context, pack *Pack) error
	FindByID(ctx context.Context, id uuid.UUID) (*Pack, error)
	FindByCode(ctx context.Context, code string) (*Pack, error)
	ListActive(ctx context.Context) ([]*Pack, error)
	Update(ctx context.Context, pack *Pack) error
}

type TransactionRepository interface {
	Create(ctx context.Context, transaction *Transaction) error
	FindByID(ctx context.Context, id uuid.UUID) (*Transaction, error)
	FindByIdempotencyKey(ctx context.Context, idempotencyKey string) (*Transaction, error)
	ListByAccountID(ctx context.Context, accountID uuid.UUID, limit int, offset int) ([]*Transaction, error)
}
