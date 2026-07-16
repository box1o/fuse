import { create } from "zustand";

type SubscriptionPlanId = "free" | "pro";

interface MockBillingState {
    planId: SubscriptionPlanId;
    usedCredits: number;
    includedCredits: number;
}

interface MockBillingStore extends MockBillingState {
    userKey: string | null;
    setUserKey: (userKey: string | null) => void;
    buyPro: () => void;
    spendCredit: (amount?: number) => void;
    resetForUser: () => void;
}

const STORAGE_PREFIX = "fuse.mock-billing";

const getDefaultState = (): MockBillingState => ({
    planId: "free",
    usedCredits: 18,
    includedCredits: 240,
});

const readStoredState = (userKey: string | null): MockBillingState => {
    if (!userKey || typeof window === "undefined") {
        return getDefaultState();
    }

    const raw = window.localStorage.getItem(`${STORAGE_PREFIX}:${userKey}`);
    if (!raw) {
        return getDefaultState();
    }

    try {
        const parsed = JSON.parse(raw) as Partial<MockBillingState>;
        return {
            planId: parsed.planId === "pro" ? "pro" : "free",
            usedCredits: Number.isFinite(parsed.usedCredits) ? Math.max(0, parsed.usedCredits ?? 0) : 18,
            includedCredits: Number.isFinite(parsed.includedCredits) ? Math.max(1, parsed.includedCredits ?? 240) : 240,
        };
    } catch {
        return getDefaultState();
    }
};

const writeStoredState = (userKey: string | null, state: MockBillingState) => {
    if (!userKey || typeof window === "undefined") {
        return;
    }

    window.localStorage.setItem(`${STORAGE_PREFIX}:${userKey}`, JSON.stringify(state));
};

const useMockBillingStore = create<MockBillingStore>((set, get) => ({
    userKey: null,
    ...getDefaultState(),
    setUserKey: (userKey) => {
        const state = readStoredState(userKey);
        set({ userKey, ...state });
    },
    buyPro: () => {
        const nextState = {
            ...get(),
            planId: "pro" as const,
            includedCredits: 2400,
            usedCredits: Math.min(get().usedCredits, 2400),
        };
        set(nextState);
        writeStoredState(get().userKey, nextState);
    },
    spendCredit: (amount = 1) => {
        const current = get();
        const nextState = {
            ...current,
            usedCredits: Math.min(current.usedCredits + amount, current.includedCredits),
        };
        set(nextState);
        writeStoredState(current.userKey, nextState);
    },
    resetForUser: () => {
        const nextState = getDefaultState();
        set(nextState);
        writeStoredState(get().userKey, nextState);
    },
}));

export { useMockBillingStore };
