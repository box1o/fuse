package compute

import (
	"encoding/json"
	"net/http"

	domainCompute "fuse/internal/domain/compute"
	"fuse/internal/interfaces/server/middleware"
	computeService "fuse/internal/services/compute"
	"fuse/pkg/errors"
	"fuse/pkg/log"

	serverWebSocket "fuse/internal/interfaces/server/websocket"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct {
	computeSvc      *computeService.Service
	webSocketServer *serverWebSocket.Server
}

type CreateNodeRequest struct {
	WorkspaceID uuid.UUID `json:"workspace_id" example:"9f99eadd-7206-45ad-b284-8031019c39bc"`
	Name        string    `json:"name" example:"development-node"`

	Capabilities domainCompute.Capabilities `json:"capabilities"`
}

type RenameNodeRequest struct {
	Name string `json:"name" example:"production-node"`
}

type UpdateNodeStatusRequest struct {
	Status domainCompute.NodeStatus `json:"status" enums:"pending,online,offline,disabled" example:"online"`
}

type UpdateNodeCapabilitiesRequest struct {
	Capabilities domainCompute.Capabilities `json:"capabilities"`
}

func NewHandler(service *computeService.Service, webSocketServer *serverWebSocket.Server) *Handler {
	return &Handler{
		computeSvc:      service,
		webSocketServer: webSocketServer,
	}
}

func (handler *Handler) RegisterRoutes(router chi.Router, authMiddleware *middleware.AuthMiddleware) {
	router.Route("/compute/node", func(router chi.Router) {
		router.Use(authMiddleware.RequireAuth)

		router.Post("/", handler.CreateNode)
		router.Get("/", handler.ListOwnerNodes)
		router.Get("/stream", handler.StreamWorkspaceNodes)
		router.Get(
			"/workspace/{workspaceID}",
			handler.ListWorkspaceNodes,
		)

		router.Get("/{nodeID}", handler.GetNode)
		router.Delete("/{nodeID}", handler.DeleteNode)

		router.Patch("/{nodeID}/name", handler.RenameNode)
		router.Patch("/{nodeID}/status", handler.UpdateNodeStatus)
		router.Patch(
			"/{nodeID}/capabilities",
			handler.UpdateNodeCapabilities,
		)
	})
}

// CreateNode godoc
//
//	@Summary		Create compute node
//	@Description	Creates a compute node owned by the authenticated user.
//	@Tags			compute
//	@Accept			json
//	@Produce		json
//	@Param			request	body		CreateNodeRequest	true	"Compute node details"
//	@Success		201		{object}	domainCompute.Node
//	@Failure		400		{object}	errors.HTTPError
//	@Failure		401		{object}	errors.HTTPError
//	@Failure		500		{object}	errors.HTTPError
//	@Router			/compute/node [post]
func (handler *Handler) CreateNode(writer http.ResponseWriter, request *http.Request) {
	ownerID, ok := authenticatedUserID(writer, request)
	if !ok {
		return
	}

	defer request.Body.Close()

	var payload CreateNodeRequest
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		log.Warn(
			"failed to decode create compute node request: %v",
			err,
		)

		errors.WriteError(
			writer,
			errors.ErrBadRequest.WithDetail(
				"invalid request payload",
			),
		)
		return
	}

	node, err := handler.computeSvc.CreateNode(
		request.Context(),
		computeService.CreateNodeInput{
			OwnerID:      ownerID,
			WorkspaceID:  payload.WorkspaceID,
			Name:         payload.Name,
			Capabilities: payload.Capabilities,
		},
	)
	if err != nil {
		log.Error("failed to create compute node: %v", err)
		errors.WriteError(writer, errors.ToHTTP(err))
		return
	}

	writeJSON(writer, http.StatusCreated, node)
}

