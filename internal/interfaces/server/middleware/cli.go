package middleware

import (
	"context"
	"net/http"
	"strings"

	"fuse/internal/domain/compute"
	computeService "fuse/internal/services/compute"
	"fuse/pkg/errors"
)

type cliContextKey string

const cliCredentialKey cliContextKey = "cli_credential"

type CLIMiddleware struct {
	service *computeService.Service
}

func NewCLIMiddleware(service *computeService.Service) *CLIMiddleware {
	return &CLIMiddleware{service: service}
}

func (m *CLIMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := strings.TrimSpace(r.Header.Get("Authorization"))
		if !strings.HasPrefix(strings.ToLower(authorization), "bearer ") {
			errors.WriteError(w, errors.ErrUnauthorized.WithDetail("CLI bearer token is required"))
			return
		}
		token := strings.TrimSpace(authorization[len("Bearer "):])
		credential, err := m.service.AuthenticateCLIToken(r.Context(), token)
		if err != nil {
			errors.WriteError(w, errors.ErrUnauthorized.WithDetail("CLI credential is invalid, expired, or revoked"))
			return
		}
		ctx := context.WithValue(r.Context(), cliCredentialKey, credential)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetCLICredential(ctx context.Context) *compute.CLICredential {
	credential, _ := ctx.Value(cliCredentialKey).(*compute.CLICredential)
	return credential
}
