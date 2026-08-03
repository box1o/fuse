import { useQuery } from "@tanstack/react-query";
import type {  WorkspaceMember } from "../../types/workspace.types";
import { WORKSPACE_QUERY_KEYS } from "../../constants/workspace.constants";
import { workspaceService } from "../../services";
import { useWorkspaceStore } from "../../store";

export const useListWorkspaceMembers = () => {
    const currentWorkspace = useWorkspaceStore(
        (state) => state.currentWorkspace,
    );

    const workspaceId = currentWorkspace?.id

    const listWsMembersQuery = useQuery<WorkspaceMember[], Error>({
        queryKey: [WORKSPACE_QUERY_KEYS.MEMBER_LIST, workspaceId],
        queryFn: async (): Promise<WorkspaceMember[]> => {
            if (!workspaceId){
                throw new Error("No workspace selected")
            }
            const response = await workspaceService.memberList(workspaceId);
            if (!response.success || !response.data) {
                throw new Error(response.error || "Failed to fetch workspace members");
            }
            return response.data;
        },
        enabled: Boolean(workspaceId),
        retry: false,
    });
    
    return {
        isLoading: listWsMembersQuery.isLoading,
        isFetching: listWsMembersQuery.isFetching,
        isSuccess: listWsMembersQuery.isSuccess,
        refetch: listWsMembersQuery.refetch,
        error: listWsMembersQuery.error,
        //data 
        members: listWsMembersQuery.data ?? [],
    };
}
