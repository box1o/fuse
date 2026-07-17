import React from "react";
// import { Badge, Button } from "@/shared/components";
// import Workspaces from "./components/workspaces/workspaces";
// import type { Workspace } from "./types";
// import { useWorkspaceStore } from "./store";
// import { useAuthStore } from "@/features/auth";


const Main: React.FC = () => {
    // const setCurrentWs = useWorkspaceStore((store) => store.setCurrentWorkspace);
    // const selectedWorkspace = useWorkspaceStore((store) => store.currentWorkspace);
    // const user = useAuthStore((state) => state.user);
    // const setUserKey = useMockBillingStore((state) => state.setUserKey);
    // const planId = useMockBillingStore((state) => state.planId);
    // const usedCredits = useMockBillingStore((state) => state.usedCredits);
    // const includedCredits = useMockBillingStore((state) => state.includedCredits);
    // const spendCredit = useMockBillingStore((state) => state.spendCredit);
    // const remainingCredits = Math.max(includedCredits - usedCredits, 0);

    // const handleRowClick = React.useCallback((workspace: Workspace) => {
    //     setCurrentWs(workspace);

    // }, [setCurrentWs]);

    // React.useEffect(() => {
    //     setUserKey(user?.id ?? user?.email ?? null);
    // }, [setUserKey, user?.email, user?.id]);

    // const handleMockSpend = React.useCallback((event: React.MouseEvent<HTMLDivElement>) => {
    //     const target = event.target as HTMLElement | null;
    //     const button = target?.closest("button");

    //     if (!button || button.hasAttribute("disabled") || button.getAttribute("aria-disabled") === "true") {
    //         return;
    //     }

    //     spendCredit(1);
    // }, [spendCredit]);

    return (
        <div className="flex h-full w-full items-center justify-center p-6">
            <h1 className="text-2xl font-semibold">
                Welcome to the Workspace Page
            </h1>
        </div>
        // <div className="w-full h-full" onClickCapture={handleMockSpend}>
        //     <div className="mx-auto mb-4 flex w-full max-w-4xl items-center justify-between rounded-2xl border bg-background/80 px-4 py-3 shadow-sm">
        //         <div>
        //             <div className="flex items-center gap-2">
        //                 <h2 className="text-sm font-semibold">Mock billing</h2>
        //                 <Badge variant="secondary">
        //                     {planId === "pro" ? "Pro" : "Free"}
        //                 </Badge>
        //             </div>
        //             <p className="text-sm text-muted-foreground">
        //                 Selected workspace: {selectedWorkspace?.name ?? "none"}
        //             </p>
        //         </div>

        //         <div className="flex items-center gap-2">
        //             <div className="text-right">
        //                 <div className="text-sm font-semibold">
        //                     {remainingCredits.toLocaleString()} credits left
        //                 </div>
        //                 <div className="text-xs text-muted-foreground">
        //                     {usedCredits.toLocaleString()} / {includedCredits.toLocaleString()} used
        //                 </div>
        //             </div>

        //             <Button type="button" variant="outline" size="sm">
        //                 Spend credit
        //             </Button>
        //         </div>
        //     </div>

        //     <Workspaces onRowClick={handleRowClick} />
        // </div>


    );
};

export default Main;
