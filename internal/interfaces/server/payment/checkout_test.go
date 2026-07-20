package payment

import (
	"bytes"
	"context"
	"encoding/json"
	stdErrors "errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"fuse/internal/interfaces/server/middleware"
	paymentService "fuse/internal/services/payment"

	"github.com/google/uuid"
)

func TestHandler_CreateCheckout_ReturnsCreatedCheckout(t *testing.T) {
	t.Parallel()

	ownerID := uuid.New()
	creditPackID := uuid.New()
	paymentID := uuid.New()

	service := &fakeCheckoutService{
		output: &paymentService.CreateCheckoutOutput{
			PaymentID:   paymentID,
			SessionID:   "cs_test_123",
			CheckoutURL: "https://checkout.stripe.com/test",
		},
	}

	handler := NewHandler(service, nil, nil)

	requestBody := CreateCheckoutRequest{
		CreditPackID: creditPackID,
		SuccessURL:   "https://example.com/payments/success",
		CancelURL:    "https://example.com/payments/cancel",
	}

	request := newCheckoutRequest(t, ownerID, requestBody)
	responseRecorder := httptest.NewRecorder()

	handler.CreateCheckout(responseRecorder, request)

	if responseRecorder.Code != http.StatusCreated {
		t.Fatalf(
			"expected status %d, got %d: %s",
			http.StatusCreated,
			responseRecorder.Code,
			responseRecorder.Body.String(),
		)
	}

	if service.receivedInput.OwnerID != ownerID {
		t.Errorf(
			"expected owner ID %s, got %s",
			ownerID,
			service.receivedInput.OwnerID,
		)
	}

	if service.receivedInput.CreditPackID != creditPackID {
		t.Errorf(
			"expected credit pack ID %s, got %s",
			creditPackID,
			service.receivedInput.CreditPackID,
		)
	}

	if service.receivedInput.SuccessURL != requestBody.SuccessURL {
		t.Errorf(
			"expected success URL %q, got %q",
			requestBody.SuccessURL,
			service.receivedInput.SuccessURL,
		)
	}

	if service.receivedInput.CancelURL != requestBody.CancelURL {
		t.Errorf(
			"expected cancel URL %q, got %q",
			requestBody.CancelURL,
			service.receivedInput.CancelURL,
		)
	}

	var response CreateCheckoutResponse

	if err := json.Unmarshal(
		responseRecorder.Body.Bytes(),
		&response,
	); err != nil {
		t.Fatalf("decode checkout response: %v", err)
	}

	if response.PaymentID != paymentID {
		t.Errorf(
			"expected payment ID %s, got %s",
			paymentID,
			response.PaymentID,
		)
	}

	if response.SessionID != "cs_test_123" {
		t.Errorf(
			"expected session ID %q, got %q",
			"cs_test_123",
			response.SessionID,
		)
	}

	if response.CheckoutURL != "https://checkout.stripe.com/test" {
		t.Errorf(
			"unexpected checkout URL %q",
			response.CheckoutURL,
		)
	}

	if contentType := responseRecorder.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf(
			"expected application/json content type, got %q",
			contentType,
		)
	}
}

func TestHandler_CreateCheckout_RejectsInvalidJSON(t *testing.T) {
	t.Parallel()

	service := &fakeCheckoutService{}
	handler := NewHandler(service, nil, nil)

	request := httptest.NewRequest(
		http.MethodPost,
		"/payments/checkout",
		bytes.NewBufferString(`{"credit_pack_id":`),
	)

	request = request.WithContext(
		context.WithValue(
			request.Context(),
			middleware.UserIDKey,
			uuid.New(),
		),
	)

	responseRecorder := httptest.NewRecorder()

	handler.CreateCheckout(responseRecorder, request)

	if responseRecorder.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			responseRecorder.Code,
		)
	}

	if service.called {
		t.Error("checkout service must not be called for invalid JSON")
	}
}

func TestHandler_CreateCheckout_RejectsUnknownFields(t *testing.T) {
	t.Parallel()

	service := &fakeCheckoutService{}
	handler := NewHandler(service, nil, nil)

	request := httptest.NewRequest(
		http.MethodPost,
		"/payments/checkout",
		bytes.NewBufferString(`{
			"credit_pack_id": "`+uuid.NewString()+`",
			"success_url": "https://example.com/success",
			"cancel_url": "https://example.com/cancel",
			"amount": 1
		}`),
	)

	request = request.WithContext(
		context.WithValue(
			request.Context(),
			middleware.UserIDKey,
			uuid.New(),
		),
	)

	responseRecorder := httptest.NewRecorder()

	handler.CreateCheckout(responseRecorder, request)

	if responseRecorder.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			responseRecorder.Code,
		)
	}

	if service.called {
		t.Error(
			"checkout service must not be called when unknown fields are supplied",
		)
	}
}

