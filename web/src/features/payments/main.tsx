import React from "react";
import { Button } from "@/shared/components";
import { useWorkspaceStore } from "../workspace/store";
import { useCreateCheckoutSession } from "./hooks";

const Main: React.FC = () => {
    const currentWorkspace = useWorkspaceStore((state) => state.currentWorkspace);
    const checkout = useCreateCheckoutSession();

    const workspaceId = currentWorkspace?.id ?? "";

    const handleCheckout = () => {
        if (!workspaceId) {
            return;
        }

        checkout.checkout({
            workspaceId,
            resourceType: "cpu",
            successUrl: `${window.location.origin}/payments`,
            cancelUrl: `${window.location.origin}/payments`,
        });
    };

    return (
        <div className="flex h-full w-full items-center justify-center p-6">
            <div className="flex w-full max-w-md flex-col gap-4 rounded-2xl border bg-background/80 p-8 shadow-sm">
                <div>
                    <h1 className="text-2xl font-semibold">Payments test</h1>
                    <p className="mt-2 text-sm text-muted-foreground">
                        One button, hardcoded checkout session, Stripe hosted UI.
                    </p>
                </div>

                <p className="text-sm">
                    Workspace: <span className="font-medium">{workspaceId || "none selected"}</span>
                </p>

                <Button type="button" onClick={handleCheckout} disabled={!workspaceId || checkout.isLoading}>
                    {checkout.isLoading ? "Opening Stripe..." : "Pay with Stripe"}
                </Button>
            </div>
        </div>
    );
};

export default Main;
