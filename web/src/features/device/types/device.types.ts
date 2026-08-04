export interface DeviceAuthorizationRequest {
    user_code: string;
    client_name: string;
    status: "pending" | "approved" | "denied";
    expires_at: string;
}

export interface DeviceDecisionResponse {
    status: "approved" | "denied";
}

export type DeviceAuthorizationState =
    | "entering-code"
    | "checking-code"
    | "ready"
    | "submitting"
    | "approved"
    | "denied"
    | "invalid";
