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

interface WorkspaceMember {
    id: string;
    user_id: string;
    workspace_id: string;
    name: string;
    mail: string;
    role: string;
    updated_at: string;
    created_at: string;
}

interface WorkspaceMemberRequest {
    workspace_id: string;
    user_mail: string;
    role: string;
}

export type { CreateWorkspaceRequest, Workspace, WorkspaceMemberRequest, WorkspaceMember }
