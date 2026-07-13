package workspace

import (
	"encoding/json"
	"fuse/internal/interfaces/server/middleware"
	"fuse/internal/services/workspace"
	"fuse/pkg/config"
	"fuse/pkg/errors"
	"fuse/pkg/log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct {
	workspaceSvc *workspace.Service
	cfg          *config.Config
}

func NewHandler(wsService *workspace.Service, cfg *config.Config) *Handler {
	return &Handler{
		workspaceSvc: wsService,
		cfg:          cfg,
	}
}

func (h *Handler) RegisterRoutes(r chi.Router, authMiddleware *middleware.AuthMiddleware) {
	r.Route("/workspaces", func(r chi.Router) {
		r.Use(authMiddleware.RequireAuth)
		r.Post("/", h.CreateWorkspace)
		r.Get("/", h.GetOwnerWorkspaces)
		r.Delete("/{workspaceID}", h.DeleteWorkspace)

	})
}

type CreateWorkspaceRequest struct {
	Name string `json:"name"`
}

// @Summary		Create a workspace
// @Description	Creates a workspace owned by the authenticated user.
// @Tags			workspaces
// @Accept			json
// @Produce		json
// @Param			request	body	CreateWorkspaceRequest	true	"Workspace details"
// @Success		201	{object}	map[string]interface{}
// @Failure		400	{object}	errors.HTTPError
// @Failure		401	{object}	errors.HTTPError
// @Failure		409	{object}	errors.HTTPError
// @Failure		500	{object}	errors.HTTPError
// @Router			/workspaces [post]
func (h *Handler) CreateWorkspace(w http.ResponseWriter, r *http.Request) {
	var req CreateWorkspaceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Warn("failed to decode create workspace request: %v", err)
		errors.WriteError(w, errors.ErrBadRequest.WithDetail("invalid request payload"))
		return
	}
	defer r.Body.Close()

	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		log.Warn("user ID not found in context")
		errors.WriteError(w, errors.ErrUnauthorized.WithDetail("user not authenticated"))
		return
	}

	ws, err := h.workspaceSvc.CreateWorkspace(r.Context(), req.Name, userID)
	if err != nil {
		// aici folosim direct metoda Is implementată de tine
		if e, ok := err.(*errors.Error); ok && e.Is(errors.ErrNameExists) {
			errors.WriteError(w, errors.ErrNameExists.WithDetail("workspace name already exists"))
			return
		}

		log.Error("failed to create workspace: %v", err)
		errors.WriteError(w, errors.ToHTTP(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(ws)
}

// @Summary		List owned workspaces
// @Description	Returns every workspace owned by the authenticated user.
// @Tags			workspaces
// @Produce		json
// @Success		200	{array}	map[string]interface{}
// @Failure		401	{object}	errors.HTTPError
// @Failure		500	{object}	errors.HTTPError
// @Router			/workspaces [get]
func (h *Handler) GetOwnerWorkspaces(w http.ResponseWriter, r *http.Request) {

	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		log.Warn("user ID not found in context")
		errors.WriteError(w, errors.ErrUnauthorized.WithDetail("user not authenticated"))
		return
	}

	workspaces, err := h.workspaceSvc.GetUserWorkspaces(r.Context(), userID)
	if err != nil {
		log.Error("failed to get workspaces: %v", err)
		errors.WriteError(w, errors.ToHTTP(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(workspaces)

}

// @Summary		Delete a workspace
// @Description	Deletes the workspace identified by its ID.
// @Tags			workspaces
// @Param			workspaceID	path	string	true	"Workspace ID"
// @Success		204
// @Failure		400	{object}	errors.HTTPError
// @Failure		404	{object}	errors.HTTPError
// @Failure		500	{object}	errors.HTTPError
// @Router			/workspaces/{workspaceID} [delete]
func (h *Handler) DeleteWorkspace(w http.ResponseWriter, r *http.Request) {
	workspaceIDStr := chi.URLParam(r, "workspaceID")
	workspaceID, err := uuid.Parse(workspaceIDStr)
	if err != nil || workspaceID == uuid.Nil {
		log.Warn("invalid workspace ID: %v", err)
		errors.WriteError(w, errors.ErrBadRequest.WithDetail("invalid workspace ID"))
		return
	}

	err = h.workspaceSvc.DeleteWorkspace(r.Context(), workspaceID)
	if err != nil {
		log.Error("failed to delete workspace: %v", err)
		errors.WriteError(w, errors.ToHTTP(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
