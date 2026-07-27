package payment

import (
	"context"
	"strings"

	domainCredit "fuse/internal/domain/credit"
	domainPayment "fuse/internal/domain/payment"
	creditService "fuse/internal/services/credit"
)

const (
	WebhookEventCheckoutCompleted     = "checkout.session.completed"
	WebhookEventAsyncPaymentSucceeded = "checkout.session.async_payment_succeeded"
)

type WebhookParser interface {
	ParseWebhook(
		payload []byte,
		signature string,
	) (*WebhookEvent, error)
}

type CreditDepositor interface {
	Deposit(
		ctx context.Context,
		input creditService.DepositInput,
	) error
}

type WebhookEvent struct {
	ID              string
	Type            string
	CheckoutSession *CompletedCheckoutSession
}

type CompletedCheckoutSession struct {
	SessionID         string
	PaymentIntentID   string
	Amount            int64
	Currency          string
	PaymentSuccessful bool
}

func (s *Service) HandleWebhook(
	ctx context.Context,
	event *WebhookEvent,
) error {
	if event == nil {
		return ErrInvalidWebhookEvent
	}

	switch event.Type {
	case WebhookEventCheckoutCompleted,
		WebhookEventAsyncPaymentSucceeded:
		return s.completeCheckout(ctx, event.CheckoutSession)

	default:
		// Unknown Stripe events are acknowledged but ignored.
		return nil
	}
}

func (s *Service) completeCheckout(
	ctx context.Context,
	session *CompletedCheckoutSession,
) error {
	if session == nil {
		return ErrInvalidWebhookEvent
	}

	sessionID := strings.TrimSpace(session.SessionID)
	if sessionID == "" {
		return domainPayment.ErrProviderSessionIDRequired
	}

	// checkout.session.completed can occur before payment settles for some
	// delayed payment methods. Credits must only be granted after payment.
	if !session.PaymentSuccessful {
		return nil
	}

	payment, err := s.payments.FindByProviderSessionID(ctx, domainPayment.ProviderStripe, sessionID)
	if err != nil {
		return err
	}

	if payment.Status == domainPayment.StatusCompleted {
		return nil
	}

	if payment.Status != domainPayment.StatusPending {
		return domainPayment.ErrPaymentNotPending
	}

	if session.Amount != payment.Amount {
		return ErrCheckoutAmountMismatch
	}

	currency := strings.ToUpper(strings.TrimSpace(session.Currency))
	if currency != payment.Currency {
		return ErrCheckoutCurrencyMismatch
	}

	creditAmount, err := domainCredit.NewAmount(payment.Credits)
	if err != nil {
		return err
	}

	if err := s.credits.Deposit(
		ctx,
		creditService.DepositInput{
			OwnerID:           payment.OwnerID,
			Amount:            creditAmount,
			ReferenceID:       payment.ID.String(),
			ExternalReference: sessionID,
			IdempotencyKey:    "stripe:checkout:" + sessionID,
		},
	); err != nil {
		return err
	}

	if err := payment.Complete(session.PaymentIntentID); err != nil {
		return err
	}

	return s.payments.Update(ctx, payment)
}
