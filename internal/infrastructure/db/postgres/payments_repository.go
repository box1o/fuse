package postgres

import (
	"context"
	"errors"
	"strings"

	"fuse/internal/domain/payments"

	paymentsM "fuse/internal/domain/payments/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const defaultPendingUsageLimit = 100

type PaymentsRepository struct {
	db *gorm.DB
}

func NewPaymentsRepository(db *gorm.DB) payments.Repository {
	return &PaymentsRepository{
		db: db,
	}
}

func (r *PaymentsRepository) CreateBillingAccount(
	ctx context.Context,
	account *payments.BillingAccount,
) error {
	if account == nil {
		return payments.ErrInvalidBillingAccount
	}

	dbAccount := paymentsM.BillingAccountFromDomain(account)

	err := r.db.WithContext(ctx).Create(dbAccount).Error
	if err != nil {
		if isUniqueConstraintError(err) {
			return payments.ErrBillingAccountAlreadyExists.WithErr(err)
		}

		return payments.ErrCreateBillingAccountFailed.WithErr(err)
	}

	*account = *dbAccount.ToDomain()

	return nil
}

func (r *PaymentsRepository) FindBillingAccountByWorkspaceID(
	ctx context.Context,
	workspaceID uuid.UUID,
) (*payments.BillingAccount, error) {
	if workspaceID == uuid.Nil {
		return nil, payments.ErrWorkspaceIDRequired
	}

	var dbAccount paymentsM.DBBillingAccount

	err := r.db.WithContext(ctx).
		First(&dbAccount, "workspace_id = ?", workspaceID).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, payments.ErrBillingAccountNotFound
		}

		return nil, payments.ErrDatabaseOperation.WithErr(err)
	}

	return dbAccount.ToDomain(), nil
}

func (r *PaymentsRepository) FindBillingAccountByStripeCustomerID(
	ctx context.Context,
	stripeCustomerID string,
) (*payments.BillingAccount, error) {
	stripeCustomerID = strings.TrimSpace(stripeCustomerID)
	if stripeCustomerID == "" {
		return nil, payments.ErrStripeCustomerIDRequired
	}

	var dbAccount paymentsM.DBBillingAccount

	err := r.db.WithContext(ctx).
		First(
			&dbAccount,
			"stripe_customer_id = ?",
			stripeCustomerID,
		).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, payments.ErrBillingAccountNotFound
		}

		return nil, payments.ErrDatabaseOperation.WithErr(err)
	}

	return dbAccount.ToDomain(), nil
}

func (r *PaymentsRepository) UpsertSubscription(
	ctx context.Context,
	subscription *payments.Subscription,
) error {
	if subscription == nil {
		return payments.ErrInvalidSubscription
	}

	dbSubscription := paymentsM.SubscriptionFromDomain(subscription)

	err := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "stripe_subscription_id"},
			},
			DoUpdates: clause.AssignmentColumns([]string{
				"workspace_id",
				"status",
				"current_period_start",
				"current_period_end",
				"cancel_at_period_end",
				"updated_at",
			}),
		}).
		Create(dbSubscription).
		Error

	if err != nil {
		return payments.ErrUpsertSubscriptionFailed.WithErr(err)
	}

	saved, err := r.FindSubscriptionByStripeID(
		ctx,
		subscription.StripeSubscriptionID,
	)
	if err != nil {
		return err
	}

	*subscription = *saved

	return nil
}

func (r *PaymentsRepository) FindSubscriptionByWorkspaceID(
	ctx context.Context,
	workspaceID uuid.UUID,
) (*payments.Subscription, error) {
	if workspaceID == uuid.Nil {
		return nil, payments.ErrWorkspaceIDRequired
	}

	var dbSubscription paymentsM.DBSubscription

	err := r.db.WithContext(ctx).
		Order("updated_at DESC").
		First(
			&dbSubscription,
			"workspace_id = ?",
			workspaceID,
		).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, payments.ErrSubscriptionNotFound
		}

		return nil, payments.ErrDatabaseOperation.WithErr(err)
	}

	return dbSubscription.ToDomain(), nil
}

func (r *PaymentsRepository) FindSubscriptionByStripeID(
	ctx context.Context,
	stripeSubscriptionID string,
) (*payments.Subscription, error) {
	stripeSubscriptionID = strings.TrimSpace(stripeSubscriptionID)
	if stripeSubscriptionID == "" {
		return nil, payments.ErrStripeSubscriptionIDRequired
	}

	var dbSubscription paymentsM.DBSubscription

	err := r.db.WithContext(ctx).
		First(
			&dbSubscription,
			"stripe_subscription_id = ?",
			stripeSubscriptionID,
		).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, payments.ErrSubscriptionNotFound
		}

		return nil, payments.ErrDatabaseOperation.WithErr(err)
	}

	return dbSubscription.ToDomain(), nil
}

