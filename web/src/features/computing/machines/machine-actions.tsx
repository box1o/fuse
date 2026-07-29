import * as React from "react";
import { MoreHorizontal, Pencil, Trash2 } from "lucide-react";
import { z } from "zod";

import {
    Button,
    Dialog,
    DropdownMenu,
    Input,
} from "@/shared/components";

import { useDeleteNode, useRenameNode } from "../hooks";
import type { ComputeNode } from "../types/node.types";

const renameNodeSchema = z.object({
    name: z
        .string()
        .trim()
        .min(1, "Node name is required")
        .min(3, "Node name must contain at least 3 characters")
        .max(100, "Node name must contain at most 100 characters"),
});

interface MachineActionsProps {
    node: ComputeNode;
}

export const MachineActions = ({ node }: MachineActionsProps) => {
    const { renameAsync, isLoading: isRenaming } = useRenameNode();
    const { deleteAsync, isLoading: isDeleting } = useDeleteNode();

    const [isRenameOpen, setIsRenameOpen] = React.useState(false);
    const [isDeleteOpen, setIsDeleteOpen] = React.useState(false);
    const [name, setName] = React.useState(node.name);
    const [validationError, setValidationError] =
        React.useState<string | null>(null);

    const isBusy = isRenaming || isDeleting;

    const handleRenameOpenChange = (open: boolean) => {
        setIsRenameOpen(open);
        setValidationError(null);

        if (open) {
            setName(node.name);
        }
    };

    const handleRename = async () => {
        const result = renameNodeSchema.safeParse({
            name,
        });

        if (!result.success) {
            setValidationError(
                result.error.issues[0]?.message ?? "Invalid node name",
            );
            return;
        }

        if (result.data.name === node.name) {
            setIsRenameOpen(false);
            return;
        }

        try {
            await renameAsync({
                nodeId: node.id,
                request: {
                    name: result.data.name,
                },
            });

            setIsRenameOpen(false);
        } catch {
            // The mutation hook already displays the error toast.
        }
    };

    const handleDelete = async () => {
        try {
            await deleteAsync(node.id);
            setIsDeleteOpen(false);
        } catch {
            // The mutation hook already displays the error toast.
        }
    };

    const handleRenameKeyDown = (
        event: React.KeyboardEvent<HTMLInputElement>,
    ) => {
        if (event.key === "Enter" && !isRenaming) {
            event.preventDefault();
            void handleRename();
        }
    };

    return (
        <>
            <DropdownMenu>
                <DropdownMenu.Trigger asChild>
                    <Button
                        type="button"
                        variant="ghost"
                        size="icon"
                        disabled={isBusy}
                        aria-label={`Actions for ${node.name}`}
                    >
                        <MoreHorizontal className="size-4" />
                    </Button>
                </DropdownMenu.Trigger>

                <DropdownMenu.Content align="end">
                    <DropdownMenu.Item
                        onSelect={() => setIsRenameOpen(true)}
                    >
                        <Pencil className="mr-2 size-4" />
                        Rename
                    </DropdownMenu.Item>

                    <DropdownMenu.Item
                        variant="destructive"
                        onSelect={() => setIsDeleteOpen(true)}
                    >
                        <Trash2 className="mr-2 size-4" />
                        Delete
                    </DropdownMenu.Item>
                </DropdownMenu.Content>
            </DropdownMenu>

            <Dialog
                open={isRenameOpen}
                onOpenChange={handleRenameOpenChange}
            >
                <Dialog.Content className="rounded-2xl sm:max-w-md">
                    <Dialog.Title>Rename compute node</Dialog.Title>

                    <Dialog.Description>
                        Change the display name for this compute node.
                    </Dialog.Description>

                    <div className="mt-5 space-y-2">
                        <label
                            htmlFor={`rename-node-${node.id}`}
                            className="text-sm font-medium"
                        >
                            Node name
                        </label>

                        <Input
                            id={`rename-node-${node.id}`}
                            value={name}
                            onChange={(event) => {
                                setName(event.target.value);
                                setValidationError(null);
                            }}
                            onKeyDown={handleRenameKeyDown}
                            disabled={isRenaming}
                            autoFocus
                        />

                        {validationError && (
                            <p className="text-sm text-destructive">
                                {validationError}
                            </p>
                        )}
                    </div>

                    <div className="mt-6 flex justify-end gap-2">
                        <Button
                            type="button"
                            variant="outline"
                            onClick={() => setIsRenameOpen(false)}
                            disabled={isRenaming}
                        >
                            Cancel
                        </Button>

                        <Button
                            type="button"
                            onClick={() => void handleRename()}
                            disabled={isRenaming}
                        >
                            {isRenaming ? "Renaming..." : "Rename"}
                        </Button>
                    </div>
                </Dialog.Content>
            </Dialog>

            <Dialog open={isDeleteOpen} onOpenChange={setIsDeleteOpen}>
                <Dialog.Content className="rounded-2xl sm:max-w-md">
                    <Dialog.Title>Delete compute node</Dialog.Title>

                    <Dialog.Description>
                        This will permanently delete{" "}
                        <span className="font-medium text-foreground">
                            {node.name}
                        </span>
                        . This action cannot be undone.
                    </Dialog.Description>

                    <div className="mt-6 flex justify-end gap-2">
                        <Button
                            type="button"
                            variant="outline"
                            onClick={() => setIsDeleteOpen(false)}
                            disabled={isDeleting}
                        >
                            Cancel
                        </Button>

                        <Button
                            type="button"
                            variant="destructive"
                            onClick={() => void handleDelete()}
                            disabled={isDeleting}
                        >
                            {isDeleting ? "Deleting..." : "Delete node"}
                        </Button>
                    </div>
                </Dialog.Content>
            </Dialog>
        </>
    );
};