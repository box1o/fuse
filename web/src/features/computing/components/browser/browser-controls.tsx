import * as React from "react";

interface BrowserControlsProps extends React.HTMLAttributes<HTMLDivElement> {
  height?: number | string;
  left?: React.ReactNode;
  right?: React.ReactNode;
}

const BrowserControls: React.FC<BrowserControlsProps> = ({
  className,
  height,
  style,
  left,
  right,
  ...props
}) => {
  return (
    <div
      className={`flex w-full items-end justify-between pb-3 ${className ?? ""}`}
      style={{
        ...style,
        height,
      }}
      {...props}
    >
      <div className="flex items-center gap-3">{left}</div>

      <div className="flex items-center gap-3">{right}</div>
    </div>
  );
};

export { BrowserControls };
