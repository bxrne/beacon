import React from "react";
import { Activity, LayoutDashboard, Settings } from "lucide-react";

export const Sidebar: React.FC = () => {
	return (
		<div className="fixed left-0 top-0 h-full w-64 bg-gray-900 text-white p-4">
			<div className="flex items-center gap-3 mb-8">
				<Activity className="h-8 w-8 text-blue-400" />
				<h1 className="text-xl font-bold">Beacon</h1>
			</div>

			<nav className="space-y-2">
				<a
					href="#"
					className="flex items-center gap-3 px-4 py-2 rounded-lg bg-gray-800 text-blue-400"
				>
					<LayoutDashboard className="h-5 w-5" />
					Dashboard
				</a>
				<a
					href="#"
					className="flex items-center gap-3 px-4 py-2 rounded-lg hover:bg-gray-800 transition-colors"
				>
					<Settings className="h-5 w-5" />
					Settings
				</a>
			</nav>
		</div>
	);
};
