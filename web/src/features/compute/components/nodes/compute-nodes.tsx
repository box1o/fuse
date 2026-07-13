import { Button, Input } from "@/shared/components";
import { useAuthStore } from "@/features/auth/store/auth.store";
import type { ColumnDef } from "@tanstack/react-table";
import { RefreshCw } from "lucide-react";
import React from "react";
import { useListComputeNodes } from "../../hooks";
import type { ComputeNode, ComputeNodeFilter } from "../../types";
import { formatDateTime, getNodeAvailability } from "../../utils";
import { DataTable } from "../tables";
import { NodeActions } from "./node-actions";
import { NodeDetails } from "./node-details";

const FILTERS: Array<{ label: string; value: ComputeNodeFilter }> = [
    { label: "All", value: "all" },
    { label: "Online", value: "online" },
    { label: "Offline", value: "offline" },
];

export const ComputeNodes = () => {
    const [globalFilter, setGlobalFilter] = React.useState("");
    const [availabilityFilter, setAvailabilityFilter] = React.useState<ComputeNodeFilter>("all");
    const [selectedNode, setSelectedNode] = React.useState<ComputeNode | null>(null);
    const user = useAuthStore((state) => state.user);
    const nodesQuery = useListComputeNodes();
    const nodes = React.useMemo(() => nodesQuery.data ?? [], [nodesQuery.data]);

    const filteredNodes = React.useMemo(() => {
        const search = globalFilter.trim().toLowerCase();

        return nodes.filter((node) => {
            const matchesAvailability =
                availabilityFilter === "all" || getNodeAvailability(node) === availabilityFilter;
            if (!matchesAvailability) return false;
            if (!search) return true;

            return [node.name, node.hostname, node.capabilities.cpu.model, node.capabilities.os.name].some((value) =>
                value.toLowerCase().includes(search),
            );
        });
    }, [availabilityFilter, globalFilter, nodes]);

    const columns = React.useMemo<ColumnDef<ComputeNode>[]>(
        () => [
            {
                accessorKey: "name",
                header: "Device",
                cell: ({ row }) => (
                    <div>
                        <div className="font-medium">{row.original.name}</div>
                        <div className="text-xs text-muted-foreground">{row.original.hostname}</div>
                    </div>
                ),
            },
            {
                id: "status",
                accessorFn: getNodeAvailability,
                header: "Status",
                cell: ({ row }) => {
                    const availability = getNodeAvailability(row.original);
                    return (
                        <span className="flex items-center gap-2 text-sm">
                            <span
                                aria-hidden="true"
                                className={`size-1.5 rounded-full ${
                                    availability === "online" ? "bg-emerald-500" : "bg-muted-foreground"
                                }`}
                            />
                            {availability}
                        </span>
                    );
                },
            },
            {
                id: "owner",
                header: "Owner",
                cell: ({ row }) => <span className="text-sm">{user?.name || row.original.owner_id}</span>,
            },
            {
                id: "platform",
                header: "Platform",
                cell: ({ row }) => (
                    <span className="text-sm">
                        {row.original.capabilities.os.name} · {row.original.capabilities.os.architecture}
                    </span>
                ),
            },
            {
                accessorKey: "updated_at",
                header: "Last updated",
                cell: ({ row }) => (
                    <span className="whitespace-nowrap text-sm">{formatDateTime(row.original.updated_at)}</span>
                ),
            },
            {
                id: "actions",
                header: () => <span className="sr-only">Actions</span>,
                cell: ({ row }) => <NodeActions node={row.original} />,
            },
        ],
        [user?.name],
    );

    return (
        <>
            <div className="mx-auto mt-[15vh] flex h-full w-[calc(100%-2rem)] flex-col gap-y-4 sm:w-3/4">
                <div className="flex items-center justify-between gap-2">
                    <Input
                        className="max-w-sm"
                        onChange={(event) => setGlobalFilter(event.target.value)}
                        placeholder="Filter devices..."
                        value={globalFilter}
                    />

                    <div className="flex items-center gap-2">
                        <div className="flex rounded-lg border p-1" role="group" aria-label="Device status filter">
                            {FILTERS.map((filter) => (
                                <Button
                                    key={filter.value}
                                    aria-pressed={availabilityFilter === filter.value}
                                    onClick={() => setAvailabilityFilter(filter.value)}
                                    size="sm"
                                    variant={availabilityFilter === filter.value ? "secondary" : "ghost"}
                                >
                                    {filter.label}
                                </Button>
                            ))}
                        </div>
                        <Button
                            aria-label="Refresh devices"
                            disabled={nodesQuery.isFetching}
                            onClick={() => nodesQuery.refetch()}
                            size="icon"
                            title="Refresh"
                            variant="ghost"
                        >
                            <RefreshCw className={nodesQuery.isFetching ? "animate-spin" : undefined} />
                        </Button>
                    </div>
                </div>

                {nodesQuery.error ? (
                    <p className="text-sm text-destructive">{nodesQuery.error.message}</p>
                ) : (
                    <DataTable
                        columns={columns}
                        data={filteredNodes}
                        emptyMessage={
                            availabilityFilter === "all"
                                ? "No compute devices are registered yet."
                                : `No ${availabilityFilter} devices.`
                        }
                        loading={nodesQuery.isLoading}
                        onRowClick={setSelectedNode}
                    />
                )}
            </div>

            <NodeDetails node={selectedNode} onOpenChange={(open) => !open && setSelectedNode(null)} />
        </>
    );
};
