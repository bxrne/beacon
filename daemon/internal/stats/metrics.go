package stats

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"
	"time"
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
	log.Println("Collecting memory usage...")
	vmStat, err := memoryMon.VirtualMemory()
	if err != nil {
		return nil, err
	}
	deviceMetrics.Metrics = append(deviceMetrics.Metrics, Metric{
		Type:       "MemoryUsedPercent",
		Value:      fmt.Sprintf("%.2f", vmStat.UsedPercent),
		Unit:       "%",
		RecordedAt: time.Now().UTC().Format(time.RFC3339),
	})
	log.Println("Memory usage collected.")

	// Collect disk usage
	log.Println("Collecting disk usage...")
	diskUsage, err := diskMon.Usage("/")
	if err != nil {
		return nil, err
	}
	deviceMetrics.Metrics = append(deviceMetrics.Metrics, Metric{
		Type:       "DiskUsedPercent",
		Value:      fmt.Sprintf("%.2f", diskUsage.UsedPercent),
		Unit:       "%",
		RecordedAt: time.Now().UTC().Format(time.RFC3339),
	})
	log.Println("Disk usage collected.")

	// Collect uptime
	log.Println("Collecting uptime...")
	uptime, err := hostMon.Uptime()
	if err != nil {
		return nil, err
	}
	deviceMetrics.Metrics = append(deviceMetrics.Metrics, Metric{
		Type:       "Uptime",
		Value:      fmt.Sprintf("%d", uptime),
		Unit:       "seconds",
		RecordedAt: time.Now().UTC().Format(time.RFC3339),
	})
	log.Println("Uptime collected.")

	return &deviceMetrics, nil
}
