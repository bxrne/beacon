import { useState } from "react";
import { mockDevices } from "@/lib/mock-data";
import { TrafficDevice } from "@/lib/types";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { TrafficStatus } from "@/components/traffic/traffic-status";

export function TrafficPage() {
  const trafficLights = mockDevices.filter((d): d is TrafficDevice => d.type === 'traffic_light');
  const [selectedDevice, setSelectedDevice] = useState<string>(trafficLights[0].id);
  
  const device = trafficLights.find(d => d.id === selectedDevice)!;

  return (
    <div className="space-y-4 p-8 pt-6">
      <div className="flex items-center justify-between space-y-2">
        <h2 className="text-3xl font-bold tracking-tight">Traffic Lights</h2>
        <Select value={selectedDevice} onValueChange={setSelectedDevice}>
          <SelectTrigger className="w-[200px]">
            <SelectValue placeholder="Select traffic light" />
          </SelectTrigger>
          <SelectContent>
            {trafficLights.map((device) => (
              <SelectItem key={device.id} value={device.id}>
                {device.id}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>

      <TrafficStatus device={device} />
    </div>
  );
}