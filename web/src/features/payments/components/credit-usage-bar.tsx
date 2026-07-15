import { cn } from "@/shared/utils";

interface CreditUsageBarProps {
    used: number;
    total: number;
    className?: string;
}

const CreditUsageBar = ({ used, total, className }: CreditUsageBarProps) => {
    const safeTotal =  Math.max(total, 1); // Ensure total is at least 1 to avoid division by zero
    const safeUsed = Math.min(Math.max(used, 0), safeTotal); // Clamp used between 0 and total
    const remaining = Math.max(safeTotal - safeUsed, 0);

    const usedPercentage = (safeUsed / safeTotal) * 100;
    const remainingPercentage = 100 - usedPercentage;

    const isEmpty = safeUsed === 0;
    const isLow = remainingPercentage <= 20 && !isEmpty; // Consider low if remaining is 20% or less and not empty

    return (
        <div className={cn("space-y-2", className)}>
            <div className={cn("flex items-center justify-between gap-4 text-sm")}>
                <span className="font-medium">
                    {remaining.toLocaleString()} credits remaining
                </span>
                <span className="shrink-0 text-muted-foreground">
                    {safeUsed.toLocaleString()} / {""}
                    {total.toLocaleString()} credits used
                </span>
            </div>

            <div
                role="progressbar"
                aria-label="compute credits usage"
                aria-valuemin={0}
                aria-valuenow={safeUsed}
                aria-valuemax={total}
                className={"h-2 overflow-hidden rounded-full bg-muted"}
            >
                <div
                    className={cn(
                        "h-full rounded-all transition-[width] duration-300",
                        isEmpty 
                            ? "bg-destructive" 
                            : isLow
                                ? "bg-yellow-500"
                                : "bg-brand"
                    )}
                    style={{ width: `${usedPercentage}%` }}
                />

            </div>
        </div>

    );
};

export { CreditUsageBar };