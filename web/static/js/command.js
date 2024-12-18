document.addEventListener("DOMContentLoaded", async () => {
	const deviceSelect = document.getElementById("deviceSelect");
	const commandForm = document.getElementById("commandForm");
	const errorAlert = document.getElementById("errorAlert");
	const successAlert = document.getElementById("successAlert");

	// Fetch and populate devices
	async function loadDevices() {
		const response = await fetch("/api/device");
		const devices = await response.json();

		devices.forEach((device) => {
			const option = document.createElement("option");
			option.value = device;
			option.textContent = device;
			deviceSelect.appendChild(option);
		});
	}

	// Handle form submission
	commandForm.addEventListener("submit", async (e) => {
		e.preventDefault();
		errorAlert.style.display = "none";
		successAlert.style.display = "none";

		const command = {
			device: deviceSelect.value,
			command: document.getElementById("commandInput").value,
		};

		try {
			const response = await fetch("/api/command", {
				method: "POST",
				headers: {
					"Content-Type": "application/json",
				},
				body: JSON.stringify(command),
			});

			if (!response.ok) {
				const error = await response.json();
				throw new Error(error.error || "Failed to send command");
			}

			successAlert.textContent = "Command sent successfully";
			successAlert.style.display = "block";
			commandForm.reset();
		} catch (err) {
			errorAlert.textContent = err.message;
			errorAlert.style.display = "block";
		}
	});

	await loadDevices();
});
