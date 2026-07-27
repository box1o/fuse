package stripe

import (
	"encoding/json"
	"strings"

	paymentService "fuse/internal/services/payment"

	stripeWebhook "github.com/stripe/stripe-go/v86/webhook"
)

type WebhookParser struct {
	secret string
}

var _ paymentService.WebhookParser = (*WebhookParser)(nil)

func NewWebhookParser(secret string) (*WebhookParser, error) {
	secret = strings.TrimSpace(secret)
	if secret == "" {
		return nil, paymentService.ErrWebhookSecretRequired
	}

	return &WebhookParser{
		secret: secret,
	}, nil
}

func (parser *WebhookParser) ParseWebhook(payload []byte, signature string) (*paymentService.WebhookEvent, error) {
	if parser == nil || strings.TrimSpace(parser.secret) == "" {
		return nil, paymentService.ErrWebhookSecretRequired
	}

	event, err := stripeWebhook.ConstructEvent(
		payload,
		strings.TrimSpace(signature),
		parser.secret,
	)
	if err != nil {
		return nil, paymentService.ErrWebhookSignatureInvalid.WithErr(err)
	}

	mappedEvent := &paymentService.WebhookEvent{
		ID:   event.ID,
		Type: string(event.Type),
	}

	switch string(event.Type) {
	case paymentService.WebhookEventCheckoutCompleted,
		paymentService.WebhookEventAsyncPaymentSucceeded:
		session, err := parseCompletedCheckoutSession(event.Data.Raw)
		if err != nil {
			return nil, paymentService.ErrInvalidWebhookEvent.WithErr(err)
		}

		mappedEvent.CheckoutSession = session
	}

	return mappedEvent, nil
}

type stripeCheckoutSessionPayload struct {
	ID            string          `json:"id"`
	PaymentIntent json.RawMessage `json:"payment_intent"`
	AmountTotal   int64           `json:"amount_total"`
	Currency      string          `json:"currency"`
	PaymentStatus string          `json:"payment_status"`
}

func parseCompletedCheckoutSession(raw json.RawMessage) (*paymentService.CompletedCheckoutSession, error) {
	var payload stripeCheckoutSessionPayload

	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, err
	}

	paymentIntentID, err := parseExpandableID(payload.PaymentIntent)
	if err != nil {
		return nil, err
	}

	return &paymentService.CompletedCheckoutSession{
		SessionID:         payload.ID,
		PaymentIntentID:   paymentIntentID,
		Amount:            payload.AmountTotal,
		Currency:          payload.Currency,
		PaymentSuccessful: payload.PaymentStatus == "paid",
	}, nil
}

func parseExpandableID(raw json.RawMessage) (string, error) {
	if len(raw) == 0 || string(raw) == "null" {
		return "", nil
	}

	var id string
	if err := json.Unmarshal(raw, &id); err == nil {
		return strings.TrimSpace(id), nil
	}

	var object struct {
		ID string `json:"id"`
	}

	if err := json.Unmarshal(raw, &object); err != nil {
		return "", err
	}

	return strings.TrimSpace(object.ID), nil
}
