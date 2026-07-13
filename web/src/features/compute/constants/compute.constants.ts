const COMPUTE_QUERY_KEYS = {
    ALL: ["compute-nodes"] as const,
    LIST: ["compute-nodes", "list"] as const,
    DETAIL: (nodeId: string) => ["compute-nodes", "detail", nodeId] as const,
} as const;

export { COMPUTE_QUERY_KEYS };
