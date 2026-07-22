package stripe

import (
	"context"
	stdErrors "errors"
	"testing"

	domainPayment "fuse/internal/domain/payment"
	paymentService "fuse/internal/services/payment"

	"github.com/google/uuid"
	stripeSDK "github.com/stripe/stripe-go/v86"
)

var errStripeAPI = stdErrors.New("forced Stripe API failure")

func TestClient_CreateCheckoutSession(t *testing.T) {
	t.Parallel()

	paymentID := uuid.New()
	ownerID := uuid.New()
	creditPackID := uuid.New()

	client := newClientForTest(
		func(
			_ context.Context,
			params *stripeSDK.CheckoutSessionCreateParams,
		) (*stripeSDK.CheckoutSession, error) {
			if params.Mode == nil ||
				*params.Mode != string(stripeSDK.CheckoutSessionModePayment) {
				t.Fatalf("expected payment mode")
			}

			if params.ClientReferenceID == nil ||
				*params.ClientReferenceID != paymentID.String() {
				t.Fatalf(
					"expected client reference ID %s",
					paymentID,
				)
			}

			if len(params.LineItems) != 1 {
				t.Fatalf(
					"expected one line item, got %d",
					len(params.LineItems),
				)
			}

			lineItem := params.LineItems[0]

			if lineItem.Price == nil ||
				*lineItem.Price != "price_test_123" {
				t.Fatalf("expected configured Stripe price")
			}

			if lineItem.Quantity == nil || *lineItem.Quantity != 1 {
				t.Fatalf("expected quantity 1")
			}

			if params.Metadata["payment_id"] != paymentID.String() {
				t.Errorf("expected payment ID metadata")
			}

			if params.Metadata["owner_id"] != ownerID.String() {
				t.Errorf("expected owner ID metadata")
			}

			if params.Metadata["credit_pack_id"] != creditPackID.String() {
				t.Errorf("expected credit pack ID metadata")
			}

			return &stripeSDK.CheckoutSession{
				ID:          "cs_test_123",
				URL:         "https://checkout.stripe.com/test",
				AmountTotal: 999,
				Currency:    stripeSDK.CurrencyUSD,
				PaymentIntent: &stripeSDK.PaymentIntent{
					ID: "pi_test_123",
				},
			}, nil
		},
	)

	result, err := client.CreateCheckoutSession(
		context.Background(),
		paymentService.CreateCheckoutSessionInput{
			PaymentID:      paymentID,
			OwnerID:        ownerID,
			CreditPackID:   creditPackID,
			PriceReference: " price_test_123 ",
			SuccessURL:     " https://example.com/payment/success ",
			CancelURL:      " https://example.com/payment/cancel ",
		},
	)
	if err != nil {
		t.Fatalf(
			"CreateCheckoutSession() returned unexpected error: %v",
			err,
		)
	}

	if result.Provider != domainPayment.ProviderStripe {
		t.Errorf(
			"expected provider %q, got %q",
			domainPayment.ProviderStripe,
			result.Provider,
		)
	}

	if result.SessionID != "cs_test_123" {
		t.Errorf(
			"expected session ID %q, got %q",
			"cs_test_123",
			result.SessionID,
		)
	}

	if result.CheckoutURL != "https://checkout.stripe.com/test" {
		t.Errorf(
			"unexpected checkout URL %q",
			result.CheckoutURL,
		)
	}

	if result.ProviderPaymentID != "pi_test_123" {
		t.Errorf(
			"expected payment intent ID %q, got %q",
			"pi_test_123",
			result.ProviderPaymentID,
		)
	}

	if result.Amount != 999 {
		t.Errorf("expected amount 999, got %d", result.Amount)
	}

	if result.Currency != "USD" {
		t.Errorf(
			"expected currency USD, got %q",
			result.Currency,
		)
	}
}

