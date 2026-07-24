import { useMutation } from "@tanstack/react-query";
import { toast } from "sonner";

import { PAYMENTS_QUERY_KEYS } from "../../constants";
import { paymentsService } from "../../services";
import type {
    CheckoutSessionResponse,
    CreateCheckoutRequest,
} from "../../types";

export const useCreateCheckoutSession = () => {
    const mutation = useMutation<
        CheckoutSessionResponse,
        Error,
        CreateCheckoutRequest
    >({
        mutationKey: [PAYMENTS_QUERY_KEYS.CHECKOUT],

        mutationFn: async (request) => {
            const creditPackId = request.creditPackId.trim();
            const successUrl = request.successUrl.trim();
            const cancelUrl = request.cancelUrl.trim();

            if (!creditPackId) {
                throw new Error("Credit pack ID is required");
            }

            if (!successUrl) {
                throw new Error("Success URL is required");
            }

            if (!cancelUrl) {
                throw new Error("Cancel URL is required");
            }

            const response =
                await paymentsService.createCheckoutSession({
                    ...request,
                    creditPackId,
                    successUrl,
                    cancelUrl,
                });

            if (!response.success || !response.data) {
                throw new Error(
                    response.error ||
                        "Failed to create checkout session",
                );
            }

            const checkoutUrl = response.data.url?.trim();

            if (!checkoutUrl) {
                throw new Error(
                    "The backend returned an empty Stripe checkout URL",
                );
            }

            return {
                ...response.data,
                url: checkoutUrl,
            };
        },

        onError: (error) => {
            toast.error(error.message);
        },
    });

    return {
        createCheckout: mutation.mutateAsync,
        checkout: mutation.data,
        paymentId: mutation.data?.payment_id,
        sessionId: mutation.data?.session_id,
        checkoutUrl: mutation.data?.url,
        isLoading: mutation.isPending,
        error: mutation.error,
        resetCheckout: mutation.reset,
    };
};