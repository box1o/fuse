package postgres

import (
	"context"
	stdErrors "errors"
	"fuse/internal/domain/credit"
	"fuse/internal/domain/credit/models"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreditPackRepository struct {
	db *gorm.DB
}

var _ credit.PackRepository = (*CreditPackRepository)(nil)

func NewCreditPackRepository(db *gorm.DB) credit.PackRepository {
	return &CreditPackRepository{
		db: db,
	}
}

func (r *CreditPackRepository) Create(ctx context.Context, pack *credit.Pack) error {
	if pack == nil {
		return credit.ErrPackNotFound
	}

	dbPack, err := models.FromDomainPack(pack)
	if err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Create(dbPack).Error; err != nil {
		return err
	}

	converted, err := dbPack.ToDomainPack()
	if err != nil {
		return err
	}

	*pack = *converted

	return nil
}

func (r *CreditPackRepository) FindByID(ctx context.Context, id uuid.UUID) (*credit.Pack, error) {
	if id == uuid.Nil {
		return nil, credit.ErrPackNotFound
	}
	var dbPack models.DBCreditPack
	err := r.db.WithContext(ctx).
		First(&dbPack, "id = ?", id).Error

	if stdErrors.Is(err, gorm.ErrRecordNotFound) {
		return nil, credit.ErrPackNotFound
	}

	if err != nil {
		return nil, err
	}

	pack, err := dbPack.ToDomainPack()
	if err != nil {
		return nil, err
	}

	return pack, nil

}

func (r *CreditPackRepository) FindByCode(ctx context.Context, code string) (*credit.Pack, error) {
	code = strings.TrimSpace(code)
	if code == "" {
		return nil, credit.ErrPackCodeRequired
	}

	var dbPack models.DBCreditPack

	err := r.db.WithContext(ctx).
		Where("code = ?", code).
		First(&dbPack).Error

	if stdErrors.Is(err, gorm.ErrRecordNotFound) {
		return nil, credit.ErrPackNotFound
	}

	if err != nil {
		return nil, err
	}

	pack, err := dbPack.ToDomainPack()
	if err != nil {
		return nil, err
	}

	return pack, nil
}

func (r *CreditPackRepository) ListActive(ctx context.Context) ([]*credit.Pack, error) {
	var dbPacks []models.DBCreditPack

	err := r.db.WithContext(ctx).
		Where(
			`active = ?
             AND stripe_price_id IS NOT NULL
             AND stripe_price_id <> ?
             AND price_amount > ?
             AND currency IS NOT NULL
             AND currency <> ?`,
			true,
			"",
			0,
			"",
		).
		Order("credits ASC").
		Find(&dbPacks).Error
	if err != nil {
		return nil, err
	}

	packs := make([]*credit.Pack, 0, len(dbPacks))

	for i := range dbPacks {
		pack, err := dbPacks[i].ToDomainPack()
		if err != nil {
			return nil, err
		}

		packs = append(packs, pack)
	}

	return packs, nil
}

func (r *CreditPackRepository) Update(ctx context.Context, pack *credit.Pack) error {
	if pack == nil || pack.ID == uuid.Nil {
		return credit.ErrPackNotFound
	}

	dbPack, err := models.FromDomainPack(pack)
	if err != nil {
		return err
	}

	result := r.db.WithContext(ctx).
		Model(&models.DBCreditPack{}).
		Where("id = ?", pack.ID).
		Updates(map[string]any{
			"code":            dbPack.Code,
			"name":            dbPack.Name,
			"credits":         dbPack.Credits,
			"active":          dbPack.Active,
			"stripe_price_id": dbPack.StripePriceID,
			"price_amount":    dbPack.PriceAmount,
			"currency":        dbPack.Currency,
			"updated_at":      dbPack.UpdatedAt,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return credit.ErrPackNotFound
	}

	return nil
}
