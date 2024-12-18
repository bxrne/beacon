package metrics

import "time"

type DeviceMetrics struct {
	Metrics []Metric `json:"metrics"`
}

type Metric struct {
	Type       string `json:"type"`
	Value      string `json:"value"`
	Unit       string `json:"unit"`
	RecordedAt string `json:"recorded_at"`
}

func NewMetric(metricType, value, unit string) Metric {
	return Metric{
		Type:       metricType,
		Value:      value,
		Unit:       unit,
		RecordedAt: time.Now().UTC().Format(time.RFC3339),
	}
}
