package stripe

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"

	domainPayment "fuse/internal/domain/payment"
	paymentService "fuse/internal/services/payment"

	stripeSDK "github.com/stripe/stripe-go/v86"
)

var _ paymentService.Provider = (*Client)(nil)

func (client *Client) CreateCheckoutSession(ctx context.Context, input paymentService.CreateCheckoutSessionInput) (*paymentService.CheckoutSession, error) {
	if err := validateCheckoutSessionInput(input); err != nil {
		return nil, err
	}

	metadata := map[string]string{
		"payment_id":     input.PaymentID.String(),
		"owner_id":       input.OwnerID.String(),
		"credit_pack_id": input.CreditPackID.String(),
	}

	params := &stripeSDK.CheckoutSessionCreateParams{
		Mode: stripeSDK.String(
			string(stripeSDK.CheckoutSessionModePayment),
		),

		SuccessURL: stripeSDK.String(
			strings.TrimSpace(input.SuccessURL),
		),

		CancelURL: stripeSDK.String(
			strings.TrimSpace(input.CancelURL),
		),

		ClientReferenceID: stripeSDK.String(
			input.PaymentID.String(),
		),

		LineItems: []*stripeSDK.CheckoutSessionCreateLineItemParams{
			{
				Price: stripeSDK.String(
					strings.TrimSpace(input.PriceReference),
				),
				Quantity: stripeSDK.Int64(1),
			},
		},

		Metadata: metadata,

		PaymentIntentData: &stripeSDK.CheckoutSessionCreatePaymentIntentDataParams{
			Metadata: metadata,
		},
	}

	session, err := client.createCheckoutSession(ctx, params)
	if err != nil {
		return nil, errors.Join(
			paymentService.ErrCheckoutSessionCreationFailed,
			err,
		)
	}

	if session == nil {
		return nil, fmt.Errorf(
			"%w: Stripe returned a nil session",
			paymentService.ErrCheckoutSessionCreationFailed,
		)
	}

	if strings.TrimSpace(session.ID) == "" {
		return nil, fmt.Errorf(
			"%w: Stripe returned an empty session ID",
			paymentService.ErrCheckoutSessionCreationFailed,
		)
	}

	if strings.TrimSpace(session.URL) == "" {
		return nil, fmt.Errorf(
			"%w: Stripe returned an empty checkout URL",
			paymentService.ErrCheckoutSessionCreationFailed,
		)
	}

	providerPaymentID := ""
	if session.PaymentIntent != nil {
		providerPaymentID = session.PaymentIntent.ID
	}

	return &paymentService.CheckoutSession{
		Provider:          domainPayment.ProviderStripe,
		SessionID:         session.ID,
		URL:               session.URL,
		ProviderPaymentID: providerPaymentID,
		Amount:            session.AmountTotal,
		Currency: strings.ToUpper(
			string(session.Currency),
		),
	}, nil
}

func validateCheckoutSessionInput(
	input paymentService.CreateCheckoutSessionInput,
) error {
	if input.PaymentID == uuid.Nil {
		return domainPayment.ErrPaymentNotFound
	}

	if input.OwnerID == uuid.Nil {
		return domainPayment.ErrOwnerIDRequired
	}

	if input.CreditPackID == uuid.Nil {
		return domainPayment.ErrCreditPackIDRequired
	}

	if strings.TrimSpace(input.PriceReference) == "" {
		return paymentService.ErrPriceReferenceRequired
	}

	if strings.TrimSpace(input.SuccessURL) == "" {
		return paymentService.ErrSuccessURLRequired
	}

	if strings.TrimSpace(input.CancelURL) == "" {
		return paymentService.ErrCancelURLRequired
	}

	return nil
}
