package compute

import (
	"context"
	"errors"
	"net/http"
	"time"

	domainCompute "fuse/internal/domain/compute"
	"fuse/internal/interfaces/server/middleware"
	serverWebSocket "fuse/internal/interfaces/server/websocket"
	"fuse/pkg/log"

	coderwebsocket "github.com/coder/websocket"
	"github.com/google/uuid"
)

const nodeStreamInterval = 2 * time.Second

func (handler *Handler) StreamWorkspaceNodes(writer http.ResponseWriter, request *http.Request) {
	ownerID := middleware.GetUserIDFromContext(request.Context())
	if ownerID == uuid.Nil {
		http.Error(
			writer,
			"authentication required",
			http.StatusUnauthorized,
		)
		return
	}

	workspaceID, err := uuid.Parse(request.URL.Query().Get("workspace_id"))
	if err != nil || workspaceID == uuid.Nil {
		http.Error(
			writer,
			"invalid workspace_id",
			http.StatusBadRequest,
		)
		return
	}

	connection, err := handler.webSocketServer.Accept(writer, request)
	if err != nil {
		log.Warn(
			"failed to accept compute websocket: %v",
			err,
		)
		return
	}

	defer func() {
		if err := connection.CloseNormal(); err != nil {
			log.Warn(
				"failed to close compute websocket: %v",
				err,
			)
		}
	}()

	err = serverWebSocket.Stream(
		request.Context(),
		connection,
		nodeStreamInterval,
		func(ctx context.Context) (any, error) {
			return handler.getOwnerWorkspaceNodes(
				ctx,
				workspaceID,
				ownerID,
			)
		},
	)
	if err != nil &&
		!errors.Is(err, context.Canceled) &&
		!errors.Is(err, context.DeadlineExceeded) &&
		coderwebsocket.CloseStatus(err) == -1 {
		log.Warn(
			"compute websocket stream stopped: %v",
			err,
		)
	}
}

func (handler *Handler) getOwnerWorkspaceNodes(ctx context.Context, workspaceID uuid.UUID, ownerID uuid.UUID) ([]*domainCompute.Node, error) {
	nodes, err := handler.computeSvc.ListNodesByWorkspaceID(ctx, workspaceID)
	if err != nil {
		return nil, err
	}

	ownerNodes := make([]*domainCompute.Node, 0, len(nodes))

	for _, node := range nodes {
		if node == nil || node.OwnerID != ownerID {
			continue
		}

		ownerNodes = append(ownerNodes, node)
	}

	return ownerNodes, nil
}
