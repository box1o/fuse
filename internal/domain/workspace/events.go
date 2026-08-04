package workspace

import "fuse/internal/domain/events"

const (
	WorkspaceAddMailEvent    = "workspace.add_member"
	WorkspaceRemoveMailEvent = "workspace.member_removed"
)

// Workspace created event
type WorkspaceAddMail struct {
	WorkspaceName string
	UserName      string
	UserEmail     string
	WorkspaceID   string
}

// Workspace remove event
type WorkspaceRemoveMail struct {
	WorkspaceName string
	UserName      string
	UserEmail     string
}

func NewWorkspaceAddMail(workspaceName, userName, userEmail, workspaceID string) *events.BaseEvent {
	return events.NewBaseEvent(WorkspaceAddMailEvent, WorkspaceAddMail{
		WorkspaceName: workspaceName,
		UserName:      userName,
		UserEmail:     userEmail,
		WorkspaceID:   workspaceID,
	})
}

func NewWorkspaceRemovedMail(workspaceName, userName, userEmail string) *events.BaseEvent {
	return events.NewBaseEvent(WorkspaceRemoveMailEvent, WorkspaceRemoveMail{
		WorkspaceName: workspaceName,
		UserName:      userName,
		UserEmail:     userEmail,
	})
}
