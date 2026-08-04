import * as React from "react";
import { z } from "zod";

import { Button, Dialog, Input } from "@/shared/components";
import { cn } from "@/shared/utils";

import { useAddWorkspaceMember } from "../../hooks";
import { useWorkspaceStore } from "../../store";
import { MemberRole, type MemberRoleSection } from "./member-role";

interface AddWorkspaceMemberModalProps {
  className?: string;
}

const AddWorkspaceMemberModal = ({ className }: AddWorkspaceMemberModalProps) => {
  const { addMember, isLoading } = useAddWorkspaceMember();

  const currentWorkspace = useWorkspaceStore(
    (state) => state.currentWorkspace,
  );

  const [email, setEmail] = React.useState("");
  const [memberRole, setMemberRole] =
    React.useState<MemberRoleSection>("member");
  const [validationError, setValidationError] = React.useState("");

  const validateAndCreate = () => {
    if (!currentWorkspace) {
      setValidationError("No workspace selected");
      return;
    }

    const result = z.email().safeParse(email.trim());

    if (!result.success) {
      setValidationError(
        result.error.issues[0]?.message ?? "Invalid email address",
      );
      return;
    }

    setValidationError("");

    addMember({
      user_mail: result.data,
      workspace_id: currentWorkspace.id,
      role: memberRole,
    });
  };

  const handleEmailChange = (
    event: React.ChangeEvent<HTMLInputElement>,
  ) => {
    setEmail(event.target.value);

    if (validationError) {
      setValidationError("");
    }
  };

  const handleKeyDown = (
    event: React.KeyboardEvent<HTMLInputElement>,
  ) => {
    if (event.key === "Enter") {
      event.preventDefault();
      validateAndCreate();
    }
  };

  return (
    <Dialog>
      <Dialog.Trigger asChild>
        <Button variant="outline" className={cn(className)}>
          Add Member
        </Button>
      </Dialog.Trigger>

      <Dialog.Content
        overlayClassName="z-[60]"
        className="z-[60] w-[90vw] rounded-2xl p-6 sm:max-w-md"
      >
        <Dialog.Title className="text-lg font-semibold">
          Add Member
        </Dialog.Title>

        <div className="mt-4 space-y-4">
          <div className="space-y-2">
            <label
              htmlFor="member-email"
              className="text-sm font-medium text-muted-foreground"
            >
              User email
            </label>

            <Input
              id="member-email"
              type="email"
              value={email}
              onChange={handleEmailChange}
              onKeyDown={handleKeyDown}
              placeholder="example@gmail.com"
              className="w-full"
              disabled={isLoading}
            />
          </div>

          <div className="space-y-2">
            <label className="text-sm font-medium text-muted-foreground">
              Member role
            </label>

            <MemberRole
              className="w-full"
              memberRole={memberRole}
              setMemberRole={setMemberRole}
            />
          </div>

          {validationError && (
            <p className="text-sm text-red-500">
              {validationError}
            </p>
          )}

          <div className="flex justify-end gap-2 pt-2">
            <Dialog.Close asChild>
              <Button variant="outline" disabled={isLoading}>
                Cancel
              </Button>
            </Dialog.Close>

            <Button
              type="button"
              onClick={validateAndCreate}
              disabled={isLoading || !currentWorkspace}
              variant="outline"
              className="h-8 rounded-md bg-brand/35"
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
