import { useMutation } from "@tanstack/react-query";
import { toast } from "sonner";
import { PAYMENTS_QUERY_KEYS } from "../../constants";
import { paymentsService } from "../../services";
import { usePaymentsStore } from "../../store";
import type {
    CheckoutSessionResponse,
    CreateCheckoutRequest,
} from "../../types";

export const useCreateCheckoutSession = () => {
    const { setCheckoutUrl, reset } = usePaymentsStore();

    const mutation = useMutation<
        CheckoutSessionResponse,
        Error,
        CreateCheckoutRequest
    >({
        mutationKey: [PAYMENTS_QUERY_KEYS.CHECKOUT],

        mutationFn: async (request) => {
            const creditPackId = request.creditPackId.trim();

            if (!creditPackId) {
                throw new Error("Credit pack ID is required");
            }

            const response =
                await paymentsService.createCheckoutSession({
                    ...request,
                    creditPackId,
                });

            if (!response.success || !response.data) {
                throw new Error(
                    response.error ||
                        "Failed to create checkout session",
                );
            }

            return response.data;
        },

        onSuccess: (checkout) => {
            setCheckoutUrl(checkout.checkout_url);
            toast.success("Checkout session created");

            window.location.assign(checkout.checkout_url);
        },

        onError: (error) => {
            reset();
            toast.error(error.message);
        },
    });

    return {
        createCheckout: mutation.mutate,
        checkoutUrl: mutation.data?.checkout_url,
        paymentId: mutation.data?.payment_id,
        sessionId: mutation.data?.session_id,
        isLoading: mutation.isPending,
        error: mutation.error,
    };
};