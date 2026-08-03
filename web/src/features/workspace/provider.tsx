import type { PropsWithChildren } from "react";
import { WorkspaceSettingsModal } from "./components/modal/workspace-settings-modal";
import { useWorkspaceStore } from "./store";

const WorkspaceSettingsModalProvider = ({ children }: PropsWithChildren) => {

  const { isworkspaceSettingsOpen, setWorkspaceSettingsOpen } = useWorkspaceStore();
  return (
    <>
      {children}
      <WorkspaceSettingsModal isOpen={isworkspaceSettingsOpen} onOpenChange={setWorkspaceSettingsOpen} />
    </>
  );
};

export { WorkspaceSettingsModalProvider };
