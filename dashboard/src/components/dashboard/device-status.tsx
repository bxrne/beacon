import { Activity, Cpu, HardDrive, Network } from "lucide-react";
import { MetricsCard } from "./metrics-card";

interface DeviceStatusProps {
  deviceId: string;
  status: "online" | "offline";
  cpuUsage: number;
  memoryUsage: number;
  networkLatency: number;
}

export function DeviceStatus({ deviceId, status, cpuUsage, memoryUsage, networkLatency }: DeviceStatusProps) {
  return (
    <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
      <MetricsCard
        title="Device Status"
        value={status}
        icon={<Activity className="h-4 w-4 text-muted-foreground" />}
        description={`Device ID: ${deviceId}`}
        className={status === "online" ? "border-green-500" : "border-red-500"}
      />
      <MetricsCard
        title="CPU Usage"
        value={`${cpuUsage}%`}
        icon={<Cpu className="h-4 w-4 text-muted-foreground" />}
      />
      <MetricsCard
        title="Memory Usage"
        value={`${memoryUsage}%`}
        icon={<HardDrive className="h-4 w-4 text-muted-foreground" />}
      />
      <MetricsCard
        title="Network Latency"
        value={`${networkLatency}ms`}
        icon={<Network className="h-4 w-4 text-muted-foreground" />}
      />
    </div>
  );
}