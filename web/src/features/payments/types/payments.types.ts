interface CheckoutSessionResponse {
    session_id: string;
    url: string;
}

interface ProjectUsageRequest {
    workspace_id: string;
    resource_type: "cpu" | "gpu" | "npu";
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

export type {
    CancelSubscriptionRequest,
    CheckoutSessionResponse,
    ProjectUsageRequest,
    WebhookRequest,
};
