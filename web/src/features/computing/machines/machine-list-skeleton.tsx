interface MachineListSkeletonProps {
    count?: number;
}

const MachineListSkeleton = ({count = 1}: MachineListSkeletonProps) => {
    return (
        <aside className="space-y-4">
            <div className="flex items-start justify-between gap-3">
                <div className="min-w-0">
                    <h2 className="font-semibold">Compute nodes</h2>
                </div>
            </div>
                <div className="space-y-3">
                    {Array.from({ length: count }).map((_, index) => (
                        <article key={index} className="border-b px-1 py-3 last:border-b-0">
                            <div className="flex items-center justify-between gap-3">
                                <div className="h-4 w-32 rounded bg-neutral-800" />
                                <div className="h-4 w-8 rounded-md bg-neutral-800" />
                            </div>

                            <div className="mt-2 flex gap-4">
                                <div className="h-3 w-16 rounded bg-neutral-800" />
                                <div className="h-3 w-20 rounded bg-neutral-800" />
                                <div className="h-3 w-16 rounded bg-neutral-800" />
                            </div>
                        </article>
                    ))}
                </div>
        </aside>
    )
}

export { MachineListSkeleton };
