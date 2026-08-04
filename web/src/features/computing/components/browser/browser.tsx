import * as React from "react";

import { BrowserViewer } from "./browser-viewer";
import { ModeSelector, type BrowserMode } from "./mode-selector";
import { ImageViewer } from "./image";
import { RdpViewer } from "./rdp";
import { VideoViewer } from "./video";



const Browser: React.FC = () => {
    const [mode, setMode] = React.useState<BrowserMode>("rdp");
    const renderViewer = () => {
        switch (mode) {
            case "rdp":
                return <RdpViewer />;

            case "video":
                return <VideoViewer />;

            case "image":
                return <ImageViewer />;

            default:
                return null;
        }
    };

    return (
        <div className="group relative aspect-video w-full overflow-hidden rounded-xl bg-card">
            <BrowserViewer
                overlay={
                    <div className="pointer-events-none absolute inset-0">
                        <div className="pointer-events-auto absolute right-4 top-4">
                            <ModeSelector
                                value={mode}
                                onValueChange={setMode}
                            />
                        </div>
                    </div>
                }
            >
                 {renderViewer()}
            </BrowserViewer>
        </div>
    );
};

export { Browser };
