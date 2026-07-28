import * as React from "react";
import { Play } from "lucide-react";

const BrowserControlLeft: React.FC = () => {
  return (
      <div className="flex justify-between gap-4">
            <Play />
            <span> 24:40 / 40:234 </span>
      </div>
  );
};

export { BrowserControlLeft };
