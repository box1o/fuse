import type { SubscriptionPlan} from "../types";

const SUBSCRIPTION_PLANS: SubscriptionPlan[] = [
    {
        id: "free",
        name: "Free",
        description: "For occasional users who want to explore the platform.",
        priceMonthlyCents: 0,
        includedCredits: 240,
        resetInterval: "monthly",
        features: [
            "240 compute credits per month",
            "CPU compute only",
            "Basic workspace management",
            "Community support",
        ],
    },
    {
        id: "pro",
        name: "Pro",
        description: "For professionals and teams who need more resources.",
        priceMonthlyCents: 3000,
        includedCredits: 2400,
        resetInterval: "monthly",
        features: [
            "2400 compute credits per month",
            "Access to CPU, GPU, and NPU compute",
            "Advanced workspace management",
            "Priority support",
        ],
        recommended: true,
    },
];

const CREDIT_COSTS = {
    cpu: 1,
    npu: 2,
    gpu: 4,
} as const;

export {
    SUBSCRIPTION_PLANS,
    CREDIT_COSTS,
};