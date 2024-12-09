import React from "react";
import { format } from "date-fns";
import { useDeviceStore } from "../stores/deviceStore";

export const MetricsTable: React.FC = () => {
	const { metrics, loading, error } = useDeviceStore();

	if (loading) {
		return <div className="text-center">Loading metrics...</div>;
	}

	if (error) {
		return <div className="text-red-500">{error}</div>;
	}

	if (!metrics || !metrics.metrics || metrics.metrics.length === 0) {
		return <div className="text-gray-500">No metrics available</div>;
	}

	return (
		<div className="overflow-x-auto">
			<table className="min-w-full bg-white rounded-lg overflow-hidden">
				<thead className="bg-gray-100">
					<tr>
						<th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
							Time
						</th>
						<th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
							Type
						</th>
						<th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
							Value
						</th>
						<th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
							Unit
						</th>
					</tr>
				</thead>
				<tbody className="divide-y divide-gray-200">
					{metrics.metrics.map((metric, index) => (
						<tr key={index} className="hover:bg-gray-50">
							<td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
								{format(new Date(metric.recorded_at), "PPpp")}
							</td>
							<td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
								{metric.type}
							</td>
							<td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
								{metric.value}
							</td>
							<td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
								{metric.unit}
							</td>
						</tr>
					))}
				</tbody>
			</table>
		</div>
	);
};
