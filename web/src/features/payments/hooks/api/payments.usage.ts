import { useMutation } from "@tanstack/react-query";
import { toast } from "sonner";
import { paymentsService } from "../../services";
import { usePaymentsStore } from "../../store";
import type { ProjectUsageRequest } from "../../types";
import { PAYMENTS_QUERY_KEYS } from "../../constants";

export const useRecordUsage = () => {
    const { setError } = usePaymentsStore();

    const mutation = useMutation<void, Error, ProjectUsageRequest>({
        mutationKey: [PAYMENTS_QUERY_KEYS.USAGE],
        mutationFn: async (request) => {
            if (!request?.workspace_id) throw new Error("Missing workspace ID");
            const response = await paymentsService.recordUsage(request);
            if (!response.success) {
                throw new Error(response.error || "Failed to record usage");
            }
        },
        onSuccess: () => {
            setError(null);
            toast.success("Usage recorded");
        },
        onError: (err) => {
            setError(err.message);
            toast.error(err.message || "Failed to record usage");
        },
    });

    return {
        recordUsage: mutation.mutate,
        isLoading: mutation.isPending,
        isSuccess: mutation.isSuccess,
        isError: mutation.isError,
        error: mutation.error,
    };
};
