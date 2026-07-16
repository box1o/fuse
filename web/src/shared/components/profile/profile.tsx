import * as React from "react";
import { User, Settings, LogOut, UserPlus, Badge, Zap } from "lucide-react";
import { Button } from "@/shared/components/ui/button";
import {
    DropdownMenu,
} from "@/shared/components/ui";
import { Avatar } from "@/shared/components/ui/avatar";
import { useAuthActions, useAuthStore } from "@/features/auth";
import { getInitials } from "@/shared/utils";
import { useMockBillingStore } from "@/features/payments/store/mock-billing.store";
import { CreditUsageBar } from "@/features/payments/components/credit-usage-bar";
import { CreditsButton } from "@/features/payments";
import { Separator } from "@radix-ui/react-dropdown-menu";


const Profile = () => {

    const { logout } = useAuthActions();
    const user = useAuthStore(state => state.user) || {
        id: "guest",
        avatar: 'https://github.com/shadcn.png',
        name: 'Guest User',
        email: 'guest@gmail.com'
    }
    const setUserKey = useMockBillingStore((state) => state.setUserKey);
    const planId = useMockBillingStore((state) => state.planId);
    const usedCredits = useMockBillingStore((state) => state.usedCredits);
    const includedCredits = useMockBillingStore((state) => state.includedCredits);
    const planStatus = planId === "pro" ? "Active" : "Free";
    const resetAt = new Date().toLocaleDateString("en-US", {
        month: "short",
        day: "numeric",
        year: "numeric",
    });

    React.useEffect(() => {
        setUserKey(user?.id ?? user?.email ?? null);
    }, [setUserKey, user?.email, user?.id]);

    return (
        <DropdownMenu>
            <DropdownMenu.Trigger asChild>
                <Button variant="ghost" size="icon" className="rounded-full">
                    <Avatar className="h-6 w-6">
                        <Avatar.Image src={user.avatar} />
                        <Avatar.Fallback className="text-xs">
                            {getInitials(user.name)}
                        </Avatar.Fallback>
                    </Avatar>
                </Button>
            </DropdownMenu.Trigger>
            <DropdownMenu.Content align="end" className="w-[14rem] p-2">
                <DropdownMenu.Label className="font-normal">
                    <div className="flex items-center gap-3 px-2 py-1.5">
                        <Avatar className="h-6 w-6">
                            <Avatar.Image src={user.avatar} />
                            <Avatar.Fallback className="text-sm">
                                {getInitials(user.name)}
                            </Avatar.Fallback>
                        </Avatar>
                        <div className="flex flex-col space-y-1 min-w-0">
                            <p className="text-xs font-medium leading-none truncate">
                                {user.name}
                            </p>
                            <p className="text-xs leading-none text-muted-foreground truncate">
                                {user.email}
                            </p>
                        </div>
                    </div>
                </DropdownMenu.Label>

                <DropdownMenu.Separator/>
                <DropdownMenu.Group>
                    <DropdownMenu.Item>
                        <User className="mr-2 h-4 w-4" />
                        <span>Projects</span>
                    </DropdownMenu.Item>
                    <DropdownMenu.Item>
                        <Settings className="mr-2 h-4 w-4" />
                        <span>Settings</span>
                    </DropdownMenu.Item>
                    <DropdownMenu.Item>
                        <UserPlus className="mr-2 h-4 w-4" />
                        <span>Share</span>
                    </DropdownMenu.Item>
                </DropdownMenu.Group>
                <DropdownMenu.Item
                variant="branded"
                >
                    <Zap className="mr-2 h-4 w-4" />
                    <span>Subscription</span>
                </DropdownMenu.Item>

                <DropdownMenu.Separator/>
            <DropdownMenu.Item
                    onClick={() => logout()}
                    variant="destructive"
                >
                    <LogOut className="mr-2 h-4 w-4" />
                    <span>Log out</span>
                </DropdownMenu.Item>
            </DropdownMenu.Content>
        </DropdownMenu>
    );
}


export { Profile };
