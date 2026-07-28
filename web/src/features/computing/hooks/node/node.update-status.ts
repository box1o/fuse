import { useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";

import type { ComputeNode, UpdateNodeStatusRequest } from "../../types/node.types";
import { COMPUTE_QUERY_KEYS } from "../../constants/computing.constants";
import { computeService } from "../../services/computing.service";
import { useNodeStore } from "../../store";

interface UpdateNodeStatusVariables {
    nodeId: string;
    request: UpdateNodeStatusRequest;
}

export const useUpdateNodeStatus = () => {
    const { updateNode, setError } = useNodeStore();
    const queryClient = useQueryClient();

    const mutation = useMutation<ComputeNode, Error, UpdateNodeStatusVariables>({
        mutationKey: [COMPUTE_QUERY_KEYS.UPDATE_STATUS],
        mutationFn: async ({ nodeId, request }) => {
            if (!nodeId) throw new Error("Missing compute node ID");

            const response = await computeService.updateStatus(nodeId, request);

            if (!response.success || !response.data) {
                throw new Error(response.error || "Failed to update compute node status");
            }

            return response.data;
        },
        onSuccess: (data) => {
            updateNode(data);
            setError(null);

            queryClient.setQueryData(
                [COMPUTE_QUERY_KEYS.GET, data.id],
                data,
            );

            queryClient.invalidateQueries({
                queryKey: [COMPUTE_QUERY_KEYS.LIST],
            });

            queryClient.invalidateQueries({
                queryKey: [COMPUTE_QUERY_KEYS.LIST_BY_WORKSPACE, data.workspace_id],
            });

            toast.success("Compute node status updated");
        },
        onError: (error) => {
            setError(error.message);
            toast.error(error.message || "Failed to update compute node status");
        },
    });

    return {
        node: mutation.data,
        isLoading: mutation.isPending,
        isSuccess: mutation.isSuccess,
        isError: mutation.isError,
        error: mutation.error,
        updateStatus: mutation.mutate,
        updateStatusAsync: mutation.mutateAsync,
    };
};
