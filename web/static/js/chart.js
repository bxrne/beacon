document.addEventListener("DOMContentLoaded", async () => {
	const deviceSelect = document.getElementById("deviceSelect");
	const gaugeContainer = document.getElementById("gaugeContainer");

	async function fetchDevices() {
		const response = await fetch("/api/device");
		const devices = await response.json();

		// Clear existing options
		deviceSelect.innerHTML = '<option value="">Select a device</option>';

		// Use a Set to avoid duplicate devices
		const uniqueDevices = new Set(devices);
		uniqueDevices.forEach((device) => {
			const option = document.createElement("option");
			option.value = device;
			option.textContent = device;
			deviceSelect.appendChild(option);
		});
	}

	async function fetchMetrics(deviceID) {
		const response = await fetch(`/api/metrics`, {
			headers: {
				"X-DeviceID": deviceID,
			},
		});
		const metrics = await response.json();
		initializeGauges(metrics);
	}

	function initializeGauges(metrics) {
		gaugeContainer.innerHTML = "";
		metrics.forEach((metric) => {
			if (metric.Unit && metric.Unit.Name === "percent") {
				const gauge = document.createElement("div");
				gauge.className = "gauge";
				gauge.innerHTML = `
					<div class="gauge-value" style="width: ${metric.Value}%;">${metric.Value}%</div>
					<div class="gauge-label">${metric.Type ? metric.Type.Name : ""}</div>
				`;
				gaugeContainer.appendChild(gauge);
			}
		});
	}

	await fetchDevices();

	deviceSelect.addEventListener("change", function () {
		const deviceID = this.value;
		if (deviceID) {
			fetchMetrics(deviceID);
		} else {
			gaugeContainer.innerHTML = "";
		}
	});
});
