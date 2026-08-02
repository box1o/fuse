import React from "react";
import type { ColumnDef } from "@tanstack/react-table";
import { Trash } from "lucide-react";

import { Badge, Button } from "@/shared/components";
import { useAlert } from "@/shared/hooks";

import {
    useDeleteWorkspaceMember,
    useListWorkspaceMembers,
} from "../../hooks";
import type { WorkspaceMember } from "../../types/workspace.types";
import { DataTable } from "../tables";
import AddWorkspaceMemberModal from "../workspaces/add-workspace-members-modal";

const WorkspaceMembers = () => {
    const alert = useAlert();

    const { members, isLoading } = useListWorkspaceMembers();
    const { delete: deleteWorkspaceMember } = useDeleteWorkspaceMember();

    const handleDeleteWorkspaceMember = React.useCallback(
        (workspaceId: string, memberId: string) => {
            alert.custom({
                title: "Delete Workspace Member",
                message:
                    "This action cannot be undone. The member will be removed from this workspace.",
                confirmText: "Delete Member",
                cancelText: "Keep Member",
                type: "warning",
                showCancel: true,
                onConfirm: async () => {
                    deleteWorkspaceMember({
                        workspaceId,
                        memberId,
                    });
                },
            });
        },
        [alert, deleteWorkspaceMember],
    );

    const columns = React.useMemo<ColumnDef<WorkspaceMember>[]>(
        () => [
            {
                accessorKey: "name",
                header: "Name",
                cell: ({ row }) => {
                    const name = row.getValue("name") as string;

                    return <div className="font-medium">{name}</div>;
                },
            },
            {
                accessorKey: "mail",
                header: "Mail",
                cell: ({ row }) => {
                    const mail = row.getValue("mail") as string;

                    return <div className="font-medium">{mail}</div>;
                },
            },
            {
                accessorKey: "role",
                header: "Role",
                cell: ({ row }) => {
                    const role = row.getValue("role") as string;
                    const isOwner = role === "owner";

                    return (
                        <Badge
                            variant={isOwner ? "default" : "outline"}
                            className="capitalize"
                        >
                            {role}
                        </Badge>
                    );
                },
            },
            {
                id: "actions",
                header: () => (
                    <div className="text-right">
                        Actions
                    </div>
                ),
                cell: ({ row }) => {
                    const member = row.original;
                    const isOwner = member.role === "owner";

                    return (
                        <div className="flex justify-end">
                            <Button
                                type="button"
                                size="icon"
                                variant="ghost"
                                disabled={isOwner}
                                aria-label={`Delete ${member.name}`}
                                onClick={() =>
                                    handleDeleteWorkspaceMember(
                                        member.workspace_id,
                                        member.id,
                                    )
                                }
                            >
                                <Trash className="size-4 text-destructive" />
                            </Button>
                        </div>
                    );
                },
            },
        ],
        [handleDeleteWorkspaceMember],
    );

    return (
        <div className="h-full w-full space-y-4">
            <div className="flex items-center">
                <span className="pl-2">
                    Workspace Members
                </span>

                <AddWorkspaceMemberModal className="ml-auto mr-2 bg-brand/20" />
            </div>

            <DataTable
                columns={columns}
                data={members}
                loading={isLoading}
            />
        </div>
    );
};

export { WorkspaceMembers };
