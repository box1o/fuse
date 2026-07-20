package credit

import (
	"context"
	stdErrors "errors"
	"testing"
	"time"

	domain "fuse/internal/domain/credit"

	"github.com/google/uuid"
)

var errPackRepositoryFailure = stdErrors.New(
	"forced pack repository failure",
)

func TestService_ListActivePacks_ReturnsActivePacks(t *testing.T) {
	t.Parallel()

	firstPack := newTestPack(t, "credits_500", "500 Credits", 500, true)
	secondPack := newTestPack(
		t,
		"credits_1500",
		"1,500 Credits",
		1500,
		true,
	)

	repository := &fakePackRepository{
		activePacks: []*domain.Pack{
			firstPack,
			secondPack,
		},
	}

	service := NewService(nil, repository)

	packs, err := service.ListActivePacks(context.Background())
	if err != nil {
		t.Fatalf(
			"ListActivePacks() returned unexpected error: %v",
			err,
		)
	}

	if len(packs) != 2 {
		t.Fatalf("expected 2 packs, got %d", len(packs))
	}

	if packs[0].ID != firstPack.ID {
		t.Errorf(
			"expected first pack ID %s, got %s",
			firstPack.ID,
			packs[0].ID,
		)
	}

	if packs[1].ID != secondPack.ID {
		t.Errorf(
			"expected second pack ID %s, got %s",
			secondPack.ID,
			packs[1].ID,
		)
	}
}

func TestService_ListActivePacks_ReturnsRepositoryError(t *testing.T) {
	t.Parallel()

	repository := &fakePackRepository{
		listActiveErr: errPackRepositoryFailure,
	}

	service := NewService(nil, repository)

	_, err := service.ListActivePacks(context.Background())
	if !stdErrors.Is(err, errPackRepositoryFailure) {
		t.Fatalf(
			"expected repository error, got %v",
			err,
		)
	}
}

func TestService_GetActivePack_ReturnsActivePack(t *testing.T) {
	t.Parallel()

	pack := newTestPack(
		t,
		"credits_500",
		"500 Credits",
		500,
		true,
	)

	repository := &fakePackRepository{
		packsByID: map[uuid.UUID]*domain.Pack{
			pack.ID: pack,
		},
	}

	service := NewService(nil, repository)

	result, err := service.GetActivePack(
		context.Background(),
		pack.ID,
	)
	if err != nil {
		t.Fatalf(
			"GetActivePack() returned unexpected error: %v",
			err,
		)
	}

	if result.ID != pack.ID {
		t.Errorf(
			"expected pack ID %s, got %s",
			pack.ID,
			result.ID,
		)
	}

	if result.Credits != pack.Credits {
		t.Errorf(
			"expected %d credits, got %d",
			pack.Credits.Value(),
			result.Credits.Value(),
		)
	}
}

func TestService_GetActivePack_RejectsInactivePack(t *testing.T) {
	t.Parallel()

	pack := newTestPack(
		t,
		"credits_500",
		"500 Credits",
		500,
		false,
	)

	repository := &fakePackRepository{
		packsByID: map[uuid.UUID]*domain.Pack{
			pack.ID: pack,
		},
	}

	service := NewService(nil, repository)

	_, err := service.GetActivePack(
		context.Background(),
		pack.ID,
	)
	if !stdErrors.Is(err, domain.ErrPackInactive) {
		t.Fatalf(
			"expected inactive pack error, got %v",
			err,
		)
	}
}

func TestService_GetActivePack_RejectsEmptyID(t *testing.T) {
	t.Parallel()

	service := NewService(nil, &fakePackRepository{})

	_, err := service.GetActivePack(
		context.Background(),
		uuid.Nil,
	)
	if !stdErrors.Is(err, domain.ErrPackNotFound) {
		t.Fatalf(
			"expected pack not found error, got %v",
			err,
		)
	}
}

func TestService_GetActivePack_ReturnsRepositoryError(t *testing.T) {
	t.Parallel()

	packID := uuid.New()

	repository := &fakePackRepository{
		findByIDErr: errPackRepositoryFailure,
	}

	service := NewService(nil, repository)

	_, err := service.GetActivePack(
		context.Background(),
		packID,
	)
	if !stdErrors.Is(err, errPackRepositoryFailure) {
		t.Fatalf(
			"expected repository error, got %v",
			err,
		)
	}
}

type fakePackRepository struct {
	packsByID     map[uuid.UUID]*domain.Pack
	activePacks   []*domain.Pack
	findByIDErr   error
	listActiveErr error
}

func (repository *fakePackRepository) Create(
	context.Context,
	*domain.Pack,
) error {
	return nil
}

func (repository *fakePackRepository) FindByID(
	_ context.Context,
	id uuid.UUID,
) (*domain.Pack, error) {
	if repository.findByIDErr != nil {
		return nil, repository.findByIDErr
	}

	pack, exists := repository.packsByID[id]
	if !exists {
		return nil, domain.ErrPackNotFound
	}

	return pack, nil
}

func (repository *fakePackRepository) FindByCode(
	context.Context,
	string,
) (*domain.Pack, error) {
	return nil, domain.ErrPackNotFound
}

func (repository *fakePackRepository) ListActive(
	context.Context,
) ([]*domain.Pack, error) {
	if repository.listActiveErr != nil {
		return nil, repository.listActiveErr
	}

	return repository.activePacks, nil
}

func (repository *fakePackRepository) Update(
	context.Context,
	*domain.Pack,
) error {
	return nil
}

func newTestPack(
	t *testing.T,
	code string,
	name string,
	credits int64,
	active bool,
) *domain.Pack {
	t.Helper()

	amount, err := domain.NewAmount(credits)
	if err != nil {
		t.Fatalf("create test amount: %v", err)
	}

	now := time.Now().UTC()

	return &domain.Pack{
		ID:        uuid.New(),
		Code:      code,
		Name:      name,
		Credits:   amount,
		Active:    active,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
