import * as React from "react";

interface BrowserViewerProps extends React.HTMLAttributes<HTMLDivElement> {
    children?: React.ReactNode;
    overlay?: React.ReactNode;
    overlayClassName?: string;
}

const BrowserViewer: React.FC<BrowserViewerProps> = ({
    children,
    overlay,
    className,
    overlayClassName,
    ...props
}) => {
    return (
        <div
            className={`relative h-full w-full overflow-hidden ${className ?? ""}`}
            {...props}
        >
            <div className="h-full w-full">{children}</div>

            {overlay && (
                <div
                    className={`absolute inset-0 h-full w-full ${overlayClassName ?? ""}`}
                >
                    {overlay}
                </div>
            )}
        </div>
    );
};

export { BrowserViewer };
