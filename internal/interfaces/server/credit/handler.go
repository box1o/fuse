package credit

import (
	"context"

	domainCredit "fuse/internal/domain/credit"
	"fuse/internal/interfaces/server/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type PackService interface {
	ListActivePacks(
		ctx context.Context,
	) ([]*domainCredit.Pack, error)
}

type BalanceService interface {
	GetBalance(
		ctx context.Context,
		ownerID uuid.UUID,
	) (domainCredit.Amount, error)
}

type Handler struct {
	packService    PackService
	balanceService BalanceService
}

func NewHandler(packService PackService, balanceService BalanceService) *Handler {
	return &Handler{
		packService:    packService,
		balanceService: balanceService,
	}
}

func (handler *Handler) RegisterRoutes(router chi.Router, authMiddleware *middleware.AuthMiddleware) {
	router.Group(func(router chi.Router) {
		router.Use(authMiddleware.RequireAuth)

		router.Get(
			"/credit-packs",
			handler.ListActivePacks,
		)

		router.Get(
			"/credits/balance",
			handler.GetBalance,
		)
	})
}
