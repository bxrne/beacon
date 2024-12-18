document.addEventListener("DOMContentLoaded", async () => {
	const deviceSelect = document.getElementById("deviceSelect");
	const gaugeContainer = document.getElementById("gaugeContainer");
	const colorContainer = document.getElementById("colorContainer");
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

	function updateMetricsDisplay(metrics) {
		gaugeContainer.innerHTML = "";
		colorContainer.innerHTML = "";

		if (!Array.isArray(metrics)) {
			console.error("Invalid metrics data received");
			return;
		}

		metrics.forEach((metric) => {
			if (!metric.Unit || !metric.Type) return;
			if (metric.Unit.Name === "percent") {
				// Display percent metrics as gauges
				const gauge = document.createElement("div");
				gauge.innerHTML = `
					<div class="gauge">
						<div class="gauge-value" style="width: ${metric.Value}%;">
							${metric.Value}%
						</div>
					</div>
					<div class="gauge-label">${metric.Type.Name}</div>
				`;
				gaugeContainer.appendChild(gauge);
			} else if (metric.Unit.Name === "color") {
				// Display color metrics as colored divs
				const colorWrapper = document.createElement("div");
				colorWrapper.className = "color-wrapper";
				colorWrapper.innerHTML = `
					<div class="color-value" style="background-color: ${metric.Value};"></div>
					<div class="color-label">${metric.Type.Name}</div>
				`;
				colorContainer.appendChild(colorWrapper);
			} else {
				// warning div to say no display for this metric
				const warning = document.createElement("div");
				warning.className = "warning";
				warning.innerHTML = `
					<div class="warning-label">No display for this metric</div>
				`;
				gaugeContainer.appendChild(warning);
			}
		});
	}

	async function fetchMetrics(deviceID) {
		const response = await fetch(`/api/metrics?view=charts`, {
			headers: {
				"X-DeviceID": deviceID,
			},
		});
		const data = await response.json();
		if (data && Array.isArray(data.metrics)) {
			updateMetricsDisplay(data.metrics);
		} else {
			console.error("Expected metrics to be an array", data);
		}
	}

	function filterMetrics() {
		const filterText = metricTypeFilter.value.toLowerCase();
		const gaugeElements = gaugeContainer.querySelectorAll(".gauge");
		const colorElements = colorContainer.querySelectorAll(".color-wrapper");

		gaugeElements.forEach((element) => {
			const label = element
				.querySelector(".gauge-label")
				.textContent.toLowerCase();
			element.style.display = label.includes(filterText) ? "" : "none";
		});

		colorElements.forEach((element) => {
			const label = element
				.querySelector(".color-label")
				.textContent.toLowerCase();
			element.style.display = label.includes(filterText) ? "" : "none";
		});
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
			clearMetrics();
			clearInterval(refreshIntervalId);
		}
	});

	metricTypeFilter.addEventListener("input", filterMetrics);

	function clearMetrics() {
		gaugeContainer.innerHTML = "";
		colorContainer.innerHTML = "";
	}
});
