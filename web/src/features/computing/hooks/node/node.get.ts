import { useQuery } from "@tanstack/react-query";

import type { ComputeNode } from "../../types/node.types";
import { COMPUTE_QUERY_KEYS } from "../../constants/computing.constants";
import { computeService } from "../../services/computing.service";
import { useNodeStore } from "../../store";

export const useNode = (nodeId: string | undefined) => {
    const { setCurrentNode, setError } = useNodeStore();

    const query = useQuery<ComputeNode, Error>({
        queryKey: [COMPUTE_QUERY_KEYS.GET, nodeId],
        enabled: Boolean(nodeId),
        queryFn: async () => {
            if (!nodeId) throw new Error("Missing compute node ID");

            const response = await computeService.getById(nodeId);

            if (!response.success || !response.data) {
                throw new Error(response.error || "Failed to get compute node");
            }

            setCurrentNode(response.data);
            setError(null);

            return response.data;
        },
    });

    return {
        node: query.data,
        isLoading: query.isPending,
        isFetching: query.isFetching,
        isSuccess: query.isSuccess,
        isError: query.isError,
        error: query.error,
        refetch: query.refetch,
    };
};
