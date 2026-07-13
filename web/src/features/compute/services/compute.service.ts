import { api } from "@/shared/services";
import type { ServiceResult } from "@/shared/types";
import type { ComputeService } from "../types";

export const COMPUTE_ROUTES = {
    BASE: "/compute",
} as const;

class ComputeServiceClient {
    async list(): Promise<ServiceResult<ComputeService[]>> {
        try {
            const { data } = await api.get<ComputeService[]>(COMPUTE_ROUTES.BASE);
            return { data, success: true };
        } catch (error: any) {
            return {
                error: error?.response?.data?.message ?? error?.message ?? "Failed to load compute services",
                success: false,
            };
        }
    }
}

export const computeService = new ComputeServiceClient();
