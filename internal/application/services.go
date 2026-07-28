package application

import (
	"fuse/internal/services/auth"
	creditService "fuse/internal/services/credit"
	"fuse/internal/services/mail"
	"fuse/internal/services/notification"
	paymentService "fuse/internal/services/payment"
	svcWorkspace "fuse/internal/services/workspace"
	computeService "fuse/internal/services/compute"
)

func (a *Application) setupServices() error {
	a.workspaceSvc = svcWorkspace.NewService(a.workspaceRepo)
	a.authSvc = auth.NewService(a.userRepo, a.sessMgr, a.workspaceSvc, a.eventManager.Bus())
	a.mailSvc = mail.NewService(a.cfg, a.eventManager)
	a.mailSvc.Setup()
	a.notificationSvc = notification.NewService(a.cfg)
	a.creditSvc = creditService.NewService(a.creditUoW, a.creditAccountRepo, a.creditPackRepo)
	a.paymentSvc = paymentService.NewService(a.paymentRepo, a.creditSvc, a.creditSvc, a.paymentPriceCatalog, a.stripeClient)
	a.computeSvc = computeService.NewService(a.computeRepo)

	return nil
}
