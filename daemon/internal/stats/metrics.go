package stats

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	metrics "github.com/bxrne/beacon/aggregator/pkg/types"
)

func GetDeviceUUID() string {
	cmd := exec.Command("hostid") // Get host identifier (available on Unix-like systems)
	output, _ := cmd.Output()     // Execute command and get output
	return string(output)
}

func CollectMetrics() (*metrics.DeviceMetrics, error) {
	var deviceMetrics metrics.DeviceMetrics

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
	deviceMetrics.Metrics = append(deviceMetrics.Metrics, metrics.Metric{
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
	deviceMetrics.Metrics = append(deviceMetrics.Metrics, metrics.Metric{
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
	deviceMetrics.Metrics = append(deviceMetrics.Metrics, metrics.Metric{
		Type:       "uptime",
		Value:      fmt.Sprintf("%d", uptime),
		Unit:       "seconds",
		RecordedAt: time.Now().UTC().Format(time.RFC3339),
	})

	return &deviceMetrics, nil
}
