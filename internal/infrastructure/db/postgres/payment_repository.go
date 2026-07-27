package postgres

import (
	"context"
	stdErrors "errors"
	"fuse/internal/domain/payment"
	"fuse/internal/domain/payment/models"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentRepository struct {
	db *gorm.DB
}

var _ payment.Repository = (*PaymentRepository)(nil)

func NewPaymentRepository(db *gorm.DB) payment.Repository {
	return &PaymentRepository{
		db: db,
	}
}

func (repository *PaymentRepository) Create(ctx context.Context, domainPayment *payment.Payment) error {
	if domainPayment == nil || domainPayment.ID == uuid.Nil {
		return payment.ErrPaymentNotFound
	}

	dbPayment, err := models.FromDomainPayment(domainPayment)
	if err != nil {
		return err
	}

	if err := repository.db.WithContext(ctx).
		Create(dbPayment).
		Error; err != nil {
		if isUniqueConstraintError(err) {
			return payment.ErrPaymentAlreadyExists
		}

		return err
	}

	convertedPayment, err := dbPayment.ToDomainPayment()
	if err != nil {
		return err
	}

	*domainPayment = *convertedPayment

	return nil
}

func (repository *PaymentRepository) FindByID(ctx context.Context, id uuid.UUID) (*payment.Payment, error) {
	if id == uuid.Nil {
		return nil, payment.ErrPaymentNotFound
	}

	var dbPayment models.DBPayment

	err := repository.db.WithContext(ctx).
		First(&dbPayment, "id = ?", id).Error

	if stdErrors.Is(err, gorm.ErrRecordNotFound) {
		return nil, payment.ErrPaymentNotFound
	}

	if err != nil {
		return nil, err
	}

	return dbPayment.ToDomainPayment()
}

func (repository *PaymentRepository) FindByProviderSessionID(ctx context.Context, provider payment.Provider, providerSessionID string) (*payment.Payment, error) {
	if !provider.IsValid() {
		return nil, payment.ErrInvalidProvider
	}

	providerSessionID = strings.TrimSpace(providerSessionID)
	if providerSessionID == "" {
		return nil, payment.ErrProviderSessionIDRequired
	}

	var dbPayment models.DBPayment

	err := repository.db.WithContext(ctx).
		Where(
			"provider = ? AND provider_session_id = ?",
			string(provider),
			providerSessionID,
		).
		First(&dbPayment).Error

	if stdErrors.Is(err, gorm.ErrRecordNotFound) {
		return nil, payment.ErrPaymentNotFound
	}

	if err != nil {
		return nil, err
	}

	return dbPayment.ToDomainPayment()
}

func (repository *PaymentRepository) Update(ctx context.Context, domainPayment *payment.Payment) error {
	if domainPayment == nil || domainPayment.ID == uuid.Nil {
		return payment.ErrPaymentNotFound
	}

	dbPayment, err := models.FromDomainPayment(domainPayment)
	if err != nil {
		return err
	}

	result := repository.db.WithContext(ctx).
		Model(&models.DBPayment{}).
		Where("id = ?", domainPayment.ID).
		Updates(map[string]any{
			"status":              dbPayment.Status,
			"provider_session_id": dbPayment.ProviderSessionID,
			"provider_payment_id": dbPayment.ProviderPaymentID,
			"completed_at":        dbPayment.CompletedAt,
			"updated_at":          dbPayment.UpdatedAt,
		})

	if result.Error != nil {
		if isUniqueConstraintError(result.Error) {
			return payment.ErrPaymentAlreadyExists
		}

		return result.Error
	}

	if result.RowsAffected == 0 {
		return payment.ErrPaymentNotFound
	}

	return nil
}
