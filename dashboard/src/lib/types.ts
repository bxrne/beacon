export type DeviceType = 'laptop' | 'traffic_light';

export interface BaseDevice {
  id: string;
  type: DeviceType;
  status: 'online' | 'offline';
}

export interface LaptopDevice extends BaseDevice {
  type: 'laptop';
  cpuUsage: number;
  memoryUsage: number;
  networkLatency: number;
}

export interface TrafficDevice extends BaseDevice {
  type: 'traffic_light';
  car_light: 'RED' | 'YELLOW' | 'GREEN';
  pedestrian_light: 'RED' | 'GREEN';
  pedestrian_button: boolean;
}

export type Device = LaptopDevice | TrafficDevice;