package workspace

import (
	"context"

	stderrors "errors"
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

	usr, err := s.userRepo.FindByID(ctx, ownerID)
	if err != nil {
		log.Error("can't find user by ownerId %s: %v", ownerID, err)
		return nil, err
	}

	wsMember := workspace.NewMember(ownerID, ws.ID, usr.Email, usr.Name, workspace.RoleOwner)

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

func (s *Service) AddWorkspaceMember(ctx context.Context, workspaceID uuid.UUID, userMail, userRole string) (*workspace.Member, error) {
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

	_, err = s.workspaceRepo.FindMember(ctx, workspaceID, user.ID)

	if err == nil {
		return nil, workspace.ErrMemberAlreadyExists
	}

	if !stderrors.Is(err, workspace.ErrMemberNotFound) {
		return nil, err
	}

	role, err := workspace.ParseRole(userRole)
	if err != nil {
		log.Error("failed to parse user role")
		return nil, err
	}

	wsMember := workspace.NewMember(user.ID, workspaceID, user.Email, user.Name, role)
	if err := s.workspaceRepo.AddMember(ctx, wsMember); err != nil {
		log.Error("failed to add member in db: %v", err)
		return nil, err
	}

	ws, err := s.workspaceRepo.GetWorkspaceByID(ctx, workspaceID)
	if err != nil {
		log.Error("can't find workspace by ID from db: %v", err)
		return nil, nil
	}

	if err := s.eventBus.Publish(ctx, workspace.NewWorkspaceAddMail(ws.Name, user.Name, user.Email)); err != nil {
		log.Error("failed to publish workspace member added event: %v", err)
	}
	return wsMember, nil
}

func (s *Service) DeleteWorkspaceMember(ctx context.Context, userMail string, workspaceID, memberID uuid.UUID) error {
	if workspaceID == uuid.Nil || memberID == uuid.Nil {
		return errors.ErrInternalServer.WithDetail("workspace ID and member ID cannot be empty")
	}

	member, err := s.workspaceRepo.FindMember(ctx, workspaceID, memberID)
	if err != nil {
		log.Error("failed to find member %s in workspace %s: %v", memberID, workspaceID, err)
		return err
	}

	ws, err := s.workspaceRepo.GetWorkspaceByID(ctx, workspaceID)
	if err != nil {
		log.Error("can't find workspace by ID from db: %v", err)
		return nil
	}

	if err := s.workspaceRepo.RemoveMember(ctx, workspaceID, memberID); err != nil {
		log.Error("failed to remove member: %s from workspace: %s, %v", memberID, workspaceID, err)
		return err
	}

	if err := s.eventBus.Publish(ctx, workspace.NewWorkspaceRemovedMail(ws.Name, member.Name, member.Mail)); err != nil {
		log.Error("failed to publish workspace remove member event: %v", err)
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
