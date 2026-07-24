package stripe

import (
	"encoding/json"
	stdErrors "errors"
	"testing"

	paymentService "fuse/internal/services/payment"
)

func TestNewWebhookParser_CreatesParser(t *testing.T) {
	t.Parallel()

	parser, err := NewWebhookParser("  whsec_test_123  ")
	if err != nil {
		t.Fatalf(
			"NewWebhookParser() returned unexpected error: %v",
			err,
		)
	}

	if parser == nil {
		t.Fatal("expected webhook parser")
	}

	if parser.secret != "whsec_test_123" {
		t.Errorf(
			"expected trimmed secret %q, got %q",
			"whsec_test_123",
			parser.secret,
		)
	}
}

func TestNewWebhookParser_RejectsEmptySecret(t *testing.T) {
	t.Parallel()

	parser, err := NewWebhookParser("   ")

	if !stdErrors.Is(
		err,
		paymentService.ErrWebhookSecretRequired,
	) {
		t.Fatalf(
			"expected webhook secret required error, got %v",
			err,
		)
	}

	if parser != nil {
		t.Error("expected nil parser")
	}
}

func TestWebhookParser_ParseWebhook_RejectsNilParser(
	t *testing.T,
) {
	t.Parallel()

	var parser *WebhookParser

	event, err := parser.ParseWebhook(
		[]byte(`{"id":"evt_test_123"}`),
		"signature",
	)

	if !stdErrors.Is(
		err,
		paymentService.ErrWebhookSecretRequired,
	) {
		t.Fatalf(
			"expected webhook secret required error, got %v",
			err,
		)
	}

	if event != nil {
		t.Error("expected nil event")
	}
}

func TestWebhookParser_ParseWebhook_RejectsInvalidSignature(
	t *testing.T,
) {
	t.Parallel()

	parser, err := NewWebhookParser("whsec_test_123")
	if err != nil {
		t.Fatalf("create webhook parser: %v", err)
	}

	event, err := parser.ParseWebhook(
		[]byte(`{"id":"evt_test_123"}`),
		"invalid-signature",
	)

	if !stdErrors.Is(
		err,
		paymentService.ErrWebhookSignatureInvalid,
	) {
		t.Fatalf(
			"expected invalid signature error, got %v",
			err,
		)
	}

	if event != nil {
		t.Error("expected nil event")
	}
}

func TestParseCompletedCheckoutSession_MapsPaidSession(
	t *testing.T,
) {
	t.Parallel()

	raw := json.RawMessage(`{
		"id": "cs_test_123",
		"payment_intent": "pi_test_123",
		"amount_total": 999,
		"currency": "usd",
		"payment_status": "paid"
	}`)

	session, err := parseCompletedCheckoutSession(raw)
	if err != nil {
		t.Fatalf(
			"parseCompletedCheckoutSession() returned error: %v",
			err,
		)
	}

	if session.SessionID != "cs_test_123" {
		t.Errorf(
			"expected session ID %q, got %q",
			"cs_test_123",
			session.SessionID,
		)
	}

	if session.PaymentIntentID != "pi_test_123" {
		t.Errorf(
			"expected payment intent ID %q, got %q",
			"pi_test_123",
			session.PaymentIntentID,
		)
	}

	if session.Amount != 999 {
		t.Errorf(
			"expected amount 999, got %d",
			session.Amount,
		)
	}

	if session.Currency != "usd" {
		t.Errorf(
			"expected currency %q, got %q",
			"usd",
			session.Currency,
		)
	}

	if !session.PaymentSuccessful {
		t.Error("expected payment to be successful")
	}
}

func TestParseCompletedCheckoutSession_MapsExpandablePaymentIntent(
	t *testing.T,
) {
	t.Parallel()

	raw := json.RawMessage(`{
		"id": "cs_test_123",
		"payment_intent": {
			"id": "pi_test_expanded"
		},
		"amount_total": 999,
		"currency": "usd",
		"payment_status": "paid"
	}`)

	session, err := parseCompletedCheckoutSession(raw)
	if err != nil {
		t.Fatalf(
			"parseCompletedCheckoutSession() returned error: %v",
			err,
		)
	}

	if session.PaymentIntentID != "pi_test_expanded" {
		t.Errorf(
			"expected expanded payment intent ID %q, got %q",
			"pi_test_expanded",
			session.PaymentIntentID,
		)
	}
}

func TestParseCompletedCheckoutSession_MapsUnpaidSession(
	t *testing.T,
) {
	t.Parallel()

	raw := json.RawMessage(`{
		"id": "cs_test_123",
		"payment_intent": null,
		"amount_total": 999,
		"currency": "usd",
		"payment_status": "unpaid"
	}`)

	session, err := parseCompletedCheckoutSession(raw)
	if err != nil {
		t.Fatalf(
			"parseCompletedCheckoutSession() returned error: %v",
			err,
		)
	}

	if session.PaymentSuccessful {
		t.Error("expected payment not to be successful")
	}

	if session.PaymentIntentID != "" {
		t.Errorf(
			"expected empty payment intent ID, got %q",
			session.PaymentIntentID,
		)
	}
}

func TestParseCompletedCheckoutSession_RejectsInvalidJSON(
	t *testing.T,
) {
	t.Parallel()

	session, err := parseCompletedCheckoutSession(
		json.RawMessage(`{"id":`),
	)

	if err == nil {
		t.Fatal("expected JSON parsing error")
	}

	if session != nil {
		t.Error("expected nil session")
	}
}

func TestParseExpandableID_ParsesStringID(t *testing.T) {
	t.Parallel()

	id, err := parseExpandableID(
		json.RawMessage(`"  pi_test_123  "`),
	)
	if err != nil {
		t.Fatalf("parseExpandableID() returned error: %v", err)
	}

	if id != "pi_test_123" {
		t.Errorf(
			"expected trimmed ID %q, got %q",
			"pi_test_123",
			id,
		)
	}
}

func TestParseExpandableID_ParsesObjectID(t *testing.T) {
	t.Parallel()

	id, err := parseExpandableID(
		json.RawMessage(`{"id":"  pi_test_123  "}`),
	)
	if err != nil {
		t.Fatalf("parseExpandableID() returned error: %v", err)
	}

	if id != "pi_test_123" {
		t.Errorf(
			"expected trimmed ID %q, got %q",
			"pi_test_123",
			id,
		)
	}
}

func TestParseExpandableID_ReturnsEmptyForNull(t *testing.T) {
	t.Parallel()

	id, err := parseExpandableID(json.RawMessage(`null`))
	if err != nil {
		t.Fatalf("parseExpandableID() returned error: %v", err)
	}

	if id != "" {
		t.Errorf("expected empty ID, got %q", id)
	}
}
