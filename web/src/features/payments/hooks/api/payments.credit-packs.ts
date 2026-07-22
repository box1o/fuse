import { useQuery } from "@tanstack/react-query";

import { PAYMENTS_QUERY_KEYS } from "../../constants";
import { paymentsService } from "../../services";
import type { CreditPack } from "../../types";

const useCreditPacks = () => {
    const query = useQuery<CreditPack[], Error>({
        queryKey: [PAYMENTS_QUERY_KEYS.CREDIT_PACKS],

        queryFn: async () => {
            const response = await paymentsService.listCreditPacks();

            if (!response.success || !response.data) {
                throw new Error(
                    response.error ?? "Failed to load credit packs",
                );
            }

            return response.data;
        },
    });

    return {
        creditPacks: query.data ?? [],
        error: query.error,
        isLoading: query.isLoading,
        isError: query.isError,
        refetch: query.refetch,
    };
};

export { useCreditPacks };