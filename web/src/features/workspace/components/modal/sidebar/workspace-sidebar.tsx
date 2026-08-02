import type { ReactNode } from "react";

import { Button, Sidebar } from "@/shared/components";

export interface WorkspaceSidebarItem {
    id: string;
    label: string;
    icon?: ReactNode;
}

interface WorkspaceSidebarProps {
    items: WorkspaceSidebarItem[];
    selectedItemId: string;
    onItemSelect: (itemId: string) => void;
    children?: ReactNode;
}

const WorkspaceSidebar = ({
    items,
    selectedItemId,
    onItemSelect,
    children,
}: WorkspaceSidebarProps) => {
    return (
        <div className="flex h-full">
            <Sidebar width="max-content">
                <Sidebar.Content>
                    {items.map((item) => {
                        const isSelected = selectedItemId === item.id;

                        return (
                            <Sidebar.Item key={item.id} asChild>
                                <Button
                                    type="button"
                                    variant={
                                        isSelected ? "secondary" : "ghost"
                                    }
                                    className="w-full justify-start gap-2 whitespace-nowrap"
                                    onClick={() => onItemSelect(item.id)}
                                    aria-pressed={isSelected}
                                >
                                    {item.icon}
                                    <span>{item.label}</span>
                                </Button>
                            </Sidebar.Item>
                        );
                    })}
                </Sidebar.Content>
            </Sidebar>

            <Sidebar.Inset className="border-0 py-2">
                {children}
            </Sidebar.Inset>
        </div>
    );
};

export { WorkspaceSidebar };