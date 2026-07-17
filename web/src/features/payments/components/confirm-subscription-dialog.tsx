import {
    Button,
    Dialog,
} from "@/shared/components";
import type { SubscriptionPlan } from "../types";

interface ConfirmSubscriptionDialogProps {
    open: boolean;
    plan: SubscriptionPlan | null;
    isLoading?: boolean;
    onOpenChange: (open: boolean) => void;
    onConfirm: () => void;
}

const formatPrice = (priceMonthlyCents: number): string => {
    return new Intl.NumberFormat("en-US", {
        style: "currency",
        currency: "USD",
    }).format(priceMonthlyCents / 100);
};

const ConfirmSubscriptionDialog = ({
    open,
    plan,
    isLoading = false,
    onOpenChange,
    onConfirm,
}: ConfirmSubscriptionDialogProps) => {
    if (!plan) {
        return null;
    }

    return (
        <Dialog
            open={open}
            onOpenChange={(nextOpen) => {
                if (!isLoading) {
                    onOpenChange(nextOpen);
                }
            }}
        >
            <Dialog.Content className="w-[calc(100vw-2rem)] max-w-md rounded-2xl">
                <Dialog.Header>
                    <Dialog.Title>
                        Confirm subscription
                    </Dialog.Title>

                    <Dialog.Description>
                        Review the subscription before continuing to checkout.
                    </Dialog.Description>
                </Dialog.Header>

                <div className="my-5 rounded-xl border bg-muted/20 p-4">
                    <div className="flex items-center justify-between gap-4">
                        <div>
                            <p className="font-semibold">
                                {plan.name} plan
                            </p>

                            <p className="mt-1 text-sm text-muted-foreground">
                                {plan.includedCredits.toLocaleString()} credits per{" "}
                                {plan.resetInterval === "monthly"
                                    ? "month"
                                    : "year"}
                            </p>
                        </div>

                        <p className="shrink-0 font-semibold">
                            {formatPrice(plan.priceMonthlyCents)}
                            <span className="text-sm font-normal text-muted-foreground">
                                /month
                            </span>
                        </p>
                    </div>
                </div>

                <Dialog.Footer>
                    <Button
                        type="button"
                        variant="outline"
                        disabled={isLoading}
                        onClick={() => onOpenChange(false)}
                    >
                        Cancel
                    </Button>

                    <Button
                        type="button"
                        disabled={isLoading}
                        onClick={onConfirm}
                    >
                        {isLoading
                            ? "Opening checkout..."
                            : "Continue to checkout"}
                    </Button>
                </Dialog.Footer>
            </Dialog.Content>
        </Dialog>
    );
};

export { ConfirmSubscriptionDialog };