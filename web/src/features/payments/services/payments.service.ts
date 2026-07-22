import { api } from "@/shared/services";
import type { ServiceResult } from "@/shared/types";
import { PAYMENTS_ROUTES } from "../constants";
import type {
    CancelSubscriptionRequest,
    CheckoutSessionResponse,
    CreateCheckoutRequest,
    CreditBalanceResponse,
    CreditPack,
    ProjectUsageRequest,
    WebhookRequest,
} from "../types";

class PaymentsService {
    async createCheckoutSession(
        request: CreateCheckoutRequest,
    ): Promise<ServiceResult<CheckoutSessionResponse>> {
        try {
            const { data } = await api.post<CheckoutSessionResponse>(
                PAYMENTS_ROUTES.CHECKOUT,
                {
                    credit_pack_id: request.creditPackId,
                    success_url: request.successUrl,
                    cancel_url: request.cancelUrl,
                },
            );

            return {
                data,
                success: true,
            };
        } catch (error: unknown) {
            return {
                error: this.handleError(error, "create checkout session"),
                success: false,
            };
        }
    }

    async recordUsage(
        request: ProjectUsageRequest,
    ): Promise<ServiceResult<void>> {
        try {
            await api.post(PAYMENTS_ROUTES.USAGE, request);

            return {
                data: undefined,
                success: true,
            };
        } catch (error: unknown) {
            return {
                error: this.handleError(error, "record usage"),
                success: false,
            };
        }
    }

    async cancelSubscription(
        request: CancelSubscriptionRequest,
    ): Promise<ServiceResult<void>> {
        try {
            await api.delete(PAYMENTS_ROUTES.SUBSCRIPTION, {
                data: request,
            });

            return {
                data: undefined,
                success: true,
            };
        } catch (error: unknown) {
            return {
                error: this.handleError(error, "cancel subscription"),
                success: false,
            };
        }
    }

    async sendWebhook(
        request: WebhookRequest,
    ): Promise<ServiceResult<void>> {
        try {
            await api.post(PAYMENTS_ROUTES.WEBHOOK, request.payload, {
                headers: {
                    "Stripe-Signature": request.signature,
                    "Content-Type": "application/json",
                },
            });

            return {
                data: undefined,
                success: true,
            };
        } catch (error: unknown) {
            return {
                error: this.handleError(error, "send webhook"),
                success: false,
            };
        }
    }

    private handleError(error: unknown, operation: string): string {
        if (
            typeof error === "object" &&
            error !== null &&
            "response" in error
        ) {
            const responseError = error as {
                response?: {
                    data?: {
                        message?: string;
                    };
                };
            };

            if (responseError.response?.data?.message) {
                return responseError.response.data.message;
            }
        }

        if (error instanceof Error) {
            return error.message;
        }

        return `Failed to ${operation}`;
    }

    async listCreditPacks(): Promise<ServiceResult<CreditPack[]>> {
        try {
            const { data } = await api.get<CreditPack[]>(
                PAYMENTS_ROUTES.CREDIT_PACKS,
            );

            return {
                data,
                success: true,
            };
        } catch (error: unknown) {
            return {
                error: this.handleError(error, "load credit packs"),
                success: false,
            };
        }
    }

    async getCreditBalance(): Promise< ServiceResult<CreditBalanceResponse>> {
        try {
            const { data } =
                await api.get<CreditBalanceResponse>(
                    PAYMENTS_ROUTES.CREDIT_BALANCE,
                );

            return {
                data,
                success: true,
            };
        } catch (error: unknown) {
            return {
                error: this.handleError(
                    error,
                    "load credit balance",
                ),
                success: false,
            };
        }
    }
}

export const paymentsService = new PaymentsService();