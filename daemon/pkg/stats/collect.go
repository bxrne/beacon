package stats

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/bxrne/beacon/pkg/config"
)

func Send(cfg *config.Config, metrics DeviceMetrics) error {
	jsonData, err := json.Marshal(metrics)
	if err != nil {
		return fmt.Errorf("failed to marshal metrics: %w", err)
	}

	hostname, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("failed to get hostname: %w", err)
	}

	req, err := http.NewRequest("POST", cfg.Telemetry.Server+"/metric", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Hostname", hostname)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send metrics: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send metrics: %s", resp.Status)
	}

	return nil
}

func Collect(cfg *config.Config, host HostMonitor, memory MemoryMonitor, disk DiskMonitor) (DeviceMetrics, error) {
	metrics := DeviceMetrics{}

	virtualMemory, err := memory.VirtualMemory()
	if err != nil {
		return metrics, fmt.Errorf("failed to get memory usage: %w", err)
	}
	metrics.Metrics = append(metrics.Metrics, Metric{
		Type:  "memory_usage",
		Value: virtualMemory.UsedPercent,
		Unit:  "percent",
	})

	for _, path := range cfg.Monitoring.DiskPaths {
		diskUsage, err := disk.Usage(path)
		if err != nil {
			fmt.Printf("Failed to get disk usage for %s: %v\n", path, err)
			continue
		}
		metrics.Metrics = append(metrics.Metrics, Metric{
			Type:  sanitizeMetricType(fmt.Sprintf("disk_usage_%s", path)),
			Value: diskUsage.UsedPercent,
			Unit:  "percent",
		})
	}

	uptime, err := host.Uptime()
	if err != nil {
		return metrics, fmt.Errorf("failed to get uptime: %w", err)
	}
	metrics.Metrics = append(metrics.Metrics, Metric{
		Type:  "uptime",
		Value: float64(uptime),
		Unit:  "seconds",
	})

	return metrics, nil
}

func sanitizeMetricType(metricType string) string {
	re := regexp.MustCompile(`[^\w]+`)
	sanitized := re.ReplaceAllString(metricType, "_")
	return strings.Trim(sanitized, "_")
}
