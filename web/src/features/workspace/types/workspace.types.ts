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


export type { CreateWorkspaceRequest, Workspace }
