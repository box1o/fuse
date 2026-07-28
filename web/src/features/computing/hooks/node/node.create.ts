import { useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";

import type { ComputeNode, CreateNodeRequest } from "../../types/node.types";
import { COMPUTE_QUERY_KEYS } from "../../constants/computing.constants";
import { computeService } from "../../services/computing.service";
import { useNodeStore } from "../../store";

export const useCreateNode = () => {
    const { addNode, setError } = useNodeStore();
    const queryClient = useQueryClient();

    const mutation = useMutation<ComputeNode, Error, CreateNodeRequest>({
        mutationKey: [COMPUTE_QUERY_KEYS.CREATE],
        mutationFn: async (request) => {
            if (!request) throw new Error("Missing create node request");

            const response = await computeService.create(request);

            if (!response.success || !response.data) {
                throw new Error(response.error || "Failed to create compute node");
            }

            return response.data;
        },
        onSuccess: (data) => {
            addNode(data);
            setError(null);

            queryClient.invalidateQueries({
                queryKey: [COMPUTE_QUERY_KEYS.LIST],
            });

            queryClient.invalidateQueries({
                queryKey: [COMPUTE_QUERY_KEYS.LIST_BY_WORKSPACE, data.workspace_id],
            });

            toast.success("Compute node created");
        },
        onError: (error) => {
            setError(error.message);
            toast.error(error.message || "Failed to create compute node");
        },
    });

    return {
        node: mutation.data,
        isLoading: mutation.isPending,
        isSuccess: mutation.isSuccess,
        isError: mutation.isError,
        error: mutation.error,
        create: mutation.mutate,
        createAsync: mutation.mutateAsync,
    };
};
