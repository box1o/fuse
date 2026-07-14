import { useMutation } from "@tanstack/react-query";
import { toast } from "sonner";
import { paymentsService } from "../../services";
import { usePaymentsStore } from "../../store";
import type { WebhookRequest } from "../../types";
import { PAYMENTS_QUERY_KEYS } from "../../constants";

export const useSendWebhook = () => {
    const { setWebhookResult, setError } = usePaymentsStore();

    const mutation = useMutation<void, Error, WebhookRequest>({
        mutationKey: [PAYMENTS_QUERY_KEYS.WEBHOOK],
        mutationFn: async (request) => {
            if (!request?.payload) throw new Error("Missing webhook payload");
            if (!request?.signature) throw new Error("Missing webhook signature");
            const response = await paymentsService.sendWebhook(request);
            if (!response.success) {
                throw new Error(response.error || "Failed to send webhook");
            }
        },
        onSuccess: () => {
            setWebhookResult("Webhook accepted");
            setError(null);
            toast.success("Webhook accepted");
        },
        onError: (err) => {
            setWebhookResult("Webhook rejected");
            setError(err.message);
            toast.error(err.message || "Failed to send webhook");
        },
    });

    return {
        sendWebhook: mutation.mutate,
        isLoading: mutation.isPending,
        isSuccess: mutation.isSuccess,
        isError: mutation.isError,
        error: mutation.error,
    };
};
