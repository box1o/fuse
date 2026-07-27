package payment

import (
	domainPayment "fuse/internal/domain/payment"
)

type Service struct {
	payments domainPayment.Repository
	packs    CreditPackReader
	credits  CreditDepositor
	prices   PriceCatalog
	provider Provider
}

func NewService(
	payments domainPayment.Repository,
	packs CreditPackReader,
	credits CreditDepositor,
	prices PriceCatalog,
	provider Provider,
) *Service {
	return &Service{
		payments: payments,
		packs:    packs,
		credits:  credits,
		prices:   prices,
		provider: provider,
	}
}
