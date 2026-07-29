import { api } from "@/shared/services";
import type { DeviceAuthorizationRequest, DeviceDecisionResponse } from "../types";

const DEVICE_ROUTES = {
    request: (userCode: string) => `/auth/device/request/${encodeURIComponent(userCode)}`,
    approve: "/auth/device/approve",
    deny: "/auth/device/deny",
};

class DeviceService {
    async getRequest(userCode: string): Promise<DeviceAuthorizationRequest> {
        const response = await api.get<DeviceAuthorizationRequest>(DEVICE_ROUTES.request(userCode));
        return response.data;
    }

    async decide(userCode: string, approve: boolean): Promise<DeviceDecisionResponse> {
        const route = approve ? DEVICE_ROUTES.approve : DEVICE_ROUTES.deny;
        const response = await api.post<DeviceDecisionResponse>(route, { user_code: userCode });
        return response.data;
    }
}

export const deviceService = new DeviceService();
