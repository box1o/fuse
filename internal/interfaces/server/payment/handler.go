package payment

import (
	"context"

	"fuse/internal/interfaces/server/middleware"
	paymentService "fuse/internal/services/payment"

	"github.com/go-chi/chi/v5"
)

type CheckoutService interface {
	CreateCheckout(
		ctx context.Context,
		input paymentService.CreateCheckoutInput,
	) (*paymentService.CreateCheckoutOutput, error)
}

type WebhookService interface {
	HandleWebhook(
		ctx context.Context,
		event *paymentService.WebhookEvent,
	) error
}

type Handler struct {
	checkoutService CheckoutService
	webhookService  WebhookService
	webhookParser   paymentService.WebhookParser
}

func NewHandler(checkoutService CheckoutService, webhookService WebhookService, webhookParser paymentService.WebhookParser) *Handler {
	return &Handler{
		checkoutService: checkoutService,
		webhookService:  webhookService,
		webhookParser:   webhookParser,
	}
}

func (h *Handler) RegisterRoutes(
	router chi.Router,
	authMiddleware *middleware.AuthMiddleware,
) {
	router.Post("/payments/webhook", h.HandleWebhook)
	router.Route("/payments", func(router chi.Router) {
		router.Use(authMiddleware.RequireAuth)

		router.Post("/checkout", h.CreateCheckout)
	})
}
