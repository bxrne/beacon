import axios from "axios";
import { API_CONFIG } from "../config/api";
import type { DeviceMetrics } from "../types/api";

const api = axios.create({
	baseURL: API_CONFIG.baseUrl,
});

export const getDevices = async (): Promise<string[]> => {
	const response = await api.get<string[]>("/device");
	return response.data;
};

export const getMetrics = async (deviceId: string): Promise<DeviceMetrics> => {
	const response = await api.get<DeviceMetrics>("/metric", {
		headers: {
			"X-DeviceID": deviceId,
		},
	});
	return response.data;
};
