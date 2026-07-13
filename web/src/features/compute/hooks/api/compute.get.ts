import { useQuery } from "@tanstack/react-query";
import { COMPUTE_QUERY_KEYS } from "../../constants";
import { computeService } from "../../services";

export const useComputeNode = (nodeId: string | undefined) =>
    useQuery({
        queryKey: COMPUTE_QUERY_KEYS.DETAIL(nodeId ?? ""),
        queryFn: async () => {
            if (!nodeId) throw new Error("Compute node ID is required");

            const response = await computeService.get(nodeId);
            if (!response.success || !response.data) {
                throw new Error(response.error ?? "Failed to load compute node");
            }
            return response.data;
        },
        enabled: Boolean(nodeId),
    });
