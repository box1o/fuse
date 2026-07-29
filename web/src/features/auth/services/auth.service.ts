import axios from "axios";
import type { HttpError, ServiceResult } from "@/shared/types";
import type { LogoutResponse, User } from "../types/auth.types";
import { api } from "@/shared/services";

const AUTH_ROUTES = {
    OAUTH: (provider: string) => `/auth/${provider}`,
    LOGOUT: "/auth/logout",
    STATUS: "/auth/status",
};

class AuthService {
    async getStatus(): Promise<ServiceResult<User>> {
        try {
            const { data } = await api.get<User>(AUTH_ROUTES.STATUS);
            return { data, success: true };
        } catch (error: unknown) {
            return {
                error: this.handleError(error, "get auth status"),
                success: false,
            };
        }
    }

    async logout(): Promise<ServiceResult<LogoutResponse>> {
        try {
            const { data } = await api.post<LogoutResponse>(AUTH_ROUTES.LOGOUT);
            return { data, success: true };
        } catch (error: unknown) {
            return {
                error: this.handleError(error, "logout"),
                success: false,
            };
        }
    }

    startOAuth(provider: string, returnTo?: string): void {
        if (!provider?.trim()) {
            throw new Error("Provider is required");
        }

        window.location.href = this.getLoginURL(provider, returnTo);
    }

    getLoginURL(provider: string, returnTo?: string): string {
        const loginURL = new URL(`${api.defaults.baseURL}${AUTH_ROUTES.OAUTH(provider)}`);
        if (returnTo) {
            loginURL.searchParams.set("return_to", returnTo);
        }

        return loginURL.toString();
    }

    private handleError(error: unknown, operation: string): string {
        if (axios.isAxiosError<HttpError>(error)) {
            return error.response?.data.message || error.response?.data.detail || error.message;
        }
        if (error instanceof Error) {
            return error.message;
        }

        return `Failed to ${operation}`;
    }
}

export const authService = new AuthService();
