const PAYMENTS_QUERY_KEYS = {
    CHECKOUT: "payments-checkout",
    CREDIT_PACKS: "credit-packs",
    CREDIT_BALANCE: "credit-balance",
} as const;

const PAYMENTS_ROUTES = {
    CHECKOUT: "/payments/checkout",
    CREDIT_PACKS: "/credit-packs",
    CREDIT_BALANCE: "/credits/balance",
    USAGE: "/payments/usage",
    SUBSCRIPTION: "/payments/subscription",
    WEBHOOK: "/payments/webhook",
} as const;

export { PAYMENTS_QUERY_KEYS, PAYMENTS_ROUTES };