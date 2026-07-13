package compute

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	domain "fuse/internal/domain/compute"
	"fuse/internal/interfaces/server/middleware"
	computeService "fuse/internal/services/compute"
	"fuse/pkg/computeapi"
	"fuse/pkg/errors"
	"fuse/pkg/log"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct {
	service *computeService.Service
	cliAuth *middleware.CLIMiddleware
}

func NewHandler(service *computeService.Service, cliAuth *middleware.CLIMiddleware) *Handler {
	return &Handler{service: service, cliAuth: cliAuth}
}

func (h *Handler) RegisterRoutes(r chi.Router, browserAuth *middleware.AuthMiddleware) {
	r.Route("/compute/nodes", func(r chi.Router) {
		r.With(h.cliAuth.RequireAuth).Post("/", h.RegisterNode)
		r.With(h.requireBrowserOrCLI(browserAuth)).Get("/", h.ListNodes)
		r.With(h.requireBrowserOrCLI(browserAuth)).Get("/{nodeID}", h.GetNode)
		r.With(h.requireBrowserOrCLI(browserAuth)).Patch("/{nodeID}", h.UpdateNode)
		r.With(h.requireBrowserOrCLI(browserAuth)).Delete("/{nodeID}", h.DeleteNode)
	})
}

func (h *Handler) requireBrowserOrCLI(browserAuth *middleware.AuthMiddleware) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(strings.ToLower(strings.TrimSpace(r.Header.Get("Authorization"))), "bearer ") {
				h.cliAuth.RequireAuth(next).ServeHTTP(w, r)
				return
			}
			browserAuth.RequireAuth(next).ServeHTTP(w, r)
		})
	}
}

func (h *Handler) RegisterNode(w http.ResponseWriter, r *http.Request) {
	ownerID := ownerIDFromCLI(r)
	var request computeapi.RegisterNodeRequest
	if err := decodeJSON(w, r, &request); err != nil {
		errors.WriteError(w, errors.ErrBadRequest.WithDetail(err.Error()))
		return
	}
	installationID, err := uuid.Parse(request.InstallationID)
	if err != nil {
		errors.WriteError(w, errors.ErrValidation.WithDetail("installation_id must be a valid UUID"))
		return
	}

	node, created, err := h.service.RegisterNode(r.Context(), ownerID, domain.RegisterNodeInput{
		InstallationID: installationID, Name: request.Name, Hostname: request.Hostname,
		AgentVersion: request.AgentVersion, Capabilities: capabilitiesToDomain(request.Capabilities),
	})
	if err != nil {
		writeDomainError(w, err)
		return
	}
	status := http.StatusOK
	if created {
		status = http.StatusCreated
	}
	writeJSON(w, status, nodeToAPI(node))
}

func (h *Handler) ListNodes(w http.ResponseWriter, r *http.Request) {
	nodes, err := h.service.ListNodes(r.Context(), ownerIDFromRequest(r))
	if err != nil {
		writeDomainError(w, err)
		return
	}
	response := make([]computeapi.Node, 0, len(nodes))
	for _, node := range nodes {
		response = append(response, nodeToAPI(node))
	}
	writeJSON(w, http.StatusOK, response)
}

func (h *Handler) GetNode(w http.ResponseWriter, r *http.Request) {
	nodeID, ok := parseNodeID(w, r)
	if !ok {
		return
	}
	node, err := h.service.GetNode(r.Context(), ownerIDFromRequest(r), nodeID)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, nodeToAPI(node))
}

func (h *Handler) UpdateNode(w http.ResponseWriter, r *http.Request) {
	nodeID, ok := parseNodeID(w, r)
	if !ok {
		return
	}
	var request computeapi.UpdateNodeRequest
	if err := decodeJSON(w, r, &request); err != nil {
		errors.WriteError(w, errors.ErrBadRequest.WithDetail(err.Error()))
		return
	}
	if request.Name == nil && request.Disabled == nil {
		errors.WriteError(w, errors.ErrBadRequest.WithDetail("name or disabled must be provided"))
		return
	}
	node, err := h.service.UpdateNode(r.Context(), ownerIDFromRequest(r), nodeID, computeService.UpdateNodeInput{
		Name: request.Name, Disabled: request.Disabled,
	})
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, nodeToAPI(node))
}

func (h *Handler) DeleteNode(w http.ResponseWriter, r *http.Request) {
	nodeID, ok := parseNodeID(w, r)
	if !ok {
		return
	}
	if err := h.service.DeleteNode(r.Context(), ownerIDFromRequest(r), nodeID); err != nil {
		writeDomainError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func ownerIDFromCLI(r *http.Request) uuid.UUID {
	credential := middleware.GetCLICredential(r.Context())
	if credential == nil {
		return uuid.Nil
	}
	return credential.OwnerID
}

func ownerIDFromRequest(r *http.Request) uuid.UUID {
	if credential := middleware.GetCLICredential(r.Context()); credential != nil {
		return credential.OwnerID
	}
	return middleware.GetUserIDFromContext(r.Context())
}

func parseNodeID(w http.ResponseWriter, r *http.Request) (uuid.UUID, bool) {
	nodeID, err := uuid.Parse(chi.URLParam(r, "nodeID"))
	if err != nil {
		errors.WriteError(w, errors.ErrValidation.WithDetail("nodeID must be a valid UUID"))
		return uuid.Nil, false
	}
	return nodeID, true
}

func decodeJSON(w http.ResponseWriter, r *http.Request, target any) error {
	defer r.Body.Close()
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return err
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return errors.ErrBadRequest.WithDetail("request body must contain exactly one JSON object")
	}
	return nil
}

func writeDomainError(w http.ResponseWriter, err error) {
	log.Warn("compute request failed: %v", err)
	errors.WriteError(w, errors.ToHTTP(err))
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Error("failed to encode compute response: %v", err)
	}
}
