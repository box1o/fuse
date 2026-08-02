import { useState } from "react";
import { Settings, Users } from "lucide-react";

import { Button, Dialog } from "@/shared/components";
import { cn } from "@/shared/utils";

import {
    WorkspaceSidebar,
    type WorkspaceSidebarItem,
} from "./sidebar";
import { WorkspaceMembers } from "../members";

interface WorkspaceSettingsModalProps {
    className?: string;
}

type WorkspaceSettingsSection = "members" | "general";

const sidebarItems: WorkspaceSidebarItem[] = [
    {
        id: "members",
        label: "Members",
        icon: <Users className="size-4 shrink-0" />,
    },
    {
        id: "general",
        label: "General",
        icon: <Settings className="size-4 shrink-0" />,
    },
];

const WorkspaceSettingsModal = ({
    className,
}: WorkspaceSettingsModalProps) => {
    const [selectedSection, setSelectedSection] =
        useState<WorkspaceSettingsSection>("members");

    const handleSectionSelect = (sectionId: string) => {
        if (sectionId === "members" || sectionId === "general") {
            setSelectedSection(sectionId);
        }
    };

    return (
        <Dialog>
            <Dialog.Trigger asChild>
                <Button
                    variant="outline"
                    className={cn(
                        "h-8 rounded-md bg-brand/35",
                        className,
                    )}
                >
                    New Workspace Member
                </Button>
            </Dialog.Trigger>

            <Dialog.Content
                showCloseButton={false}
                className="!h-[75vh] !w-[75vw] !max-w-[75vw] !rounded-2xl p-2 sm:!h-[75vh] sm:!w-[75vw] sm:!max-w-[75vw]"
            >
                <WorkspaceSidebar
                    items={sidebarItems}
                    selectedItemId={selectedSection}
                    onItemSelect={handleSectionSelect}
                >
                    {selectedSection === "members" && (
                        <WorkspaceMembers/>
                    )}

                    {selectedSection === "general" && (
                        <div></div>
                    )}
                </WorkspaceSidebar>
            </Dialog.Content>
        </Dialog>
    );
};

export { WorkspaceSettingsModal };