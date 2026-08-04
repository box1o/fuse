import { useMutation, useQueryClient } from "@tanstack/react-query";
import type {  WorkspaceMember, WorkspaceMemberRequest } from "../../types/workspace.types";
import { WORKSPACE_QUERY_KEYS } from "../../constants/workspace.constants";
import { workspaceService } from "../../services/workspace.service";
import { toast } from "sonner";

export const useAddWorkspaceMember = () => {
    // const { addWorkspace, reset } = useWorkspaceStore();
    const queryClient = useQueryClient();

    const mutation = useMutation<WorkspaceMember, Error, WorkspaceMemberRequest>({
        mutationKey: [WORKSPACE_QUERY_KEYS.MEMBER_ADD],
        mutationFn: async (request) => {
            if (!request) throw new Error("Missing add workspace member request");
            const response = await workspaceService.addWorkspaceMember(request);
            if (!response.success || !response.data) {
                throw new Error(response.error || "Failed to add workspace member");
            }
            return response.data;
        },
        onSuccess: (_, variables) => {
            queryClient.invalidateQueries({ queryKey: [WORKSPACE_QUERY_KEYS.MEMBER_LIST, variables.workspace_id] });
            toast.success("Workspace member added");
        },
        onError: (err) => {
            toast.error(err.message || "Failed to add workspace member");
            toast.error("Please try again with a different email");
        },
    });

    return {
        workspace: mutation.data,
        isLoading: mutation.isPending,
        isSuccess: mutation.isSuccess,
        isError: mutation.isError,
        error: mutation.error,
        addMember: mutation.mutate, 
    }
};
