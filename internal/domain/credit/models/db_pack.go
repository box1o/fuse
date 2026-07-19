package models

import (
	"fmt"
	"fuse/internal/domain/credit"
	"fuse/internal/infrastructure/db"
)

type DBCreditPack struct {
	db.Model

	Code    string `gorm:"not null;size:100;uniqueIndex" json:"code"`
	Name    string `gorm:"not null;size:255" json:"name"`
	Credits int64  `gorm:"not null;check:credits > 0" json:"credits"`
	Active  bool   `gorm:"not null;default:true;index" json:"active"`
}

func (DBCreditPack) TableName() string {
	return "credit_packs"
}

func FromDomainPack(pack *credit.Pack) (*DBCreditPack, error) {
	if pack == nil {
		return nil, credit.ErrPackNotFound
	}

	return &DBCreditPack{
		Model:   db.Model{ID: pack.ID, CreatedAt: pack.CreatedAt, UpdatedAt: pack.UpdatedAt},
		Code:    pack.Code,
		Name:    pack.Name,
		Credits: pack.Credits.Value(),
		Active:  pack.Active,
	}, nil
}

func (pack *DBCreditPack) ToDomainPack() (*credit.Pack, error) {
	if pack == nil {
		return nil, credit.ErrPackNotFound
	}

	credits, err := credit.NewAmount(pack.Credits)
	if err != nil {
		return nil, fmt.Errorf("create credit pack amount: %w", err)
	}

	return &credit.Pack{
		ID:        pack.ID,
		Code:      pack.Code,
		Name:      pack.Name,
		Credits:   credits,
		Active:    pack.Active,
		CreatedAt: pack.CreatedAt,
		UpdatedAt: pack.UpdatedAt,
	}, nil
}
