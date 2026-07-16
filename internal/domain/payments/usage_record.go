package payments

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// UsageQuantityUnit defines what one usage quantity means.
//
// One quantity unit equals one compute-second.
const UsageQuantityUnit = "compute-second"

type ResourceType string

const (
	ResourceTypeCPU ResourceType = "cpu"
	ResourceTypeGPU ResourceType = "gpu"
	ResourceTypeNPU ResourceType = "npu"
)

func (r ResourceType) IsValid() bool {
	switch r {
	case ResourceTypeCPU, ResourceTypeGPU, ResourceTypeNPU:
		return true
	default:
		return false
	}
}

type UsageStatus string

const (
	UsageStatusPending  UsageStatus = "pending"
	UsageStatusReported UsageStatus = "reported"
	UsageStatusFailed   UsageStatus = "failed"
)

func (s UsageStatus) IsValid() bool {
	switch s {
	case UsageStatusPending, UsageStatusReported, UsageStatusFailed:
		return true
	default:
		return false
	}
}

type UsageRecord struct {
	ID             uuid.UUID    `json:"id"`
	UserID         uuid.UUID    `json:"user_id"`
	ResourceType   ResourceType `json:"resource_type"`
	Quantity       int64        `json:"quantity"`
	OccurredAt     time.Time    `json:"occurred_at"`
	IdempotencyKey string       `json:"idempotency_key"`
	StripeEventID  string       `json:"stripe_event_id,omitempty"`
	Status         UsageStatus  `json:"status"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}

func NewUsageRecord(
	userID uuid.UUID,
	resourceType ResourceType,
	quantity int64,
	occurredAt time.Time,
	idempotencyKey string,
) (*UsageRecord, error) {
	if userID == uuid.Nil {
		return nil, ErrUserIDRequired
	}

	if !resourceType.IsValid() {
		return nil, ErrInvalidResourceType
	}

	if quantity <= 0 {
		return nil, ErrInvalidUsageQuantity
	}

	if occurredAt.IsZero() {
		return nil, ErrOccurredAtRequired
	}

	idempotencyKey = strings.TrimSpace(idempotencyKey)
	if idempotencyKey == "" {
		return nil, ErrIdempotencyKeyRequired
	}

	now := time.Now().UTC()

	return &UsageRecord{
		ID:             uuid.New(),
		UserID:         userID,
		ResourceType:   resourceType,
		Quantity:       quantity,
		OccurredAt:     occurredAt.UTC(),
		IdempotencyKey: idempotencyKey,
		Status:         UsageStatusPending,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

func (u *UsageRecord) MarkReported(stripeEventID string) error {
	if u == nil {
		return ErrInvalidUsageRecord
	}

	stripeEventID = strings.TrimSpace(stripeEventID)
	if stripeEventID == "" {
		return ErrStripeEventIDRequired
	}

	if u.Status == UsageStatusReported {
		return ErrUsageAlreadyReported
	}

	u.StripeEventID = stripeEventID
	u.Status = UsageStatusReported
	u.UpdatedAt = time.Now().UTC()

	return nil
}

func (u *UsageRecord) MarkFailed() error {
	if u == nil {
		return ErrInvalidUsageRecord
	}

	if u.Status == UsageStatusReported {
		return ErrUsageAlreadyReported
	}

	u.Status = UsageStatusFailed
	u.UpdatedAt = time.Now().UTC()

	return nil
}

func (u *UsageRecord) MarkPending() error {
	if u == nil {
		return ErrInvalidUsageRecord
	}

	if u.Status == UsageStatusReported {
		return ErrUsageAlreadyReported
	}

	u.Status = UsageStatusPending
	u.UpdatedAt = time.Now().UTC()

	return nil
}
