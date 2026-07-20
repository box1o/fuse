package payment

import (
	"context"
	stdErrors "errors"
	"testing"

	domainCredit "fuse/internal/domain/credit"
	domainPayment "fuse/internal/domain/payment"

	"github.com/google/uuid"
)

var (
	errPaymentRepository = stdErrors.New(
		"forced payment repository failure",
	)
	errPaymentProvider = stdErrors.New(
		"forced payment provider failure",
	)
)

func TestService_CreateCheckout_CreatesPaymentAndSession(
	t *testing.T,
) {
	t.Parallel()

	ownerID := uuid.New()
	pack := newCheckoutTestPack(t)

	paymentRepository := &fakePaymentRepository{}
	provider := &fakeProvider{
		createSession: func(
			_ context.Context,
			input CreateCheckoutSessionInput,
		) (*CheckoutSession, error) {
			if input.OwnerID != ownerID {
				t.Errorf(
					"expected owner ID %s, got %s",
					ownerID,
					input.OwnerID,
				)
			}

			if input.CreditPackID != pack.ID {
				t.Errorf(
					"expected pack ID %s, got %s",
					pack.ID,
					input.CreditPackID,
				)
			}

			if input.PriceReference != "price_test_500" {
				t.Errorf(
					"expected price reference %q, got %q",
					"price_test_500",
					input.PriceReference,
				)
			}

			return &CheckoutSession{
				Provider:    domainPayment.ProviderStripe,
				SessionID:   "cs_test_123",
				CheckoutURL: "https://checkout.stripe.com/test",
				Amount:      999,
				Currency:    "USD",
			}, nil
		},
	}

	service := NewService(
		paymentRepository,
		&fakeCreditPackReader{pack: pack},
		nil,
		&fakePriceCatalog{
			price: &Price{
				Reference: "price_test_500",
				Amount:    999,
				Currency:  "USD",
			},
		},
		provider,
	)

	output, err := service.CreateCheckout(
		context.Background(),
		CreateCheckoutInput{
			OwnerID:      ownerID,
			CreditPackID: pack.ID,
			SuccessURL:   "https://example.com/payment/success",
			CancelURL:    "https://example.com/payment/cancel",
		},
	)
	if err != nil {
		t.Fatalf(
			"CreateCheckout() returned unexpected error: %v",
			err,
		)
	}

	if output.PaymentID == uuid.Nil {
		t.Error("expected payment ID")
	}

	if output.SessionID != "cs_test_123" {
		t.Errorf(
			"expected session ID %q, got %q",
			"cs_test_123",
			output.SessionID,
		)
	}

	if output.CheckoutURL != "https://checkout.stripe.com/test" {
		t.Errorf(
			"unexpected checkout URL %q",
			output.CheckoutURL,
		)
	}

	if paymentRepository.created == nil {
		t.Fatal("expected payment to be created")
	}

	if paymentRepository.created.Credits != 500 {
		t.Errorf(
			"expected 500 credits, got %d",
			paymentRepository.created.Credits,
		)
	}

	if paymentRepository.created.Amount != 999 {
		t.Errorf(
			"expected amount 999, got %d",
			paymentRepository.created.Amount,
		)
	}

	if paymentRepository.updated == nil {
		t.Fatal("expected payment to be updated")
	}

	if paymentRepository.updated.ProviderSessionID != "cs_test_123" {
		t.Errorf(
			"expected attached provider session ID, got %q",
			paymentRepository.updated.ProviderSessionID,
		)
	}
}

func TestService_CreateCheckout_ReturnsProviderErrorAndFailsPayment(
	t *testing.T,
) {
	t.Parallel()

	pack := newCheckoutTestPack(t)
	repository := &fakePaymentRepository{}

	service := NewService(
		repository,
		&fakeCreditPackReader{pack: pack},
		nil,
		&fakePriceCatalog{
			price: &Price{
				Reference: "price_test_500",
				Amount:    999,
				Currency:  "USD",
			},
		},
		&fakeProvider{
			err: errPaymentProvider,
		},
	)

	_, err := service.CreateCheckout(
		context.Background(),
		CreateCheckoutInput{
			OwnerID:      uuid.New(),
			CreditPackID: pack.ID,
			SuccessURL:   "https://example.com/payment/success",
			CancelURL:    "https://example.com/payment/cancel",
		},
	)
	if !stdErrors.Is(err, errPaymentProvider) {
		t.Fatalf(
			"expected provider error, got %v",
			err,
		)
	}

	if repository.updated == nil {
		t.Fatal("expected failed payment to be persisted")
	}

	if repository.updated.Status != domainPayment.StatusFailed {
		t.Errorf(
			"expected failed status, got %q",
			repository.updated.Status,
		)
	}
}

