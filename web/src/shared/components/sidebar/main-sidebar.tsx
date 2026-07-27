import { Sidebar } from "@/shared/components/ui/sidebar";
import { Button } from "@/shared/components/ui/button";
import { Computer, Folder, Layers} from "lucide-react";
import React from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { ROUTES } from "@/shared/constants";

interface Route {
    name: string;
    path: string;
    icon: React.ReactNode;
}

const routes: Route[] = [
    { name: "Projects", path: ROUTES.PROJECTS, icon: <Folder size={8} /> },
    { name: "Workspace", path: ROUTES.WORKSPACE, icon: <Layers size={8} /> },
    { name: "Compute", path: ROUTES.COMPUTE, icon: <Computer size={8} /> },
];

const MainSidebar = () => {
    const location = useLocation();
    const navigate = useNavigate();

    return (
        <Sidebar.Content>
            {routes.map((route) => {
                const isActive = location.pathname === route.path;
                return (
                    <Sidebar.Item key={route.name} asChild>
                        <Button
                            variant={isActive ? "secondary" : "ghost"}
                            onClick={() => navigate(route.path)}
                            size="icon"
                        >
                            {route.icon}
                        </Button>
                    </Sidebar.Item>
                );
            })}
        </Sidebar.Content>
    );
};

export { MainSidebar };