func (r *PaymentsRepository) CreateUsageRecord(
	ctx context.Context,
	record *payments.UsageRecord,
) error {
	if record == nil {
		return payments.ErrInvalidUsageRecord
	}

	dbRecord := paymentsM.UsageRecordFromDomain(record)

	err := r.db.WithContext(ctx).Create(dbRecord).Error
	if err != nil {
		if isUniqueConstraintError(err) {
			return payments.ErrIdempotencyKeyAlreadyExists.WithErr(err)
		}

		return payments.ErrCreateUsageRecordFailed.WithErr(err)
	}

	*record = *dbRecord.ToDomain()

	return nil
}

func (r *PaymentsRepository) FindUsageRecordByID(
	ctx context.Context,
	id uuid.UUID,
) (*payments.UsageRecord, error) {
	if id == uuid.Nil {
		return nil, payments.ErrInvalidUsageRecord
	}

	var dbRecord paymentsM.DBUsageRecord

	err := r.db.WithContext(ctx).
		First(&dbRecord, "id = ?", id).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, payments.ErrUsageRecordNotFound
		}

		return nil, payments.ErrDatabaseOperation.WithErr(err)
	}

	return dbRecord.ToDomain(), nil
}

func (r *PaymentsRepository) ListPendingUsage(
	ctx context.Context,
	limit int,
) ([]*payments.UsageRecord, error) {
	if limit <= 0 {
		limit = defaultPendingUsageLimit
	}

	var dbRecords []paymentsM.DBUsageRecord

	err := r.db.WithContext(ctx).
		Where("status = ?", string(payments.UsageStatusPending)).
		Order("occurred_at ASC").
		Limit(limit).
		Find(&dbRecords).
		Error

	if err != nil {
		return nil, payments.ErrDatabaseOperation.WithErr(err)
	}

	records := make([]*payments.UsageRecord, len(dbRecords))
	for i := range dbRecords {
		records[i] = dbRecords[i].ToDomain()
	}

	return records, nil
}

func (r *PaymentsRepository) UpdateUsageRecord(
	ctx context.Context,
	record *payments.UsageRecord,
) error {
	if record == nil {
		return payments.ErrInvalidUsageRecord
	}

	dbRecord := paymentsM.UsageRecordFromDomain(record)

	result := r.db.WithContext(ctx).
		Model(&paymentsM.DBUsageRecord{}).
		Where("id = ?", record.ID).
		Updates(map[string]any{
			"stripe_event_id": dbRecord.StripeEventID,
			"status":          dbRecord.Status,
			"updated_at":      dbRecord.UpdatedAt,
		})

	if result.Error != nil {
		return payments.ErrUpdateUsageRecordFailed.WithErr(result.Error)
	}

	if result.RowsAffected == 0 {
		return payments.ErrUsageRecordNotFound
	}

	return nil
}

func (r *PaymentsRepository) WebhookEventExists(
	ctx context.Context,
	stripeEventID string,
) (bool, error) {
	stripeEventID = strings.TrimSpace(stripeEventID)
	if stripeEventID == "" {
		return false, payments.ErrStripeEventIDRequired
	}

	var count int64

	err := r.db.WithContext(ctx).
		Model(&paymentsM.DBWebhookEvent{}).
		Where("stripe_event_id = ?", stripeEventID).
		Count(&count).
		Error

	if err != nil {
		return false, payments.ErrDatabaseOperation.WithErr(err)
	}

	return count > 0, nil
}

func (r *PaymentsRepository) CreateWebhookEvent(
	ctx context.Context,
	event *payments.WebhookEvent,
) error {
	if event == nil {
		return payments.ErrInvalidWebhookEvent
	}

	dbEvent := paymentsM.WebhookEventFromDomain(event)

	err := r.db.WithContext(ctx).Create(dbEvent).Error
	if err != nil {
		if isUniqueConstraintError(err) {
			return payments.ErrWebhookEventAlreadyProcessed.WithErr(err)
		}

		return payments.ErrCreateWebhookEventFailed.WithErr(err)
	}

	return nil
}

func isUniqueConstraintError(err error) bool {
	if err == nil {
		return false
	}

	message := strings.ToLower(err.Error())

	return strings.Contains(message, "duplicate key") ||
		strings.Contains(message, "unique constraint") ||
		strings.Contains(message, "sqlstate 23505")
}
