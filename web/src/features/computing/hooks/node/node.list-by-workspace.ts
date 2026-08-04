import { useQuery } from "@tanstack/react-query";

import type { ComputeNode } from "../../types/node.types";
import { COMPUTE_QUERY_KEYS } from "../../constants/computing.constants";
import { computeService } from "../../services/computing.service";
import { useNodeStore } from "../../store";

export const useWorkspaceNodes = (workspaceId: string | undefined) => {
    const { setNodes, setError } = useNodeStore();

    const query = useQuery<ComputeNode[], Error>({
        queryKey: [COMPUTE_QUERY_KEYS.LIST_BY_WORKSPACE, workspaceId],
        enabled: Boolean(workspaceId),
        queryFn: async () => {
            if (!workspaceId) throw new Error("Missing workspace ID");

            const response = await computeService.listByWorkspace(workspaceId);

            if (!response.success || !response.data) {
                throw new Error(response.error || "Failed to list workspace compute nodes");
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
