package payments

import (
	"strings"
	"time"
)

type WebhookEvent struct {
	StripeEventID string    `json:"stripe_event_id"`
	EventType     string    `json:"event_type"`
	ProcessedAt   time.Time `json:"processed_at"`
}

func NewWebhookEvent(
	stripeEventID string,
	eventType string,
) (*WebhookEvent, error) {
	stripeEventID = strings.TrimSpace(stripeEventID)
	if stripeEventID == "" {
		return nil, ErrStripeEventIDRequired
	}

	eventType = strings.TrimSpace(eventType)
	if eventType == "" {
		return nil, ErrWebhookEventTypeRequired
	}

	return &WebhookEvent{
		StripeEventID: stripeEventID,
		EventType:     eventType,
		ProcessedAt:   time.Now().UTC(),
	}, nil
}
