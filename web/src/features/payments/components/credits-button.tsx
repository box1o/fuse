import * as React from "react";
import { Zap } from "lucide-react";
import { toast } from "sonner";

import {
    Button,
    Tooltip,
    TooltipContent,
    TooltipTrigger,
} from "@/shared/components";
import { useWorkspaceStore } from "@/features/workspace";
import { useCreateCheckoutSession } from "../hooks";
import type {
    CreditBalance,
    SubscriptionPlanId,
} from "../types";
import { SubscriptionDialog } from "./subscription-dialog";

const MOCK_BALANCE: CreditBalance = {
    planId: "free",
    status: "free",
    usedCredits: 18,
    includedCredits: 240,
    remainingCredits: 222,
    resetAt: new Date().toISOString(),
    nextResetDate: new Date(
        Date.now() + 30 * 24 * 60 * 60 * 1000
    ).toISOString(),
};

const CreditsButton = () => {
    const [open, setOpen] = React.useState(false);

    const currentWorkspace = useWorkspaceStore(
        (state) => state.currentWorkspace,
    );

    const checkout = useCreateCheckoutSession();

    const balance = MOCK_BALANCE;

    const remaining = Math.max(
        balance.includedCredits - balance.usedCredits,
        0,
    );

    const handleSelectPlan = (
        planId: SubscriptionPlanId,
    ) => {
        if (planId === balance.planId) {
            return;
        }

        if (!currentWorkspace?.id) {
            toast.error(
                "Select a workspace before upgrading",
            );
            return;
        }

        /*
         * Your existing backend checkout accepts a resource type,
         * not a plan ID.
         *
         * This temporarily starts CPU checkout.
         * The backend should later accept:
         *
         * plan_id: "pro"
         */
        checkout.checkout({
            workspaceId: currentWorkspace.id,
            resourceType: "cpu",
            successUrl: `${window.location.origin}/payments?checkout=success`,
            cancelUrl: `${window.location.origin}/payments?checkout=canceled`,
        });
    };

    const handleManageSubscription = () => {
        toast.info(
            "Stripe Customer Portal endpoint is not connected yet",
        );
    };

    return (
        <>
            <Tooltip>
                <TooltipTrigger asChild>
                    <Button
                        type="button"
                        size="sm"
                        variant="outline"
                        className="h-8 max-w-[9rem] gap-1.5 rounded-full px-2.5"
                        onClick={() => setOpen(true)}
                    >
                        <Zap className="size-3.5 fill-brand text-brand" />

                        <span className="hidden truncate text-xs sm:inline">
                            {balance.planId === "pro"
                                ? `Pro · ${remaining.toLocaleString()}`
                                : `${remaining} credits`}
                        </span>

                        <span className="sr-only">
                            Open subscription plans
                        </span>
                    </Button>
                </TooltipTrigger>

                <TooltipContent>
                    {remaining.toLocaleString()} of{" "}
                    {balance.includedCredits.toLocaleString()}{" "}
                    credits remaining
                </TooltipContent>
            </Tooltip>

            <SubscriptionDialog
                open={open}
                onOpenChange={setOpen}
                balance={balance}
                isUpgrading={checkout.isLoading}
                onSelectPlan={handleSelectPlan}
                onManageSubscription={
                    handleManageSubscription
                }
            />
        </>
    );
};

export { CreditsButton };