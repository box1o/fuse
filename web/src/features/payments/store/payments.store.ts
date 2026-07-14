import { create } from "zustand";

interface PaymentsStoreProps {
    error: string | null;
    checkoutUrl: string | null;
    lastWebhookResult: string | null;
    setCheckoutUrl: (url: string | null) => void;
    setWebhookResult: (result: string | null) => void;
    setError: (error: string | null) => void;
    reset: () => void;
}

const usePaymentsStore = create<PaymentsStoreProps>((set) => ({
    error: null,
    checkoutUrl: null,
    lastWebhookResult: null,
    setCheckoutUrl: (checkoutUrl) => set({ checkoutUrl }),
    setWebhookResult: (lastWebhookResult) => set({ lastWebhookResult }),
    setError: (error) => set({ error }),
    reset: () => set({ error: null, checkoutUrl: null, lastWebhookResult: null }),
}));

export default usePaymentsStore;
