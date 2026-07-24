package payment

import (
	"context"
	stdErrors "errors"
	"testing"
	"time"

	domainPayment "fuse/internal/domain/payment"
	creditService "fuse/internal/services/credit"

	"github.com/google/uuid"
)

var (
	errWebhookPaymentRepository = stdErrors.New(
		"forced webhook payment repository failure",
	)
	errWebhookCreditDeposit = stdErrors.New(
		"forced webhook credit deposit failure",
	)
)

func TestService_HandleWebhook_CompletesPaymentAndDepositsCredits(
	t *testing.T,
) {
	t.Parallel()

	payment := newWebhookTestPayment(t)
	repository := &fakeWebhookPaymentRepository{
		payment: payment,
	}
	depositor := &fakeCreditDepositor{}

	service := NewService(
		repository,
		nil,
		depositor,
		nil,
		nil,
	)

	err := service.HandleWebhook(
		context.Background(),
		&WebhookEvent{
			ID:   "evt_test_123",
			Type: WebhookEventCheckoutCompleted,
			CheckoutSession: &CompletedCheckoutSession{
				SessionID:         payment.ProviderSessionID,
				PaymentIntentID:   "pi_test_123",
				Amount:            payment.Amount,
				Currency:          "usd",
				PaymentSuccessful: true,
			},
		},
	)
	if err != nil {
		t.Fatalf("HandleWebhook() returned unexpected error: %v", err)
	}

	if depositor.calls != 1 {
		t.Fatalf(
			"expected one credit deposit, got %d",
			depositor.calls,
		)
	}

	if depositor.input.OwnerID != payment.OwnerID {
		t.Errorf(
			"expected owner ID %s, got %s",
			payment.OwnerID,
			depositor.input.OwnerID,
		)
	}

	if depositor.input.Amount.Value() != payment.Credits {
		t.Errorf(
			"expected %d deposited credits, got %d",
			payment.Credits,
			depositor.input.Amount.Value(),
		)
	}

	if depositor.input.ReferenceID != payment.ID.String() {
		t.Errorf(
			"expected reference ID %q, got %q",
			payment.ID.String(),
			depositor.input.ReferenceID,
		)
	}

	if depositor.input.ExternalReference != payment.ProviderSessionID {
		t.Errorf(
			"expected external reference %q, got %q",
			payment.ProviderSessionID,
			depositor.input.ExternalReference,
		)
	}

	expectedIdempotencyKey := "stripe:checkout:" +
		payment.ProviderSessionID

	if depositor.input.IdempotencyKey != expectedIdempotencyKey {
		t.Errorf(
			"expected idempotency key %q, got %q",
			expectedIdempotencyKey,
			depositor.input.IdempotencyKey,
		)
	}

	if repository.findProvider != domainPayment.ProviderStripe {
		t.Errorf(
			"expected Stripe provider lookup, got %q",
			repository.findProvider,
		)
	}

	if repository.findSessionID != payment.ProviderSessionID {
		t.Errorf(
			"expected session lookup %q, got %q",
			payment.ProviderSessionID,
			repository.findSessionID,
		)
	}

	if repository.updated == nil {
		t.Fatal("expected completed payment to be persisted")
	}

	if repository.updated.Status != domainPayment.StatusCompleted {
		t.Errorf(
			"expected completed status, got %q",
			repository.updated.Status,
		)
	}

	if repository.updated.ProviderPaymentID != "pi_test_123" {
		t.Errorf(
			"expected provider payment ID %q, got %q",
			"pi_test_123",
			repository.updated.ProviderPaymentID,
		)
	}

	if repository.updated.CompletedAt == nil {
		t.Error("expected payment completion timestamp")
	}
}

