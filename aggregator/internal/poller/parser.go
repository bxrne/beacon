package poller

import (
	"strings"
	"time"

	metrics "github.com/bxrne/beacon/aggregator/pkg/types"
)

func parseMetrics(response []byte) (*metrics.DeviceMetrics, error) {
	// Find the payload after HTTP headers
	parts := strings.Split(string(response), "\r\n\r\n")

	var payload string
	if len(parts) > 1 {
		payload = parts[1]
	} else {
		payload = parts[0]
	}
	// Remove STX, length byte, and ETX
	payload = payload[2 : len(payload)-1]

	// Remove the end byte if it exists
	if len(payload) > 0 && payload[len(payload)-1] == 0x03 {
		payload = payload[:len(payload)-1]
	}

	// Parse the key-value pairs
	pairs := strings.Split(payload, ", ")
	result := &metrics.DeviceMetrics{
		Metrics: make([]metrics.Metric, 0),
	}

	var recordedAt string
	for _, pair := range pairs {
		parts := strings.Split(pair, ": ")
		if len(parts) != 2 {
			continue
		}

		key := parts[0]
		value := parts[1]

		// Handle recorded_at separately
		if key == "recorded_at" {
			recordedAt = value
			continue
		}

		// Create metric based on the type
		metric := metrics.Metric{
			Type:  key,
			Value: value,
			Unit:  determineUnit(key),
		}

		result.Metrics = append(result.Metrics, metric)
	}

	// Set recorded_at for all metrics
	if recordedAt == "" {
		recordedAt = time.Now().UTC().Format(time.RFC3339)
	}
	for i := range result.Metrics {
		result.Metrics[i].RecordedAt = recordedAt
	}

	return result, nil
}

func determineUnit(metricType string) string {
	switch metricType {
	case "memory_used":
		return "percent"
	case "disk_used":
		return "percent"
	case "uptime":
		return "seconds"
	case "car_light":
		return "color"
	case "ped_light":
		return "color"
	default:
		return "unknown"
	}
}
