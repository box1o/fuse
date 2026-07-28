interface MachineListSkeletonProps {
    count?: number;
}

const MachineListSkeleton = ({count = 1}: MachineListSkeletonProps) => {
    return (
         <div className="space-y-3">
            <aside className="space-y-4">
                {Array.from({ length: count }).map((_, index) => (
                    <div
                        key={index}
                        className="grid grid-cols-[160px_1fr] gap-3"
                    >
                        <div className="aspect-video rounded-lg bg-neutral-800" />

                        <div className="space-y-2">
                            <div className="h-4 rounded bg-neutral-800" />
                            <div className="h-4 w-3/4 rounded bg-neutral-800" />
                            <div className="h-3 w-1/2 rounded bg-neutral-800" />
                        </div>
                    </div>
                ))}
            </aside>
        </div>
    )
}

export { MachineListSkeleton };
