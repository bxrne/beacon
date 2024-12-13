document.addEventListener("DOMContentLoaded", async () => {
	const deviceSelect = document.getElementById("deviceSelect");
	const gaugeContainer = document.getElementById("gaugeContainer");
	let eventSource = null;

	async function fetchDevices() {
		const response = await fetch("/api/device");
		const devices = await response.json();
		devices.forEach((device) => {
			const option = document.createElement("option");
			option.value = device;
			option.textContent = device;
			deviceSelect.appendChild(option);
		});
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

	function startEventSource(deviceID) {
		if (eventSource) {
			eventSource.close();
		}
		eventSource = new EventSource(`/metrics/stream?deviceID=${deviceID}`);

		eventSource.onmessage = function (event) {
			const metrics = JSON.parse(event.data);
			initializeGauges(metrics);
		};

		eventSource.onerror = function () {
			console.error("EventSource failed. Reconnecting...");
			eventSource.close();
			setTimeout(() => startEventSource(deviceID), 5000);
		};
	}

	await fetchDevices();

	deviceSelect.addEventListener("change", function () {
		const deviceID = this.value;
		if (deviceID) {
			startEventSource(deviceID);
		} else {
			gaugeContainer.innerHTML = "";
			if (eventSource) {
				eventSource.close();
			}
		}
	});
});
