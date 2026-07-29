package deviceauth

import (
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"fuse/internal/interfaces/server/middleware"
	"fuse/internal/services/deviceauth"
	"fuse/pkg/deviceauthapi"
	appErrors "fuse/pkg/errors"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct {
	service *deviceauth.Service
}

func NewHandler(service *deviceauth.Service) *Handler {
	return &Handler{service: service}
}

func (handler *Handler) RegisterRoutes(router chi.Router, browserAuth *middleware.AuthMiddleware) {
	router.Post("/auth/device/code", handler.CreateCode)
	router.Post("/auth/device/token", handler.ExchangeCode)

	router.With(browserAuth.RequireAuth).Get("/auth/device/request/{userCode}", handler.GetRequest)
	router.With(browserAuth.RequireAuth).Post("/auth/device/approve", handler.Approve)
	router.With(browserAuth.RequireAuth).Post("/auth/device/deny", handler.Deny)

	router.Get("/auth/cli/status", handler.Status)
	router.Post("/auth/cli/logout", handler.Logout)
}

// CreateCode starts a browser-based CLI authorization flow.
//
//	@Summary		Create a CLI device code
//	@Tags			CLI authentication
//	@Accept			json
//	@Produce		json
//	@Param			request	body		deviceauthapi.CreateCodeRequest	true	"CLI client details"
//	@Success		200		{object}	deviceauthapi.CodeResponse
//	@Failure		400		{object}	errors.HTTPError
//	@Failure		429		{object}	errors.HTTPError
//	@Router			/auth/device/code [post]
func (handler *Handler) CreateCode(writer http.ResponseWriter, request *http.Request) {
	if !handler.service.Allow(request.Context(), "create:"+clientIP(request), 10, 10*time.Minute) {
		appErrors.WriteError(writer, appErrors.NewHTTP(http.StatusTooManyRequests, "RATE_LIMITED", "Too many device authorization requests"))
		return
	}

	var input deviceauthapi.CreateCodeRequest
	if err := decodeJSON(writer, request, &input); err != nil {
		appErrors.WriteError(writer, appErrors.ErrBadRequest.WithDetail(err.Error()))
		return
	}

	response, err := handler.service.Create(request.Context(), input.ClientName)
	if err != nil {
		if errors.Is(err, deviceauth.ErrInvalidCode) {
			appErrors.WriteError(writer, appErrors.ErrBadRequest.WithDetail(err.Error()))
		} else {
			appErrors.WriteError(writer, appErrors.ErrInternalServer.WithDetail("failed to create device authorization"))
		}
		return
	}

	writeJSON(writer, http.StatusOK, response)
}

// ExchangeCode polls a device authorization and returns its CLI token when approved.
//
//	@Summary		Exchange a CLI device code
//	@Tags			CLI authentication
//	@Accept			json
//	@Produce		json
//	@Param			request	body		deviceauthapi.TokenRequest	true	"Device code"
//	@Success		200		{object}	deviceauthapi.TokenResponse
//	@Failure		202		{object}	map[string]string
//	@Failure		403		{object}	map[string]string
//	@Failure		410		{object}	map[string]string
//	@Failure		429		{object}	map[string]string
//	@Router			/auth/device/token [post]
func (handler *Handler) ExchangeCode(writer http.ResponseWriter, request *http.Request) {
	var input deviceauthapi.TokenRequest
	if err := decodeJSON(writer, request, &input); err != nil {
		appErrors.WriteError(writer, appErrors.ErrBadRequest.WithDetail(err.Error()))
		return
	}
	if strings.TrimSpace(input.DeviceCode) == "" {
		appErrors.WriteError(writer, appErrors.ErrBadRequest.WithDetail("device_code is required"))
		return
	}

	if !handler.service.Allow(request.Context(), "poll:"+input.DeviceCode, 20, time.Minute) {
		writeJSON(writer, http.StatusTooManyRequests, map[string]string{"error": "slow_down"})
		return
	}

	response, err := handler.service.Exchange(request.Context(), input.DeviceCode)
	if err != nil {
		handler.writeExchangeError(writer, err)
		return
	}

	writeJSON(writer, http.StatusOK, response)
}

// GetRequest returns a pending device request to an authenticated browser session.
//
//	@Summary		Inspect a CLI device request
//	@Tags			CLI authentication
//	@Produce		json
//	@Param			userCode	path		string	true	"One-time user code"
//	@Success		200			{object}	deviceauthapi.DeviceRequest
//	@Failure		401			{object}	errors.HTTPError
//	@Failure		404			{object}	errors.HTTPError
//	@Router			/auth/device/request/{userCode} [get]
func (handler *Handler) GetRequest(writer http.ResponseWriter, request *http.Request) {
	state, err := handler.service.Inspect(request.Context(), chi.URLParam(request, "userCode"))
	if err != nil {
		appErrors.WriteError(writer, appErrors.ErrNotFound.WithDetail("device authorization not found or expired"))
		return
	}

	writeJSON(writer, http.StatusOK, deviceauthapi.DeviceRequest{
		UserCode:   state.UserCode,
		ClientName: state.ClientName,
		Status:     state.Status,
		ExpiresAt:  state.ExpiresAt,
	})
}

// Approve authorizes a pending CLI device request.
//
//	@Summary		Approve a CLI device request
//	@Tags			CLI authentication
//	@Accept			json
//	@Produce		json
//	@Param			request	body		deviceauthapi.DecisionRequest	true	"One-time user code"
//	@Success		200		{object}	map[string]string
//	@Failure		401		{object}	errors.HTTPError
//	@Failure		409		{object}	errors.HTTPError
//	@Failure		410		{object}	errors.HTTPError
//	@Router			/auth/device/approve [post]
func (handler *Handler) Approve(writer http.ResponseWriter, request *http.Request) {
	handler.handleDecision(writer, request, true)
}

// Deny rejects a pending CLI device request.
//
//	@Summary		Deny a CLI device request
//	@Tags			CLI authentication
//	@Accept			json
//	@Produce		json
//	@Param			request	body		deviceauthapi.DecisionRequest	true	"One-time user code"
//	@Success		200		{object}	map[string]string
//	@Failure		401		{object}	errors.HTTPError
//	@Failure		409		{object}	errors.HTTPError
//	@Failure		410		{object}	errors.HTTPError
//	@Router			/auth/device/deny [post]
func (handler *Handler) Deny(writer http.ResponseWriter, request *http.Request) {
	handler.handleDecision(writer, request, false)
}

// Status validates a CLI bearer token and returns its owner.
//
//	@Summary		Get CLI authentication status
//	@Tags			CLI authentication
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	deviceauthapi.StatusResponse
//	@Failure		401	{object}	errors.HTTPError
//	@Router			/auth/cli/status [get]
func (handler *Handler) Status(writer http.ResponseWriter, request *http.Request) {
	accessToken, ok := bearerToken(request)
	if !ok {
		appErrors.WriteError(writer, appErrors.ErrUnauthorized.WithDetail("CLI bearer token is required"))
		return
	}

	credential, owner, err := handler.service.ValidateToken(request.Context(), accessToken)
	if err != nil {
		appErrors.WriteError(writer, appErrors.ErrUnauthorized.WithDetail("CLI credential is invalid, expired, or revoked"))
		return
	}

	writeJSON(writer, http.StatusOK, deviceauthapi.StatusResponse{
		Authenticated: true,
		OwnerID:       owner.ID.String(),
		OwnerName:     owner.Name,
		OwnerEmail:    owner.Email,
		ExpiresAt:     credential.ExpiresAt,
	})
}

// Logout revokes the current CLI bearer token.
//
//	@Summary		Revoke a CLI credential
//	@Tags			CLI authentication
//	@Security		BearerAuth
//	@Success		204
//	@Failure		401	{object}	errors.HTTPError
//	@Router			/auth/cli/logout [post]
func (handler *Handler) Logout(writer http.ResponseWriter, request *http.Request) {
	accessToken, ok := bearerToken(request)
	if !ok {
		appErrors.WriteError(writer, appErrors.ErrUnauthorized.WithDetail("CLI bearer token is required"))
		return
	}

	if err := handler.service.RevokeToken(request.Context(), accessToken); err != nil {
		appErrors.WriteError(writer, appErrors.ErrUnauthorized.WithDetail("CLI credential is invalid, expired, or revoked"))
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}

func (handler *Handler) handleDecision(writer http.ResponseWriter, request *http.Request, approve bool) {
	ownerID := middleware.GetUserIDFromContext(request.Context())
	if ownerID == uuid.Nil {
		appErrors.WriteError(writer, appErrors.ErrUnauthorized)
		return
	}

	var input deviceauthapi.DecisionRequest
	if err := decodeJSON(writer, request, &input); err != nil {
		appErrors.WriteError(writer, appErrors.ErrBadRequest.WithDetail(err.Error()))
		return
	}
	if strings.TrimSpace(input.UserCode) == "" {
		appErrors.WriteError(writer, appErrors.ErrBadRequest.WithDetail("user_code is required"))
		return
	}

	var err error
	status := "denied"
	if approve {
		err = handler.service.Approve(request.Context(), input.UserCode, ownerID)
		status = "approved"
	} else {
		err = handler.service.Deny(request.Context(), input.UserCode, ownerID)
	}
	if err != nil {
		switch {
		case errors.Is(err, deviceauth.ErrAuthorizationExpired):
			appErrors.WriteError(writer, appErrors.NewHTTP(http.StatusGone, "DEVICE_AUTHORIZATION_EXPIRED", "Device authorization expired"))
		case errors.Is(err, deviceauth.ErrAlreadyHandled):
			appErrors.WriteError(writer, appErrors.ErrConflict.WithDetail("device authorization was already handled"))
		case errors.Is(err, deviceauth.ErrInvalidCode):
			appErrors.WriteError(writer, appErrors.ErrBadRequest.WithDetail(err.Error()))
		default:
			appErrors.WriteError(writer, appErrors.ErrInternalServer.WithDetail("failed to update device authorization"))
		}
		return
	}

	writeJSON(writer, http.StatusOK, map[string]string{"status": status})
}

func (handler *Handler) writeExchangeError(writer http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, deviceauth.ErrAuthorizationPending):
		writeJSON(writer, http.StatusAccepted, map[string]string{"error": "authorization_pending"})
	case errors.Is(err, deviceauth.ErrAuthorizationDenied):
		writeJSON(writer, http.StatusForbidden, map[string]string{"error": "access_denied"})
	case errors.Is(err, deviceauth.ErrAuthorizationExpired), errors.Is(err, deviceauth.ErrInvalidCode):
		writeJSON(writer, http.StatusGone, map[string]string{"error": "expired_token"})
	default:
		appErrors.WriteError(writer, appErrors.ErrInternalServer.WithDetail("failed to exchange device code"))
	}
}

func decodeJSON(writer http.ResponseWriter, request *http.Request, target any) error {
	defer request.Body.Close()

	decoder := json.NewDecoder(http.MaxBytesReader(writer, request.Body, 1<<20))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(target); err != nil {
		return err
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return appErrors.ErrBadRequest.WithDetail("request body must contain exactly one JSON object")
	}

	return nil
}

func bearerToken(request *http.Request) (string, bool) {
	authorization := strings.TrimSpace(request.Header.Get("Authorization"))
	parts := strings.Fields(authorization)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] == "" {
		return "", false
	}

	return parts[1], true
}

func writeJSON(writer http.ResponseWriter, status int, value any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(value)
}

func clientIP(request *http.Request) string {
	host, _, err := net.SplitHostPort(request.RemoteAddr)
	if err == nil {
		return host
	}

	return request.RemoteAddr
}
