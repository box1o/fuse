package application

import (
	"context"

	"fuse/internal/infrastructure/db/postgres"
)

func (a *Application) setupRepositories() error {
	a.userRepo = postgres.NewUserRepository(a.db.DB)
	a.workspaceRepo = postgres.NewWorkspaceRepository(a.db.DB)
	a.creditUoW = postgres.NewCreditUnitOfWork(a.db.DB)
	a.creditPackRepo = postgres.NewCreditPackRepository(a.db.DB)
	if err := syncConfiguredCreditPacks(context.Background(), a.creditPackRepo, a.cfg.Stripe.CreditPrices); err != nil {
		return err
	}
	a.paymentRepo = postgres.NewPaymentRepository(a.db.DB)
	a.paymentPriceCatalog = postgres.NewPaymentPriceCatalog(a.creditPackRepo)
	a.creditAccountRepo = postgres.NewCreditAccountRepository(a.db.DB)
	a.computeRepo = postgres.NewComputeRepository(a.db.DB)

	return nil
}