func TestClient_CreateCheckoutSession_ValidatesInput(t *testing.T) {
	t.Parallel()

	validInput := paymentService.CreateCheckoutSessionInput{
		PaymentID:      uuid.New(),
		OwnerID:        uuid.New(),
		CreditPackID:   uuid.New(),
		PriceReference: "price_test_123",
		SuccessURL:     "https://example.com/payment/success",
		CancelURL:      "https://example.com/payment/cancel",
	}

	tests := []struct {
		name          string
		changeInput   func(*paymentService.CreateCheckoutSessionInput)
		expectedError error
	}{
		{
			name: "missing payment ID",
			changeInput: func(
				input *paymentService.CreateCheckoutSessionInput,
			) {
				input.PaymentID = uuid.Nil
			},
			expectedError: domainPayment.ErrPaymentNotFound,
		},
		{
			name: "missing owner ID",
			changeInput: func(
				input *paymentService.CreateCheckoutSessionInput,
			) {
				input.OwnerID = uuid.Nil
			},
			expectedError: domainPayment.ErrOwnerIDRequired,
		},
		{
			name: "missing credit pack ID",
			changeInput: func(
				input *paymentService.CreateCheckoutSessionInput,
			) {
				input.CreditPackID = uuid.Nil
			},
			expectedError: domainPayment.ErrCreditPackIDRequired,
		},
		{
			name: "missing price reference",
			changeInput: func(
				input *paymentService.CreateCheckoutSessionInput,
			) {
				input.PriceReference = " "
			},
			expectedError: paymentService.ErrPriceReferenceRequired,
		},
		{
			name: "missing success URL",
			changeInput: func(
				input *paymentService.CreateCheckoutSessionInput,
			) {
				input.SuccessURL = " "
			},
			expectedError: paymentService.ErrSuccessURLRequired,
		},
		{
			name: "missing cancel URL",
			changeInput: func(
				input *paymentService.CreateCheckoutSessionInput,
			) {
				input.CancelURL = " "
			},
			expectedError: paymentService.ErrCancelURLRequired,
		},
	}

	client := newClientForTest(
		func(
			context.Context,
			*stripeSDK.CheckoutSessionCreateParams,
		) (*stripeSDK.CheckoutSession, error) {
			t.Fatal("Stripe API must not be called for invalid input")
			return nil, nil
		},
	)

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			input := validInput
			test.changeInput(&input)

			_, err := client.CreateCheckoutSession(
				context.Background(),
				input,
			)

			if !stdErrors.Is(err, test.expectedError) {
				t.Fatalf(
					"expected error %v, got %v",
					test.expectedError,
					err,
				)
			}
		})
	}
}

func TestClient_CreateCheckoutSession_ReturnsProviderError(
	t *testing.T,
) {
	t.Parallel()

	client := newClientForTest(
		func(
			context.Context,
			*stripeSDK.CheckoutSessionCreateParams,
		) (*stripeSDK.CheckoutSession, error) {
			return nil, errStripeAPI
		},
	)

	_, err := client.CreateCheckoutSession(
		context.Background(),
		paymentService.CreateCheckoutSessionInput{
			PaymentID:      uuid.New(),
			OwnerID:        uuid.New(),
			CreditPackID:   uuid.New(),
			PriceReference: "price_test_123",
			SuccessURL:     "https://example.com/payment/success",
			CancelURL:      "https://example.com/payment/cancel",
		},
	)

	if !stdErrors.Is(
		err,
		paymentService.ErrCheckoutSessionCreationFailed,
	) {
		t.Fatalf(
			"expected checkout creation error, got %v",
			err,
		)
	}
}

func TestClient_CreateCheckoutSession_RejectsInvalidResponse(
	t *testing.T,
) {
	t.Parallel()

	client := newClientForTest(
		func(
			context.Context,
			*stripeSDK.CheckoutSessionCreateParams,
		) (*stripeSDK.CheckoutSession, error) {
			return &stripeSDK.CheckoutSession{}, nil
		},
	)

	_, err := client.CreateCheckoutSession(
		context.Background(),
		paymentService.CreateCheckoutSessionInput{
			PaymentID:      uuid.New(),
			OwnerID:        uuid.New(),
			CreditPackID:   uuid.New(),
			PriceReference: "price_test_123",
			SuccessURL:     "https://example.com/payment/success",
			CancelURL:      "https://example.com/payment/cancel",
		},
	)

	if !stdErrors.Is(
		err,
		paymentService.ErrCheckoutSessionCreationFailed,
	) {
		t.Fatalf(
			"expected checkout creation error, got %v",
			err,
		)
	}
}

func TestNewClient_RequiresSecretKey(t *testing.T) {
	t.Parallel()

	_, err := NewClient(" ")

	if !stdErrors.Is(err, ErrSecretKeyRequired) {
		t.Fatalf(
			"expected secret key error, got %v",
			err,
		)
	}
}
