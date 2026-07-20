package application

import (
	paymentH "fuse/internal/interfaces/server/payment"
	"fuse/internal/services/auth"
	computeSvc "fuse/internal/services/compute"
	creditService "fuse/internal/services/credit"
	deviceAuthSvc "fuse/internal/services/deviceauth"
	"fuse/internal/services/mail"
	"fuse/internal/services/notification"
	paymentService "fuse/internal/services/payment"
	svcWorkspace "fuse/internal/services/workspace"
)

func (a *Application) setupServices() error {
	a.workspaceSvc = svcWorkspace.NewService(a.workspaceRepo)
	a.authSvc = auth.NewService(a.userRepo, a.sessMgr, a.workspaceSvc, a.eventManager.Bus())
	a.mailSvc = mail.NewService(a.cfg, a.eventManager)
	a.mailSvc.Setup()
	a.notificationSvc = notification.NewService(a.cfg)
	a.computeSvc = computeSvc.NewService(a.computeRepo)
	a.deviceAuthSvc = deviceAuthSvc.NewService(a.cfg, a.redis, a.computeSvc, a.userRepo)
	a.creditSvc = creditService.NewService(a.creditUoW, a.creditPackRepo)
	a.paymentSvc = paymentService.NewService(a.paymentRepo, a.creditSvc, a.creditSvc, a.stripePriceCatalog, a.stripeClient)
	a.paymentHandler = paymentH.NewHandler(a.paymentSvc, a.paymentSvc, a.stripeWebhookParser)
	return nil
}
