package application

import (
	"fuse/internal/services/auth"
	computeSvc "fuse/internal/services/compute"
	deviceAuthSvc "fuse/internal/services/deviceauth"
	"fuse/internal/services/mail"
	"fuse/internal/services/notification"
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
	return nil
}
