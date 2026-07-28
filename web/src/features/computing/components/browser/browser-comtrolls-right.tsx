import * as React from "react";
import { Maximize, Captions, Settings } from "lucide-react";

const BrowserControlRight: React.FC = () => {
  return (
    <div className="flex justify-between gap-4">
      <Captions />
      <Settings />
      <Maximize />
    </div>
  );
};

export { BrowserControlRight };
