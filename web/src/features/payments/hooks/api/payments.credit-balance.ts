import { useQuery } from "@tanstack/react-query";

import { PAYMENTS_QUERY_KEYS } from "../../constants";
import { paymentsService } from "../../services";

const useCreditBalance = () => {
    const query = useQuery<number, Error>({
        queryKey: [
            PAYMENTS_QUERY_KEYS.CREDIT_BALANCE,
        ],
        queryFn: async () => {
            const response =
                await paymentsService.getCreditBalance();

            if (!response.success || !response.data) {
                throw new Error(
                    response.error ??
                        "Failed to load credit balance",
                );
            }

            return response.data.balance;
        },
    });

    return {
        balance: query.data ?? 0,
        error: query.error,
        isLoading: query.isLoading,
        isError: query.isError,
        refetch: query.refetch,
    };
};

export { useCreditBalance };