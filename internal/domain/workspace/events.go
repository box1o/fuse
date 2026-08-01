package workspace

import "fuse/internal/domain/events"

const (
	WorkspaceAddMailEvent = "workspace.add_member"
)

// Workspace created event
type WorkspaceAddMail struct {
	WorkspaceName string
	UserName      string
	UserEmail     string
}

func NewWorkspaceAddMail(workspaceName, userName, userEmail string) *events.BaseEvent {
	return events.NewBaseEvent(WorkspaceAddMailEvent, WorkspaceAddMail{
		WorkspaceName: workspaceName,
		UserName:      userName,
		UserEmail:     userEmail,
	})
}
