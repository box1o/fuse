import { Cpu, MemoryStick, Microchip } from "lucide-react";

import { useComputing } from "../hooks";
import type { ComputeNode } from "../types/node.types";
import { CreateNodeModal } from "../components/create-node-modal";
import { MachineListSkeleton } from "./machine-list-skeleton";
import { MachineActions } from "./machine-actions";

const MachineCard = ({ node }: { node: ComputeNode }) => (
    <article className="px-1 py-3 border-b last:border-b-0">
        <div className="flex items-center justify-between">
            <h3 className="min-w-0 truncate text-sm font-medium">
                {node.name}
            </h3>

            <MachineActions node={node} />
        </div>

        <div className="mt-2 flex flex-wrap gap-4 text-xs text-muted-foreground">
            <span className="flex items-center gap-1.5">
                <Cpu className="size-3.5" />
                {node.capabilities.cpu_cores} CPU
            </span>

            <span className="flex items-center gap-1.5">
                <MemoryStick className="size-3.5" />
                {node.capabilities.memory_mb} MB
            </span>

            <span className="flex items-center gap-1.5">
                <Microchip className="size-3.5" />
                {node.capabilities.gpu_count} GPU
            </span>
        </div>
    </article>
);

export const MachineList = () => {
    const {
        currentWorkspace,
        nodes,
        isLoadingNodes,
        nodesError,
    } = useComputing();

    if (!currentWorkspace) {
        return (
            <aside>
                <h2 className="font-semibold">Compute nodes</h2>

                <p className="mt-2 text-sm text-muted-foreground">
                    Select a workspace to view and create compute nodes.
                </p>
            </aside>
        );
    }

    if (isLoadingNodes) {
        return <MachineListSkeleton/>;
    }

    return (
        <aside className="rounded-2xl bp-5">
            <div className="flex items-start justify-between gap-3">
                <div className="min-w-0">
                    <h2 className="font-semibold">Compute nodes</h2>
                </div>

                <CreateNodeModal />
            </div>

            {nodesError && (
                <p className="mt-4 text-sm text-destructive">
                    {nodesError.message}
                </p>
            )}

            {!nodesError && nodes.length === 0 && (
                <div className="mt-6 p-6 text-center">
                    <p className="text-sm font-medium">
                        No compute nodes
                    </p>

                    <p className="mt-1 text-xs text-muted-foreground">
                        Create the first node for this workspace.
                    </p>
                </div>
            )}

            {nodes.length > 0 && (
                <div className="mt-4">
                    {nodes.map((node) => (
                        <MachineCard key={node.id} node={node} />
                        
                    ))}
                </div>
            )}
        </aside>
        
    );
};