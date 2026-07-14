package application

import (
	"fuse/internal/services/auth"
	"fuse/internal/services/mail"
	paymentsSvc "fuse/internal/services/payments"
	"fuse/internal/services/notification"
	svcWorkspace "fuse/internal/services/workspace"
)

func (a *Application) setupServices() error {
	a.workspaceSvc = svcWorkspace.NewService(a.workspaceRepo)
	a.authSvc = auth.NewService(a.userRepo, a.sessMgr, a.workspaceSvc, a.eventManager.Bus())
	a.mailSvc = mail.NewService(a.cfg, a.eventManager)
	a.mailSvc.Setup()
	a.paymentsSvc = paymentsSvc.NewService(a.cfg, a.paymentsRepo)
	a.notificationSvc = notification.NewService(a.cfg)
	return nil
}
