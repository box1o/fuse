import { useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";

import { COMPUTE_QUERY_KEYS } from "../../constants/computing.constants";
import { computeService } from "../../services/computing.service";
import { useNodeStore } from "../../store";

export const useDeleteNode = () => {
    const { deleteNode, setError } = useNodeStore();
    const queryClient = useQueryClient();

    const mutation = useMutation<void, Error, string>({
        mutationKey: [COMPUTE_QUERY_KEYS.DELETE],
        mutationFn: async (nodeId) => {
            if (!nodeId) throw new Error("Missing compute node ID");

            const response = await computeService.delete(nodeId);

            if (!response.success) {
                throw new Error(response.error || "Failed to delete compute node");
            }
        },
        onSuccess: (_, nodeId) => {
            deleteNode(nodeId);
            setError(null);

            queryClient.removeQueries({
                queryKey: [COMPUTE_QUERY_KEYS.GET, nodeId],
            });

            queryClient.invalidateQueries({
                queryKey: [COMPUTE_QUERY_KEYS.LIST],
            });

            queryClient.invalidateQueries({
                queryKey: [COMPUTE_QUERY_KEYS.LIST_BY_WORKSPACE],
            });

            toast.success("Compute node deleted");
        },
        onError: (error) => {
            setError(error.message);
            toast.error(error.message || "Failed to delete compute node");
        },
    });

    return {
        isLoading: mutation.isPending,
        isSuccess: mutation.isSuccess,
        isError: mutation.isError,
        error: mutation.error,
        delete: mutation.mutate,
        deleteAsync: mutation.mutateAsync,
    };
};