func TestService_HandleWebhook_IsIdempotentForCompletedPayment(
	t *testing.T,
) {
	t.Parallel()

	payment := newWebhookTestPayment(t)
	completedAt := time.Now().UTC()

	payment.Status = domainPayment.StatusCompleted
	payment.ProviderPaymentID = "pi_existing"
	payment.CompletedAt = &completedAt

	repository := &fakeWebhookPaymentRepository{
		payment: payment,
	}
	depositor := &fakeCreditDepositor{}

	service := NewService(
		repository,
		nil,
		depositor,
		nil,
		nil,
	)

	err := service.HandleWebhook(
		context.Background(),
		&WebhookEvent{
			ID:   "evt_duplicate",
			Type: WebhookEventCheckoutCompleted,
			CheckoutSession: &CompletedCheckoutSession{
				SessionID:         payment.ProviderSessionID,
				PaymentIntentID:   "pi_existing",
				Amount:            payment.Amount,
				Currency:          payment.Currency,
				PaymentSuccessful: true,
			},
		},
	)
	if err != nil {
		t.Fatalf("HandleWebhook() returned unexpected error: %v", err)
	}

	if depositor.calls != 0 {
		t.Errorf(
			"expected no duplicate deposit, got %d calls",
			depositor.calls,
		)
	}

	if repository.updateCalls != 0 {
		t.Errorf(
			"expected no payment update, got %d calls",
			repository.updateCalls,
		)
	}
}

func TestService_HandleWebhook_IgnoresUnpaidCheckout(
	t *testing.T,
) {
	t.Parallel()

	payment := newWebhookTestPayment(t)
	repository := &fakeWebhookPaymentRepository{
		payment: payment,
	}
	depositor := &fakeCreditDepositor{}

	service := NewService(
		repository,
		nil,
		depositor,
		nil,
		nil,
	)

	err := service.HandleWebhook(
		context.Background(),
		&WebhookEvent{
			ID:   "evt_unpaid",
			Type: WebhookEventCheckoutCompleted,
			CheckoutSession: &CompletedCheckoutSession{
				SessionID:         payment.ProviderSessionID,
				Amount:            payment.Amount,
				Currency:          payment.Currency,
				PaymentSuccessful: false,
			},
		},
	)
	if err != nil {
		t.Fatalf("HandleWebhook() returned unexpected error: %v", err)
	}

	if repository.findCalls != 0 {
		t.Errorf(
			"expected no payment lookup, got %d calls",
			repository.findCalls,
		)
	}

	if depositor.calls != 0 {
		t.Errorf(
			"expected no credit deposit, got %d calls",
			depositor.calls,
		)
	}

	if repository.updateCalls != 0 {
		t.Errorf(
			"expected no payment update, got %d calls",
			repository.updateCalls,
		)
	}
}

func TestService_HandleWebhook_IgnoresUnknownEvent(t *testing.T) {
	t.Parallel()

	repository := &fakeWebhookPaymentRepository{}
	depositor := &fakeCreditDepositor{}

	service := NewService(
		repository,
		nil,
		depositor,
		nil,
		nil,
	)

	err := service.HandleWebhook(
		context.Background(),
		&WebhookEvent{
			ID:   "evt_unknown",
			Type: "customer.created",
		},
	)
	if err != nil {
		t.Fatalf("HandleWebhook() returned unexpected error: %v", err)
	}

	if repository.findCalls != 0 {
		t.Errorf(
			"expected no payment lookup, got %d calls",
			repository.findCalls,
		)
	}

	if depositor.calls != 0 {
		t.Errorf(
			"expected no credit deposit, got %d calls",
			depositor.calls,
		)
	}
}

func TestService_HandleWebhook_RejectsNilEvent(t *testing.T) {
	t.Parallel()

	service := NewService(
		&fakeWebhookPaymentRepository{},
		nil,
		&fakeCreditDepositor{},
		nil,
		nil,
	)

	err := service.HandleWebhook(context.Background(), nil)

	if !stdErrors.Is(err, ErrInvalidWebhookEvent) {
		t.Fatalf(
			"expected invalid webhook event error, got %v",
			err,
		)
	}
}

func TestService_HandleWebhook_RejectsMissingCheckoutSession(
	t *testing.T,
) {
	t.Parallel()

	service := NewService(
		&fakeWebhookPaymentRepository{},
		nil,
		&fakeCreditDepositor{},
		nil,
		nil,
	)

	err := service.HandleWebhook(
		context.Background(),
		&WebhookEvent{
			ID:   "evt_test_123",
			Type: WebhookEventCheckoutCompleted,
		},
	)

	if !stdErrors.Is(err, ErrInvalidWebhookEvent) {
		t.Fatalf(
			"expected invalid webhook event error, got %v",
			err,
		)
	}
}

