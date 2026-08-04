export interface HttpError {
    status: number;
    message: string;
    detail?: string;
    details?: string;
    error?: string;
    timestamp?: string;
}
