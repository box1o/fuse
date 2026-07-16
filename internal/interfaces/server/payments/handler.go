package payments

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	domain "fuse/internal/domain/payments"
	"fuse/internal/interfaces/server/middleware"
	service "fuse/internal/services/payments"
	"fuse/pkg/config"
	"fuse/pkg/errors"
	"fuse/pkg/log"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct {
	cfg *config.Config
	svc *service.Service
}

func NewHandler(cfg *config.Config, svc *service.Service) *Handler {
	return &Handler{cfg: cfg, svc: svc}
}

func (h *Handler) RegisterRoutes(r chi.Router, authMiddleware *middleware.AuthMiddleware) {
	r.Route("/payments", func(r chi.Router) {
		r.Use(authMiddleware.RequireAuth)
		r.Post("/checkout", h.CreateCheckoutSession)
		r.Post("/usage", h.RecordUsage)
		r.Delete("/subscription", h.CancelSubscription)
	})

	r.Post("/payments/webhook", h.HandleWebhook)
}

type checkoutRequest struct {
	PlanID     string `json:"plan_id"`
	SuccessURL string `json:"success_url"`
	CancelURL  string `json:"cancel_url"`
}

func (h *Handler) CreateCheckoutSession(w http.ResponseWriter, r *http.Request) {
	var req checkoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.WriteError(w, errors.ErrBadRequest.WithDetail("invalid request payload"))
		return
	}

	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		errors.WriteError(w, errors.ErrUnauthorized.WithDetail("user not authenticated"))
		return
	}

	priceID, err := h.resolvePriceID(strings.TrimSpace(req.PlanID))
	if err != nil {
		errors.WriteError(w, err)
		return
	}

	result, svcErr := h.svc.CreateCheckoutSession(r.Context(), userID, req.SuccessURL, req.CancelURL, priceID)
	if svcErr != nil {
		log.Warn("create checkout session failed: %v", svcErr)
		errors.WriteError(w, errors.ToHTTP(svcErr))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(result)
}

type usageRequest struct {
	ResourceType   string    `json:"resource_type"`
	Quantity       int64     `json:"quantity"`
	OccurredAt     time.Time `json:"occurred_at"`
	IdempotencyKey string    `json:"idempotency_key"`
}

func (h *Handler) RecordUsage(w http.ResponseWriter, r *http.Request) {
	var req usageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.WriteError(w, errors.ErrBadRequest.WithDetail("invalid request payload"))
		return
	}

	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		errors.WriteError(w, errors.ErrUnauthorized.WithDetail("user not authenticated"))
		return
	}

	record, err := h.svc.RecordUsage(
		r.Context(),
		userID,
		domain.ResourceType(strings.TrimSpace(req.ResourceType)),
		req.Quantity,
		req.OccurredAt,
		req.IdempotencyKey,
	)
	if err != nil {
		log.Warn("record usage failed: %v", err)
		errors.WriteError(w, errors.ToHTTP(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(record)
}

func (h *Handler) CancelSubscription(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		errors.WriteError(w, errors.ErrUnauthorized.WithDetail("user not authenticated"))
		return
	}

	if err := h.svc.CancelSubscription(r.Context(), userID); err != nil {
		log.Warn("cancel subscription failed: %v", err)
		errors.WriteError(w, errors.ToHTTP(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		errors.WriteError(w, errors.ErrBadRequest.WithDetail("failed to read request body"))
		return
	}
	defer r.Body.Close()

	signature := r.Header.Get("Stripe-Signature")
	if strings.TrimSpace(signature) == "" {
		errors.WriteError(w, errors.ErrBadRequest.WithDetail("missing Stripe signature"))
		return
	}

	if err := h.svc.HandleWebhook(r.Context(), body, signature); err != nil {
		log.Warn("stripe webhook failed: %v", err)
		errors.WriteError(w, errors.ToHTTP(err))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) resolvePriceID(planID string) (string, *errors.HTTPError) {
	if strings.TrimSpace(planID) == "" {
		return "", errors.ErrBadRequest.WithDetail("plan_id is required")
	}

	proPriceID := strings.TrimSpace(h.cfg.Stripe.ProPriceID)
	if proPriceID == "" {
		return "", errors.ErrBadRequest.WithDetail("stripe pro price id is not configured")
	}

	if strings.Contains(strings.ToLower(proPriceID), "placeholder") {
		return "", errors.ErrBadRequest.WithDetail("stripe pro price id is still a placeholder")
	}

	return proPriceID, nil
}
