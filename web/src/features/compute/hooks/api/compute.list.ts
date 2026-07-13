import { useQuery } from "@tanstack/react-query";
import { COMPUTE_QUERY_KEYS } from "../../constants";
import { computeService } from "../../services";

export const useListComputeNodes = () =>
    useQuery({
        queryKey: COMPUTE_QUERY_KEYS.LIST,
        queryFn: async () => {
            const response = await computeService.list();
            if (!response.success || !response.data) {
                throw new Error(response.error ?? "Failed to load compute nodes");
            }
            return response.data;
        },
        refetchInterval: 30_000,
    });
