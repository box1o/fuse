import { useEffect, useState } from "react";
import { Button, Dialog, Skeleton } from "@/shared/components/ui";
import { cn } from "@/shared/utils";


import {
    useCreateCheckoutSession,
    useCreditPacks,
} from "../hooks";
import type { CreditPack } from "../types";

interface CreditPurchaseModalProps {
    open: boolean;
    onOpenChange: (open: boolean) => void;
}

const CreditPurchaseModal = ({open,onOpenChange,}: CreditPurchaseModalProps) => {
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
        resetCheckout,
    } = useCreateCheckoutSession();

    useEffect(() => {
        if (!open) {
            setSelectedPackId(null);
            resetCheckout();
        }
    }, [open, resetCheckout]);

    useEffect(() => {
        if (
            open &&
            creditPacks.length > 0 &&
            selectedPackId === null
        ) {
            const defaultPackIndex =
                creditPacks.length >= 3 ? 1 : 0;

            setSelectedPackId(
                creditPacks[defaultPackIndex]?.id ?? null,
            );
        }
    }, [creditPacks, open, selectedPackId]);

    const handlePurchase = async (pack: CreditPack) => {
        if (isCreatingCheckout) {
            return;
        }

        setSelectedPackId(pack.id);

        try {
            const checkout = await createCheckout({
                creditPackId: pack.id,
                successUrl:
                    `${window.location.origin}/payments/success` +
                    "?session_id={CHECKOUT_SESSION_ID}",
                cancelUrl: `${window.location.origin}/payments/cancel`,
            });

            window.location.assign(checkout.url);
        } catch {
            // The mutation displays the error through the hook.
        }
    };
    
    return (
        <Dialog open={open} onOpenChange={onOpenChange} >
            <Dialog.Content 
                showCloseButton={false}
                overlayClassName="bg-black/30 backdrop-blur-sm"
                className="h-[70vh] w-[90vw] max-w-noneo verflow-visible border-0 bg-transparent p-0 shadow-none sm:max-w-6xl"
            >
                <div className="h-full px-10 py-10 sm:px-7 sm:py-8">
                    {isLoadingCreditPacks && <CreditPackSkeletons />}

                    {!isLoadingCreditPacks && creditPacksError && (
                        <ErrorMessage
                            message={creditPacksError.message}
                        />
                    )}

                    {!isLoadingCreditPacks &&
                        !creditPacksError &&
                        creditPacks.length === 0 && (
                            <EmptyCreditPacks />
                        )}

                    {!isLoadingCreditPacks &&
                        !creditPacksError &&
                        creditPacks.length > 0 && (
                            <>
                                <div className="grid h-full gap-4 md:grid-cols-3">
                                    {creditPacks.map((pack, index) => {
                                        const isFeatured =
                                            creditPacks.length >= 3 && index === 1;

                                        return (
                                            <CreditPackCard
                                                key={pack.id}
                                                pack={pack}
                                                isFeatured={isFeatured}
                                                isSelected={selectedPackId === pack.id}
                                                isLoading={
                                                    isCreatingCheckout &&
                                                    selectedPackId === pack.id
                                                }
                                                isDisabled={
                                                    isCreatingCheckout &&
                                                    selectedPackId !== pack.id
                                                }
                                                onSelect={() => setSelectedPackId(pack.id)}
                                                onPurchase={() => handlePurchase(pack)}
                                            />
                                        );
                                    })}
                                </div>
                            </>
                        )}

                    {checkoutError && (
                        <div className="mt-4">
                            <ErrorMessage
                                message={checkoutError.message}
                            />
                        </div>
                    )}

                </div>
            </Dialog.Content>
        </Dialog>
    );
};

interface CreditPackCardProps {
    pack: CreditPack;
    isFeatured: boolean;
    isSelected: boolean;
    isLoading: boolean;
    isDisabled: boolean;
    onSelect: () => void;
    onPurchase: () => void;
}

