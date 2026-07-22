type ResourceType = "cpu" | "gpu" | "npu";

type SubscriptionPlanId = "free" | "pro";

type SubscriptionStatus =
    | "free"
    | "active"
    | "canceled"
    | "past_due"
    | "unpaid"
    | "incomplete"
    | "incomplete_expired"
    | "trialing";

interface CreateCheckoutRequest {
    creditPackId: string;
    successUrl: string;
    cancelUrl: string;
}

interface CheckoutSessionResponse {
    payment_id: string;
    session_id: string;
    checkout_url: string;
}

interface ProjectUsageRequest {
    workspace_id: string;
    resource_type: ResourceType;
    quantity: number;
    occurred_at: string;
    idempotency_key: string;
}

interface CancelSubscriptionRequest {
    workspace_id: string;
}

interface WebhookRequest {
    payload: string;
    signature: string;
}

interface SubscriptionPlan {
    id: SubscriptionPlanId;
    name: string;
    description: string;
    priceMonthlyCents: number;
    includedCredits: number;
    resetInterval: "monthly" | "yearly";
    features: string[];
    recommended?: boolean;
}

interface CreditBalance {
    planId: SubscriptionPlanId;
    status: SubscriptionStatus;
    usedCredits: number;
    includedCredits: number;
    remainingCredits: number;
    resetAt: string;
    nextResetDate: string;
}

interface CreditPack {
    id: string;
    code: string;
    name: string;
    credits: number;
    price_amount: number;
    currency: string;
}

interface CreditBalanceResponse {
    balance: number;
}

export type {
    CancelSubscriptionRequest,
    CheckoutSessionResponse,
    CreateCheckoutRequest,
    ProjectUsageRequest,
    WebhookRequest,
    SubscriptionPlan,
    CreditBalance,
    CreditBalanceResponse,
    CreditPack,
    ResourceType,
    SubscriptionStatus,
    SubscriptionPlanId,
};