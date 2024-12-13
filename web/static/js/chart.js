document.addEventListener("DOMContentLoaded", async () => {
	const deviceSelect = document.getElementById("deviceSelect");
	const gaugeContainer = document.getElementById("gaugeContainer");
	const metricTypeFilter = document.getElementById("metricTypeFilter");
	let refreshIntervalId = null;

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
		const data = await response.json();
		if (data && Array.isArray(data.metrics)) {
			updateGauges(data.metrics);
		} else {
			console.error("Expected metrics to be an array", data);
		}
	}

	function updateGauges(metrics) {
		gaugeContainer.innerHTML = "";
		const latestMetrics = {};
		metrics.forEach((metric) => {
			latestMetrics[metric.Type.Name] = metric;
		});
		Object.values(latestMetrics).forEach((metric) => {
			const gauge = document.createElement("div");
			gauge.className = "gauge";
			gauge.innerHTML = `
				<div class="gauge-value" style="width: ${metric.Value}%;">${metric.Value}%</div>
				<div class="gauge-label">${metric.Type ? metric.Type.Name : ""}</div>
			`;
			gaugeContainer.appendChild(gauge);
		});
	}

	function filterGauges() {
		const filterText = metricTypeFilter.value.toLowerCase();
		const gauges = gaugeContainer.getElementsByClassName("gauge");
		for (let gauge of gauges) {
			const label = gauge
				.getElementsByClassName("gauge-label")[0]
				.textContent.toLowerCase();
			if (label.includes(filterText)) {
				gauge.style.display = "";
			} else {
				gauge.style.display = "none";
			}
		}
	}

	function startAutoRefresh(deviceID) {
		const interval = 5000; // 5 seconds
		clearInterval(refreshIntervalId);
		refreshIntervalId = setInterval(() => {
			fetchMetrics(deviceID);
		}, interval);
	}

	await fetchDevices();

	deviceSelect.addEventListener("change", function () {
		const deviceID = this.value;
		if (deviceID) {
			fetchMetrics(deviceID);
			startAutoRefresh(deviceID);
		} else {
			gaugeContainer.innerHTML = "";
			clearInterval(refreshIntervalId);
		}
	});

	metricTypeFilter.addEventListener("input", filterGauges);
});
