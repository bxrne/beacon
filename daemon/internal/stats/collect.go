package stats

import (
	"fmt"
	"time"

	"github.com/bxrne/beacon/daemon/internal/config"
)

// Collect collects metrics from the system
func Collect(cfg *config.Config, host HostMonitor, memory MemoryMonitor, disk DiskMonitor) (DeviceMetrics, error) {
	var metrics []Metric
	currentTime := time.Now().UTC()

	// Collect host metrics
	hostUptime, err := host.Uptime()
	if err != nil {
		return DeviceMetrics{}, fmt.Errorf("failed to collect host metrics: %w", err)
	}
	hostMetrics := []Metric{
		{
			Type:       "uptime",
			Unit:       "seconds",
			Value:      fmt.Sprintf("%d", hostUptime),
			RecordedAt: currentTime.Format(time.RFC3339),
		},
	}
	metrics = append(metrics, hostMetrics...)

	// Collect memory metrics
	memoryMetrics, err := memory.VirtualMemory()
	if err != nil {
		return DeviceMetrics{}, fmt.Errorf("failed to collect memory metrics: %w", err)
	}
	metrics = append(metrics, Metric{
		Type:       "cpu_usage",
		Unit:       "percent",
		Value:      fmt.Sprintf("%.2f", memoryMetrics.UsedPercent),
		RecordedAt: currentTime.Format(time.RFC3339),
	})

	// Collect disk metrics
	diskMetrics, err := disk.Usage("/")
	if err != nil {
		return DeviceMetrics{}, fmt.Errorf("failed to collect disk metrics: %w", err)
	}
	metrics = append(metrics, Metric{
		Type:       "disk_usage",
		Unit:       "percent",
		Value:      fmt.Sprintf("%.2f", diskMetrics.UsedPercent),
		RecordedAt: currentTime.Format(time.RFC3339),
	})

	return DeviceMetrics{Metrics: metrics}, nil
}
