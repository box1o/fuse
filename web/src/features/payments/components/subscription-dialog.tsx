import * as React from "react";
import {
    ChevronDown,
    ExternalLink,
} from "lucide-react";

import {
    Button,
    Collapsible,
    Dialog,
} from "@/shared/components";
import { SUBSCRIPTION_PLANS } from "../constants";
import type {
    CreditBalance,
    SubscriptionPlanId,
} from "../types";
import { CreditUsageBar } from "./credit-usage-bar";
import { PlanCard } from "./plan-card";

interface SubscriptionDialogProps {
    open: boolean;
    onOpenChange: (open: boolean) => void;
    balance: CreditBalance;
    isUpgrading?: boolean;
    isManaging?: boolean;
    onSelectPlan: (planId: SubscriptionPlanId) => void;
    onManageSubscription: () => void;
}

const formatResetDate = (value: string): string => {
    const date = new Date(value);

    if (Number.isNaN(date.getTime())) {
        return "Unknown";
    }

    return new Intl.DateTimeFormat("en-US", {
        month: "short",
        day: "numeric",
        year: "numeric",
    }).format(date);
};

const formatStatus = (status: CreditBalance["status"]): string => {
    if (status === "free") {
        return "Free";
    }

    return status
        .replaceAll("_", " ")
        .replace(/\b\w/g, (character) =>
            character.toUpperCase(),
        );
};

const SubscriptionDialog = ({
    open,
    onOpenChange,
    balance,
    isUpgrading = false,
    isManaging = false,
    onSelectPlan,
    onManageSubscription,
}: SubscriptionDialogProps) => {
    const [creditsOpen, setCreditsOpen] =
        React.useState(false);

    const currentPlan =
        SUBSCRIPTION_PLANS.find(
            (plan) => plan.id === balance.planId,
        ) ?? SUBSCRIPTION_PLANS[0];
    const isPaidPlan = balance.planId === "pro";

    return (
        <Dialog
            open={open}
            onOpenChange={onOpenChange}
        >
            <Dialog.Content className="max-h-[calc(100dvh-2rem)] w-[calc(100vw-2rem)] sm:max-w-2xl overflow-y-auto rounded-2xl p-0">
                <div className="border-b p-5 pr-12 sm:p-6 sm:pr-12">
                    <Dialog.Title>
                        Your plan
                    </Dialog.Title>

                    <Dialog.Description className="mt-1">
                        Manage your compute credits and
                        subscription.
                    </Dialog.Description>
                </div>

                <div className="space-y-6 p-5 sm:p-6">
                    <section className="rounded-2xl border bg-muted/20 p-4 sm:p-5">
                        <div className="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
                            <div className="min-w-0">
                                <div className="flex flex-wrap items-center gap-2">
                                    <h2 className="font-semibold">
                                        {currentPlan.name} plan
                                    </h2>

                                    <span className="rounded-full bg-brand/15 px-2 py-0.5 text-xs font-medium">
                                        {formatStatus(
                                            balance.status,
                                        )}
                                    </span>
                                </div>

                                <p className="mt-1 text-sm text-muted-foreground">
                                    Credits reset on{" "}
                                    {formatResetDate(
                                        balance.resetAt,
                                    )}
                                    .
                                </p>
                            </div>

                            {isPaidPlan && (
                                <Button
                                    type="button"
                                    size="sm"
                                    variant="outline"
                                    className="w-full shrink-0 sm:w-auto"
                                    disabled={isManaging}
                                    onClick={
                                        onManageSubscription
                                    }
                                >
                                    {isManaging
                                        ? "Opening..."
                                        : "Manage subscription"}

                                    <ExternalLink className="size-3.5" />
                                </Button>
                            )}
                        </div>

                        <CreditUsageBar
                            className="mt-5"
                            used={balance.usedCredits}
                            total={balance.includedCredits}
                        />
                    </section>

                    <section>
                        <h2 className="font-semibold">
                            Choose a plan
                        </h2>

                        <p className="mt-1 text-sm text-muted-foreground">
                            Free is the baseline plan. Pro unlocks more compute credits and accelerated resources.
                        </p>

                        <div className="mt-4 grid grid-cols-1 gap-4 md:grid-cols-2">
                            <p>
                                SUB here will be the plan cards for each subscription plan, allowing users to select or upgrade their plan. Each card will display the plan name, description, price, and included credits. Users can click on a plan to select it, which will trigger the onSelectPlan callback with the selected plan's ID.
                            </p>
                        </div>
                    </section>

                    <Collapsible
                        open={creditsOpen}
                        onOpenChange={setCreditsOpen}
                    >
                        <Collapsible.Trigger asChild>
                            <Button
                                type="button"
                                variant="ghost"
                                className="w-full justify-between rounded-xl px-3"
                            >
                                How compute credits work

                                <ChevronDown
                                    className={
                                        creditsOpen
                                            ? "size-4 rotate-180 transition-transform"
                                            : "size-4 transition-transform"
                                    }
                                />
                            </Button>
                        </Collapsible.Trigger>

                        <Collapsible.Content>
                            <div className="mt-2 space-y-2 rounded-xl border bg-muted/20 p-4 text-sm">
                                <div className="flex justify-between gap-4">
                                    <span>
                                        1 CPU minute
                                    </span>
                                    <span className="font-medium">
                                        1 credit
                                    </span>
                                </div>

                                <div className="flex justify-between gap-4">
                                    <span>
                                        1 NPU minute
                                    </span>
                                    <span className="font-medium">
                                        2 credits
                                    </span>
                                </div>

                                <div className="flex justify-between gap-4">
                                    <span>
                                        1 GPU minute
                                    </span>
                                    <span className="font-medium">
                                        4 credits
                                    </span>
                                </div>
                            </div>
                        </Collapsible.Content>
                    </Collapsible>
                </div>
            </Dialog.Content>
        </Dialog>
    );
};

export { SubscriptionDialog };
