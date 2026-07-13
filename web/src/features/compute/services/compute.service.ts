import { api } from "@/shared/services";
import type { ServiceResult } from "@/shared/types";
import type {
    ComputeNode,
    RegisterComputeNodeRequest,
    UpdateComputeNodeRequest,
} from "../types";

export const COMPUTE_ROUTES = {
    NODES: "/compute/nodes/",
    NODE: (nodeId: string) => `/compute/nodes/${nodeId}`,
} as const;

class ComputeService {
    async register(request: RegisterComputeNodeRequest): Promise<ServiceResult<ComputeNode>> {
        try {
            const { data } = await api.post<ComputeNode>(COMPUTE_ROUTES.NODES, request);
            return { data, success: true };
        } catch (error: unknown) {
            return this.failure(error, "register compute node");
        }
    }

    async list(): Promise<ServiceResult<ComputeNode[]>> {
        try {
            const { data } = await api.get<ComputeNode[]>(COMPUTE_ROUTES.NODES);
            return { data, success: true };
        } catch (error: unknown) {
            return this.failure(error, "load compute nodes");
        }
    }

    async get(nodeId: string): Promise<ServiceResult<ComputeNode>> {
        try {
            const { data } = await api.get<ComputeNode>(COMPUTE_ROUTES.NODE(nodeId));
            return { data, success: true };
        } catch (error: unknown) {
            return this.failure(error, "load compute node");
        }
    }

    async update(
        nodeId: string,
        request: UpdateComputeNodeRequest,
    ): Promise<ServiceResult<ComputeNode>> {
        try {
            const { data } = await api.patch<ComputeNode>(COMPUTE_ROUTES.NODE(nodeId), request);
            return { data, success: true };
        } catch (error: unknown) {
            return this.failure(error, "update compute node");
        }
    }

    async delete(nodeId: string): Promise<ServiceResult<void>> {
        try {
            await api.delete(COMPUTE_ROUTES.NODE(nodeId));
            return { data: undefined, success: true };
        } catch (error: unknown) {
            return this.failure(error, "delete compute node");
        }
    }

    private failure<T>(error: unknown, operation: string): ServiceResult<T> {
        const apiError = error as {
            message?: string;
            response?: { data?: { message?: string; detail?: string } };
        };

        return {
            error:
                apiError.response?.data?.detail ??
                apiError.response?.data?.message ??
                apiError.message ??
                `Failed to ${operation}`,
            success: false,
        };
    }
}

export const computeService = new ComputeService();