func TestHandler_CreateCheckout_RejectsUnauthenticatedRequest(
	t *testing.T,
) {
	t.Parallel()

	service := &fakeCheckoutService{}
	handler := NewHandler(service, nil, nil)

	request := newCheckoutRequest(
		t,
		uuid.Nil,
		CreateCheckoutRequest{
			CreditPackID: uuid.New(),
			SuccessURL:   "https://example.com/success",
			CancelURL:    "https://example.com/cancel",
		},
	)

	responseRecorder := httptest.NewRecorder()

	handler.CreateCheckout(responseRecorder, request)

	if responseRecorder.Code != http.StatusUnauthorized {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusUnauthorized,
			responseRecorder.Code,
		)
	}

	if service.called {
		t.Error(
			"checkout service must not be called for unauthenticated requests",
		)
	}
}

func TestHandler_CreateCheckout_ReturnsServiceError(t *testing.T) {
	t.Parallel()

	service := &fakeCheckoutService{
		err: paymentService.ErrPriceNotFound,
	}

	handler := NewHandler(service, nil, nil)

	request := newCheckoutRequest(
		t,
		uuid.New(),
		CreateCheckoutRequest{
			CreditPackID: uuid.New(),
			SuccessURL:   "https://example.com/success",
			CancelURL:    "https://example.com/cancel",
		},
	)

	responseRecorder := httptest.NewRecorder()

	handler.CreateCheckout(responseRecorder, request)

	if responseRecorder.Code != http.StatusNotFound {
		t.Fatalf(
			"expected status %d, got %d: %s",
			http.StatusNotFound,
			responseRecorder.Code,
			responseRecorder.Body.String(),
		)
	}
}

func TestHandler_CreateCheckout_RejectsMultipleJSONObjects(
	t *testing.T,
) {
	t.Parallel()

	service := &fakeCheckoutService{}
	handler := NewHandler(service, nil, nil)

	request := httptest.NewRequest(
		http.MethodPost,
		"/payments/checkout",
		bytes.NewBufferString(`{
			"credit_pack_id": "`+uuid.NewString()+`",
			"success_url": "https://example.com/success",
			"cancel_url": "https://example.com/cancel"
		}
		{
			"credit_pack_id": "`+uuid.NewString()+`"
		}`),
	)

	request = request.WithContext(
		context.WithValue(
			request.Context(),
			middleware.UserIDKey,
			uuid.New(),
		),
	)

	responseRecorder := httptest.NewRecorder()

	handler.CreateCheckout(responseRecorder, request)

	if responseRecorder.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			responseRecorder.Code,
		)
	}

	if service.called {
		t.Error(
			"checkout service must not be called for multiple JSON objects",
		)
	}
}

type fakeCheckoutService struct {
	output        *paymentService.CreateCheckoutOutput
	err           error
	called        bool
	receivedInput paymentService.CreateCheckoutInput
}

func (service *fakeCheckoutService) CreateCheckout(
	_ context.Context,
	input paymentService.CreateCheckoutInput,
) (*paymentService.CreateCheckoutOutput, error) {
	service.called = true
	service.receivedInput = input

	if service.err != nil {
		return nil, service.err
	}

	if service.output == nil {
		return nil, stdErrors.New("fake checkout output not configured")
	}

	return service.output, nil
}

func newCheckoutRequest(
	t *testing.T,
	ownerID uuid.UUID,
	body CreateCheckoutRequest,
) *http.Request {
	t.Helper()

	requestBody, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("encode checkout request: %v", err)
	}

	request := httptest.NewRequest(
		http.MethodPost,
		"/payments/checkout",
		bytes.NewReader(requestBody),
	)
	request.Header.Set("Content-Type", "application/json")

	if ownerID == uuid.Nil {
		return request
	}

	ctx := context.WithValue(
		request.Context(),
		middleware.UserIDKey,
		ownerID,
	)

	return request.WithContext(ctx)
}
