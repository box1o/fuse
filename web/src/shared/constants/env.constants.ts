type NodeEnv = "development" | "production";

interface Environment {
    readonly API_URL: string;
    readonly API_TIMEOUT?: number;
    readonly NODE_ENV: NodeEnv;
    readonly APPLICATION_NAME?: string;
    readonly STRIPE_PUBLISHABLE_KEY: string;
}

const createEnv = (): Environment => {
    const nodeEnv = (import.meta.env.VITE_NODE_ENV as NodeEnv) || "development";
    const apiUrl = import.meta.env.VITE_API_URL || "http://localhost:3000";
    const apiTimeout = import.meta.env.VITE_API_TIMEOUT
        ? parseInt(import.meta.env.VITE_API_TIMEOUT, 10)
        : undefined;
    const applicationName = import.meta.env.VITE_APPLICATION_NAME || "Fuse";
    const stripePublishableKey = import.meta.env.VITE_STRIPE_PUBLISHABLE_KEY || "";

    return {
        API_URL: apiUrl,
        NODE_ENV: nodeEnv,
        API_TIMEOUT: apiTimeout,
        APPLICATION_NAME: applicationName,
        STRIPE_PUBLISHABLE_KEY: stripePublishableKey,
    } as const;
};

export const env = createEnv();
