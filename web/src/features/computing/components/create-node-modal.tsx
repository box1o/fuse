import z from "zod";
import { useComputing } from "../hooks";
import { Cpu, Plus } from "lucide-react";
import { Button, Dialog, Input } from "@/shared/components";
import React from "react";


const createNodeSchema = z.object({
    name: z
        .string()
        .trim()
        .min(1, "Node name is required")
        .min(3, "Node name must contain at least 3 characters")
        .max(100, "Node name must contain at most 100 characters"),
    cpuCores: z.coerce
        .number()
        .int("CPU cores must be a whole number")
        .min(1, "At least one CPU core is required"),
    memoryMB: z.coerce
        .number()
        .int("Memory must be a whole number")
        .min(128, "Memory must be at least 128 MB"),
    gpuCount: z.coerce
        .number()
        .int("GPU count must be a whole number")
        .min(0, "GPU count cannot be negative"),
});

interface CreateNodeFormState {
    name: string;
    cpuCores: string;
    memoryMB: string;
    gpuCount: string;
    npuSupported: boolean;
}

const INITIAL_FORM_STATE: CreateNodeFormState = {
    name: "",
    cpuCores: "2",
    memoryMB: "4096",
    gpuCount: "0",
    npuSupported: false,
}

export const CreateNodeModal = () => {
    const {
        currentWorkspace,
        createNode,
        isCreatingNode,
        canCreateNode,

    } = useComputing();

    const [isOpen, SetIsOpen] = React.useState(false);
    const [form, setForm] = React.useState<CreateNodeFormState>(INITIAL_FORM_STATE);
    const [ValidationError, setValidationError] = React.useState<string | null> (null);

    const updateField = (field: keyof CreateNodeFormState, value: string | boolean) => {
        setForm((current) => ({
            ...current,
            [field]: value,
        }));

        setValidationError(null);
    }

    const resetForm = () => {
        setForm(INITIAL_FORM_STATE);
        setValidationError(null);
    };

    const handleOpenChange = (open: boolean) => {
        SetIsOpen(open);

        if (!open) {
            resetForm();
        }
    };

    const handleSubmit = async () => {
        const result = createNodeSchema.safeParse({
            name: form.name,
            cpuCores: form.cpuCores,
            memoryMB: form.memoryMB,
            gpuCount: form.gpuCount,
        });

        if (!result.success) {
            setValidationError(result.error.issues[0]?.message ?? "Invalid node");
            return;
        }

        try {
            await createNode({
                name: result.data.name,
                capabilities: {
                    cpu_cores: result.data.cpuCores,
                    memory_mb: result.data.memoryMB,
                    gpu_count: result.data.gpuCount,
                    npu_supported: form.npuSupported,
                },
            });

            SetIsOpen(false);
            resetForm();

        } catch (error) {
            setValidationError(
                error instanceof Error
                ? error.message
                : "Failed to create compute node"
            );
        }
    };

    const handleKeyDown = (event: React.KeyboardEvent<HTMLInputElement>) => {
        if (event.key === "Enter" && !isCreatingNode) {
            void handleSubmit();
        }
    }

    return(
        <Dialog open={isOpen} onOpenChange={handleOpenChange}>
            <Dialog.Trigger asChild>
                <Button
                    type="button"
                    disabled={!canCreateNode}
                    className="gap-2"
                >
                    <Plus className="size-4" />
                    Create node
                </Button>
            </Dialog.Trigger>
            <Dialog.Content className="rounded-2xl sm:max-w-lg">
                <Dialog.Title className="flex items-center gap-2">
                    Create compute node
                </Dialog.Title>
                <div className="grid grid-row gap-2">
                    <div className="space-y-2">                    
                        <label
                            htmlFor="compute-node-name"
                            className="text-sm font-medium"
                        >
                            Node name
                        </label>

                        <Input
                            id="compute-node-name"
                            value={form.name}
                            onChange={(event) =>
                                updateField("name", event.target.value)
                            }
                            onKeyDown={handleKeyDown}
                            placeholder="Development machine"
                            disabled={isCreatingNode}
                            autoFocus
                        />
                    </div>

                    <div className="grid grid-row gap-2">
                        <div className="space-y-2">
                            <label
                                htmlFor="compute-node-cpu"
                                className="text-sm font-medium"
                            >
                                CPU cores
                            </label>

                            <Input
                                id="compute-node-cpu"
                                type="number"
                                min={1}
                                step={1}
                                value={form.cpuCores}
                                onChange={(event) =>
                                    updateField(
                                        "cpuCores",
                                        event.target.value,
                                    )
                                }
                                onKeyDown={handleKeyDown}
                                disabled={isCreatingNode}
                            />
                        </div>

                        <div className="space-y-2">
                            <label
                                htmlFor="compute-node-memory"
                                className="text-sm font-medium"
                            >
                                Memory, MB
                            </label>

                            <Input
                                id="compute-node-memory"
                                type="number"
                                min={128}
                                step={128}
                                value={form.memoryMB}
                                onChange={(event) =>
                                    updateField(
                                        "memoryMB",
                                        event.target.value,
                                    )
                                }
                                onKeyDown={handleKeyDown}
                                disabled={isCreatingNode}
                            />
                        </div>

                        <div className="space-y-2">
                            <label
                                htmlFor="compute-node-gpu"
                                className="text-sm font-medium"
                            >
                                GPU count
                            </label>

                            <Input
                                id="compute-node-gpu"
                                type="number"
                                min={0}
                                step={1}
                                value={form.gpuCount}
                                onChange={(event) =>
                                    updateField(
                                        "gpuCount",
                                        event.target.value,
                                    )
                                }
                                onKeyDown={handleKeyDown}
                                disabled={isCreatingNode}
                            />
                        </div>

                        <label className="flex cursor-pointer items-center gap-6 self-end rounded-md border px-3 py-2">
                            <input
                                type="checkbox"
                                checked={form.npuSupported}
                                onChange={(event) =>
                                    updateField(
                                        "npuSupported",
                                        event.target.checked,
                                    )
                                }
                                disabled={isCreatingNode}
                                className="size-4"
                            />

                            <span className="text-sm font-medium">
                                NPU supported
                            </span>
                        </label>
                    </div>

                    {ValidationError && (
                        <p className="text-sm text-destructive">
                            {ValidationError}
                        </p>
                    )}

                    <div className="flex justify-end gap-2 pt-2">
                        <Dialog.Close asChild>
                            <Button
                                type="button"
                                variant="outline"
                                disabled={isCreatingNode}
                            >
                                Cancel
                            </Button>
                        </Dialog.Close>

                        <Button
                            type="button"
                            onClick={() => void handleSubmit()}
                            disabled={!canCreateNode || isCreatingNode}
                        >
                            {isCreatingNode ? "Creating..." : "Create node"}
                        </Button>
                    </div>
                </div>
            </Dialog.Content>
        </Dialog>
    )



    

}