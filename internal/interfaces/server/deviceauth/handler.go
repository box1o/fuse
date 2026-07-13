package deviceauth

import (
	"encoding/json"
	stdErrors "errors"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"fuse/internal/interfaces/server/middleware"
	computeService "fuse/internal/services/compute"
	deviceService "fuse/internal/services/deviceauth"
	"fuse/pkg/deviceauthapi"
	appErrors "fuse/pkg/errors"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct {
	service    *deviceService.Service
	compute    *computeService.Service
	computeCLI *middleware.CLIMiddleware
}

func NewHandler(service *deviceService.Service, compute *computeService.Service, cliMiddleware *middleware.CLIMiddleware) *Handler {
	return &Handler{service: service, compute: compute, computeCLI: cliMiddleware}
}

func (h *Handler) RegisterRoutes(r chi.Router, browserAuth *middleware.AuthMiddleware) {
	r.Post("/auth/device/code", h.CreateCode)
	r.Post("/auth/device/token", h.ExchangeCode)
	r.With(browserAuth.RequireAuth).Get("/auth/device/request/{userCode}", h.GetRequest)
	r.With(browserAuth.RequireAuth).Post("/auth/device/approve", h.Approve)
	r.With(browserAuth.RequireAuth).Post("/auth/device/deny", h.Deny)
	r.With(h.computeCLI.RequireAuth).Get("/auth/cli/status", h.Status)
	r.With(h.computeCLI.RequireAuth).Post("/auth/cli/logout", h.Logout)
}

type codeRequest struct {
	UserCode string `json:"user_code"`
}

func (h *Handler) CreateCode(w http.ResponseWriter, r *http.Request) {
	if !h.service.Allow(r.Context(), "create:"+clientIP(r), 10, time.Minute) {
		appErrors.WriteError(w, appErrors.NewHTTP(http.StatusTooManyRequests, "RATE_LIMITED", "Too many device authorization requests"))
		return
	}
	var request deviceauthapi.CreateCodeRequest
	if err := decodeJSON(w, r, &request); err != nil {
		appErrors.WriteError(w, appErrors.ErrBadRequest.WithDetail(err.Error()))
		return
	}
	response, err := h.service.Create(r.Context(), request.ClientName)
	if err != nil {
		appErrors.WriteError(w, appErrors.ErrInternalServer.WithDetail("failed to create device authorization"))
		return
	}
	writeJSON(w, http.StatusOK, response)
}

func (h *Handler) ExchangeCode(w http.ResponseWriter, r *http.Request) {
	var request deviceauthapi.TokenRequest
	if err := decodeJSON(w, r, &request); err != nil || strings.TrimSpace(request.DeviceCode) == "" {
		appErrors.WriteError(w, appErrors.ErrBadRequest.WithDetail("device_code is required"))
		return
	}
	if !h.service.Allow(r.Context(), "poll:"+request.DeviceCode, 20, time.Minute) {
		appErrors.WriteError(w, appErrors.NewHTTP(http.StatusTooManyRequests, "SLOW_DOWN", "Device authorization is being polled too frequently"))
		return
	}
	response, err := h.service.Exchange(r.Context(), request.DeviceCode)
	if err != nil {
		switch {
		case stdErrors.Is(err, deviceService.ErrAuthorizationPending):
			writeJSON(w, http.StatusAccepted, map[string]string{"error": "authorization_pending"})
		case stdErrors.Is(err, deviceService.ErrAuthorizationDenied):
			writeJSON(w, http.StatusForbidden, map[string]string{"error": "access_denied"})
		case stdErrors.Is(err, deviceService.ErrAuthorizationExpired), stdErrors.Is(err, deviceService.ErrInvalidCode):
			writeJSON(w, http.StatusGone, map[string]string{"error": "expired_token"})
		default:
			appErrors.WriteError(w, appErrors.ErrInternalServer.WithDetail("failed to exchange device code"))
		}
		return
	}
	writeJSON(w, http.StatusOK, response)
}

func (h *Handler) GetRequest(w http.ResponseWriter, r *http.Request) {
	state, err := h.service.Inspect(r.Context(), chi.URLParam(r, "userCode"))
	if err != nil {
		appErrors.WriteError(w, appErrors.ErrNotFound.WithDetail("device authorization not found or expired"))
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"user_code": state.UserCode, "client_name": state.ClientName,
		"status": state.Status, "expires_at": state.ExpiresAt,
	})
}

func (h *Handler) Approve(w http.ResponseWriter, r *http.Request) {
	h.handleDecision(w, r, true)
}

func (h *Handler) Deny(w http.ResponseWriter, r *http.Request) {
	h.handleDecision(w, r, false)
}

func (h *Handler) handleDecision(w http.ResponseWriter, r *http.Request, approve bool) {
	ownerID := middleware.GetUserIDFromContext(r.Context())
	if ownerID == uuid.Nil {
		appErrors.WriteError(w, appErrors.ErrUnauthorized)
		return
	}
	var request codeRequest
	if err := decodeJSON(w, r, &request); err != nil || strings.TrimSpace(request.UserCode) == "" {
		appErrors.WriteError(w, appErrors.ErrBadRequest.WithDetail("user_code is required"))
		return
	}
	var err error
	if approve {
		err = h.service.Approve(r.Context(), request.UserCode, ownerID)
	} else {
		err = h.service.Deny(r.Context(), request.UserCode, ownerID)
	}
	if err != nil {
		appErrors.WriteError(w, appErrors.ErrBadRequest.WithDetail(err.Error()))
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": map[bool]string{true: "approved", false: "denied"}[approve]})
}

func (h *Handler) Status(w http.ResponseWriter, r *http.Request) {
	credential := middleware.GetCLICredential(r.Context())
	owner, err := h.service.Owner(r.Context(), credential.OwnerID)
	if err != nil {
		appErrors.WriteError(w, appErrors.ErrInternalServer.WithDetail("failed to load CLI owner"))
		return
	}
	writeJSON(w, http.StatusOK, deviceauthapi.StatusResponse{
		Authenticated: true, OwnerID: credential.OwnerID.String(), CredentialName: credential.Name,
		OwnerName: owner.Name, OwnerEmail: owner.Email, ExpiresAt: credential.ExpiresAt,
	})
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	credential := middleware.GetCLICredential(r.Context())
	if err := h.compute.RevokeCLICredential(r.Context(), credential); err != nil {
		appErrors.WriteError(w, appErrors.ErrInternalServer.WithDetail("failed to revoke CLI credential"))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func decodeJSON(w http.ResponseWriter, r *http.Request, target any) error {
	defer r.Body.Close()
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return err
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return appErrors.ErrBadRequest.WithDetail("request body must contain exactly one JSON object")
	}
	return nil
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}
	return r.RemoteAddr
}
