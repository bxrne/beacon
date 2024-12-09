import { create } from "zustand";
import { getDevices, getMetrics } from "../services/api";
import toast from "react-hot-toast";
import type { DeviceMetrics } from "../types/api";

interface DeviceStore {
	devices: string[];
	selectedDevice: string;
	metrics: DeviceMetrics | null;
	loading: boolean;
	error: string | null;
	fetchDevices: () => Promise<void>;
	setSelectedDevice: (deviceId: string) => void;
	fetchMetrics: (deviceId: string) => Promise<void>;
}

export const useDeviceStore = create<DeviceStore>((set) => ({
	devices: [],
	selectedDevice: "",
	metrics: null,
	loading: false,
	error: null,
	fetchDevices: async () => {
		set({ loading: true, error: null });
		try {
			const devices = await getDevices();
			set({ devices, loading: false });
			toast.success("Devices loaded successfully");
		} catch {
			const message = "Failed to fetch devices";
			set({ error: message, loading: false });
			toast.error(message);
		}
	},
	setSelectedDevice: (deviceId: string) => {
		set({ selectedDevice: deviceId });
		if (deviceId) {
			set({ loading: true, error: null });
			getMetrics(deviceId)
				.then((metrics) => set({ metrics, loading: false }))
				.catch(() => {
					const message = "Failed to fetch metrics";
					set({ error: message, loading: false });
					toast.error(message);
				});
		}
	},
	fetchMetrics: async (deviceId: string) => {
		set({ loading: true, error: null });
		try {
			const metrics = await getMetrics(deviceId);
			set({ metrics, loading: false });
		} catch {
			const message = "Failed to fetch metrics";
			set({ error: message, loading: false });
			toast.error(message);
		}
	},
}));
