package metrics

import (
	"fmt"
	"sort"
	"time"
)

type Metric struct {
	Type       string `json:"type"`
	Value      string `json:"value"`
	Unit       string `json:"unit"`
	RecordedAt string `json:"recorded_at"`
}

type DeviceMetrics struct {
	Metrics  []Metric `json:"metrics"`
	Hostname string   `json:"hostname"`
}

func (d *DeviceMetrics) String() string {
	sort.Slice(d.Metrics, func(i, j int) bool {
		return d.Metrics[i].Type < d.Metrics[j].Type
	})

	metricsStr := ""
	// type: value, type: value, recorded_at: time
	for _, metric := range d.Metrics {
		metricsStr += fmt.Sprintf("%s: %s, ", metric.Type, metric.Value)
	}
	if len(metricsStr) > 0 {
		metricsStr = metricsStr[:len(metricsStr)-2] // Remove trailing comma and space
	}
	metricsStr += fmt.Sprintf(", recorded_at: %s", d.Metrics[0].RecordedAt) // WARN: Only first metric's recorded_at is used

	return metricsStr
}
func NewMetric(metricType, value, unit string) Metric {
	return Metric{
		Type:       metricType,
		Value:      value,
		Unit:       unit,
		RecordedAt: time.Now().UTC().Format(time.RFC3339),
	}
}
