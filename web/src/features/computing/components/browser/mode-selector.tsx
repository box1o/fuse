import * as React from "react";
import { Check, ChevronDown, Monitor, Video, ImageIcon } from "lucide-react";

import { Button } from "@/shared/components/ui/button";
import { DropdownMenu } from "@/shared/components/ui";

type BrowserMode = "video" | "rdp" | "image";

interface ModeSelectorProps {
    value: BrowserMode;
    onValueChange: (value: BrowserMode) => void;
    className?: string;
}

const MODE_OPTIONS: Array<{
    value: BrowserMode;
    label: string;
    icon: React.ComponentType<React.SVGProps<SVGSVGElement>>;
}> = [
    {
        value: "video",
        label: "Video",
        icon: Video,
    },
    {
        value: "rdp",
        label: "RDP",
        icon: Monitor,
    },
    {
        value: "image",
        label: "Image",
        icon: ImageIcon,
    }
];

const ModeSelector: React.FC<ModeSelectorProps> = ({
    value,
    onValueChange,
    // className,
}) => {
    const selectedMode =
        MODE_OPTIONS.find((option) => option.value === value) ?? MODE_OPTIONS[0];

    const SelectedIcon = selectedMode.icon;

    return (
        <DropdownMenu>
            <DropdownMenu.Trigger asChild>
                <Button
                    type="button"
                    variant="ghost"
                >
                    <SelectedIcon className="h-4 w-4" />

                    <span>{selectedMode.label}</span>

                    <ChevronDown className="h-4 w-4 opacity-70" />
                </Button>
            </DropdownMenu.Trigger>

            <DropdownMenu.Content
                className="min-w-36"
            >
                {MODE_OPTIONS.map((option) => {
                    const Icon = option.icon;
                    const isSelected = option.value === value;

                    return (
                        <DropdownMenu.Item
                            key={option.value}
                            onSelect={() => onValueChange(option.value)}
                            className="flex items-center gap-2"
                        >
                            <Icon className="h-4 w-4" />

                            <span>{option.label}</span>

                            {isSelected && (
                                <Check className="ml-auto h-4 w-4" />
                            )}
                        </DropdownMenu.Item>
                    );
                })}
            </DropdownMenu.Content>
        </DropdownMenu>
    );
};

export type { BrowserMode };
export { ModeSelector };