// GetNode godoc
//
//	@Summary		Get compute node
//	@Description	Returns a compute node owned by the authenticated user.
//	@Tags			compute
//	@Produce		json
//	@Param			nodeID	path		string	true	"Compute node ID"	Format(uuid)
//	@Success		200		{object}	domainCompute.Node
//	@Failure		400		{object}	errors.HTTPError
//	@Failure		401		{object}	errors.HTTPError
//	@Failure		404		{object}	errors.HTTPError
//	@Failure		500		{object}	errors.HTTPError
//	@Router			/compute/node/{nodeID} [get]
func (handler *Handler) GetNode(writer http.ResponseWriter, request *http.Request) {
	ownerID, ok := authenticatedUserID(writer, request)
	if !ok {
		return
	}

	nodeID, ok := nodeIDFromRequest(writer, request)
	if !ok {
		return
	}

	node, err := handler.computeSvc.GetNode(
		request.Context(),
		nodeID,
	)
	if err != nil {
		log.Error("failed to get compute node: %v", err)
		errors.WriteError(writer, errors.ToHTTP(err))
		return
	}

	if node.OwnerID != ownerID {
		writeNodeNotFound(writer)
		return
	}

	writeJSON(writer, http.StatusOK, node)
}

// ListOwnerNodes godoc
//
//	@Summary		List owned compute nodes
//	@Description	Returns all compute nodes owned by the authenticated user.
//	@Tags			compute
//	@Produce		json
//	@Success		200	{array}		domainCompute.Node
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/compute/node [get]
func (handler *Handler) ListOwnerNodes(writer http.ResponseWriter, request *http.Request) {
	ownerID, ok := authenticatedUserID(writer, request)
	if !ok {
		return
	}

	nodes, err := handler.computeSvc.ListNodesByOwnerID(
		request.Context(),
		ownerID,
	)
	if err != nil {
		log.Error("failed to list compute nodes: %v", err)
		errors.WriteError(writer, errors.ToHTTP(err))
		return
	}

	writeJSON(writer, http.StatusOK, nodes)
}

// ListWorkspaceNodes godoc
//
//	@Summary		List workspace compute nodes
//	@Description	Returns compute nodes in a workspace that are owned by the authenticated user.
//	@Tags			compute
//	@Produce		json
//	@Param			workspaceID	path		string	true	"Workspace ID"	Format(uuid)
//	@Success		200			{array}		domainCompute.Node
//	@Failure		400			{object}	errors.HTTPError
//	@Failure		401			{object}	errors.HTTPError
//	@Failure		500			{object}	errors.HTTPError
//	@Router			/compute/node/workspace/{workspaceID} [get]
func (handler *Handler) ListWorkspaceNodes(writer http.ResponseWriter, request *http.Request) {
	ownerID, ok := authenticatedUserID(writer, request)
	if !ok {
		return
	}

	workspaceID, err := uuid.Parse(chi.URLParam(request, "workspaceID"))
	if err != nil || workspaceID == uuid.Nil {
		errors.WriteError(
			writer,
			errors.ErrBadRequest.WithDetail(
				"invalid workspace ID",
			),
		)
		return
	}

	nodes, err := handler.computeSvc.ListNodesByWorkspaceID(request.Context(), workspaceID)
	if err != nil {
		log.Error(
			"failed to list workspace compute nodes: %v",
			err,
		)

		errors.WriteError(writer, errors.ToHTTP(err))
		return
	}

	ownerNodes := make([]*domainCompute.Node, 0, len(nodes))

	for _, node := range nodes {
		if node != nil && node.OwnerID == ownerID {
			ownerNodes = append(ownerNodes, node)
		}
	}

	writeJSON(writer, http.StatusOK, ownerNodes)
}

