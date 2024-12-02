package stats

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os/exec"
	"sort"
)

type Metric struct {
	Type  string  `json:"type"`
	Value float64 `json:"value"`
	Unit  string  `json:"unit"`
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
		metricsStr += fmt.Sprintf("%s: %.2f %s | ", metric.Type, metric.Value, metric.Unit)
	}

	return metricsStr
}

func GetDeviceUUID() string {
	cmd := exec.Command("hostid")       // Get host identifier (available on Unix-like systems)
	output, _ := cmd.Output()           // Execute command and get output
	hash := sha256.Sum256(output)       // Hash the output
	return hex.EncodeToString(hash[:8]) // Return first 8 bytes as hex string
}
