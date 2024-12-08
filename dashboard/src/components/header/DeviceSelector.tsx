import React from "react";
import { useDeviceStore } from "../../stores/deviceStore";

export const DeviceSelector: React.FC = () => {
	const { devices, selectedDevice, setSelectedDevice } = useDeviceStore();

	const handleDeviceChange = (event: React.ChangeEvent<HTMLSelectElement>) => {
		const deviceId = event.target.value;
		setSelectedDevice(deviceId);
	};

	return (
		<select
			className="block w-48 px-4 py-2 rounded-lg border border-gray-300 bg-white focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
			value={selectedDevice || ""}
			onChange={handleDeviceChange}
		>
			<option value="">Select device</option>
			{devices.map((device) => (
				<option key={device} value={device}>
					{device}
				</option>
			))}
		</select>
	);
};
