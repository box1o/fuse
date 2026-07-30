import { useWorkspaceStore } from "@/features/workspace";

import { Browser } from "./components/browser/browser";
import { ComputeLogs } from "./components/compute-logs";
import { useComputeWebSocket } from "./hooks/use-compute-websocket";
import { MachineList } from "./machines";

const Main = () => {
    const currentWorkspace = useWorkspaceStore(
        (state) => state.currentWorkspace,
    );

    const webSocket = useComputeWebSocket(
        currentWorkspace?.id,
    );

    return (
        <div className="h-full w-full overflow-y-auto">
            <div className="grid min-h-full grid-cols-1 gap-6 p-6 lg:grid-cols-[minmax(0,1fr)_360px]">
                <main className="min-w-0">
                    <Browser />
                </main>
                    <MachineList />

                    <ComputeLogs
                        isConnected={webSocket.isConnected}
                        isEnabled={webSocket.isEnabled}
                        onStart={webSocket.start}
                        onStop={webSocket.stop}
                    />
            </div>
        </div>
    );
};

export default Main;
