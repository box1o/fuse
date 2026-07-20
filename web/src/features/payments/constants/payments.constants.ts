const PAYMENTS_QUERY_KEYS = {
    CHECKOUT: "payments-checkout",
} as const;

const PAYMENTS_ROUTES = {
    CHECKOUT: "/payments/checkout",
    USAGE: "/payments/usage",
    SUBSCRIPTION: "/payments/subscription",
    WEBHOOK: "/payments/webhook",
} as const;

export { PAYMENTS_QUERY_KEYS, PAYMENTS_ROUTES };