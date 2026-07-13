package httpapi

import (
	"context"
	"net/http"

	"fuse/pkg/computeapi"
)

type NodeGateway struct{ client *Client }

func NewNodeGateway(client *Client) *NodeGateway { return &NodeGateway{client: client} }

func (g *NodeGateway) Register(ctx context.Context, token string, request computeapi.RegisterNodeRequest) (computeapi.Node, bool, error) {
	var node computeapi.Node
	status, err := g.client.Do(ctx, http.MethodPost, "/compute/nodes/", token, request, &node)
	return node, status == http.StatusCreated, err
}

func (g *NodeGateway) List(ctx context.Context, token string) ([]computeapi.Node, error) {
	var nodes []computeapi.Node
	_, err := g.client.Do(ctx, http.MethodGet, "/compute/nodes/", token, nil, &nodes)
	return nodes, err
}

func (g *NodeGateway) Get(ctx context.Context, token, nodeID string) (computeapi.Node, error) {
	var node computeapi.Node
	_, err := g.client.Do(ctx, http.MethodGet, "/compute/nodes/"+nodeID, token, nil, &node)
	return node, err
}

func (g *NodeGateway) Update(ctx context.Context, token, nodeID string, request computeapi.UpdateNodeRequest) (computeapi.Node, error) {
	var node computeapi.Node
	_, err := g.client.Do(ctx, http.MethodPatch, "/compute/nodes/"+nodeID, token, request, &node)
	return node, err
}

func (g *NodeGateway) Delete(ctx context.Context, token, nodeID string) error {
	_, err := g.client.Do(ctx, http.MethodDelete, "/compute/nodes/"+nodeID, token, nil, nil)
	return err
}
