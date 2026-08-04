import { useMutation, useQueryClient } from "@tanstack/react-query";
import { WORKSPACE_QUERY_KEYS } from "../../constants/workspace.constants";
import { workspaceService } from "../../services/workspace.service";
import { toast } from "sonner";

interface DeleteWorkspaceMemberRequest {
    workspaceId: string;
    memberId: string;
}

export const useDeleteWorkspaceMember = () => {
    const queryClient = useQueryClient();

    const mutation = useMutation<void, Error, DeleteWorkspaceMemberRequest>({
        mutationKey: [WORKSPACE_QUERY_KEYS.MEMBER_DELETE],
        mutationFn: async (request) => {
            if (!request) throw new Error("Missing delete workspace member request");
            const response = await workspaceService.deleteMember(request.workspaceId, request.memberId);
            if (!response.success) {
                throw new Error(response.error || "Failed to delete workspace member");
            }
        },
        onSuccess: (_, variables) => {
            queryClient.invalidateQueries({ queryKey: [WORKSPACE_QUERY_KEYS.MEMBER_LIST, variables.workspaceId] });
            toast.success("Workspace member deleted");
        },
        onError: (err) => {
            toast.error(err.message || "Failed to delete workspace member");
        },
    });

    return {
        isLoading: mutation.isPending,
        isSuccess: mutation.isSuccess,
        isError: mutation.isError,
        error: mutation.error,
        delete: mutation.mutate,
    };
};
