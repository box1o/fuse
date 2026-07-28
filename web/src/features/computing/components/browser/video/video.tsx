import { BrowserControlsLeft } from "../browser-controls-left";
import { BrowserControlsRight } from "../browser-controls-right";
import { BrowserControls } from "../browser-controls";

const VideoViewer = () => {
    return (
        <div className="relative h-full w-full bg-card">
            <div className="flex h-full items-center justify-center">
                Video viewer
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
                left={<BrowserControlsLeft />}
                right={<BrowserControlsRight />}
            />
        </div>
    );
};

export { VideoViewer };