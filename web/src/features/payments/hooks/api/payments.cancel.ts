import { useMutation } from "@tanstack/react-query";
import { toast } from "sonner";
import { paymentsService } from "../../services";
import { usePaymentsStore } from "../../store";
import type { CancelSubscriptionRequest } from "../../types";
import { PAYMENTS_QUERY_KEYS } from "../../constants";

export const useCancelSubscription = () => {
    const { setError } = usePaymentsStore();

    const mutation = useMutation<void, Error, CancelSubscriptionRequest>({
        mutationKey: [PAYMENTS_QUERY_KEYS.CANCEL],
        mutationFn: async (request) => {
            if (!request?.workspace_id) throw new Error("Missing workspace ID");
            const response = await paymentsService.cancelSubscription(request);
            if (!response.success) {
                throw new Error(response.error || "Failed to cancel subscription");
            }
        },
        onSuccess: () => {
            setError(null);
            toast.success("Subscription cancellation requested");
        },
        onError: (err) => {
            setError(err.message);
            toast.error(err.message || "Failed to cancel subscription");
        },
    });

    return {
        cancelSubscription: mutation.mutate,
        isLoading: mutation.isPending,
        isSuccess: mutation.isSuccess,
        isError: mutation.isError,
        error: mutation.error,
    };
};
