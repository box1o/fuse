import { Button } from "@/shared/components";
import { useAlert } from "@/shared/hooks";
import { Trash2 } from "lucide-react";
import { useDeleteComputeNode, useUpdateComputeNode } from "../../hooks";
import type { ComputeNode } from "../../types";

interface NodeActionsProps {
    node: ComputeNode;
}

export const NodeActions = ({ node }: NodeActionsProps) => {
    const alert = useAlert();
    const updateNode = useUpdateComputeNode();
    const deleteNode = useDeleteComputeNode();
    const disabled = node.status === "disabled";
    const busy = updateNode.isPending || deleteNode.isPending;

    const toggleDisabled = () => {
        updateNode.mutate({
            nodeId: node.id,
            request: { disabled: !disabled },
        });
    };

    const confirmDelete = () => {
        alert.custom({
            title: "Delete compute node",
            message: `Delete ${node.name}? The CLI will need to register this machine again.`,
            confirmText: "Delete",
            cancelText: "Cancel",
            type: "warning",
            showCancel: true,
            onConfirm: async () => {
                await deleteNode.mutateAsync(node.id);
            },
        });
    };

    return (
        <div className="flex justify-end gap-2" onClick={(event) => event.stopPropagation()}>
            <Button disabled={busy} onClick={toggleDisabled} size="sm" variant="ghost">
                {disabled ? "Enable" : "Disable"}
            </Button>
            <Button
                aria-label={`Delete ${node.name}`}
                disabled={busy}
                onClick={confirmDelete}
                size="icon"
                title="Delete"
                className="text-destructive hover:text-destructive"
                variant="ghost"
            >
                <Trash2 />
            </Button>
        </div>
    );
};
