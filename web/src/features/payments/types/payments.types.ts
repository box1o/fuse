type ResourceType = "cpu" | "gpu" | "npu";
type SubscriptionPlanId = "free" | "pro";
type SubscriptionStatus =  "free" | "active" | "canceled" | "past_due" | "unpaid" | "incomplete" | "incomplete_expired" | "trialing";

interface CheckoutSessionResponse {
    session_id: string;
    url: string;
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

export type {
    CancelSubscriptionRequest,
    CheckoutSessionResponse,
    ProjectUsageRequest,
    WebhookRequest,
    SubscriptionPlan,
    CreditBalance,
    ResourceType,
    SubscriptionStatus,
    SubscriptionPlanId,
};
