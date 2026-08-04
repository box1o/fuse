import { useRef } from "react";
import {
    Play,
    Square,
    Trash2,
} from "lucide-react";

import { Button } from "@/shared/components";

import { useNodeStore } from "../store";

interface ComputeLogsProps {
    isConnected: boolean;
    isEnabled: boolean;
    onStart: () => void;
    onStop: () => void;
}

const ComputeLogs = ({isConnected, isEnabled, onStart, onStop}: ComputeLogsProps) => {
    const logs = useNodeStore((state) => state.logs);
    const clearLogs = useNodeStore(
        (state) => state.clearLogs,
    );

    const bottomRef = useRef<HTMLDivElement | null>(null);

    return (
        <section className="flex h-full min-h-128 flex-col overflow-hidden rounded-xl bg-card">
            <header className="flex items-center justify-between border-b px-4 py-3">
                <div>
                    <div className="flex items-center gap-2">
                        <h2 className="text-sm font-medium">
                            Compute logs
                        </h2>

                        <span
                            className={
                                isConnected
                                    ? "text-xs text-emerald-500"
                                    : "text-xs text-muted-foreground"
                            }
                        >
                            {isConnected
                                ? "Connected"
                                : isEnabled
                                  ? "Connecting"
                                  : "Stopped"}
                        </span>
                    </div>
                </div>

                <div className="flex items-center gap-1">
                    {isEnabled ? (
                        <Button
                            type="button"
                            variant="ghost"
                            size="icon"
                            onClick={onStop}
                            aria-label="Stop WebSocket"
                            title="Stop WebSocket"
                        >
                            <Square className="size-4" />
                        </Button>
                    ) : (
                        <Button
                            type="button"
                            variant="ghost"
                            size="icon"
                            onClick={onStart}
                            aria-label="Start WebSocket"
                            title="Start WebSocket"
                        >
                            <Play className="size-4" />
                        </Button>
                    )}

                    <Button
                        type="button"
                        variant="ghost"
                        size="icon"
                        onClick={clearLogs}
                        disabled={logs.length === 0}
                        aria-label="Clear logs"
                        title="Clear logs"
                    >
                        <Trash2 className="size-4" />
                    </Button>
                </div>
            </header>

            <div className="min-h-0 flex-1 overflow-y-auto p-4 font-mono text-xs">
                {logs.length === 0 ? (
                    <p className="text-muted-foreground">
                        Waiting for WebSocket activity...
                    </p>
                ) : (
                    <div className="space-y-1">
                        {logs.map((log) => (
                            <div
                                key={log.id}
                                className="flex gap-3"
                            >
                                <span className="shrink-0 text-muted-foreground">
                                    [{log.timestamp}]
                                </span>

                                <span
                                    className={
                                        log.level === "error"
                                            ? "text-destructive"
                                            : "text-foreground"
                                    }
                                >
                                    {log.message}
                                </span>
                            </div>
                        ))}

                        <div ref={bottomRef} />
                    </div>
                )}
            </div>
        </section>
    );
};

export { ComputeLogs };