import * as React from "react";
import { Zap } from "lucide-react";
import { toast } from "sonner";

import {
    Button,
    Tooltip,
    TooltipContent,
    TooltipTrigger,
} from "@/shared/components";
import { useAuthStore } from "@/features/auth";
import { useCreateCheckoutSession } from "../hooks";
import type {
    CreditBalance,
    SubscriptionPlanId,
} from "../types";
import { useMockBillingStore } from "../store/mock-billing.store";
import { SubscriptionDialog } from "./subscription-dialog";

const CreditsButton = () => {
    const [open, setOpen] = React.useState(false);

    const checkout = useCreateCheckoutSession();
    const user = useAuthStore((state) => state.user);
    const setUserKey = useMockBillingStore((state) => state.setUserKey);
    const planId = useMockBillingStore((state) => state.planId);
    const usedCredits = useMockBillingStore((state) => state.usedCredits);
    const includedCredits = useMockBillingStore((state) => state.includedCredits);
    const buyPro = useMockBillingStore((state) => state.buyPro);

    React.useEffect(() => {
        setUserKey(user?.id ?? user?.email ?? null);
    }, [setUserKey, user?.email, user?.id]);

    const balance: CreditBalance = {
        planId,
        status: planId === "pro" ? "active" : "free",
        usedCredits,
        includedCredits,
        remainingCredits: Math.max(includedCredits - usedCredits, 0),
        resetAt: new Date().toISOString(),
        nextResetDate: new Date(
            Date.now() + 30 * 24 * 60 * 60 * 1000
        ).toISOString(),
    };

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

        buyPro();

        checkout.checkout({
            planId: "pro",
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
                        className="px-2 py-0 max-w-[9rem] gap-1.5 rounded-full "
                        onClick={() => setOpen(true)}
                    >
                        <Zap className="size-3.5 fill-brand text-brand" />
                        <span className="hidden truncate text-xs sm:inline">
                            {balance.planId === "pro" ? "Pro(Manage)" : "Free"}
                        </span>
                    </Button>
                </TooltipTrigger>

                <TooltipContent>
                    {balance.planId === "pro"
                        ? `Pro active · ${remaining.toLocaleString()} credits remaining`
                        : `${remaining.toLocaleString()} of ${balance.includedCredits.toLocaleString()} credits remaining`}
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
