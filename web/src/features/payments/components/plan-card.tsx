import { Check } from "lucide-react";

import { Badge, Button } from "@/shared/components";
import { cn } from "@/shared/utils";
import type {
    SubscriptionPlan,
    SubscriptionPlanId,
} from "../types";

interface PlanCardProps {
    plan: SubscriptionPlan;
    currentPlanId: SubscriptionPlanId;
    isLoading?: boolean;
    onSelect: (planId: SubscriptionPlanId) => void;
}

const formatPrice = (priceMonthlyCents: number): string => {
    if (priceMonthlyCents === 0) {
        return "$0";
    }

    return new Intl.NumberFormat("en-US", {
        style: "currency",
        currency: "USD",
        maximumFractionDigits: 0,
    }).format(priceMonthlyCents / 100);
};

const PlanCard = ({
    plan,
    currentPlanId,
    isLoading = false,
    onSelect,
}: PlanCardProps) => {
    const isCurrent = plan.id === currentPlanId;
    const isFree = plan.priceMonthlyCents === 0;

    return (
        <article
            className={cn(
                "relative flex min-w-0 flex-col rounded-2xl border bg-background p-5",
                plan.recommended && "border-brand/70 shadow-sm",
            )}
        >
            {plan.recommended && (
                <Badge
                    variant="brand"
                    className="absolute right-4 top-4 rounded-full px-2 py-0.5"
                >
                    Recommended
                </Badge>
            )}

            <div className="pr-24">
                <h3 className="text-lg font-semibold">
                    {plan.name}
                </h3>

                <p className="mt-1 text-sm text-muted-foreground">
                    {plan.description}
                </p>
            </div>

            <div className="mt-6 flex items-baseline gap-1">
                <span className="text-3xl font-semibold tracking-tight">
                    {formatPrice(plan.priceMonthlyCents)}
                </span>

                {!isFree && (
                    <span className="text-sm text-muted-foreground">
                        / month
                    </span>
                )}
            </div>

            <p className="mt-2 text-sm font-medium">
                {plan.includedCredits.toLocaleString()} credits{" "}
                {plan.resetInterval === "monthly"
                    ? "per month"
                    : "per year"}
            </p>

            <ul className="mt-6 flex-1 space-y-3">
                {plan.features.map((feature) => (
                    <li
                        key={feature}
                        className="flex items-start gap-2 text-sm"
                    >
                        <Check className="mt-0.5 size-4 shrink-0 text-brand" />
                        <span>{feature}</span>
                    </li>
                ))}
            </ul>

            <Button
                type="button"
                variant={isCurrent ? "outline" : "default"}
                className={cn(
                    "mt-6 w-full",
                    !isCurrent &&
                        plan.recommended &&
                        "bg-brand text-black hover:bg-brand/90",
                )}
                disabled={isCurrent || isLoading}
                onClick={() => onSelect(plan.id)}
            >
                {isCurrent
                    ? "Current plan"
                    : isLoading
                      ? "Opening checkout..."
                      : `Upgrade to ${plan.name}`}
            </Button>
        </article>
    );
};

export { PlanCard };