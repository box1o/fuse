import { useWorkspaceStore } from "@/features/workspace";

import { useCreateNode } from "./node/node.create";
import { useWorkspaceNodes } from "./node/node.list-by-workspace";
import type { Capabilities } from "../types/node.types";

interface CreateActiveWorkspaceNodeInput {
    name: string;
    capabilities: Capabilities;
}

export const useComputing = () => {
    const currentWorkspace = useWorkspaceStore(
        (state) => state.currentWorkspace,
    );

    const workspaceNodesQuery = useWorkspaceNodes(currentWorkspace?.id);
    const createNodeMutation = useCreateNode();

    const createNode = async ({
        name,
        capabilities,
    }: CreateActiveWorkspaceNodeInput) => {
        if (!currentWorkspace) {
            throw new Error("Select an active workspace before creating a node");
        }

        return createNodeMutation.createAsync({
            workspace_id: currentWorkspace.id,
            name,
            capabilities,
        });
    };

    return {
        currentWorkspace,

        nodes: workspaceNodesQuery.nodes,
        isLoadingNodes: workspaceNodesQuery.isLoading,
        isFetchingNodes: workspaceNodesQuery.isFetching,
        nodesError: workspaceNodesQuery.error,
        refetchNodes: workspaceNodesQuery.refetch,

        createNode,
        isCreatingNode: createNodeMutation.isLoading,
        createNodeError: createNodeMutation.error,

        canCreateNode: currentWorkspace !== null,
    };
};