func TestService_CreateCheckout_RejectsAmountMismatch(
	t *testing.T,
) {
	t.Parallel()

	pack := newCheckoutTestPack(t)

	service := NewService(
		&fakePaymentRepository{},
		&fakeCreditPackReader{pack: pack},
		nil,
		&fakePriceCatalog{
			price: &Price{
				Reference: "price_test_500",
				Amount:    999,
				Currency:  "USD",
			},
		},
		&fakeProvider{
			session: &CheckoutSession{
				Provider:    domainPayment.ProviderStripe,
				SessionID:   "cs_test_123",
				CheckoutURL: "https://checkout.stripe.com/test",
				Amount:      1999,
				Currency:    "USD",
			},
		},
	)

	_, err := service.CreateCheckout(
		context.Background(),
		CreateCheckoutInput{
			OwnerID:      uuid.New(),
			CreditPackID: pack.ID,
			SuccessURL:   "https://example.com/payment/success",
			CancelURL:    "https://example.com/payment/cancel",
		},
	)

	if !stdErrors.Is(err, ErrCheckoutAmountMismatch) {
		t.Fatalf(
			"expected amount mismatch error, got %v",
			err,
		)
	}
}

type fakePaymentRepository struct {
	created   *domainPayment.Payment
	updated   *domainPayment.Payment
	createErr error
	updateErr error
}

func (repository *fakePaymentRepository) Create(
	_ context.Context,
	payment *domainPayment.Payment,
) error {
	if repository.createErr != nil {
		return repository.createErr
	}

	copy := *payment
	repository.created = &copy

	return nil
}

func (repository *fakePaymentRepository) FindByID(
	context.Context,
	uuid.UUID,
) (*domainPayment.Payment, error) {
	return nil, domainPayment.ErrPaymentNotFound
}

func (repository *fakePaymentRepository) FindByProviderSessionID(
	context.Context,
	domainPayment.Provider,
	string,
) (*domainPayment.Payment, error) {
	return nil, domainPayment.ErrPaymentNotFound
}

func (repository *fakePaymentRepository) Update(
	_ context.Context,
	payment *domainPayment.Payment,
) error {
	if repository.updateErr != nil {
		return repository.updateErr
	}

	copy := *payment
	repository.updated = &copy

	return nil
}

type fakeCreditPackReader struct {
	pack *domainCredit.Pack
	err  error
}

func (reader *fakeCreditPackReader) GetActivePack(
	context.Context,
	uuid.UUID,
) (*domainCredit.Pack, error) {
	if reader.err != nil {
		return nil, reader.err
	}

	return reader.pack, nil
}

type fakePriceCatalog struct {
	price *Price
	err   error
}

func (catalog *fakePriceCatalog) FindByPackCode(
	context.Context,
	string,
) (*Price, error) {
	if catalog.err != nil {
		return nil, catalog.err
	}

	return catalog.price, nil
}

type fakeProvider struct {
	session       *CheckoutSession
	err           error
	createSession func(
		context.Context,
		CreateCheckoutSessionInput,
	) (*CheckoutSession, error)
}

func (provider *fakeProvider) CreateCheckoutSession(
	ctx context.Context,
	input CreateCheckoutSessionInput,
) (*CheckoutSession, error) {
	if provider.createSession != nil {
		return provider.createSession(ctx, input)
	}

	if provider.err != nil {
		return nil, provider.err
	}

	return provider.session, nil
}

func newCheckoutTestPack(t *testing.T) *domainCredit.Pack {
	t.Helper()

	credits, err := domainCredit.NewAmount(500)
	if err != nil {
		t.Fatalf("create credit amount: %v", err)
	}

	return &domainCredit.Pack{
		ID:      uuid.New(),
		Code:    "credits_500",
		Name:    "500 Credits",
		Credits: credits,
		Active:  true,
	}
}
