import React from "react";
import { useDeviceStore } from "../stores/deviceStore";

export const DeviceSelector: React.FC = () => {
	const { devices, selectedDevice, setSelectedDevice } = useDeviceStore();

	const handleDeviceChange = (event: React.ChangeEvent<HTMLSelectElement>) => {
		const deviceId = event.target.value;
		setSelectedDevice(deviceId);
	};

	return (
		<div className="w-full max-w-xs">
			<select
				className="block w-full px-4 py-2 rounded-lg border border-gray-300 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
				value={selectedDevice || ""}
				onChange={handleDeviceChange}
			>
				<option value="">Select a device</option>
				{devices.map((device) => (
					<option key={device} value={device}>
						{device}
					</option>
				))}
			</select>
		</div>
	);
};
