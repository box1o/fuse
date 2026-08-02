import { useQuery, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import type {  WorkspaceMember } from "../../types/workspace.types";
import { WORKSPACE_QUERY_KEYS } from "../../constants/workspace.constants";
import { workspaceService } from "../../services";
import React from "react";
import { useWorkspaceStore } from "../../store";

export const useListWorkspaceMembers = () => {
    const queryClient = useQueryClient();
    const {currentWorkspace} = useWorkspaceStore();

    const listWsMemberQuery = useQuery<WorkspaceMember[], Error>({
        queryKey: [WORKSPACE_QUERY_KEYS.MEMBER_LIST],
        queryFn: async (): Promise<WorkspaceMember[]> => {
            const response = await workspaceService.memberList(currentWorkspace?.id ?? "");
            if (!response.success || !response.data) {
                throw new Error(response.error || "Failed to fetch workspace members");
            }
            return response.data;
        },
        retry: false
    });

    React.useEffect(() => {
        if (listWsMemberQuery.isError && listWsMemberQuery.error) {
            toast.error(listWsMemberQuery.error.message || "Failed to fetch workspace member");
            queryClient.removeQueries({ queryKey: [WORKSPACE_QUERY_KEYS.MEMBER_LIST] });
        }
    }, [listWsMemberQuery.isError, listWsMemberQuery.error, queryClient]);

    return {
        isLoading: listWsMemberQuery.isLoading,
        isFetching: listWsMemberQuery.isFetching,
        isSuccess: listWsMemberQuery.isSuccess,
        refetch: listWsMemberQuery.refetch,
        error: listWsMemberQuery.error,
        //data 
        members: listWsMemberQuery.data ?? [],
    };
}
