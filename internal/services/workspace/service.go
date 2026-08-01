package workspace

import (
	"context"

	"fuse/internal/domain/user"
	"fuse/internal/domain/workspace"
	"fuse/internal/infrastructure/events"
	"fuse/pkg/errors"
	"fuse/pkg/log"

	"github.com/google/uuid"
)

type Service struct {
	workspaceRepo workspace.Repository
	userRepo      user.Repository
	eventBus      events.EventBus
}

func NewService(wsRepo workspace.Repository, userRepo user.Repository, eventBus events.EventBus) *Service {
	return &Service{
		workspaceRepo: wsRepo,
		userRepo:      userRepo,
		eventBus:      eventBus,
	}
}

func (s *Service) CreateWorkspace(ctx context.Context, name string, ownerID uuid.UUID) (*workspace.Workspace, error) {
	ws := workspace.NewWorkspace(name, ownerID)
	wsMember := workspace.NewMember(ownerID, ws.ID, workspace.RoleOwner)

	// existing, _ := s.workspaceRepo.FindByName(ctx, name)
	// if existing != nil {
	// 	return nil, errors.ErrNameExists.WithDetail("workspace name already exists in the system")
	// }
	//
	if err := s.workspaceRepo.CreateWorkspaceWithOwner(ctx, ws, wsMember); err != nil {
		log.Error("failed to create workspace in db: %v", err)
		return nil, err
	}
	return ws, nil

}

func (s *Service) GetUserWorkspaces(ctx context.Context, ownerID uuid.UUID) ([]*workspace.Workspace, error) {
	if ownerID == uuid.Nil {
		return nil, errors.ErrInternalServer.WithDetail("owner ID cannot be empty")
	}

	workspaces, err := s.workspaceRepo.GetUserWorkspaces(ctx, ownerID)
	if err != nil {
		log.Error("failed to retrieve workspaces for owner %s: %v", ownerID, err)
		return nil, err
	}
	return workspaces, nil
}

func (s *Service) DeleteWorkspace(ctx context.Context, wsID uuid.UUID) error {
	if wsID == uuid.Nil {
		return errors.ErrInternalServer.WithDetail("workspace ID cannot be empty")
	}

	if err := s.workspaceRepo.Delete(ctx, wsID); err != nil {
		log.Error("failed to delete workspace %s: %v", wsID, err)
		return err
	}

	return nil
}

func (s *Service) CreateWorkspaceMember(ctx context.Context, workspaceID uuid.UUID, userMail string) (*workspace.Member, error) {

	if workspaceID == uuid.Nil {
		return nil, errors.ErrInternalServer.WithDetail("workspace ID cannot be empty")
	}

	if userMail == "" {
		return nil, errors.ErrInternalServer.WithDetail("userMail cannot be empty")
	}

	user, err := s.userRepo.FindByEmail(ctx, userMail)
	if err != nil {
		log.Error("can't find user by email: %v", err)
		return nil, err
	}

	wsMember := workspace.NewMember(user.ID, workspaceID, workspace.RoleMember)

	if err := s.workspaceRepo.AddMember(ctx, wsMember); err != nil {
		log.Error("failed to add member in db: %v", err)
		return nil, err
	}

	ws, err := s.workspaceRepo.GetWorkspaceByID(ctx, workspaceID)
	if err != nil {
		log.Error("can't find workspace by ID from db: %v", err)
		return nil, err
	}

	if err := s.eventBus.Publish(ctx, workspace.NewWorkspaceAddMail(ws.Name, user.Name, user.Email)); err != nil {
		log.Error("failed to publish workspace member added event: %v", err)
	}
	return wsMember, nil
}

func (s *Service) DeleteWorkspaceMember(ctx context.Context, workspaceID, userID uuid.UUID) error {
	if workspaceID == uuid.Nil || userID == uuid.Nil {
		return errors.ErrInternalServer.WithDetail("workspace ID and user ID cannot be empty")
	}

	if err := s.workspaceRepo.RemoveMember(ctx, workspaceID, userID); err != nil {
		log.Error("failed to remove member: %s from workspace: %s, %v", userID, workspaceID, err)
		return err
	}

	return nil
}

func (s *Service) ListWorkspaceMembers(ctx context.Context, workspaceID uuid.UUID) ([]*workspace.Member, error) {
	if workspaceID == uuid.Nil {
		return nil, errors.ErrInternalServer.WithDetail("workspace ID cannot be empty")
	}

	members, err := s.workspaceRepo.ListMembers(ctx, workspaceID)
	if err != nil {
		log.Error("failed to list members for workspace %s: %v", workspaceID, err)
		return nil, err
	}
	return members, nil
}
