import { Button } from "@/shared/components";
import { CheckCircle2, LoaderCircle, XCircle } from "lucide-react";
import { DeviceCodeInput } from "./components";
import { useDeviceAuthorization } from "./hooks";

const DevicePage = () => {
    const authorization = useDeviceAuthorization();
    const busy = authorization.state === "checking-code" || authorization.state === "submitting";
    const complete = authorization.state === "approved" || authorization.state === "denied";

    return (
        <main className="flex min-h-screen items-center justify-center bg-background px-6 py-12 text-foreground">
            <section className="w-full max-w-md rounded-2xl border bg-card p-8 shadow-sm">
                {!complete && (
                    <DeviceCodeInput
                        disabled={busy}
                        onChange={authorization.updateCode}
                        value={authorization.code}
                    />
                )}

                {authorization.authenticated === false && !complete && (
                    <div className="mt-8 rounded-xl border bg-muted/30 p-4">
                        <p className="text-sm text-muted-foreground">
                            Sign in before authorizing this device. You will return to this page afterward.
                        </p>
                        <Button className="mt-4 w-full" onClick={authorization.signIn} size="lg">
                            Sign in with Google
                        </Button>
                    </div>
                )}

                {authorization.request && authorization.state === "ready" && (
                    <div className="mt-8 flex justify-center gap-3">
                        <Button onClick={authorization.deny} variant="ghost">
                            Deny
                        </Button>
                        <Button onClick={authorization.approve}>Authorize</Button>
                    </div>
                )}

                {busy && (
                    <p className="mt-6 flex items-center gap-2 text-sm text-muted-foreground">
                        <LoaderCircle className="size-4 animate-spin" aria-hidden="true" />
                        {authorization.state === "checking-code" ? "Checking code…" : "Authorizing CLI…"}
                    </p>
                )}

                {authorization.state === "invalid" && (
                    <p className="mt-6 text-sm text-destructive">
                        This code is invalid, expired, or has already been used.
                    </p>
                )}

                {authorization.state === "approved" && (
                    <div className="text-center">
                        <h2 className="text-lg font-medium">CLI authorized</h2>
                        <CheckCircle2 className="mx-auto mt-4 size-10 text-emerald-500" aria-hidden="true" />
                        <p className="mt-2 text-sm text-muted-foreground">
                            Return to your terminal to continue. You can safely close this page.
                        </p>
                    </div>
                )}

                {authorization.state === "denied" && (
                    <div className="text-center">
                        <XCircle className="mx-auto size-10 text-muted-foreground" aria-hidden="true" />
                        <h2 className="mt-4 text-lg font-medium">Authorization denied</h2>
                        <p className="mt-2 text-sm text-muted-foreground">
                            No CLI credential was created. You can safely close this page.
                        </p>
                    </div>
                )}
            </section>
        </main>
    );
};

export const Component = DevicePage;
