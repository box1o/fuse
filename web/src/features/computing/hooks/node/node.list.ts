import { useQuery } from "@tanstack/react-query";

import type { ComputeNode } from "../../types/node.types";
import { COMPUTE_QUERY_KEYS } from "../../constants/computing.constants";
import { computeService } from "../../services/computing.service";
import { useNodeStore } from "../../store";

export const useNodes = () => {
    const { setNodes, setError } = useNodeStore();

    const query = useQuery<ComputeNode[], Error>({
        queryKey: [COMPUTE_QUERY_KEYS.LIST],
        queryFn: async () => {
            const response = await computeService.list();

            if (!response.success || !response.data) {
                throw new Error(response.error || "Failed to list compute nodes");
            }

            setNodes(response.data);
            setError(null);

            return response.data;
        },
    });

    return {
        nodes: query.data ?? [],
        isLoading: query.isPending,
        isFetching: query.isFetching,
        isSuccess: query.isSuccess,
        isError: query.isError,
        error: query.error,
        refetch: query.refetch,
    };
};
