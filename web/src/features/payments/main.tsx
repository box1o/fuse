import { type FormEvent, useState } from "react";
import { Button, Input } from "@/shared/components";
import { useCreateCheckoutSession } from "./hooks";

const Main = () => {
    const [creditPackId, setCreditPackId] = useState("");
    const checkout = useCreateCheckoutSession();

    const handleSubmit = (event: FormEvent<HTMLFormElement>) => {
        event.preventDefault();

        checkout.createCheckout({
            creditPackId,
            successUrl: `${window.location.origin}/payments?checkout=success`,
            cancelUrl: `${window.location.origin}/payments?checkout=cancelled`,
        });
    };

    return (
        <div className="flex h-full w-full items-center justify-center p-6">
            <form
                onSubmit={handleSubmit}
                className="flex w-full max-w-md flex-col gap-5 rounded-2xl border bg-background p-8 shadow-sm"
            >
                <div>
                    <h1 className="text-2xl font-semibold">
                        Credit checkout test
                    </h1>

                    <p className="mt-2 text-sm text-muted-foreground">
                        Enter an active credit pack UUID and open Stripe
                        Checkout.
                    </p>
                </div>

                <div className="flex flex-col gap-2">
                    <label
                        htmlFor="credit-pack-id"
                        className="text-sm font-medium"
                    >
                        Credit pack ID
                    </label>

                    <Input
                        id="credit-pack-id"
                        value={creditPackId}
                        onChange={(event) =>
                            setCreditPackId(event.target.value)
                        }
                        placeholder="xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
                        autoComplete="off"
                        required
                    />
                </div>

                {checkout.error && (
                    <p className="text-sm text-destructive">
                        {checkout.error.message}
                    </p>
                )}

                <Button
                    type="submit"
                    disabled={
                        checkout.isLoading ||
                        creditPackId.trim().length === 0
                    }
                >
                    {checkout.isLoading
                        ? "Opening Stripe..."
                        : "Buy credits"}
                </Button>
            </form>
        </div>
    );
};

export default Main;