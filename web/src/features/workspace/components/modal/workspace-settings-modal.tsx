import { useState } from "react";
import { Settings, Users } from "lucide-react";

import { Button, Dialog } from "@/shared/components";
import { cn } from "@/shared/utils";

import { WorkspaceSidebar, type WorkspaceSidebarItem } from "./sidebar";
import { WorkspaceMembers } from "../members";

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

interface WorkspaceSettingsModalProps {
  className?: string;
  isOpen?: boolean;
  onOpenChange?: (isOpen: boolean) => void;
  selectedImplicitSection?: WorkspaceSettingsSection;
}

const WorkspaceSettingsModal = ({
  className,
  isOpen,
  onOpenChange,
  selectedImplicitSection,
}: WorkspaceSettingsModalProps) => {
  const [selectedSection, setSelectedSection] =
    useState<WorkspaceSettingsSection>(selectedImplicitSection || "members");

  const handleSectionSelect = (sectionId: string) => {
    if (sectionId === "members" || sectionId === "general") {
      setSelectedSection(sectionId);
    }
  };

  return (
    <Dialog open={isOpen} onOpenChange={onOpenChange}>
      <Dialog.Content
        showCloseButton={false}
        className={cn(
            "!h-[75vh] !w-[75vw] !max-w-[75vw] !rounded-2xl p-2 sm:!h-[75vh] sm:!w-[75vw] sm:!max-w-[75vw] overflow-auto",
          className,
        )}
      >
        <WorkspaceSidebar
          items={sidebarItems}
          selectedItemId={selectedSection}
          onItemSelect={handleSectionSelect}
        >
          {selectedSection === "members" && <WorkspaceMembers />}
          {selectedSection === "general" && <div />}
        </WorkspaceSidebar>
      </Dialog.Content>
    </Dialog>
  );
};

export { WorkspaceSettingsModal };