func TestService_HandleWebhook_RejectsMissingSessionID(
	t *testing.T,
) {
	t.Parallel()

	service := NewService(
		&fakeWebhookPaymentRepository{},
		nil,
		&fakeCreditDepositor{},
		nil,
		nil,
	)

	err := service.HandleWebhook(
		context.Background(),
		&WebhookEvent{
			ID:   "evt_test_123",
			Type: WebhookEventCheckoutCompleted,
			CheckoutSession: &CompletedCheckoutSession{
				SessionID:         "   ",
				PaymentSuccessful: true,
			},
		},
	)

	if !stdErrors.Is(
		err,
		domainPayment.ErrProviderSessionIDRequired,
	) {
		t.Fatalf(
			"expected provider session ID error, got %v",
			err,
		)
	}
}

func TestService_HandleWebhook_RejectsAmountMismatch(
	t *testing.T,
) {
	t.Parallel()

	payment := newWebhookTestPayment(t)
	repository := &fakeWebhookPaymentRepository{
		payment: payment,
	}
	depositor := &fakeCreditDepositor{}

	service := NewService(
		repository,
		nil,
		depositor,
		nil,
		nil,
	)

	err := service.HandleWebhook(
		context.Background(),
		&WebhookEvent{
			ID:   "evt_amount_mismatch",
			Type: WebhookEventCheckoutCompleted,
			CheckoutSession: &CompletedCheckoutSession{
				SessionID:         payment.ProviderSessionID,
				PaymentIntentID:   "pi_test_123",
				Amount:            payment.Amount + 1,
				Currency:          payment.Currency,
				PaymentSuccessful: true,
			},
		},
	)

	if !stdErrors.Is(err, ErrCheckoutAmountMismatch) {
		t.Fatalf(
			"expected checkout amount mismatch, got %v",
			err,
		)
	}

	if depositor.calls != 0 {
		t.Errorf(
			"expected no credit deposit, got %d calls",
			depositor.calls,
		)
	}

	if repository.updateCalls != 0 {
		t.Errorf(
			"expected no payment update, got %d calls",
			repository.updateCalls,
		)
	}
}

func TestService_HandleWebhook_RejectsCurrencyMismatch(
	t *testing.T,
) {
	t.Parallel()

	payment := newWebhookTestPayment(t)
	repository := &fakeWebhookPaymentRepository{
		payment: payment,
	}
	depositor := &fakeCreditDepositor{}

	service := NewService(
		repository,
		nil,
		depositor,
		nil,
		nil,
	)

	err := service.HandleWebhook(
		context.Background(),
		&WebhookEvent{
			ID:   "evt_currency_mismatch",
			Type: WebhookEventCheckoutCompleted,
			CheckoutSession: &CompletedCheckoutSession{
				SessionID:         payment.ProviderSessionID,
				PaymentIntentID:   "pi_test_123",
				Amount:            payment.Amount,
				Currency:          "eur",
				PaymentSuccessful: true,
			},
		},
	)

	if !stdErrors.Is(err, ErrCheckoutCurrencyMismatch) {
		t.Fatalf(
			"expected checkout currency mismatch, got %v",
			err,
		)
	}

	if depositor.calls != 0 {
		t.Errorf(
			"expected no credit deposit, got %d calls",
			depositor.calls,
		)
	}

	if repository.updateCalls != 0 {
		t.Errorf(
			"expected no payment update, got %d calls",
			repository.updateCalls,
		)
	}
}

func TestService_HandleWebhook_ReturnsRepositoryLookupError(
	t *testing.T,
) {
	t.Parallel()

	repository := &fakeWebhookPaymentRepository{
		findErr: errWebhookPaymentRepository,
	}

	service := NewService(
		repository,
		nil,
		&fakeCreditDepositor{},
		nil,
		nil,
	)

	err := service.HandleWebhook(
		context.Background(),
		&WebhookEvent{
			ID:   "evt_repository_error",
			Type: WebhookEventCheckoutCompleted,
			CheckoutSession: &CompletedCheckoutSession{
				SessionID:         "cs_test_123",
				Amount:            999,
				Currency:          "USD",
				PaymentSuccessful: true,
			},
		},
	)

	if !stdErrors.Is(err, errWebhookPaymentRepository) {
		t.Fatalf(
			"expected repository error, got %v",
			err,
		)
	}
}

