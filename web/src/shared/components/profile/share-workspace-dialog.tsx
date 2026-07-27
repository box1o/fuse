import { Copy } from "lucide-react";
import { toast } from "sonner";

import { useWorkspaceStore } from "@/features/workspace";
import { Button, Dialog, Input } from "@/shared/components/ui";
import { ROUTES } from "@/shared/constants";

interface ShareWorkspaceDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

const ShareWorkspaceDialog = ({
  open,
  onOpenChange,
}: ShareWorkspaceDialogProps) => {
  const currentWorkspace = useWorkspaceStore((state) => state.currentWorkspace);

  const workspaceLink = (() => {
    if (!currentWorkspace) {
      return "";
    }

    const url = new URL(ROUTES.WORKSPACE, window.location.origin);
    url.searchParams.set("workspace", currentWorkspace.id);
    return url.toString();
  })();

  const handleCopy = async () => {
    if (!workspaceLink) {
      return;
    }

    try {
      await navigator.clipboard.writeText(workspaceLink);
      toast.success("Workspace link copied");
    } catch {
      toast.error("Could not copy workspace link");
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <Dialog.Content>
        <Dialog.Header>
          <Dialog.Title>Share workspace</Dialog.Title>
          <Dialog.Description>
            Copy a link to open this workspace in Fuse.
          </Dialog.Description>
        </Dialog.Header>

        <div className="space-y-4">
          <div className="space-y-2">
            <label className="text-sm font-medium" htmlFor="workspace-access">
              Access
            </label>
            <select
              id="workspace-access"
              className="border-input bg-background h-9 w-full rounded-md border px-3 text-sm"
              defaultValue="link"
            >
              <option value="link">Anyone with the link</option>
            </select>
          </div>

          <div className="space-y-2">
            <label className="text-sm font-medium" htmlFor="workspace-link">
              Workspace link
            </label>
            <div className="flex gap-2">
              <Input
                id="workspace-link"
                value={workspaceLink}
                placeholder="Select a workspace first"
                readOnly
              />
              <Button onClick={handleCopy} disabled={!workspaceLink}>
                <Copy />
                Copy
              </Button>
            </div>
          </div>
        </div>
      </Dialog.Content>
    </Dialog>
  );
};

export { ShareWorkspaceDialog };
