import { Device, LaptopDevice, TrafficDevice } from './types';

export const mockDevices: Device[] = [
  {
    id: 'laptop-001',
    type: 'laptop',
    status: 'online',
    cpuUsage: 45,
    memoryUsage: 60,
    networkLatency: 20,
  },
  {
    id: 'laptop-002',
    type: 'laptop',
    status: 'online',
    cpuUsage: 30,
    memoryUsage: 45,
    networkLatency: 15,
  },
  {
    id: 'traffic-001',
    type: 'traffic_light',
    status: 'online',
    car_light: 'GREEN',
    pedestrian_light: 'RED',
    pedestrian_button: false,
  },
  {
    id: 'traffic-002',
    type: 'traffic_light',
    status: 'online',
    car_light: 'RED',
    pedestrian_light: 'GREEN',
    pedestrian_button: true,
  },
];

export function generateMockTimeseriesData(hours: number = 24) {
  const data = [];
  const now = new Date();
  
  for (let i = 0; i < hours * 60; i++) {
    const timestamp = new Date(now.getTime() - (i * 60000)).toISOString();
    data.unshift({
      timestamp,
      cpuUsage: Math.floor(Math.random() * 40) + 20,
      memoryUsage: Math.floor(Math.random() * 30) + 40,
      networkLatency: Math.floor(Math.random() * 50) + 10,
    });
  }
  
  return data;
}