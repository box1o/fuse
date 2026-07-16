import { api } from "@/shared/services";
import type { ServiceResult } from "@/shared/types";
import type {
    CancelSubscriptionRequest,
    CheckoutSessionResponse,
    ProjectUsageRequest,
    WebhookRequest,
} from "../types";
import { PAYMENTS_ROUTES } from "../constants";

class PaymentsService {
    async createCheckoutSession(planId: "free" | "pro", successUrl: string, cancelUrl: string): Promise<ServiceResult<CheckoutSessionResponse>> {
        try {
            const { data } = await api.post<CheckoutSessionResponse>(PAYMENTS_ROUTES.CHECKOUT, {
                plan_id: planId,
                success_url: successUrl,
                cancel_url: cancelUrl,
            });
            return { data, success: true };
        } catch (error: any) {
            return { error: this.handleError(error, "create checkout session"), success: false };
        }
    }

    async recordUsage(request: ProjectUsageRequest): Promise<ServiceResult<void>> {
        try {
            await api.post(PAYMENTS_ROUTES.USAGE, request);
            return { data: undefined, success: true };
        } catch (error: any) {
            return { error: this.handleError(error, "record usage"), success: false };
        }
    }

    async cancelSubscription(request: CancelSubscriptionRequest): Promise<ServiceResult<void>> {
        try {
            await api.delete(PAYMENTS_ROUTES.SUBSCRIPTION, {
                data: request,
            });
            return { data: undefined, success: true };
        } catch (error: any) {
            return { error: this.handleError(error, "cancel subscription"), success: false };
        }
    }

    async sendWebhook(request: WebhookRequest): Promise<ServiceResult<void>> {
        try {
            await api.post(PAYMENTS_ROUTES.WEBHOOK, request.payload, {
                headers: {
                    "Stripe-Signature": request.signature,
                    "Content-Type": "application/json",
                },
            });
            return { data: undefined, success: true };
        } catch (error: any) {
            return { error: this.handleError(error, "send webhook"), success: false };
        }
    }

    private handleError(error: any, operation: string): string {
        return error?.response?.data?.message ||
            error?.message ||
            `Failed to ${operation}`;
    }
}

export const paymentsService = new PaymentsService();
