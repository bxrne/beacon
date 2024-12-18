package stats

import (
	"fmt"
	"time"

	metric_types "github.com/bxrne/beacon/aggregator/pkg/types"
	"github.com/bxrne/beacon/daemon/internal/config"
)

// Collect collects metrics from the system
func Collect(cfg *config.Config, host HostMonitor, memory MemoryMonitor, disk DiskMonitor) (metric_types.DeviceMetrics, error) {
	var metrics []metric_types.Metric
	currentTime := time.Now().UTC()

	// Collect host metrics
	hostUptime, err := host.Uptime()
	if err != nil {
		return metric_types.DeviceMetrics{}, fmt.Errorf("failed to collect host metrics: %w", err)
	}
	hostMetrics := []metric_types.Metric{
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
		return metric_types.DeviceMetrics{}, fmt.Errorf("failed to collect memory metrics: %w", err)
	}
	metrics = append(metrics, metric_types.Metric{
		Type:       "cpu_usage",
		Unit:       "percent",
		Value:      fmt.Sprintf("%.2f", memoryMetrics.UsedPercent),
		RecordedAt: currentTime.Format(time.RFC3339),
	})

	// Collect disk metrics
	diskMetrics, err := disk.Usage("/")
	if err != nil {
		return metric_types.DeviceMetrics{}, fmt.Errorf("failed to collect disk metrics: %w", err)
	}
	metrics = append(metrics, metric_types.Metric{
		Type:       "disk_usage",
		Unit:       "percent",
		Value:      fmt.Sprintf("%.2f", diskMetrics.UsedPercent),
		RecordedAt: currentTime.Format(time.RFC3339),
	})

	return metric_types.DeviceMetrics{Metrics: metrics}, nil
}
