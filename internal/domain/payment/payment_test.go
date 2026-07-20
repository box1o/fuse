package payment

import (
	stdErrors "errors"
	"testing"

	"github.com/google/uuid"
)

func TestNewPayment_CreatesPendingPayment(t *testing.T) {
	t.Parallel()

	ownerID := uuid.New()
	packID := uuid.New()

	payment, err := NewPayment(NewPaymentInput{
		OwnerID:      ownerID,
		CreditPackID: packID,
		Credits:      500,
		Amount:       999,
		Currency:     " usd ",
		Provider:     ProviderStripe,
	})
	if err != nil {
		t.Fatalf("NewPayment() returned unexpected error: %v", err)
	}

	if payment.ID == uuid.Nil {
		t.Error("expected payment ID to be generated")
	}

	if payment.OwnerID != ownerID {
		t.Errorf(
			"expected owner ID %s, got %s",
			ownerID,
			payment.OwnerID,
		)
	}

	if payment.CreditPackID != packID {
		t.Errorf(
			"expected credit pack ID %s, got %s",
			packID,
			payment.CreditPackID,
		)
	}

	if payment.Credits != 500 {
		t.Errorf("expected 500 credits, got %d", payment.Credits)
	}

	if payment.Amount != 999 {
		t.Errorf("expected amount 999, got %d", payment.Amount)
	}

	if payment.Currency != "USD" {
		t.Errorf("expected currency USD, got %q", payment.Currency)
	}

	if payment.Status != StatusPending {
		t.Errorf(
			"expected status %q, got %q",
			StatusPending,
			payment.Status,
		)
	}

	if payment.Provider != ProviderStripe {
		t.Errorf(
			"expected provider %q, got %q",
			ProviderStripe,
			payment.Provider,
		)
	}

	if payment.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}

	if payment.UpdatedAt.IsZero() {
		t.Error("expected UpdatedAt to be set")
	}

	if payment.CompletedAt != nil {
		t.Error("expected CompletedAt to be nil")
	}
}

func TestNewPayment_ValidatesInput(t *testing.T) {
	t.Parallel()

	validInput := NewPaymentInput{
		OwnerID:      uuid.New(),
		CreditPackID: uuid.New(),
		Credits:      500,
		Amount:       999,
		Currency:     "USD",
		Provider:     ProviderStripe,
	}

	tests := []struct {
		name          string
		changeInput   func(*NewPaymentInput)
		expectedError error
	}{
		{
			name: "missing owner ID",
			changeInput: func(input *NewPaymentInput) {
				input.OwnerID = uuid.Nil
			},
			expectedError: ErrOwnerIDRequired,
		},
		{
			name: "missing credit pack ID",
			changeInput: func(input *NewPaymentInput) {
				input.CreditPackID = uuid.Nil
			},
			expectedError: ErrCreditPackIDRequired,
		},
		{
			name: "zero credits",
			changeInput: func(input *NewPaymentInput) {
				input.Credits = 0
			},
			expectedError: ErrCreditsMustBePositive,
		},
		{
			name: "negative credits",
			changeInput: func(input *NewPaymentInput) {
				input.Credits = -1
			},
			expectedError: ErrCreditsMustBePositive,
		},
		{
			name: "zero amount",
			changeInput: func(input *NewPaymentInput) {
				input.Amount = 0
			},
			expectedError: ErrAmountMustBePositive,
		},
		{
			name: "invalid currency length",
			changeInput: func(input *NewPaymentInput) {
				input.Currency = "US"
			},
			expectedError: ErrInvalidCurrency,
		},
		{
			name: "invalid currency characters",
			changeInput: func(input *NewPaymentInput) {
				input.Currency = "U1D"
			},
			expectedError: ErrInvalidCurrency,
		},
		{
			name: "invalid provider",
			changeInput: func(input *NewPaymentInput) {
				input.Provider = Provider("unknown")
			},
			expectedError: ErrInvalidProvider,
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			input := validInput
			test.changeInput(&input)

			_, err := NewPayment(input)
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

func TestPayment_AttachProviderSession(t *testing.T) {
	t.Parallel()

	payment := newTestPayment(t)

	err := payment.AttachProviderSession(" cs_test_123 ")
	if err != nil {
		t.Fatalf(
			"AttachProviderSession() returned unexpected error: %v",
			err,
		)
	}

	if payment.ProviderSessionID != "cs_test_123" {
		t.Errorf(
			"expected provider session ID %q, got %q",
			"cs_test_123",
			payment.ProviderSessionID,
		)
	}
}

func TestPayment_AttachProviderSession_RejectsDuplicateSession(
	t *testing.T,
) {
	t.Parallel()

	payment := newTestPayment(t)

	if err := payment.AttachProviderSession("cs_test_123"); err != nil {
		t.Fatalf("attach initial provider session: %v", err)
	}

	err := payment.AttachProviderSession("cs_test_456")
	if !stdErrors.Is(err, ErrProviderSessionAlreadyAttached) {
		t.Fatalf(
			"expected duplicate session error, got %v",
			err,
		)
	}
}

func TestPayment_Complete(t *testing.T) {
	t.Parallel()

	payment := newTestPayment(t)

	err := payment.Complete(" pi_test_123 ")
	if err != nil {
		t.Fatalf("Complete() returned unexpected error: %v", err)
	}

	if payment.Status != StatusCompleted {
		t.Errorf(
			"expected status %q, got %q",
			StatusCompleted,
			payment.Status,
		)
	}

	if payment.ProviderPaymentID != "pi_test_123" {
		t.Errorf(
			"expected provider payment ID %q, got %q",
			"pi_test_123",
			payment.ProviderPaymentID,
		)
	}

	if payment.CompletedAt == nil {
		t.Error("expected CompletedAt to be set")
	}
}

func TestPayment_CannotCompleteTwice(t *testing.T) {
	t.Parallel()

	payment := newTestPayment(t)

	if err := payment.Complete("pi_test_123"); err != nil {
		t.Fatalf("complete payment: %v", err)
	}

	err := payment.Complete("pi_test_123")
	if !stdErrors.Is(err, ErrPaymentNotPending) {
		t.Fatalf(
			"expected payment not pending error, got %v",
			err,
		)
	}
}

func TestPayment_Fail(t *testing.T) {
	t.Parallel()

	payment := newTestPayment(t)

	err := payment.Fail()
	if err != nil {
		t.Fatalf("Fail() returned unexpected error: %v", err)
	}

	if payment.Status != StatusFailed {
		t.Errorf(
			"expected status %q, got %q",
			StatusFailed,
			payment.Status,
		)
	}
}

func TestPayment_Cancel(t *testing.T) {
	t.Parallel()

	payment := newTestPayment(t)

	err := payment.Cancel()
	if err != nil {
		t.Fatalf("Cancel() returned unexpected error: %v", err)
	}

	if payment.Status != StatusCanceled {
		t.Errorf(
			"expected status %q, got %q",
			StatusCanceled,
			payment.Status,
		)
	}
}

func newTestPayment(t *testing.T) *Payment {
	t.Helper()

	payment, err := NewPayment(NewPaymentInput{
		OwnerID:      uuid.New(),
		CreditPackID: uuid.New(),
		Credits:      500,
		Amount:       999,
		Currency:     "USD",
		Provider:     ProviderStripe,
	})
	if err != nil {
		t.Fatalf("create test payment: %v", err)
	}

	return payment
}
