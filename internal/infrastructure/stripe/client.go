package stripe

import (
	"context"
	"strings"

	stripeSDK "github.com/stripe/stripe-go/v83"
)

type checkoutSessionCreator func(ctx context.Context, params *stripeSDK.CheckoutSessionCreateParams) (*stripeSDK.CheckoutSession, error)

type Client struct{ createCheckoutSession checkoutSessionCreator }

func NewClient(secretKey string) (*Client, error) {
	secretKey = strings.TrimSpace(secretKey)
	if secretKey == "" {
		return nil, ErrSecretKeyRequired
	}

	stripeClient := stripeSDK.NewClient(secretKey)

	return &Client{
		createCheckoutSession: stripeClient.V1CheckoutSessions.Create,
	}, nil
}

func newClientForTest(
	createCheckoutSession checkoutSessionCreator,
) *Client {
	return &Client{
		createCheckoutSession: createCheckoutSession,
	}
}
