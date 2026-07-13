import { Dialog } from "@/shared/components";
import type { ComputeNode } from "../../types";
import { formatBytes, formatDateTime, getNodeAvailability } from "../../utils";

interface NodeDetailsProps {
    node: ComputeNode | null;
    onOpenChange: (open: boolean) => void;
}

interface DetailProps {
    label: string;
    value: string;
}

const Detail = ({ label, value }: DetailProps) => (
    <div>
        <dt className="text-xs text-muted-foreground">{label}</dt>
        <dd className="mt-1 break-words text-sm">{value}</dd>
    </div>
);

export const NodeDetails = ({ node, onOpenChange }: NodeDetailsProps) => {
    if (!node) return null;

    const { capabilities } = node;
    const availability = getNodeAvailability(node);
    const accelerators = capabilities.accelerators.length
        ? capabilities.accelerators
              .map((accelerator) => {
                  const memory = accelerator.memory_bytes ? ` (${formatBytes(accelerator.memory_bytes)})` : "";
                  return `${accelerator.vendor} ${accelerator.model}${memory}`;
              })
              .join(", ")
        : "None";
    const runtime = capabilities.container_runtime.available
        ? `${capabilities.container_runtime.name}${
              capabilities.container_runtime.version ? ` ${capabilities.container_runtime.version}` : ""
          }`
        : "Unavailable";

    return (
        <Dialog onOpenChange={onOpenChange} open>
            <Dialog.Content className="max-h-[calc(100vh-2rem)] overflow-y-auto sm:max-w-2xl">
                <Dialog.Header>
                    <Dialog.Title>{node.name}</Dialog.Title>
                    <Dialog.Description>{node.hostname}</Dialog.Description>
                </Dialog.Header>

                <dl className="grid gap-x-8 gap-y-5 sm:grid-cols-2">
                    <Detail label="Status" value={availability} />
                    <Detail label="Agent version" value={node.agent_version} />
                    <Detail label="Operating system" value={`${capabilities.os.name} ${capabilities.os.version ?? ""}`.trim()} />
                    <Detail label="Architecture" value={capabilities.os.architecture} />
                    <Detail label="CPU" value={capabilities.cpu.model} />
                    <Detail label="CPU cores" value={`${capabilities.cpu.physical_cores} physical / ${capabilities.cpu.logical_cores} logical`} />
                    <Detail label="Memory" value={formatBytes(capabilities.memory.total_bytes)} />
                    <Detail label="Storage" value={formatBytes(capabilities.storage.total_bytes)} />
                    <Detail label="Container runtime" value={runtime} />
                    <Detail label="Accelerators" value={accelerators} />
                    <Detail label="Registered" value={formatDateTime(node.registered_at)} />
                    <Detail label="Last updated" value={formatDateTime(node.updated_at)} />
                </dl>
            </Dialog.Content>
        </Dialog>
    );
};
