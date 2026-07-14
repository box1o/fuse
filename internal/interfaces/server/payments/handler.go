package payments

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	domain "fuse/internal/domain/payments"
	"fuse/internal/interfaces/server/middleware"
	service "fuse/internal/services/payments"
	"fuse/internal/services/workspace"
	"fuse/pkg/config"
	"fuse/pkg/errors"
	"fuse/pkg/log"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct {
	cfg         *config.Config
	svc         *service.Service
	workspaceSvc *workspace.Service
}

func NewHandler(cfg *config.Config, svc *service.Service, workspaceSvc *workspace.Service) *Handler {
	return &Handler{cfg: cfg, svc: svc, workspaceSvc: workspaceSvc}
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
	WorkspaceID uuid.UUID `json:"workspace_id"`
	ResourceType string   `json:"resource_type"`
	SuccessURL  string    `json:"success_url"`
	CancelURL   string    `json:"cancel_url"`
}

func (h *Handler) CreateCheckoutSession(w http.ResponseWriter, r *http.Request) {
	var req checkoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.WriteError(w, errors.ErrBadRequest.WithDetail("invalid request payload"))
		return
	}

	if err := h.ensureWorkspaceAccess(r.Context(), req.WorkspaceID); err != nil {
		errors.WriteError(w, err)
		return
	}

	priceID, err := h.resolvePriceID(strings.TrimSpace(req.ResourceType))
	if err != nil {
		errors.WriteError(w, err)
		return
	}

	result, svcErr := h.svc.CreateCheckoutSession(r.Context(), req.WorkspaceID, req.SuccessURL, req.CancelURL, priceID)
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
	WorkspaceID    uuid.UUID `json:"workspace_id"`
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

	if err := h.ensureWorkspaceAccess(r.Context(), req.WorkspaceID); err != nil {
		errors.WriteError(w, err)
		return
	}

	record, err := h.svc.RecordUsage(
		r.Context(),
		req.WorkspaceID,
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

type cancelRequest struct {
	WorkspaceID uuid.UUID `json:"workspace_id"`
}

func (h *Handler) CancelSubscription(w http.ResponseWriter, r *http.Request) {
	var req cancelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.WriteError(w, errors.ErrBadRequest.WithDetail("invalid request payload"))
		return
	}

	if err := h.ensureWorkspaceAccess(r.Context(), req.WorkspaceID); err != nil {
		errors.WriteError(w, err)
		return
	}

	if err := h.svc.CancelSubscription(r.Context(), req.WorkspaceID); err != nil {
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

func (h *Handler) ensureWorkspaceAccess(ctx context.Context, workspaceID uuid.UUID) *errors.HTTPError {
	if workspaceID == uuid.Nil {
		return errors.ErrBadRequest.WithDetail("workspace_id is required")
	}

	userID := middleware.GetUserIDFromContext(ctx)
	if userID == uuid.Nil {
		return errors.ErrUnauthorized.WithDetail("user not authenticated")
	}

	workspaces, err := h.workspaceSvc.GetUserWorkspaces(ctx, userID)
	if err != nil {
		log.Warn("workspace lookup failed: %v", err)
		return errors.ErrInternalServer.WithDetail("failed to verify workspace access")
	}

	for _, ws := range workspaces {
		if ws != nil && ws.ID == workspaceID {
			return nil
		}
	}

	return errors.ErrForbidden.WithDetail("workspace access denied")
}

func (h *Handler) resolvePriceID(resourceType string) (string, *errors.HTTPError) {
	switch resourceType {
	case "cpu":
		return h.cfg.Stripe.CPUPriceID, nil
	case "gpu":
		return h.cfg.Stripe.GPUPriceID, nil
	case "npu":
		return h.cfg.Stripe.NPUPriceID, nil
	default:
		return "", errors.ErrBadRequest.WithDetail("invalid resource_type")
	}
}
