import { useEffect } from "react";
import { Toaster } from "react-hot-toast";
import { Sidebar } from "./components/layout/Sidebar";
import { DeviceSelector } from "./components/header/DeviceSelector";
import { MetricsGrid } from "./components/metrics/MetricsGrid";
import { useDeviceStore } from "./stores/deviceStore";

function App() {
	const { fetchDevices } = useDeviceStore();

	useEffect(() => {
		fetchDevices();
	}, [fetchDevices]);

	return (
		<div className="min-h-screen bg-gray-50">
			<Sidebar />
			<Toaster position="top-right" />

			<main className="ml-64 p-8">
				<div className="mb-8 flex justify-between items-center">
					<h2 className="text-2xl font-bold text-gray-900">Dashboard</h2>
					<DeviceSelector />
				</div>

				<MetricsGrid />
			</main>
		</div>
	);
}

export default App;
