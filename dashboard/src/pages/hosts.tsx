import { useState } from "react";
import { DeviceStatus } from "@/components/dashboard/device-status";
import { MetricsChart } from "@/components/dashboard/metrics-chart";
import { generateMockTimeseriesData, mockDevices } from "@/lib/mock-data";
import { LaptopDevice } from "@/lib/types";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";

export function HostsPage() {
  const laptops = mockDevices.filter((d): d is LaptopDevice => d.type === 'laptop');
  const [selectedDevice, setSelectedDevice] = useState<string>(laptops[0].id);
  const [data] = useState(generateMockTimeseriesData(4));
  
  const device = laptops.find(d => d.id === selectedDevice)!;
  const latestMetrics = data[data.length - 1];

  return (
    <div className="space-y-4 p-8 pt-6">
      <div className="flex items-center justify-between space-y-2">
        <h2 className="text-3xl font-bold tracking-tight">Hosts</h2>
        <Select value={selectedDevice} onValueChange={setSelectedDevice}>
          <SelectTrigger className="w-[200px]">
            <SelectValue placeholder="Select device" />
          </SelectTrigger>
          <SelectContent>
            {laptops.map((device) => (
              <SelectItem key={device.id} value={device.id}>
                {device.id}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>

      <DeviceStatus
        deviceId={device.id}
        status={device.status}
        cpuUsage={latestMetrics.cpuUsage}
        memoryUsage={latestMetrics.memoryUsage}
        networkLatency={latestMetrics.networkLatency}
      />

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        <MetricsChart
          title="CPU Usage"
          data={data}
          dataKey="cpuUsage"
          color="hsl(var(--chart-1))"
        />
        <MetricsChart
          title="Memory Usage"
          data={data}
          dataKey="memoryUsage"
          color="hsl(var(--chart-2))"
        />
        <MetricsChart
          title="Network Latency"
          data={data}
          dataKey="networkLatency"
          color="hsl(var(--chart-3))"
        />
      </div>
    </div>
  );
}