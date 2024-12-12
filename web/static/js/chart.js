import {
	fetchDevices,
	populateDeviceSelect,
	fetchMetricsForDevice,
} from "./utils.js";

document.addEventListener("DOMContentLoaded", async () => {
	const deviceSelect = document.getElementById("deviceSelect");
	const updateIntervalSelect = document.getElementById("updateInterval");
	const gaugeContainer = document.getElementById("gaugeContainer");
	let refreshIntervalId = null;

	async function initializeGauges(deviceID) {
		const metrics = await fetchMetricsForDevice(deviceID);
		gaugeContainer.innerHTML = "";
		metrics.forEach((metric) => {
			if (metric.Unit.Name === "percent") {
				const gauge = document.createElement("div");
				gauge.className = "gauge";
				gauge.innerHTML = `
					<div class="gauge-value" style="width: ${metric.Value}%;">${metric.Value}%</div>
					<div class="gauge-label">${metric.Type.Name}</div>
				`;
				gaugeContainer.appendChild(gauge);
			}
		});
	}

	function startAutoRefresh(deviceID) {
		const interval = parseInt(updateIntervalSelect.value);
		clearInterval(refreshIntervalId);
		if (interval > 0) {
			refreshIntervalId = setInterval(() => {
				initializeGauges(deviceID);
			}, interval);
		}
	}

	await fetchDevices().then((devices) =>
		populateDeviceSelect(deviceSelect, devices)
	);

	deviceSelect.addEventListener("change", function () {
		const deviceID = this.value;
		if (deviceID) {
			initializeGauges(deviceID);
			startAutoRefresh(deviceID);
		} else {
			gaugeContainer.innerHTML = "";
			clearInterval(refreshIntervalId);
		}
	});

	updateIntervalSelect.addEventListener("change", function () {
		const deviceID = deviceSelect.value;
		if (deviceID) {
			startAutoRefresh(deviceID);
		}
	});

	// Set "No update" as the default active option
	updateIntervalSelect.value = "0";
});
