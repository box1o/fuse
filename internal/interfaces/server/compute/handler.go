package compute

import (
	"encoding/json"
	"fuse/internal/interfaces/server/middleware"

	//"fuse/internal/services/workspace"
	"fuse/pkg/config"
	"fuse/pkg/errors"
	"fuse/pkg/log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ComputeService struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Status string    `json:"status"`
}

type Handler struct {
	cfg *config.Config
}

func NewHandler(cfg *config.Config) *Handler {
	return &Handler{
		cfg: cfg,
	}
}

func (h *Handler) RegisterRoutes(r chi.Router, authMiddleware *middleware.AuthMiddleware) {
	r.Route("/compute", func(r chi.Router) {
		r.Use(authMiddleware.RequireAuth)
		// r.Post("/", h.CreateWorkspace)
		r.Get("/", h.GetComputeServices)
		// r.Delete("/{workspaceID}", h.DeleteWorkspace)

	})
}

// type CreateWorkspaceRequest struct {
// 	Name string `json:"name"`
// }

// func (h *Handler) CreateWorkspace(w http.ResponseWriter, r *http.Request) {
// 	var req CreateWorkspaceRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		log.Warn("failed to decode create workspace request: %v", err)
// 		errors.WriteError(w, errors.ErrBadRequest.WithDetail("invalid request payload"))
// 		return
// 	}
// 	defer r.Body.Close()

// 	userID := middleware.GetUserIDFromContext(r.Context())
// 	if userID == uuid.Nil {
// 		log.Warn("user ID not found in context")
// 		errors.WriteError(w, errors.ErrUnauthorized.WithDetail("user not authenticated"))
// 		return
// 	}

// 	ws, err := h.workspaceSvc.CreateWorkspace(r.Context(), req.Name, userID)
// 	if err != nil {
// 		// aici folosim direct metoda Is implementată de tine
// 		if e, ok := err.(*errors.Error); ok && e.Is(errors.ErrNameExists) {
// 			errors.WriteError(w, errors.ErrNameExists.WithDetail("workspace name already exists"))
// 			return
// 		}

// 		log.Error("failed to create workspace: %v", err)
// 		errors.WriteError(w, errors.ToHTTP(err))
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusCreated)
// 	_ = json.NewEncoder(w).Encode(ws)
// }

// @Summary		List compute services
// @Description	Returns the compute services available to the authenticated user.
// @Tags			compute
// @Produce		json
// @Success		200	{array}	ComputeService
// @Failure		401	{object}	errors.HTTPError
// @Router			/compute [get]
func (h *Handler) GetComputeServices(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		log.Warn("user ID not found in context")
		errors.WriteError(w, errors.ErrUnauthorized.WithDetail("user not authenticated"))
		return
	}

	services := []ComputeService{
		{
			ID:     uuid.New(),
			Name:   "Inference Engine",
			Status: "running",
		},
		{
			ID:     uuid.New(),
			Name:   "Training Cluster",
			Status: "stopped",
		},
		{
			ID:     uuid.New(),
			Name:   "Batch Processor",
			Status: "available",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(services)

}

// func (h *Handler) DeleteWorkspace(w http.ResponseWriter, r *http.Request) {
// 	workspaceIDStr := chi.URLParam(r, "workspaceID")
// 	workspaceID, err := uuid.Parse(workspaceIDStr)
// 	if err != nil || workspaceID == uuid.Nil {
// 		log.Warn("invalid workspace ID: %v", err)
// 		errors.WriteError(w, errors.ErrBadRequest.WithDetail("invalid workspace ID"))
// 		return
// 	}

// 	err = h.workspaceSvc.DeleteWorkspace(r.Context(), workspaceID)
// 	if err != nil {
// 		log.Error("failed to delete workspace: %v", err)
// 		errors.WriteError(w, errors.ToHTTP(err))
// 		return
// 	}

// 	w.WriteHeader(http.StatusNoContent)
// }
