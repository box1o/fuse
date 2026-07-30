type ComputeLogLevel = "info" | "error";

interface ComputeLogEntry {
    id: string;
    level: ComputeLogLevel;
    message: string;
    timestamp: string;
}

export type { ComputeLogEntry, ComputeLogLevel };