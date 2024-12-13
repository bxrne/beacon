document.addEventListener("DOMContentLoaded", async () => {
	const deviceSelect = document.getElementById("deviceSelect");
	const updateIntervalSelect = document.getElementById("updateInterval");
	const metricsTable = document.getElementById("metrics");
	const pagination = document.getElementById("pagination");
	const metricTypeFilter = document.getElementById("metricTypeFilter");
	let refreshIntervalId = null;

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

	async function fetchMetrics(deviceID, page = 1) {
		const metricType = metricTypeFilter.value;
		const response = await fetch(
			`/api/metrics?page=${page}&type=${metricType}`,
			{
				headers: {
					"X-DeviceID": deviceID,
				},
			}
		);
		const data = await response.json();
		metricsTable.innerHTML = "";
		if (data.metrics) {
			data.metrics.forEach((metric) => {
				const row = document.createElement("tr");
				row.innerHTML = `
                <td>${metric.Type ? metric.Type.Name : ""}</td>
                <td>${metric.Value}</td>
                <td>${metric.Unit ? metric.Unit.Name : ""}</td>
                <td>${new Date(metric.RecordedAt).toLocaleString()}</td>
            `;
				metricsTable.appendChild(row);
			});
		}
		renderPagination(data.totalPages, data.currentPage);
	}

	function renderPagination(totalPages, currentPage) {
		pagination.innerHTML = "";

		const createPageItem = (page, text = page, disabled = false) => {
			const pageItem = document.createElement("li");
			pageItem.className =
				"page-item" +
				(page === currentPage ? " active" : "") +
				(disabled ? " disabled" : "");
			pageItem.innerHTML = `<a class="page-link" href="#" data-page="${page}">${text}</a>`;
			if (!disabled && page !== currentPage) {
				pageItem.addEventListener("click", function (event) {
					event.preventDefault();
					fetchMetrics(deviceSelect.value, page);
				});
			}
			return pageItem;
		};

		// Previous page link
		if (currentPage > 1) {
			pagination.appendChild(createPageItem(currentPage - 1, "‹ Prev"));
		}

		// Page numbers
		const startPage = Math.max(1, currentPage - 3);
		const endPage = Math.min(totalPages, currentPage + 3);

		// Ellipsis before
		if (startPage > 1) {
			pagination.appendChild(createPageItem(null, "...", true));
		}

		// Page links
		for (let i = startPage; i <= endPage; i++) {
			pagination.appendChild(createPageItem(i));
		}

		// Ellipsis after
		if (endPage < totalPages) {
			pagination.appendChild(createPageItem(null, "...", true));
		}

		// Next page link
		if (currentPage < totalPages) {
			pagination.appendChild(createPageItem(currentPage + 1, "Next ›"));
		}
	}

	function startAutoRefresh(deviceID) {
		const interval = parseInt(updateIntervalSelect.value);
		clearInterval(refreshIntervalId);
		if (interval > 0) {
			refreshIntervalId = setInterval(() => {
				fetchMetrics(deviceID);
			}, interval);
			pagination.style.display = "none"; // Hide pagination
		} else {
			pagination.style.display = ""; // Show pagination
		}
	}

	function clearMetrics() {
		metricsTable.innerHTML = "";
		pagination.innerHTML = "";
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

	updateIntervalSelect.addEventListener("change", function () {
		const deviceID = deviceSelect.value;
		if (deviceID) {
			startAutoRefresh(deviceID);
		}
	});

	metricTypeFilter.addEventListener("input", function () {
		const deviceID = deviceSelect.value;
		if (deviceID) {
			fetchMetrics(deviceID);
		}
	});

	// Set "No update" as the default active option
	updateIntervalSelect.value = "0";
});
