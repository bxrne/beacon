export interface DeviceMetrics {
  metrics: Metric[];
}

export interface Metric {
  recorded_at: string;
  type: string;
  unit: string;
  value: string;
}

export interface ErrorResponse {
  error: string;
}

export interface HealthResponse {
  status: string;
}