package application

import (
	"fuse/internal/services/auth"
	computeService "fuse/internal/services/compute"
	creditService "fuse/internal/services/credit"
	deviceAuthService "fuse/internal/services/deviceauth"
	"fuse/internal/services/mail"
	"fuse/internal/services/notification"
	paymentService "fuse/internal/services/payment"
	svcWorkspace "fuse/internal/services/workspace"
)

func (a *Application) setupServices() error {
	a.workspaceSvc = svcWorkspace.NewService(a.workspaceRepo, a.userRepo, a.eventManager.Bus())
	a.authSvc = auth.NewService(a.userRepo, a.sessMgr, a.workspaceSvc, a.eventManager.Bus())
	a.deviceAuthSvc = deviceAuthService.NewService(a.cfg, a.redis, a.cliCredentialRepo, a.userRepo)
	a.mailSvc = mail.NewService(a.cfg, a.eventManager)
	if err := a.mailSvc.Setup(); err != nil {
		return err
	}
	a.notificationSvc = notification.NewService(a.cfg)
	a.creditSvc = creditService.NewService(a.creditUoW, a.creditAccountRepo, a.creditPackRepo)
	a.paymentSvc = paymentService.NewService(a.paymentRepo, a.creditSvc, a.creditSvc, a.paymentPriceCatalog, a.stripeClient)
	a.computeSvc = computeService.NewService(a.computeRepo)

	return nil
}
