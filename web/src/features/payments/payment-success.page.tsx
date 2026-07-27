import { CheckCircle2 } from "lucide-react";
import { Link } from "react-router-dom";

import { Button } from "@/shared/components/ui";
import { ROUTES } from "@/shared/constants/routes.constants";

const PaymentSuccessPage = () => (
    <main className="flex min-h-screen items-center justify-center bg-background px-6">
        <section className="w-full max-w-md rounded-2xl border border-emerald-500/20 bg-card p-8 text-center">
            <CheckCircle2 className="mx-auto size-12 text-emerald-400" />

            <h1 className="mt-5 text-2xl font-semibold">
                Payment successful
            </h1>

            <p className="mt-3 text-sm text-muted-foreground">
                Stripe confirmed the payment. Your credits will be added
                after the webhook is processed.
            </p>

            <Button asChild className="mt-6">
                <Link to={ROUTES.PROJECTS}>Return to Fuse</Link>
            </Button>
        </section>
    </main>
);

export { PaymentSuccessPage as Component };