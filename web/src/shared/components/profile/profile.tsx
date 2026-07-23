import { useState } from "react";
import {
    Coins,
    CreditCard,
    LogOut,
    Settings,
    User,
    UserPlus,
} from "lucide-react";

import { useAuthActions, useAuthStore } from "@/features/auth";
import {
    CreditPurchaseModal,
    useCreditBalance,
} from "@/features/payments";
import { Avatar } from "@/shared/components/ui/avatar";
import { Button } from "@/shared/components/ui/button";
import { DropdownMenu } from "@/shared/components/ui";
import { getInitials } from "@/shared/utils";

const Profile = () => {
    const [isProfileMenuOpen, setIsProfileMenuOpen] = useState(false);
    const [isCreditModalOpen, setIsCreditModalOpen] = useState(false);

    const { logout } = useAuthActions();
    //NOTE: Fallback user in case the is null
    const user = useAuthStore(state => state.user) || {
        avatar: 'https://github.com/shadcn.png',
        name: 'Guest User',
        email: 'guest@gmail.com'
    }

    const handleOpenCreditModal = () => {
        setIsProfileMenuOpen(false);

        // Radix must finish closing the dropdown before opening
        // another modal layer.
        window.setTimeout(() => {
            setIsCreditModalOpen(true);
        }, 0);
    };

    const {balance, isLoading: isLoadingBalance,} = useCreditBalance();

    return (
        <>
        <DropdownMenu open={isProfileMenuOpen} onOpenChange={setIsProfileMenuOpen}>
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
            
            <DropdownMenu.Content align="end" className="max-w-64" sideOffset={8}>
                <DropdownMenu.Label className="font-normal">
                    <div className="flex items-center gap-3 p-2">
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
                    <div className="mx-1 my-2 rounded-xl border bg-card p-3">
                        <div className="flex items-center gap-3">
                            <div className="flex size-9 shrink-0 items-center justify-center rounded-lg bg-muted">
                                <Coins className="size-4 text-emerald-400" />
                            </div>

                            <div className="min-w-0">
                                <p className="text-xs font-medium text-muted-foreground">
                                    Balance
                                </p>

                                <div className="mt-0.5 flex items-baseline gap-1.5">
                                    <p className="truncate text-base font-semibold tracking-tight">
                                        {isLoadingBalance
                                            ? "Loading..."
                                            : balance.toLocaleString()}
                                    </p>

                                    {!isLoadingBalance && (
                                        <span className="text-xs text-muted-foreground">
                                            available
                                        </span>
                                    )}
                                </div>
                            </div>
                        </div>

                        <Button
                            type="button"
                            size="sm"
                            variant="outline"
                            className="mt-3 w-full justify-center hover:border-emerald-500/40 hover:bg-emerald-500/5"
                            onClick={handleOpenCreditModal}
                        >
                            <CreditCard className="size-4" />
                            Buy
                        </Button>
                    </div>

                <DropdownMenu.Separator />

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
                <DropdownMenu.Separator />
                <DropdownMenu.Item
                    onClick={() => logout()}
                    variant="destructive"
                >
                    <LogOut className="mr-2 h-4 w-4" />
                    <span>Log out</span>
                </DropdownMenu.Item>
            </DropdownMenu.Content>
        </DropdownMenu>
         <CreditPurchaseModal
            open={isCreditModalOpen}
            onOpenChange={setIsCreditModalOpen}
        />
        </>
    );
}


export { Profile };
