import React from "react";
import { useSearchParams } from "react-router-dom";
import Workspaces from "./components/workspaces/workspaces";
import type { Workspace } from "./types";
import { useWorkspaceStore } from "./store";
import AddWorkspaceMemberModal from "./components/workspaces/add-workspace-members-modal";
import { WorkspaceSettingsModal } from "./components/modal/workspace-settings-modal";


const Main: React.FC = () => {

    const setCurrentWs = useWorkspaceStore(store => store.setCurrentWorkspace);
    const workspaces = useWorkspaceStore(store => store.workspaces);
    const [searchParams] = useSearchParams();
    const workspaceId = searchParams.get("workspace");

    React.useEffect(() => {
        if (!workspaceId) {
            return;
        }

        const workspace = workspaces.find(item => item.id === workspaceId);
        if (workspace) {
            setCurrentWs(workspace);
        }
    }, [workspaceId, workspaces, setCurrentWs]);

    const handleRowClick = React.useCallback((workspace: Workspace) => {
        setCurrentWs(workspace);

    }, [setCurrentWs]);

    return (
        <div className="w-full h-full">
            <AddWorkspaceMemberModal/>
            <WorkspaceSettingsModal/>
            <Workspaces onRowClick={handleRowClick} />
        </div>


    );
};

export default Main;
