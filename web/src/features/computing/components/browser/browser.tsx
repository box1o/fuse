import * as React from "react";

import { BrowserControlLeft } from "./browser-comtrolls-left";
import { BrowserControlRight } from "./browser-comtrolls-right";
import { BrowserViewer } from "./browser-viewer";
import { ModeSelector, type BrowserMode } from "./mode-selector";
import { BrowserControls } from "./browser-controlls";

interface BrowserOverlayProps {
    mode: BrowserMode;
    onModeChange: (mode: BrowserMode) => void;
}

const BrowserOverlay: React.FC<BrowserOverlayProps> = ({
    mode,
    onModeChange,
}) => {
    return (
        <div className="pointer-events-none absolute inset-0">
            <div className="pointer-events-auto absolute right-4 top-4">
                <ModeSelector
                    value={mode}
                    onValueChange={onModeChange}
                />
            </div>

            <BrowserControls
                height={72}
                className="
                    pointer-events-auto
                    absolute inset-x-0 bottom-0
                    bg-gradient-to-t
                    from-black/80
                    via-black/40
                    to-transparent
                    px-4
                    text-white
                    opacity-0
                    backdrop-blur-sm
                    transition-opacity
                    duration-300
                    group-hover:opacity-100
                    group-focus-within:opacity-100
                "
                left={<BrowserControlLeft />}
                right={<BrowserControlRight />}
            />
        </div>
    );
};


const BrowserRDP: React.FC = () => {
    return (
        <div className="flex h-full w-full items-center justify-center bg-card">
            RDP viewer
        </div>
    );
};

const BrowserVideo: React.FC = () => {
    return (
        <div className="flex h-full w-full items-center justify-center bg-card">
            Video viewer
        </div>
    );
};

const Browser: React.FC = () => {
    const [mode, setMode] = React.useState<BrowserMode>("rdp");

    return (
        <div className="group relative aspect-video w-full overflow-hidden rounded-xl bg-card">
            <BrowserViewer
                overlay={
                    <BrowserOverlay
                        mode={mode}
                        onModeChange={setMode}
                    />
                }
            >
                {mode === "rdp" ? <BrowserRDP /> : <BrowserVideo />}
            </BrowserViewer>
        </div>
    );
};

export { Browser };