// RenameNode godoc
//
//	@Summary		Rename compute node
//	@Description	Changes the name of a compute node owned by the authenticated user.
//	@Tags			compute
//	@Accept			json
//	@Produce		json
//	@Param			nodeID	path		string				true	"Compute node ID"	Format(uuid)
//	@Param			request	body		RenameNodeRequest	true	"New compute node name"
//	@Success		200		{object}	domainCompute.Node
//	@Failure		400		{object}	errors.HTTPError
//	@Failure		401		{object}	errors.HTTPError
//	@Failure		404		{object}	errors.HTTPError
//	@Failure		500		{object}	errors.HTTPError
//	@Router			/compute/node/{nodeID}/name [patch]
func (handler *Handler) RenameNode(writer http.ResponseWriter, request *http.Request) {
	ownerID, ok := authenticatedUserID(writer, request)
	if !ok {
		return
	}

	nodeID, ok := nodeIDFromRequest(writer, request)
	if !ok {
		return
	}

	_, ok = handler.getOwnedNode(writer, request, nodeID, ownerID)
	if !ok {
		return
	}

	defer request.Body.Close()

	var payload RenameNodeRequest
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		errors.WriteError(
			writer,
			errors.ErrBadRequest.WithDetail(
				"invalid request payload",
			),
		)
		return
	}

	node, err := handler.computeSvc.RenameNode(request.Context(), nodeID, payload.Name)
	if err != nil {
		log.Error("failed to rename compute node: %v", err)
		errors.WriteError(writer, errors.ToHTTP(err))
		return
	}

	writeJSON(writer, http.StatusOK, node)
}

// UpdateNodeStatus godoc
//
//	@Summary		Update compute node status
//	@Description	Changes the status of a compute node owned by the authenticated user.
//	@Tags			compute
//	@Accept			json
//	@Produce		json
//	@Param			nodeID	path		string						true	"Compute node ID"	Format(uuid)
//	@Param			request	body		UpdateNodeStatusRequest	true	"New compute node status"
//	@Success		200		{object}	domainCompute.Node
//	@Failure		400		{object}	errors.HTTPError
//	@Failure		401		{object}	errors.HTTPError
//	@Failure		404		{object}	errors.HTTPError
//	@Failure		500		{object}	errors.HTTPError
//	@Router			/compute/node/{nodeID}/status [patch]
func (handler *Handler) UpdateNodeStatus(writer http.ResponseWriter, request *http.Request) {
	ownerID, ok := authenticatedUserID(writer, request)
	if !ok {
		return
	}

	nodeID, ok := nodeIDFromRequest(writer, request)
	if !ok {
		return
	}

	_, ok = handler.getOwnedNode(writer, request, nodeID, ownerID)
	if !ok {
		return
	}

	defer request.Body.Close()

	var payload UpdateNodeStatusRequest
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		errors.WriteError(
			writer,
			errors.ErrBadRequest.WithDetail(
				"invalid request payload",
			),
		)
		return
	}

	node, err := handler.computeSvc.UpdateNodeStatus(request.Context(), nodeID, payload.Status)
	if err != nil {
		log.Error(
			"failed to update compute node status: %v",
			err,
		)

		errors.WriteError(writer, errors.ToHTTP(err))
		return
	}

	writeJSON(writer, http.StatusOK, node)
}

// UpdateNodeCapabilities godoc
//
//	@Summary		Update compute node capabilities
//	@Description	Changes the hardware capabilities of a compute node owned by the authenticated user.
//	@Tags			compute
//	@Accept			json
//	@Produce		json
//	@Param			nodeID	path		string								true	"Compute node ID"	Format(uuid)
//	@Param			request	body		UpdateNodeCapabilitiesRequest	true	"New compute node capabilities"
//	@Success		200		{object}	domainCompute.Node
//	@Failure		400		{object}	errors.HTTPError
//	@Failure		401		{object}	errors.HTTPError
//	@Failure		404		{object}	errors.HTTPError
//	@Failure		500		{object}	errors.HTTPError
//	@Router			/compute/node/{nodeID}/capabilities [patch]
func (handler *Handler) UpdateNodeCapabilities(writer http.ResponseWriter, request *http.Request) {
	ownerID, ok := authenticatedUserID(writer, request)
	if !ok {
		return
	}

	nodeID, ok := nodeIDFromRequest(writer, request)
	if !ok {
		return
	}

	_, ok = handler.getOwnedNode(writer, request, nodeID, ownerID)
	if !ok {
		return
	}

	defer request.Body.Close()

	var payload UpdateNodeCapabilitiesRequest
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		errors.WriteError(
			writer,
			errors.ErrBadRequest.WithDetail(
				"invalid request payload",
			),
		)
		return
	}

	node, err := handler.computeSvc.UpdateNodeCapabilities(request.Context(), nodeID, payload.Capabilities)
	if err != nil {
		log.Error(
			"failed to update compute node capabilities: %v",
			err,
		)

		errors.WriteError(writer, errors.ToHTTP(err))
		return
	}

	writeJSON(writer, http.StatusOK, node)
}

