interface CreateWorkspaceRequest {
    name: string;
}


interface CreateWorkspaceRequest {
    name: string;
}

interface Workspace {
    id: string;
    name: string;
    owner_id: string;
    updated_at: string;
    created_at: string;
}

interface WorkspaceMember{
    id:  string;
	user_id: string;
	workspace_id: string;
	role: string;
	updated_at: string;
	created_at: string;
}

interface WorkspaceMemberRequest {
	user_mail: string; 
}

export type { CreateWorkspaceRequest, Workspace, WorkspaceMemberRequest, WorkspaceMember }
