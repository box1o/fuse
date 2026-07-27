package payment

import (
	"bytes"
	"context"
	stdErrors "errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	paymentService "fuse/internal/services/payment"
)

var errWebhookService = stdErrors.New(
	"forced webhook service failure",
)

func TestHandler_HandleWebhook_ReturnsNoContent(t *testing.T) {
	t.Parallel()

	expectedEvent := &paymentService.WebhookEvent{
		ID:   "evt_test_123",
		Type: paymentService.WebhookEventCheckoutCompleted,
		CheckoutSession: &paymentService.CompletedCheckoutSession{
			SessionID:         "cs_test_123",
			PaymentIntentID:   "pi_test_123",
			Amount:            999,
			Currency:          "USD",
			PaymentSuccessful: true,
		},
	}

	parser := &fakeWebhookParser{
		event: expectedEvent,
	}
	service := &fakeWebhookService{}

	handler := NewHandler(nil, service, parser)

	payload := []byte(`{"id":"evt_test_123"}`)
	request := httptest.NewRequest(
		http.MethodPost,
		"/payments/webhook",
		bytes.NewReader(payload),
	)
	request.Header.Set(
		"Stripe-Signature",
		"t=123,v1=test-signature",
	)

	responseRecorder := httptest.NewRecorder()

	handler.HandleWebhook(responseRecorder, request)

	if responseRecorder.Code != http.StatusNoContent {
		t.Fatalf(
			"expected status %d, got %d: %s",
			http.StatusNoContent,
			responseRecorder.Code,
			responseRecorder.Body.String(),
		)
	}

	if !bytes.Equal(parser.payload, payload) {
		t.Errorf(
			"expected parser payload %q, got %q",
			payload,
			parser.payload,
		)
	}

	if parser.signature != "t=123,v1=test-signature" {
		t.Errorf(
			"unexpected Stripe signature %q",
			parser.signature,
		)
	}

	if service.event != expectedEvent {
		t.Error("expected parsed event to be passed to service")
	}

	if service.calls != 1 {
		t.Errorf(
			"expected webhook service to be called once, got %d",
			service.calls,
		)
	}
}

func TestHandler_HandleWebhook_RejectsEmptyPayload(t *testing.T) {
	t.Parallel()

	parser := &fakeWebhookParser{}
	service := &fakeWebhookService{}

	handler := NewHandler(nil, service, parser)

	request := httptest.NewRequest(
		http.MethodPost,
		"/payments/webhook",
		strings.NewReader("   "),
	)
	responseRecorder := httptest.NewRecorder()

	handler.HandleWebhook(responseRecorder, request)

	if responseRecorder.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d: %s",
			http.StatusBadRequest,
			responseRecorder.Code,
			responseRecorder.Body.String(),
		)
	}

	if parser.calls != 0 {
		t.Errorf(
			"expected parser not to be called, got %d calls",
			parser.calls,
		)
	}

	if service.calls != 0 {
		t.Errorf(
			"expected service not to be called, got %d calls",
			service.calls,
		)
	}
}

func TestHandler_HandleWebhook_RejectsOversizedPayload(
	t *testing.T,
) {
	t.Parallel()

	parser := &fakeWebhookParser{}
	service := &fakeWebhookService{}

	handler := NewHandler(nil, service, parser)

	payload := strings.Repeat(
		"a",
		maxWebhookRequestBodySize+1,
	)

	request := httptest.NewRequest(
		http.MethodPost,
		"/payments/webhook",
		strings.NewReader(payload),
	)
	responseRecorder := httptest.NewRecorder()

	handler.HandleWebhook(responseRecorder, request)

	if responseRecorder.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d: %s",
			http.StatusBadRequest,
			responseRecorder.Code,
			responseRecorder.Body.String(),
		)
	}

	if parser.calls != 0 {
		t.Errorf(
			"expected parser not to be called, got %d calls",
			parser.calls,
		)
	}

	if service.calls != 0 {
		t.Errorf(
			"expected service not to be called, got %d calls",
			service.calls,
		)
	}
}

func TestHandler_HandleWebhook_ReturnsParserError(t *testing.T) {
	t.Parallel()

	parser := &fakeWebhookParser{
		err: paymentService.ErrWebhookSignatureInvalid,
	}
	service := &fakeWebhookService{}

	handler := NewHandler(nil, service, parser)

	request := httptest.NewRequest(
		http.MethodPost,
		"/payments/webhook",
		strings.NewReader(`{"id":"evt_test_123"}`),
	)
	request.Header.Set(
		"Stripe-Signature",
		"invalid-signature",
	)

	responseRecorder := httptest.NewRecorder()

	handler.HandleWebhook(responseRecorder, request)

	if responseRecorder.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d: %s",
			http.StatusBadRequest,
			responseRecorder.Code,
			responseRecorder.Body.String(),
		)
	}

	if service.calls != 0 {
		t.Errorf(
			"expected service not to be called, got %d calls",
			service.calls,
		)
	}
}

func TestHandler_HandleWebhook_ReturnsServiceError(t *testing.T) {
	t.Parallel()

	event := &paymentService.WebhookEvent{
		ID:   "evt_test_123",
		Type: paymentService.WebhookEventCheckoutCompleted,
	}

	parser := &fakeWebhookParser{
		event: event,
	}
	service := &fakeWebhookService{
		err: errWebhookService,
	}

	handler := NewHandler(nil, service, parser)

	request := httptest.NewRequest(
		http.MethodPost,
		"/payments/webhook",
		strings.NewReader(`{"id":"evt_test_123"}`),
	)
	responseRecorder := httptest.NewRecorder()

	handler.HandleWebhook(responseRecorder, request)

	if responseRecorder.Code != http.StatusInternalServerError {
		t.Fatalf(
			"expected status %d, got %d: %s",
			http.StatusInternalServerError,
			responseRecorder.Code,
			responseRecorder.Body.String(),
		)
	}

	if service.calls != 1 {
		t.Errorf(
			"expected service to be called once, got %d",
			service.calls,
		)
	}
}

type fakeWebhookParser struct {
	event     *paymentService.WebhookEvent
	payload   []byte
	signature string
	calls     int
	err       error
}

func (parser *fakeWebhookParser) ParseWebhook(
	payload []byte,
	signature string,
) (*paymentService.WebhookEvent, error) {
	parser.calls++
	parser.payload = append([]byte(nil), payload...)
	parser.signature = signature

	if parser.err != nil {
		return nil, parser.err
	}

	return parser.event, nil
}

type fakeWebhookService struct {
	event *paymentService.WebhookEvent
	calls int
	err   error
}

func (service *fakeWebhookService) HandleWebhook(
	_ context.Context,
	event *paymentService.WebhookEvent,
) error {
	service.calls++
	service.event = event

	return service.err
}