const CreditPackCard = ({ pack, isFeatured, isSelected, isLoading, isDisabled, onSelect, onPurchase}: CreditPackCardProps) => {
    return (
        <article
            onClick={onSelect}
            className={cn(
                "relative flex min-h-[360px] cursor-pointer flex-col overflow-hidden rounded-2xl border bg-card p-5 transition-all duration-200",
                "hover:-translate-y-1 hover:border-emerald-500/40 hover:shadow-lg",
                isSelected &&
                    "border-emerald-500/70 shadow-[0_0_0_1px_rgba(34,197,94,0.2),0_18px_50px_rgba(34,197,94,0.10)]",
            )}
        >
            {isFeatured && (
                <FeaturedBadge />
            )}

            <div className="mt-7 text-center">
                <div className="mt-3 flex items-center justify-center gap-2">
                    <p className="text-4xl font-semibold tracking-tight">
                        {pack.credits.toLocaleString()}
                    </p>
                    <p>
                        Credits
                    </p>
                </div>
            </div>

            <div className="my-6 h-px bg-border/70" />

            <div className="text-center">
                <p className="text-sm leading-6 text-muted-foreground">
                    {getPackDescription(pack.credits)}
                </p>
            </div>

            <div className="mt-auto pt-7 text-center">
                <p className="text-3xl font-semibold tracking-tight">
                    {formatPrice(
                        pack.price_amount,
                        pack.currency,
                    )}
                </p>

                <p className="mt-1 text-xs text-muted-foreground">
                    One-time payment
                </p>

                <Button
                    type="button"
                    className={cn(
                        "mt-5 w-full font-medium",
                        "bg-emerald-500 text-black hover:bg-emerald-400",
                        "shadow-[0_8px_24px_rgba(34,197,94,0.18)]",
                    )}
                    disabled={isDisabled || isLoading}
                    onClick={(event) => {
                        event.stopPropagation();
                        onPurchase();
                    }}
                >
                    {isLoading ? (
                        "Preparing checkout..."
                    ) : (
                        <>
                            Buy now
                        </>
                    )}
                </Button>
            </div>
        </article>
    );
};

const FeaturedBadge = () => (
    <div className="absolute left-1/2 top-0 -translate-x-1/2 -translate-y-px">
        <div className="flex items-center gap-1.5 rounded-b-xl border-x border-b border-emerald-500/40 bg-emerald-500 px-3 py-1 text-[10px] font-bold uppercase tracking-wider text-black">
            Best value
        </div>
    </div>
);

const CreditPackSkeletons = () => (
    <div className="grid gap-4 md:grid-cols-3">
        {[0, 1, 2].map((item) => (
            <div
                key={item}
                className="flex min-h-[360px] flex-col rounded-2xl border bg-card/70 p-5"
            >
                <div className="flex justify-between">
                    <Skeleton className="size-12 rounded-2xl" />
                    <Skeleton className="size-6 rounded-full" />
                </div>

                <div className="mt-7 space-y-3 text-center">
                    <Skeleton className="mx-auto h-3 w-20" />
                    <Skeleton className="mx-auto h-10 w-32" />
                    <Skeleton className="mx-auto h-4 w-16" />
                </div>

                <Skeleton className="my-6 h-px w-full" />
                <Skeleton className="mx-auto h-12 w-40" />

                <div className="mt-auto space-y-4 pt-7">
                    <Skeleton className="mx-auto h-9 w-24" />
                    <Skeleton className="h-10 w-full" />
                </div>
            </div>
        ))}
    </div>
);

interface ErrorMessageProps {
    message: string;
}

const ErrorMessage = ({ message }: ErrorMessageProps) => (
    <p className="rounded-xl border border-destructive/40 bg-destructive/5 p-4 text-sm text-destructive">
        {message}
    </p>
);

const EmptyCreditPacks = () => (
    <div className="rounded-2xl border border-dashed p-10 text-center">

        <p className="mt-4 font-medium">
            No credit packs available
        </p>

        <p className="mt-1 text-sm text-muted-foreground">
            Available credit packs will appear here.
        </p>
    </div>
);

const formatPrice = (
    amount: number,
    currency: string,
): string =>
    new Intl.NumberFormat(undefined, {
        style: "currency",
        currency,
    }).format(amount / 100);

const getPackDescription = (credits: number): string => {
    if (credits <= 500) {
        return "Suitable for light usage, testing, and smaller workloads.";
    }

    if (credits <= 3000) {
        return "A balanced option for regular projects and ongoing work.";
    }

    return "Designed for larger workloads and frequent compute usage.";
};

export { CreditPurchaseModal };