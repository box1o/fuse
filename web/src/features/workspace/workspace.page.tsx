import Main from "./main";
import { WorkspaceSettingsModalProvider } from "./provider";

const WorkspacePage = () => {
  return (
    <WorkspaceSettingsModalProvider>
      <Main />
    </WorkspaceSettingsModalProvider>
  );
};

export const Component = WorkspacePage;
