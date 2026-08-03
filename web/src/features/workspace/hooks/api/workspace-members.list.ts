import { useQuery, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import type {  WorkspaceMember } from "../../types/workspace.types";
import { WORKSPACE_QUERY_KEYS } from "../../constants/workspace.constants";
import { workspaceService } from "../../services";
import React from "react";
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

    // React.useEffect(() => {
    //     if (listWsMembersQuery.isError && listWsMembersQuery.error) {
    //         toast.error(listWsMembersQuery.error.message || "Failed to fetch workspace member");
    //         queryClient.removeQueries({ queryKey: [WORKSPACE_QUERY_KEYS.MEMBER_LIST] });
    //     }
    // }, [listWsMembersQuery.isError, listWsMembersQuery.error, queryClient]);

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
