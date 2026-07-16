import { useMutation } from "@tanstack/react-query";
import { toast } from "sonner";
import { paymentsService } from "../../services";
import type { CheckoutSessionResponse, SubscriptionPlanId } from "../../types";
import { usePaymentsStore } from "../../store";
import { PAYMENTS_QUERY_KEYS } from "../../constants";

interface CreateCheckoutRequest {
    planId: SubscriptionPlanId;
    successUrl: string;
    cancelUrl: string;
}

export const useCreateCheckoutSession = () => {
    const { setCheckoutUrl, reset } = usePaymentsStore();

    const mutation = useMutation<CheckoutSessionResponse, Error, CreateCheckoutRequest>({
        mutationKey: [PAYMENTS_QUERY_KEYS.CHECKOUT],
        mutationFn: async (request) => {
            if (!request?.planId) throw new Error("Missing plan ID");
            const response = await paymentsService.createCheckoutSession(
                request.planId,
                request.successUrl,
                request.cancelUrl,
            );
            if (!response.success || !response.data) {
                throw new Error(response.error || "Failed to create checkout session");
            }
            return response.data;
        },
        onSuccess: (data) => {
            setCheckoutUrl(data.url);
            toast.success("Checkout session created");
            window.location.href = data.url;
        },
        onError: (err) => {
            reset();
            toast.error(err.message || "Failed to create checkout session");
        },
    });

    return {
        checkout: mutation.mutate,
        checkoutUrl: mutation.data?.url,
        sessionId: mutation.data?.session_id,
        isLoading: mutation.isPending,
        isSuccess: mutation.isSuccess,
        isError: mutation.isError,
        error: mutation.error,
    };
};
