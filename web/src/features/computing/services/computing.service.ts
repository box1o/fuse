import type { ServiceResult } from "@/shared/types";
import { api } from "@/shared/services";

import type {
    ComputeNode,
    CreateNodeRequest,
    RenameNodeRequest,
    UpdateNodeCapabilitiesRequest,
    UpdateNodeStatusRequest,
} from "../types/node.types";

export const COMPUTE_ROUTES = {
    BASE: "/compute/node",
    BY_ID: (nodeId: string) => `/compute/node/${nodeId}`,
    BY_WORKSPACE: (workspaceId: string) => `/compute/node/workspace/${workspaceId}`,
    NAME: (nodeId: string) => `/compute/node/${nodeId}/name`,
    STATUS: (nodeId: string) => `/compute/node/${nodeId}/status`,
    CAPABILITIES: (nodeId: string) => `/compute/node/${nodeId}/capabilities`,
};

class ComputeService {
    async create(request: CreateNodeRequest): Promise<ServiceResult<ComputeNode>> {
        try {
            const { data } = await api.post<ComputeNode>(COMPUTE_ROUTES.BASE, request);
            return { data, success: true };
        } catch (error: unknown) {
            return {
                error: this.handleError(error, "create compute node"),
                success: false,
            };
        }
    }

    async list(): Promise<ServiceResult<ComputeNode[]>> {
        try {
            const { data } = await api.get<ComputeNode[]>(COMPUTE_ROUTES.BASE);
            return { data, success: true };
        } catch (error: unknown) {
            return {
                error: this.handleError(error, "list compute nodes"),
                success: false,
            };
        }
    }

    async getById(nodeId: string): Promise<ServiceResult<ComputeNode>> {
        try {
            const { data } = await api.get<ComputeNode>(COMPUTE_ROUTES.BY_ID(nodeId));
            return { data, success: true };
        } catch (error: unknown) {
            return {
                error: this.handleError(error, "get compute node"),
                success: false,
            };
        }
    }

    async listByWorkspace(workspaceId: string): Promise<ServiceResult<ComputeNode[]>> {
        try {
            const { data } = await api.get<ComputeNode[]>(COMPUTE_ROUTES.BY_WORKSPACE(workspaceId));
            return { data, success: true };
        } catch (error: unknown) {
            return {
                error: this.handleError(error, "list workspace compute nodes"),
                success: false,
            };
        }
    }

    async rename(nodeId: string, request: RenameNodeRequest): Promise<ServiceResult<ComputeNode>> {
        try {
            const { data } = await api.patch<ComputeNode>(COMPUTE_ROUTES.NAME(nodeId), request);
            return { data, success: true };
        } catch (error: unknown) {
            return {
                error: this.handleError(error, "rename compute node"),
                success: false,
            };
        }
    }

    async updateStatus(nodeId: string, request: UpdateNodeStatusRequest): Promise<ServiceResult<ComputeNode>> {
        try {
            const { data } = await api.patch<ComputeNode>(COMPUTE_ROUTES.STATUS(nodeId), request);
            return { data, success: true };
        } catch (error: unknown) {
            return {
                error: this.handleError(error, "update compute node status"),
                success: false,
            };
        }
    }

    async updateCapabilities(nodeId: string, request: UpdateNodeCapabilitiesRequest): Promise<ServiceResult<ComputeNode>> {
        try {
            const { data } = await api.patch<ComputeNode>(COMPUTE_ROUTES.CAPABILITIES(nodeId), request);
            return { data, success: true };
        } catch (error: unknown) {
            return {
                error: this.handleError(error, "update compute node capabilities"),
                success: false,
            };
        }
    }

    async delete(nodeId: string): Promise<ServiceResult<void>> {
        try {
            await api.delete(COMPUTE_ROUTES.BY_ID(nodeId));
            return { data: undefined, success: true };
        } catch (error: unknown) {
            return {
                error: this.handleError(error, "delete compute node"),
                success: false,
            };
        }
    }

    private handleError(error: unknown, operation: string): string {
        const responseError = error as {
            response?: {
                data?: {
                    message?: string;
                };
            };
            message?: string;
        };

        return responseError?.response?.data?.message || responseError?.message || `Failed to ${operation}`; }
}

export const computeService = new ComputeService();
