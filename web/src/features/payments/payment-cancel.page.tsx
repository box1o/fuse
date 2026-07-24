import { XCircle } from "lucide-react";
import { Link } from "react-router-dom";

import { Button } from "@/shared/components/ui";
import { ROUTES } from "@/shared/constants/routes.constants";

const PaymentCancelPage = () => (
    <main className="flex min-h-screen items-center justify-center bg-background px-6">
        <section className="w-full max-w-md rounded-2xl border border-border bg-card p-8 text-center">
            <XCircle className="mx-auto size-12 text-muted-foreground" />

            <h1 className="mt-5 text-2xl font-semibold">
                Payment canceled
            </h1>

            <p className="mt-3 text-sm text-muted-foreground">
                No payment was completed and no credits were added.
            </p>

            <Button asChild className="mt-6">
                <Link to={ROUTES.PROJECTS}>Return to Fuse</Link>
            </Button>
        </section>
    </main>
);

export { PaymentCancelPage as Component };