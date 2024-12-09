package stats

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/bxrne/beacon/daemon/pkg/config"
)

// Send sends the collected metrics to the server
func Send(cfg *config.Config, metrics DeviceMetrics) error {
	url := cfg.Telemetry.Server + "/metric"
	data, err := json.Marshal(metrics)
	if err != nil {
		return fmt.Errorf("failed to marshal metrics: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-DeviceID", cfg.Labels.Service)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-OK response: %s", resp.Status)
	}

	return nil
}

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
			RecordedAt: fmt.Sprintf("%s", currentTime.Format(time.RFC3339)),
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
		RecordedAt: fmt.Sprintf("%s", currentTime.Format(time.RFC3339)),
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
		RecordedAt: fmt.Sprintf("%s", currentTime.Format(time.RFC3339)),
	})

	return DeviceMetrics{Metrics: metrics}, nil
}

// sanitizeMetricType sanitizes the metric type string
func sanitizeMetricType(metricType string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9_]+`)
	return strings.ToLower(re.ReplaceAllString(metricType, "_"))
}
