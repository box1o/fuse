package application

import (
	authH "fuse/internal/interfaces/server/auth"
	computeH "fuse/internal/interfaces/server/compute"
	healthH "fuse/internal/interfaces/server/health"
	mailH "fuse/internal/interfaces/server/mail"
	authMW "fuse/internal/interfaces/server/middleware"
	wsH "fuse/internal/interfaces/server/workspace"

	"fuse/internal/interfaces/server"
	"github.com/go-chi/chi/v5"
)

func (a *Application) setupHandlers() error {
	a.healthHandler = healthH.NewHandler(a.cfg)
	a.authMW = authMW.NewAuthMiddleware(a.authSvc, a.cfg)
	a.authHandler = authH.NewHandler(a.authSvc, a.cfg)
	a.workspaceHandler = wsH.NewHandler(a.workspaceSvc, a.cfg)
	a.computeHandler = computeH.NewHandler(a.cfg)
	a.mailHandler = mailH.NewHandler(a.cfg, a.mailSvc)
	return nil
}

func (a *Application) setupServer() error {
	opts := []server.ServerOption{
		server.WithRoutes(func(r chi.Router) {
			a.healthHandler.RegisterRoutes(r)
			a.authHandler.RegisterRoutes(r)
			a.workspaceHandler.RegisterRoutes(r, a.authMW)
			a.computeHandler.RegisterRoutes(r, a.authMW)
			a.mailHandler.RegisterRoutes(r, a.authMW)
		}),
	}
	a.srv = server.NewServer(a.cfg, opts...)
	return nil
}
