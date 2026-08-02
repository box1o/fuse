import React from "react";
import { z } from "zod";

import { Button, Dialog, Input } from "@/shared/components";
import { useAddWorkspaceMember } from "../../hooks";
import { cn } from "@/shared/utils";
import { useWorkspaceStore } from "../../store";

interface AddWorkspaceMemberModalProps {
    className?: string;
}

const AddWorkspaceMemberModal: React.FC<AddWorkspaceMemberModalProps> = ({ className }) => {
    const { create, isLoading } = useAddWorkspaceMember();
    const [name, setName] = React.useState<string>("");
    const [validationError, setValidationError] = React.useState<string>("");
    const {currentWorkspace} = useWorkspaceStore();

    const validateAndCreate = () => {
        try {
            const validatedData = z.email().parse(name.trim());
            setValidationError("");
            create({ user_mail: validatedData , workspace_id: currentWorkspace?.id ?? ""});
        } catch (error) {
            if (error instanceof z.ZodError) {
                const messages = error.issues.map(i => i.message).join(", ");
                setValidationError(messages || "Invalid input");
            }
        }
    };

    const handleNameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setName(e.target.value);
        if (validationError) setValidationError("");
    };

    const handleKeyPress = (e: React.KeyboardEvent) => {
        if (e.key === "Enter") {
            validateAndCreate();
        }
    };

    return (
        <Dialog>
            <Dialog.Trigger asChild>
                <Button
                    variant="outline"
                    className={cn(
                        className
                    )}
                >
                    Add Member
                </Button>
            </Dialog.Trigger>
            <Dialog.Content
                overlayClassName="z-[60]"
                className="z-[60] sm:max-w-md w-[90vw] rounded-2xl p-6"
            >
                <Dialog.Title className="text-lg font-semibold">
                    Add Member
                </Dialog.Title>
                <Dialog.Description className="text-sm text-muted-foreground">
                    Member Mails
                </Dialog.Description>

                <div className="mt-6 space-y-4">
                    <div className="space-y-4">
                        <label htmlFor="workspace-name" className="text-sm font-medium">
                            Mail
                        </label>
                        <Input
                            id="workspace-name"
                            value={name}
                            onChange={handleNameChange}
                            onKeyPress={handleKeyPress}
                            placeholder="Enter user email"
                            className="w-full"
                            disabled={isLoading}
                        />
                    </div>

                    {validationError && (
                        <div className="text-red-500 text-sm p-2 rounded-md ">
                            {validationError}
                        </div>
                    )}

                    <div className="flex justify-end gap-2 pt-4">
                        <Dialog.Close asChild>
                            <Button
                                variant="outline"
                                className="rounded-lg px-4"
                                disabled={isLoading}
                            >
                                Cancel
                            </Button>
                        </Dialog.Close>


                        <Button
                            onClick={validateAndCreate}
                            disabled={isLoading}
                            variant="outline" className="rounded-md h-8 bg-brand/35"
                        >
                            {isLoading ? "Creating..." : "Create"}
                        </Button>
                    </div>
                </div>
            </Dialog.Content>
        </Dialog>
    );
};

export default AddWorkspaceMemberModal;
