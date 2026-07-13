import { useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { COMPUTE_QUERY_KEYS } from "../../constants";
import { computeService } from "../../services";
import type { UpdateComputeNodeRequest } from "../../types";

interface UpdateComputeNodeVariables {
    nodeId: string;
    request: UpdateComputeNodeRequest;
}

export const useUpdateComputeNode = () => {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async ({ nodeId, request }: UpdateComputeNodeVariables) => {
            const response = await computeService.update(nodeId, request);
            if (!response.success || !response.data) {
                throw new Error(response.error ?? "Failed to update compute node");
            }
            return response.data;
        },
        onSuccess: async (node) => {
            queryClient.setQueryData(COMPUTE_QUERY_KEYS.DETAIL(node.id), node);
            await queryClient.invalidateQueries({ queryKey: COMPUTE_QUERY_KEYS.LIST });
            toast.success("Compute node updated");
        },
        onError: (error) => toast.error(error.message),
    });
};
