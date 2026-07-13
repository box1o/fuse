import { useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { COMPUTE_QUERY_KEYS } from "../../constants";
import { computeService } from "../../services";

export const useDeleteComputeNode = () => {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async (nodeId: string) => {
            const response = await computeService.delete(nodeId);
            if (!response.success) {
                throw new Error(response.error ?? "Failed to delete compute node");
            }
            return nodeId;
        },
        onSuccess: async (nodeId) => {
            queryClient.removeQueries({ queryKey: COMPUTE_QUERY_KEYS.DETAIL(nodeId) });
            await queryClient.invalidateQueries({ queryKey: COMPUTE_QUERY_KEYS.LIST });
            toast.success("Compute node deleted");
        },
        onError: (error) => toast.error(error.message),
    });
};
