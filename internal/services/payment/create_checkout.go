package payment

import (
	"context"
	"errors"
	"fmt"
	"strings"

	domainPayment "fuse/internal/domain/payment"

	"github.com/google/uuid"
)

type CreateCheckoutInput struct {
	OwnerID      uuid.UUID
	CreditPackID uuid.UUID
	SuccessURL   string
	CancelURL    string
}

type CreateCheckoutOutput struct {
	PaymentID uuid.UUID
	SessionID string
	URL       string
}

func (s *Service) CreateCheckout(
	ctx context.Context,
	input CreateCheckoutInput,
) (*CreateCheckoutOutput, error) {
	if input.OwnerID == uuid.Nil {
		return nil, domainPayment.ErrOwnerIDRequired
	}

	if input.CreditPackID == uuid.Nil {
		return nil, domainPayment.ErrCreditPackIDRequired
	}

	successURL := strings.TrimSpace(input.SuccessURL)
	if successURL == "" {
		return nil, ErrSuccessURLRequired
	}

	cancelURL := strings.TrimSpace(input.CancelURL)
	if cancelURL == "" {
		return nil, ErrCancelURLRequired
	}

	pack, err := s.packs.GetActivePack(ctx, input.CreditPackID)
	if err != nil {
		return nil, err
	}

	price, err := s.prices.FindByPackCode(ctx, pack.Code)
	if err != nil {
		return nil, err
	}

	if err := validatePrice(price); err != nil {
		return nil, err
	}

	payment, err := domainPayment.NewPayment(
		domainPayment.NewPaymentInput{
			OwnerID:      input.OwnerID,
			CreditPackID: pack.ID,
			Credits:      pack.Credits.Value(),
			Amount:       price.Amount,
			Currency:     price.Currency,
			Provider:     domainPayment.ProviderStripe,
		},
	)
	if err != nil {
		return nil, err
	}

	if err := s.payments.Create(ctx, payment); err != nil {
		return nil, err
	}

	session, err := s.provider.CreateCheckoutSession(
		ctx,
		CreateCheckoutSessionInput{
			PaymentID:      payment.ID,
			OwnerID:        payment.OwnerID,
			CreditPackID:   payment.CreditPackID,
			PriceReference: price.Reference,
			SuccessURL:     successURL,
			CancelURL:      cancelURL,
		},
	)
	if err != nil {
		return nil, s.failPayment(
			ctx,
			payment,
			fmt.Errorf("create checkout session: %w", err),
		)
	}

	if err := validateCheckoutSession(payment, session); err != nil {
		return nil, s.failPayment(
			ctx,
			payment,
			fmt.Errorf("validate checkout session: %w", err),
		)
	}

	if err := payment.AttachProviderSession(session.SessionID); err != nil {
		return nil, err
	}

	if err := s.payments.Update(ctx, payment); err != nil {
		return nil, err
	}

	return &CreateCheckoutOutput{
		PaymentID: payment.ID,
		SessionID: session.SessionID,
		URL:       session.URL,
	}, nil
}

func validatePrice(price *Price) error {
	if price == nil {
		return ErrPriceNotFound
	}

	if strings.TrimSpace(price.Reference) == "" {
		return ErrPriceReferenceRequired
	}

	if price.Amount <= 0 {
		return ErrInvalidPrice
	}

	currency := strings.ToUpper(strings.TrimSpace(price.Currency))
	if len(currency) != 3 {
		return ErrInvalidPrice
	}

	return nil
}

func validateCheckoutSession(
	payment *domainPayment.Payment,
	session *CheckoutSession,
) error {
	if session == nil {
		return ErrCheckoutSessionCreationFailed
	}

	if session.Provider != payment.Provider {
		return domainPayment.ErrInvalidProvider
	}

	if strings.TrimSpace(session.SessionID) == "" {
		return ErrCheckoutSessionCreationFailed
	}

	if strings.TrimSpace(session.URL) == "" {
		return ErrCheckoutSessionCreationFailed
	}

	if session.Amount != payment.Amount {
		return ErrCheckoutAmountMismatch
	}

	sessionCurrency := strings.ToUpper(
		strings.TrimSpace(session.Currency),
	)

	if sessionCurrency != payment.Currency {
		return ErrCheckoutCurrencyMismatch
	}

	return nil
}

func (s *Service) failPayment(
	ctx context.Context,
	payment *domainPayment.Payment,
	cause error,
) error {
	if failErr := payment.Fail(); failErr != nil {
		return errors.Join(
			cause,
			fmt.Errorf("mark payment failed: %w", failErr),
		)
	}

	if updateErr := s.payments.Update(ctx, payment); updateErr != nil {
		return errors.Join(
			cause,
			fmt.Errorf(
				"persist failed payment %s: %w",
				payment.ID,
				updateErr,
			),
		)
	}

	return cause
}
