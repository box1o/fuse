package application

import (
	"fmt"

	"fuse/pkg/config"
	"fuse/pkg/log"
	"fuse/pkg/shutdown"

	"fuse/internal/infrastructure/db/postgres"
	"fuse/internal/infrastructure/provider"
	"fuse/internal/infrastructure/redis"
	"fuse/internal/infrastructure/session"
	stripeInfrastructure "fuse/internal/infrastructure/stripe"

	"fuse/internal/services/auth"
	computeSvc "fuse/internal/services/compute"
	creditService "fuse/internal/services/credit"
	deviceAuthSvc "fuse/internal/services/deviceauth"
	eventsSvc "fuse/internal/services/events"
	"fuse/internal/services/mail"
	svcNotification "fuse/internal/services/notification"
	paymentSvc "fuse/internal/services/payment"
	svcWorkspace "fuse/internal/services/workspace"

	"fuse/internal/interfaces/server"
	authH "fuse/internal/interfaces/server/auth"
	computeH "fuse/internal/interfaces/server/compute"
	deviceAuthH "fuse/internal/interfaces/server/deviceauth"
	healthH "fuse/internal/interfaces/server/health"
	mailH "fuse/internal/interfaces/server/mail"
	authMW "fuse/internal/interfaces/server/middleware"
	paymentH "fuse/internal/interfaces/server/payment"
	wsH "fuse/internal/interfaces/server/workspace"

	"fuse/internal/domain/compute"
	domainCredit "fuse/internal/domain/credit"
	domainPayment "fuse/internal/domain/payment"
	"fuse/internal/domain/user"
	"fuse/internal/domain/workspace"
)

type Application struct {
	// Core
	cfg          *config.Config
	srv          *server.Server
	eventManager *eventsSvc.Service

	// Infrastructure
	db                  *postgres.PostgresDB
	redis               *redis.RedisClient
	authProv            *provider.AuthProvider
	sessMgr             *session.Manager
	stripeClient        *stripeInfrastructure.Client
	stripePriceCatalog  *stripeInfrastructure.ConfigPriceCatalog
	stripeWebhookParser *stripeInfrastructure.WebhookParser

	// Repositories
	userRepo       user.Repository
	workspaceRepo  workspace.Repository
	computeRepo    compute.Repository
	creditPackRepo domainCredit.PackRepository
	creditUoW      *postgres.CreditUnitOfWork
	paymentRepo    domainPayment.Repository

	// Services
	authSvc         *auth.Service
	workspaceSvc    *svcWorkspace.Service
	mailSvc         *mail.Service
	notificationSvc *svcNotification.Service
	computeSvc      *computeSvc.Service
	deviceAuthSvc   *deviceAuthSvc.Service
	creditSvc       *creditService.Service
	paymentSvc      *paymentSvc.Service

	// Middleware
	authMW *authMW.AuthMiddleware
	cliMW  *authMW.CLIMiddleware

	// Handlers
	healthHandler     *healthH.Handler
	authHandler       *authH.Handler
	workspaceHandler  *wsH.Handler
	computeHandler    *computeH.Handler
	deviceAuthHandler *deviceAuthH.Handler
	mailHandler       *mailH.Handler
	paymentHandler    *paymentH.Handler
}

func NewApplication() (*Application, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrConfigLoad, err)
	}
	return &Application{cfg: cfg}, nil
}

func (a *Application) Run() error {
	if err := a.initialize(); err != nil {
		return fmt.Errorf("application initialization failed: %w", err)
	}

	if err := a.srv.Start(); err != nil {
		return fmt.Errorf("HTTP server failed to start: %w", err)
	}

	shutdown.GracefulShutdown(a.srv, a.db)
	return nil
}

func (a *Application) initialize() error {
	steps := []struct {
		name string
		fn   func() error
	}{
		{"logger", a.setupLogger},
		{"database", a.setupDatabase},
		{"infrastructure", a.setupInfrastructure},
		{"repositories", a.setupRepositories},
		{"services", a.setupServices},
		{"handlers", a.setupHandlers},
		{"server", a.setupServer},
		{"dispatchEvents", a.dispatchEvents},
	}

	if a.cfg.Environment != "production" {
		log.Warn("Running in non-production mode: %s", a.cfg.Environment)
	}

	for _, s := range steps {
		if err := s.fn(); err != nil {
			return fmt.Errorf("%s initialization failed: %w", s.name, err)
		}
	}
	return nil
}

func (a *Application) setupLogger() error {
	if err := log.Setup(log.Console, ""); err != nil {
		return fmt.Errorf("%w: %v", ErrLoggerSetup, err)
	}
	return nil
}
