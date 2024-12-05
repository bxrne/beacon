import { CircleDot } from "lucide-react";
import { MetricsCard } from "../dashboard/metrics-card";
import { TrafficDevice } from "@/lib/types";

interface TrafficStatusProps {
  device: TrafficDevice;
}

export function TrafficStatus({ device }: TrafficStatusProps) {
  return (
    <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
      <MetricsCard
        title="Traffic Light Status"
        value={device.car_light}
        icon={<CircleDot className="h-4 w-4" style={{ 
          color: device.car_light === 'RED' ? 'red' : 
                 device.car_light === 'YELLOW' ? 'yellow' : 'green' 
        }} />}
        description={`Device ID: ${device.id}`}
      />
      <MetricsCard
        title="Pedestrian Light"
        value={device.pedestrian_light}
        icon={<CircleDot className="h-4 w-4" style={{ 
          color: device.pedestrian_light === 'RED' ? 'red' : 'green'
        }} />}
      />
      <MetricsCard
        title="Pedestrian Button"
        value={device.pedestrian_button ? "Pressed" : "Not Pressed"}
        icon={<CircleDot className="h-4 w-4" />}
      />
    </div>
  );
}