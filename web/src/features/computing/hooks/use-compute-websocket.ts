import { useCallback, useEffect, useRef, useState } from "react";

import { env } from "@/shared/constants/env.constants";

import { useNodeStore } from "../store";
import type {
    ComputeLogLevel,
    ComputeNode,
} from "../types";

const RECONNECT_DELAY_MILLISECONDS = 2_000;

interface UseComputeWebSocketResult {
    isConnected: boolean;
    isEnabled: boolean;
    start: () => void;
    stop: () => void;
}

const createWebSocketUrl = (workspaceId: string): string => {
    const url = new URL(env.API_URL);

    url.protocol = url.protocol === "https:" ? "wss:" : "ws:";
    url.pathname = "/compute/node/stream";
    url.searchParams.set("workspace_id", workspaceId);

    return url.toString();
};

export const useComputeWebSocket = (workspaceId: string | undefined): UseComputeWebSocketResult => {
    const setNodes = useNodeStore((state) => state.setNodes);
    const setError = useNodeStore((state) => state.setError);
    const addLog = useNodeStore((state) => state.addLog);

    const socketRef = useRef<WebSocket | null>(null);
    const reconnectTimeoutRef =useRef<ReturnType<typeof setTimeout> | null>(null);

    const [isConnected, setIsConnected] = useState(false);
    const [isEnabled, setIsEnabled] = useState(false);

    const addWebSocketLog = useCallback(
        (
            level: ComputeLogLevel,
            message: string,
        ): void => {
            addLog({
                id: crypto.randomUUID(),
                level,
                message,
                timestamp: new Date().toLocaleTimeString(),
            });
        },
        [addLog],
    );

    const clearReconnectTimeout = useCallback((): void => {
        if (reconnectTimeoutRef.current === null) {
            return;
        }

        clearTimeout(reconnectTimeoutRef.current);
        reconnectTimeoutRef.current = null;
    }, []);

    const closeSocket = useCallback(
        (
            code = 1000,
            reason = "WebSocket stopped",
        ): void => {
            clearReconnectTimeout();

            const socket = socketRef.current;
            socketRef.current = null;

            if (
                socket?.readyState === WebSocket.OPEN ||
                socket?.readyState === WebSocket.CONNECTING
            ) {
                socket.close(code, reason);
            }

            setIsConnected(false);
        },
        [clearReconnectTimeout],
    );

    const stop = useCallback((): void => {
        setIsEnabled(false);
        closeSocket(1000, "Stopped by user");
        addWebSocketLog("info", "Compute WebSocket stopped");
    }, [
        addWebSocketLog,
        closeSocket,
    ]);

    const start = useCallback((): void => {
        if (!workspaceId) {
            addWebSocketLog(
                "error",
                "Select a workspace before starting the WebSocket",
            );
            return;
        }

        setIsEnabled(true);
    }, [
        addWebSocketLog,
        workspaceId,
    ]);

    useEffect(() => {
        if (!workspaceId || !isEnabled) {
            closeSocket();
            return;
        }

        let isDisposed = false;

        const connect = (): void => {
            if (
                isDisposed ||
                !isEnabled ||
                socketRef.current !== null
            ) {
                return;
            }

            addWebSocketLog(
                "info",
                "Connecting to compute WebSocket...",
            );

            const socket = new WebSocket(
                createWebSocketUrl(workspaceId),
            );

            socketRef.current = socket;

            socket.addEventListener("open", () => {
                if (isDisposed) {
                    socket.close();
                    return;
                }

                setIsConnected(true);
                setError(null);

                addWebSocketLog(
                    "info",
                    "Compute WebSocket connected",
                );
            });

            socket.addEventListener("message", (event) => {
                try {
                    const parsedData: unknown = JSON.parse(
                        String(event.data),
                    );

                    if (!Array.isArray(parsedData)) {
                        throw new Error(
                            "WebSocket payload is not an array",
                        );
                    }

                    const nodes = parsedData as ComputeNode[];

                    setNodes(nodes);
                    setError(null);

                    addWebSocketLog(
                        "info",
                        `Received ${nodes.length} compute ${
                            nodes.length === 1 ? "node" : "nodes"
                        }`,
                    );
                } catch (error) {
                    const message =
                        error instanceof Error
                            ? error.message
                            : "Received invalid WebSocket data";

                    setError(message);
                    addWebSocketLog("error", message);
                }
            });

            socket.addEventListener("error", () => {
                if (isDisposed || !isEnabled) {
                    return;
                }

                const message = "Compute WebSocket connection error";

                setError(message);
                addWebSocketLog("error", message);
            });

            socket.addEventListener("close", (event) => {
                if (socketRef.current === socket) {
                    socketRef.current = null;
                }

                setIsConnected(false);

                if (isDisposed || !isEnabled) {
                    return;
                }

                addWebSocketLog(
                    event.wasClean ? "info" : "error",
                    event.reason
                        ? `Compute WebSocket closed: ${event.reason}`
                        : `Compute WebSocket closed with code ${event.code}`,
                );

                addWebSocketLog(
                    "info",
                    "Reconnecting to compute WebSocket...",
                );

                reconnectTimeoutRef.current = setTimeout(
                    connect,
                    RECONNECT_DELAY_MILLISECONDS,
                );
            });
        };

        connect();

        return () => {
            isDisposed = true;
            closeSocket(1000, "WebSocket effect stopped");
        };
    }, [
        workspaceId,
        isEnabled,
        addWebSocketLog,
        closeSocket,
        setError,
        setNodes,
    ]);

    return {
        isConnected,
        isEnabled,
        start,
        stop,
    };
};