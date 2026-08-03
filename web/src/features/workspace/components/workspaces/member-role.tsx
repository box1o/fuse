import { Check, ChevronDown } from "lucide-react";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@radix-ui/react-dropdown-menu";

import { Button } from "@/shared/components";
import { cn } from "@/shared/utils";

type MemberRoleSection = "owner" | "admin" | "member";

interface MemberRoleProps {
  className?: string;
  memberRole?: MemberRoleSection;
  setMemberRole?: (role: MemberRoleSection) => void;
}

const roleOptions: Array<{
  value: MemberRoleSection;
  label: string;
}> = [
  {
    value: "owner",
    label: "Owner",
  },
  {
    value: "admin",
    label: "Admin",
  },
  {
    value: "member",
    label: "Member",
  },
];

const MemberRole = ({
  className,
  memberRole = "member",
  setMemberRole,
}: MemberRoleProps) => {
  const selectedRole =
    roleOptions.find((role) => role.value === memberRole) ?? roleOptions[2];

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button
          type="button"
          variant="outline"
          className={cn("w-full justify-between", className)}
        >
          <span>{selectedRole.label}</span>
          <ChevronDown className="size-4 opacity-70" />
        </Button>
      </DropdownMenuTrigger>

      <DropdownMenuContent
        align="start"
        className="z-[70] min-w-[180px] rounded-md border bg-popover p-1 shadow-md"
      >
        {roleOptions.map((role) => {
          const isSelected = role.value === memberRole;

          return (
            <DropdownMenuItem
              key={role.value}
              onSelect={() => setMemberRole?.(role.value)}
              className="flex cursor-pointer items-center gap-2 rounded-sm px-2 py-1.5 text-sm outline-none focus:bg-accent"
            >
              <span>{role.label}</span>

              {isSelected && <Check className="ml-auto size-4" />}
            </DropdownMenuItem>
          );
        })}
      </DropdownMenuContent>
    </DropdownMenu>
  );
};

export { MemberRole, type MemberRoleSection };
