import type { ServiceResult } from "@/shared/types";
import { api } from "@/shared/services";
import type { CreateWorkspaceRequest, Workspace, WorkspaceMember, WorkspaceMemberRequest } from "../types/workspace.types";

const BASE = "/workspaces";

export const WORKSPACE_ROUTES = {
    BASE,
    MEMBERS: `${BASE}/members`,
    SEARCH: `${BASE}/search`,
    BY_ID: (workspaceId: string) => `${BASE}/${workspaceId}`,
    BY_VS_MB_ID: (workspaceId: string, memberId: string) => `${BASE}/${workspaceId}/members/${memberId}`,
};
class WorkspaceService {

    async create(request: CreateWorkspaceRequest): Promise<ServiceResult<Workspace>> {
        try {
            const { data } = await api.post<Workspace>(WORKSPACE_ROUTES.BASE, request);
            return { data, success: true };
        } catch (error: any) {
            return {
                error: this.handleError(error, "create workspace"),
                success: false,
            };
        }
    }


    async list(): Promise<ServiceResult<Workspace[]>> {
        try {
            const { data } = await api.get<Workspace[]>(WORKSPACE_ROUTES.BASE);
            return { data, success: true };
        } catch (error: any) {
            return {
                error: this.handleError(error, "list workspaces"),
                success: false,
            };
        }
    }



    async delete(workspaceId: string): Promise<ServiceResult<void>> {
        try {
            await api.delete(WORKSPACE_ROUTES.BY_ID(workspaceId));
            return { data: undefined, success: true };
        } catch (error: any) {
            return {
                error: this.handleError(error, "delete workspace"),
                success: false,
            };
        }
    }


    // members
    async addWorkspaceMember(request: WorkspaceMemberRequest  ): Promise<ServiceResult<WorkspaceMember>> {
        try {
            const { data } = await api.post<WorkspaceMember>(WORKSPACE_ROUTES.MEMBERS, request);
            return { data, success: true };
        } catch (error: any) {
            return {
                error: this.handleError(error, "add workspace member"),
                success: false,
            };
        }
    }


    async memberList(): Promise<ServiceResult<WorkspaceMember[]>> {
        try {
            const { data } = await api.get<WorkspaceMember[]>(WORKSPACE_ROUTES.MEMBERS);
            return { data, success: true };
        } catch (error: any) {
            return {
                error: this.handleError(error, "list members"),
                success: false,
            };
        }
    }


    async deleteMember(workspaceId: string, memberId: string): Promise<ServiceResult<void>> {
        try {
            await api.delete(WORKSPACE_ROUTES.BY_VS_MB_ID(workspaceId, memberId));
            return { data: undefined, success: true };
        } catch (error: any) {
            return {
                error: this.handleError(error, "delete members"),
                success: false,
            };
        }
    }

    
    private handleError(error: any, operation: string): string {
        return error?.response?.data?.message ||
            error?.message ||
            `Failed to ${operation}`;
    }
}

export const workspaceService = new WorkspaceService();
