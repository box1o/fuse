import { authService } from "@/features/auth/services";
import { useCallback, useEffect, useState } from "react";
import { useSearchParams } from "react-router-dom";
import { deviceService } from "../services";
import type { DeviceAuthorizationRequest, DeviceAuthorizationState } from "../types";
import { formatDeviceCode, normalizeDeviceCode } from "../utils";

const useDeviceAuthorization = () => {
    const [searchParams] = useSearchParams();
    const [code, setCode] = useState(() => normalizeDeviceCode(searchParams.get("code") ?? ""));
    const [authenticated, setAuthenticated] = useState<boolean | null>(null);
    const [request, setRequest] = useState<DeviceAuthorizationRequest | null>(null);
    const [state, setState] = useState<DeviceAuthorizationState>("entering-code");

    useEffect(() => {
        let active = true;

        authService.getStatus().then((result) => {
            if (active) setAuthenticated(result.success);
        });

        return () => {
            active = false;
        };
    }, []);

    useEffect(() => {
        if (!authenticated || code.length !== 8) {
            setRequest(null);
            setState("entering-code");
            return;
        }

        let active = true;
        setState("checking-code");

        deviceService
            .getRequest(formatDeviceCode(code))
            .then((deviceRequest) => {
                if (!active) return;
                setRequest(deviceRequest);
                setState("ready");
            })
            .catch(() => {
                if (!active) return;
                setRequest(null);
                setState("invalid");
            });

        return () => {
            active = false;
        };
    }, [authenticated, code]);

    const updateCode = useCallback((value: string) => {
        setCode(normalizeDeviceCode(value));
        setRequest(null);
        setState("entering-code");
    }, []);

    const signIn = useCallback(() => {
        authService.startOAuth("google", window.location.href);
    }, []);

    const decide = useCallback(
        async (approve: boolean) => {
            if (!request) return;

            setState("submitting");
            try {
                await deviceService.decide(request.user_code, approve);
                setRequest(null);
                setState(approve ? "approved" : "denied");
            } catch {
                setState("invalid");
            }
        },
        [request],
    );

    return {
        authenticated,
        code,
        request,
        state,
        updateCode,
        signIn,
        approve: () => decide(true),
        deny: () => decide(false),
    };
};

export { useDeviceAuthorization };
