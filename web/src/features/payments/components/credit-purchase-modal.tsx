import { useEffect, useState } from "react";
import { Check, Coins } from "lucide-react";

import { Button, Dialog, Skeleton } from "@/shared/components/ui";
import { cn } from "@/shared/utils";
import {
    useCreateCheckoutSession,
    useCreditPacks,
} from "../hooks";


interface CreditPurchaseModalProps {
    open: boolean;
    onOpenChange: (open: boolean) => void;
}

const CreditPurchaseModal = ({
    open,
    onOpenChange,
}: CreditPurchaseModalProps) => {
    const [selectedPackId, setSelectedPackId] = useState<string | null>(
        null,
    );

    const {
        creditPacks,
        error: creditPacksError,
        isLoading: isLoadingCreditPacks,
    } = useCreditPacks();

    const {
        createCheckout,
        error: checkoutError,
        isLoading: isCreatingCheckout,
    } = useCreateCheckoutSession();

    const selectedPack = creditPacks.find(
        (pack) => pack.id === selectedPackId,
    );

    useEffect(() => {
        if (!open) {
            setSelectedPackId(null);
        }
    }, [open]);

    const handleContinue = () => {
        if (!selectedPack) {
            return;
        }

        createCheckout({
            creditPackId: selectedPack.id,
            successUrl: `${window.location.origin}/payments?checkout=success`,
            cancelUrl: `${window.location.origin}/payments?checkout=cancelled`,
        });
    };

    const formatPrice = (
    amount: number,
    currency: string,
): string =>
    new Intl.NumberFormat(undefined, {
        style: "currency",
        currency,
    }).format(amount / 100);

    const isCheckoutDisabled =
        !selectedPack ||
        isLoadingCreditPacks ||
        isCreatingCheckout;

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <Dialog.Content>
                <Dialog.Header>
                    <Dialog.Title>Buy credits</Dialog.Title>

                    <Dialog.Description>
                        Select the credit pack you want to purchase.
                    </Dialog.Description>
                </Dialog.Header>

                <div className="grid gap-3 py-2">
                    {isLoadingCreditPacks && <CreditPackSkeletons />}

                    {!isLoadingCreditPacks && creditPacksError && (
                        <p className="rounded-lg border border-destructive/40 p-4 text-sm text-destructive">
                            {creditPacksError.message}
                        </p>
                    )}

                    {!isLoadingCreditPacks &&
                        !creditPacksError &&
                        creditPacks.length === 0 && (
                            <p className="rounded-lg border p-4 text-sm text-muted-foreground">
                                No credit packs are currently available.
                            </p>
                        )}

                    {!isLoadingCreditPacks &&
                        !creditPacksError &&
                        creditPacks.map((pack) => {
                        const isSelected = pack.id === selectedPackId;

                        return (
                            <button
                                key={pack.id}
                                type="button"
                                disabled={isCreatingCheckout}
                                onClick={() => setSelectedPackId(pack.id)}
                                className={cn(
                                    "flex w-full items-center gap-4 rounded-xl border p-4 text-left transition-colors",
                                    "hover:bg-accent",
                                    "disabled:cursor-not-allowed disabled:opacity-60",
                                    isSelected &&
                                        "border-primary bg-accent",
                                )}
                            >
                                <div className="flex size-10 shrink-0 items-center justify-center rounded-full bg-muted">
                                    <Coins className="size-5" />
                                </div>

                                <div className="min-w-0 flex-1">
                                    <p className="font-medium">
                                        {pack.name}
                                    </p>

                                    <p className="text-sm text-muted-foreground">
                                        {pack.credits.toLocaleString()} credits
                                    </p>
                                </div>
                                <div className="shrink-0 text-right">
                                    <p className="font-medium">
                                        {formatPrice(
                                            pack.price_amount,
                                            pack.currency,
                                        )}
                                    </p>
                                </div>

                                <div
                                    className={cn(
                                        "flex size-5 items-center justify-center rounded-full border",
                                        isSelected &&
                                            "border-primary bg-primary text-primary-foreground",
                                    )}
                                >
                                    {isSelected && (
                                        <Check className="size-3" />
                                    )}
                                </div>
                            </button>
                        );
                    })}

                    {checkoutError && (
                        <p className="rounded-lg border border-destructive/40 p-4 text-sm text-destructive">
                            {checkoutError.message}
                        </p>
                    )}
                </div>

                <Dialog.Footer>
                    <Button
                        type="button"
                        variant="outline"
                        disabled={isCreatingCheckout}
                        onClick={() => onOpenChange(false)}
                    >
                        Cancel
                    </Button>

                    <Button
                        type="button"
                        disabled={isCheckoutDisabled}
                        onClick={handleContinue}
                    >
                        {isCreatingCheckout
                            ? "Opening Stripe..."
                            : selectedPack
                              ? `Continue with ${selectedPack.credits.toLocaleString()} credits`
                              : "Select a pack"}
                    </Button>
                </Dialog.Footer>
            </Dialog.Content>
        </Dialog>
    );
};

const CreditPackSkeletons = () => (
    <>
        {[0, 1, 2].map((item) => (
            <div
                key={item}
                className="flex items-center gap-4 rounded-xl border p-4"
            >
                <Skeleton className="size-10 rounded-full" />

                <div className="flex-1 space-y-2">
                    <Skeleton className="h-4 w-32" />
                    <Skeleton className="h-3 w-24" />
                </div>
            </div>
        ))}
    </>
);

export { CreditPurchaseModal };