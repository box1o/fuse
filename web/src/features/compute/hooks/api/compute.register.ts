import { useMutation, useQueryClient } from "@tanstack/react-query";
import { COMPUTE_QUERY_KEYS } from "../../constants";
import { computeService } from "../../services";
import type { RegisterComputeNodeRequest } from "../../types";

export const useRegisterComputeNode = () => {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async (request: RegisterComputeNodeRequest) => {
            const response = await computeService.register(request);
            if (!response.success || !response.data) {
                throw new Error(response.error ?? "Failed to register compute node");
            }
            return response.data;
        },
        onSuccess: async (node) => {
            queryClient.setQueryData(COMPUTE_QUERY_KEYS.DETAIL(node.id), node);
            await queryClient.invalidateQueries({ queryKey: COMPUTE_QUERY_KEYS.ALL });
        },
    });
};
