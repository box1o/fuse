package application

import "fuse/internal/infrastructure/db/postgres"

func (a *Application) setupRepositories() error {
	a.userRepo = postgres.NewUserRepository(a.db.DB)
	a.workspaceRepo = postgres.NewWorkspaceRepository(a.db.DB)
	a.computeRepo = postgres.NewComputeRepository(a.db.DB)
	a.creditUoW = postgres.NewCreditUnitOfWork(a.db.DB)
	a.creditPackRepo = postgres.NewCreditPackRepository(a.db.DB)

	return nil
}
