package httpapi

import (
	"context"
	"net/http"

	"fuse/pkg/deviceauthapi"
)

type DeviceAuthGateway struct{ client *Client }

func NewDeviceAuthGateway(client *Client) *DeviceAuthGateway {
	return &DeviceAuthGateway{client: client}
}

func (g *DeviceAuthGateway) CreateCode(ctx context.Context, clientName string) (deviceauthapi.CodeResponse, error) {
	var response deviceauthapi.CodeResponse
	_, err := g.client.Do(ctx, http.MethodPost, "/auth/device/code", "", deviceauthapi.CreateCodeRequest{ClientName: clientName}, &response)
	return response, err
}

func (g *DeviceAuthGateway) ExchangeCode(ctx context.Context, deviceCode string) (deviceauthapi.TokenResponse, bool, error) {
	var response deviceauthapi.TokenResponse
	status, err := g.client.Do(ctx, http.MethodPost, "/auth/device/token", "", deviceauthapi.TokenRequest{DeviceCode: deviceCode}, &response)
	if status == http.StatusAccepted {
		return deviceauthapi.TokenResponse{}, true, nil
	}
	return response, false, err
}

func (g *DeviceAuthGateway) Status(ctx context.Context, token string) (deviceauthapi.StatusResponse, error) {
	var response deviceauthapi.StatusResponse
	_, err := g.client.Do(ctx, http.MethodGet, "/auth/cli/status", token, nil, &response)
	return response, err
}

func (g *DeviceAuthGateway) Logout(ctx context.Context, token string) error {
	_, err := g.client.Do(ctx, http.MethodPost, "/auth/cli/logout", token, nil, nil)
	return err
}
