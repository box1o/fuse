import { create } from "zustand";

interface PaymentsStore {
    checkoutUrl: string | null;
    error: string | null;
    setCheckoutUrl: (checkoutUrl: string | null) => void;
    setError: (error: string | null) => void;
    reset: () => void;
}

const usePaymentsStore = create<PaymentsStore>((set) => ({
    checkoutUrl: null,
    error: null,

    setCheckoutUrl: (checkoutUrl) => {
        set({ checkoutUrl });
    },

    setError: (error) => {
        set({ error });
    },

    reset: () => {
        set({
            checkoutUrl: null,
            error: null,
        });
    },
}));

export { usePaymentsStore };