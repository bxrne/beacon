import { LayoutDashboard, Monitor, Gauge } from "lucide-react";
import { Link, useLocation } from "react-router-dom";
import { cn } from "@/lib/utils";

const navigation = [
  { name: 'Hosts', href: '/', icon: Monitor },
  { name: 'Traffic', href: '/traffic', icon: Gauge },
];

export function Navbar() {
  const location = useLocation();

  return (
    <div className="border-b">
      <div className="flex h-16 items-center px-4">
        <LayoutDashboard className="mr-2 h-6 w-6" />
        <h2 className="text-lg font-semibold mr-8">Beacon Dashboard</h2>
        <nav className="flex space-x-4">
          {navigation.map((item) => {
            const Icon = item.icon;
            return (
              <Link
                key={item.name}
                to={item.href}
                className={cn(
                  "inline-flex items-center px-3 py-2 text-sm font-medium rounded-md",
                  location.pathname === item.href
                    ? "bg-secondary text-secondary-foreground"
                    : "text-muted-foreground hover:text-primary hover:bg-secondary/50"
                )}
              >
                <Icon className="mr-2 h-4 w-4" />
                {item.name}
              </Link>
            );
          })}
        </nav>
      </div>
    </div>
  );
}