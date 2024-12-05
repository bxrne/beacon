package stats

import (
	"fmt"
	"os/exec"
	"sort"
)

type Metric struct {
	Type       string `json:"type"`
	Value      string `json:"value"` // Changed from float64 to string
	Unit       string `json:"unit"`
	RecordedAt string `json:"recorded_at"` // WARN: Daemon only field
}

type DeviceMetrics struct {
	Metrics  []Metric `json:"metrics"`
	Hostname string   `json:"hostname"` // WARN: Daemon only field

}

func (d *DeviceMetrics) String() string {
	sort.Slice(d.Metrics, func(i, j int) bool {
		return d.Metrics[i].Type < d.Metrics[j].Type
	})

	metricsStr := ""
	for _, metric := range d.Metrics {
		metricsStr += fmt.Sprintf("%s: %s %s | ", metric.Type, metric.Value, metric.Unit) // Adjusted formatting
	}

	return metricsStr
}

func GetDeviceUUID() string {
	cmd := exec.Command("hostid") // Get host identifier (available on Unix-like systems)
	output, _ := cmd.Output()     // Execute command and get output
	return string(output)
}
