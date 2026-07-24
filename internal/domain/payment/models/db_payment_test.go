package models

import (
	"testing"

	"fuse/internal/domain/payment"

	"github.com/google/uuid"
)

func TestPaymentConversion_RoundTrip(t *testing.T) {
	t.Parallel()

	original, err := payment.NewPayment(payment.NewPaymentInput{
		OwnerID:      uuid.New(),
		CreditPackID: uuid.New(),
		Credits:      500,
		Amount:       999,
		Currency:     "USD",
		Provider:     payment.ProviderStripe,
	})
	if err != nil {
		t.Fatalf("create payment: %v", err)
	}

	if err := original.AttachProviderSession("cs_test_123"); err != nil {
		t.Fatalf("attach provider session: %v", err)
	}

	dbPayment, err := FromDomainPayment(original)
	if err != nil {
		t.Fatalf("convert payment to database model: %v", err)
	}

	restored, err := dbPayment.ToDomainPayment()
	if err != nil {
		t.Fatalf("restore payment from database model: %v", err)
	}

	if restored.ID != original.ID {
		t.Errorf("expected ID %s, got %s", original.ID, restored.ID)
	}

	if restored.OwnerID != original.OwnerID {
		t.Errorf(
			"expected owner ID %s, got %s",
			original.OwnerID,
			restored.OwnerID,
		)
	}

	if restored.CreditPackID != original.CreditPackID {
		t.Errorf(
			"expected credit pack ID %s, got %s",
			original.CreditPackID,
			restored.CreditPackID,
		)
	}

	if restored.Credits != original.Credits {
		t.Errorf(
			"expected credits %d, got %d",
			original.Credits,
			restored.Credits,
		)
	}

	if restored.Amount != original.Amount {
		t.Errorf(
			"expected amount %d, got %d",
			original.Amount,
			restored.Amount,
		)
	}

	if restored.Currency != original.Currency {
		t.Errorf(
			"expected currency %q, got %q",
			original.Currency,
			restored.Currency,
		)
	}

	if restored.ProviderSessionID != original.ProviderSessionID {
		t.Errorf(
			"expected session ID %q, got %q",
			original.ProviderSessionID,
			restored.ProviderSessionID,
		)
	}
}

func TestFromDomainPayment_StoresEmptyProviderIDsAsNull(
	t *testing.T,
) {
	t.Parallel()

	domainPayment, err := payment.NewPayment(payment.NewPaymentInput{
		OwnerID:      uuid.New(),
		CreditPackID: uuid.New(),
		Credits:      500,
		Amount:       999,
		Currency:     "USD",
		Provider:     payment.ProviderStripe,
	})
	if err != nil {
		t.Fatalf("create payment: %v", err)
	}

	dbPayment, err := FromDomainPayment(domainPayment)
	if err != nil {
		t.Fatalf("convert payment: %v", err)
	}

	if dbPayment.ProviderSessionID != nil {
		t.Error("expected empty provider session ID to be stored as NULL")
	}

	if dbPayment.ProviderPaymentID != nil {
		t.Error("expected empty provider payment ID to be stored as NULL")
	}
}