func TestService_HandleWebhook_ReturnsDepositErrorWithoutCompletingPayment(
	t *testing.T,
) {
	t.Parallel()

	payment := newWebhookTestPayment(t)
	repository := &fakeWebhookPaymentRepository{
		payment: payment,
	}
	depositor := &fakeCreditDepositor{
		err: errWebhookCreditDeposit,
	}

	service := NewService(
		repository,
		nil,
		depositor,
		nil,
		nil,
	)

	err := service.HandleWebhook(
		context.Background(),
		&WebhookEvent{
			ID:   "evt_deposit_error",
			Type: WebhookEventCheckoutCompleted,
			CheckoutSession: &CompletedCheckoutSession{
				SessionID:         payment.ProviderSessionID,
				PaymentIntentID:   "pi_test_123",
				Amount:            payment.Amount,
				Currency:          payment.Currency,
				PaymentSuccessful: true,
			},
		},
	)

	if !stdErrors.Is(err, errWebhookCreditDeposit) {
		t.Fatalf(
			"expected deposit error, got %v",
			err,
		)
	}

	if repository.updateCalls != 0 {
		t.Errorf(
			"expected no payment update, got %d calls",
			repository.updateCalls,
		)
	}

	if payment.Status != domainPayment.StatusPending {
		t.Errorf(
			"expected payment to remain pending, got %q",
			payment.Status,
		)
	}
}

func TestService_HandleWebhook_ReturnsPaymentUpdateError(
	t *testing.T,
) {
	t.Parallel()

	payment := newWebhookTestPayment(t)
	repository := &fakeWebhookPaymentRepository{
		payment:   payment,
		updateErr: errWebhookPaymentRepository,
	}

	service := NewService(
		repository,
		nil,
		&fakeCreditDepositor{},
		nil,
		nil,
	)

	err := service.HandleWebhook(
		context.Background(),
		&WebhookEvent{
			ID:   "evt_update_error",
			Type: WebhookEventCheckoutCompleted,
			CheckoutSession: &CompletedCheckoutSession{
				SessionID:         payment.ProviderSessionID,
				PaymentIntentID:   "pi_test_123",
				Amount:            payment.Amount,
				Currency:          payment.Currency,
				PaymentSuccessful: true,
			},
		},
	)

	if !stdErrors.Is(err, errWebhookPaymentRepository) {
		t.Fatalf(
			"expected update error, got %v",
			err,
		)
	}
}

func newWebhookTestPayment(t *testing.T) *domainPayment.Payment {
	t.Helper()

	payment, err := domainPayment.NewPayment(
		domainPayment.NewPaymentInput{
			OwnerID:      uuid.New(),
			CreditPackID: uuid.New(),
			Credits:      500,
			Amount:       999,
			Currency:     "USD",
			Provider:     domainPayment.ProviderStripe,
		},
	)
	if err != nil {
		t.Fatalf("create test payment: %v", err)
	}

	if err := payment.AttachProviderSession("cs_test_123"); err != nil {
		t.Fatalf("attach provider session: %v", err)
	}

	return payment
}

type fakeWebhookPaymentRepository struct {
	payment *domainPayment.Payment

	findProvider  domainPayment.Provider
	findSessionID string
	findCalls     int
	updateCalls   int
	updated       *domainPayment.Payment

	findErr   error
	updateErr error
}

func (repository *fakeWebhookPaymentRepository) Create(
	context.Context,
	*domainPayment.Payment,
) error {
	return nil
}

func (repository *fakeWebhookPaymentRepository) FindByID(
	context.Context,
	uuid.UUID,
) (*domainPayment.Payment, error) {
	return nil, domainPayment.ErrPaymentNotFound
}

func (repository *fakeWebhookPaymentRepository) FindByProviderSessionID(
	_ context.Context,
	provider domainPayment.Provider,
	sessionID string,
) (*domainPayment.Payment, error) {
	repository.findCalls++
	repository.findProvider = provider
	repository.findSessionID = sessionID

	if repository.findErr != nil {
		return nil, repository.findErr
	}

	if repository.payment == nil {
		return nil, domainPayment.ErrPaymentNotFound
	}

	return repository.payment, nil
}

func (repository *fakeWebhookPaymentRepository) Update(
	_ context.Context,
	payment *domainPayment.Payment,
) error {
	repository.updateCalls++

	if repository.updateErr != nil {
		return repository.updateErr
	}

	copy := *payment
	repository.updated = &copy

	return nil
}

type fakeCreditDepositor struct {
	input creditService.DepositInput
	calls int
	err   error
}

func (depositor *fakeCreditDepositor) Deposit(
	_ context.Context,
	input creditService.DepositInput,
) error {
	depositor.calls++
	depositor.input = input

	return depositor.err
}
