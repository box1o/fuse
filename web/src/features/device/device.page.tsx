import { api } from "@/shared/services";
import { Button, InputOTP } from "@/shared/components";
import { useEffect, useState } from "react";

interface DeviceRequest {
    user_code: string;
    client_name: string;
    status: string;
    expires_at: string;
}

const normalizeCode = (value: string) => value.toUpperCase().replace(/[^A-Z0-9]/g, "").slice(0, 8);
const formatCode = (value: string) => `${value.slice(0, 4)}-${value.slice(4)}`;

const DevicePage = () => {
    const [code, setCode] = useState("");
    const [authenticated, setAuthenticated] = useState<boolean | null>(null);
    const [request, setRequest] = useState<DeviceRequest | null>(null);
    const [message, setMessage] = useState("");
    const [busy, setBusy] = useState(false);

    useEffect(() => {
        api.get("/auth/status")
            .then(() => setAuthenticated(true))
            .catch(() => setAuthenticated(false));
    }, []);

    useEffect(() => {
        if (!authenticated || code.length !== 8) return;
        setBusy(true);
        setMessage("");
        api.get<DeviceRequest>(`/auth/device/request/${encodeURIComponent(formatCode(code))}`)
            .then(({ data }) => setRequest(data))
            .catch(() => {
                setRequest(null);
                setMessage("This code is invalid or has expired.");
            })
            .finally(() => setBusy(false));
    }, [authenticated, code]);

    const signIn = () => {
        const returnTo = `${window.location.origin}/device`;
        window.location.href = `${api.defaults.baseURL}/auth/google?return_to=${encodeURIComponent(returnTo)}`;
    };

    const decide = async (approve: boolean) => {
        setBusy(true);
        try {
            await api.post(approve ? "/auth/device/approve" : "/auth/device/deny", { user_code: formatCode(code) });
            setRequest(null);
            setMessage(approve ? "Fuse CLI has been authorized. You can close this page." : "Authorization was denied.");
        } catch {
            setMessage("The authorization request could not be updated.");
        } finally {
            setBusy(false);
        }
    };

    return (
        <main className="flex min-h-screen items-center justify-center bg-background px-6 text-foreground">
            <section className="w-full max-w-sm">
                <h1 className="text-2xl font-semibold tracking-tight">Authorize Fuse CLI</h1>
                <p className="mt-2 text-sm text-muted-foreground">
                    Enter the code displayed in your terminal.
                </p>

                <InputOTP
                    aria-label="Device authorization code"
                    autoFocus
                    className="sr-only"
                    containerClassName="mt-8 justify-center"
                    maxLength={8}
                    onChange={(value) => {
                        setCode(normalizeCode(value));
                        setRequest(null);
                        setMessage("");
                    }}
                    value={code}
                >
                    <InputOTP.Group>
                        <InputOTP.Slot index={0} />
                        <InputOTP.Slot index={1} />
                        <InputOTP.Slot index={2} />
                        <InputOTP.Slot index={3} />
                    </InputOTP.Group>
                    <InputOTP.Separator />
                    <InputOTP.Group>
                        <InputOTP.Slot index={4} />
                        <InputOTP.Slot index={5} />
                        <InputOTP.Slot index={6} />
                        <InputOTP.Slot index={7} />
                    </InputOTP.Group>
                </InputOTP>

                {authenticated === false && (
                    <Button className="mt-5 w-full" onClick={signIn} size="lg">
                        Sign in with Google
                    </Button>
                )}

                {request && (
                    <div className="mt-6">
                        <p className="text-sm">
                            <span className="font-medium">{request.client_name}</span> wants permission to register and manage your compute nodes.
                        </p>
                        <div className="mt-6 flex justify-end gap-2">
                            <Button disabled={busy} onClick={() => decide(false)} variant="ghost">
                                Deny
                            </Button>
                            <Button disabled={busy} onClick={() => decide(true)}>
                                Authorize
                            </Button>
                        </div>
                    </div>
                )}

                {busy && <p className="mt-5 text-sm text-muted-foreground">Checking authorization…</p>}
                {message && <p className="mt-5 text-sm">{message}</p>}
            </section>
        </main>
    );
};

export const Component = DevicePage;
