import React from "react";
import { useDeviceStore } from "../../stores/deviceStore";

export const MetricsGrid: React.FC = () => {
	const { metrics, loading, error } = useDeviceStore();

	if (loading) {
		return <div>Loading metrics...</div>;
	}

	if (error) {
		return <div className="text-red-500">{error}</div>;
	}
	console.log(metrics);

	if (!metrics) {
		return <div>No device selected</div>;
	}

	if (!metrics || Object.keys(metrics).length === 0) {
		return <div>No metrics available</div>;
	}

	return (
		<div className="overflow-x-auto">
			<table className="min-w-full bg-white">
				<thead>
					<tr>
						<th className="py-2 px-4 border-b">Recorded At</th>
						<th className="py-2 px-4 border-b">Type</th>
						<th className="py-2 px-4 border-b">Unit</th>
						<th className="py-2 px-4 border-b">Value</th>
					</tr>
				</thead>
				<tbody>
					{metrics.map((metric, index) => (
						<tr key={index}>
							<td className="py-2 px-4 border-b">{metric.RecordedAt}</td>
							<td className="py-2 px-4 border-b">{metric.Type.Name}</td>
							<td className="py-2 px-4 border-b">{metric.Unit.Name}</td>
							<td className="py-2 px-4 border-b">{metric.Value}</td>
						</tr>
					))}
				</tbody>
			</table>
		</div>
	);
};