// DeleteNode godoc
//
//	@Summary		Delete compute node
//	@Description	Permanently deletes a compute node owned by the authenticated user.
//	@Tags			compute
//	@Param			nodeID	path	string	true	"Compute node ID"	Format(uuid)
//	@Success		204
//	@Failure		400	{object}	errors.HTTPError
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		404	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/compute/node/{nodeID} [delete]
func (handler *Handler) DeleteNode(writer http.ResponseWriter, request *http.Request) {
	ownerID, ok := authenticatedUserID(writer, request)
	if !ok {
		return
	}

	nodeID, ok := nodeIDFromRequest(writer, request)
	if !ok {
		return
	}

	_, ok = handler.getOwnedNode(
		writer,
		request,
		nodeID,
		ownerID,
	)
	if !ok {
		return
	}

	if err := handler.computeSvc.DeleteNode(request.Context(), nodeID); err != nil {
		log.Error("failed to delete compute node: %v", err)
		errors.WriteError(writer, errors.ToHTTP(err))
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}

func (handler *Handler) getOwnedNode(writer http.ResponseWriter, request *http.Request, nodeID uuid.UUID, ownerID uuid.UUID) (*domainCompute.Node, bool) {
	node, err := handler.computeSvc.GetNode(request.Context(), nodeID)
	if err != nil {
		log.Error("failed to get compute node: %v", err)
		errors.WriteError(writer, errors.ToHTTP(err))
		return nil, false
	}

	if node == nil || node.OwnerID != ownerID {
		writeNodeNotFound(writer)
		return nil, false
	}

	return node, true
}

func authenticatedUserID(writer http.ResponseWriter, request *http.Request) (uuid.UUID, bool) {
	userID := middleware.GetUserIDFromContext(
		request.Context(),
	)

	if userID == uuid.Nil {
		log.Warn("user ID not found in context")

		errors.WriteError(
			writer,
			errors.ErrUnauthorized.WithDetail(
				"user not authenticated",
			),
		)

		return uuid.Nil, false
	}

	return userID, true
}

func nodeIDFromRequest(writer http.ResponseWriter, request *http.Request) (uuid.UUID, bool) {

	nodeID, err := uuid.Parse(
		chi.URLParam(request, "nodeID"),
	)
	if err != nil || nodeID == uuid.Nil {
		errors.WriteError(
			writer,
			errors.ErrBadRequest.WithDetail(
				"invalid compute node ID",
			),
		)

		return uuid.Nil, false
	}

	return nodeID, true
}

func writeNodeNotFound(writer http.ResponseWriter) {
	errors.WriteError(
		writer,
		errors.ToHTTP(domainCompute.ErrNodeNotFound),
	)
}

func writeJSON(
	writer http.ResponseWriter,
	statusCode int,
	value any,
) {
	writer.Header().Set(
		"Content-Type",
		"application/json",
	)

	writer.WriteHeader(statusCode)

	if err := json.NewEncoder(writer).Encode(value); err != nil {
		log.Error("failed to encode response: %v", err)
	}
}
