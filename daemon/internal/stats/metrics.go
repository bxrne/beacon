package stats

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"time"
)

type Metric struct {
	Type       string `json:"type"`
	Value      string `json:"value"` // Changed from float64 to string on read
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

func GetDeviceUUID() string {
	cmd := exec.Command("hostid") // Get host identifier (available on Unix-like systems)
	output, _ := cmd.Output()     // Execute command and get output
	return string(output)
}

func CollectMetrics() (*DeviceMetrics, error) {
	var deviceMetrics DeviceMetrics

	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	deviceMetrics.Hostname = hostname

	// Initialize monitors
	var memoryMon MemoryMonitor = MemoryMon{}
	var diskMon DiskMon = DiskMon{}
	var hostMon HostMon = HostMon{}

	// Collect memory usage
	vmStat, err := memoryMon.VirtualMemory()
	if err != nil {
		return nil, err
	}
	deviceMetrics.Metrics = append(deviceMetrics.Metrics, Metric{
		Type:       "memory_used",
		Value:      fmt.Sprintf("%.2f", vmStat.UsedPercent),
		Unit:       "%",
		RecordedAt: time.Now().UTC().Format(time.RFC3339),
	})

	// Collect disk usage
	diskUsage, err := diskMon.Usage("/")
	if err != nil {
		return nil, err
	}
	deviceMetrics.Metrics = append(deviceMetrics.Metrics, Metric{
		Type:       "disk_used",
		Value:      fmt.Sprintf("%.2f", diskUsage.UsedPercent),
		Unit:       "%",
		RecordedAt: time.Now().UTC().Format(time.RFC3339),
	})

	// Collect uptime
	uptime, err := hostMon.Uptime()
	if err != nil {
		return nil, err
	}
	deviceMetrics.Metrics = append(deviceMetrics.Metrics, Metric{
		Type:       "uptime",
		Value:      fmt.Sprintf("%d", uptime),
		Unit:       "seconds",
		RecordedAt: time.Now().UTC().Format(time.RFC3339),
	})

	return &deviceMetrics, nil
}
