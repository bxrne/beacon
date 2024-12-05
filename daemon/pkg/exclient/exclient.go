package exclient

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/bxrne/beacon/daemon/pkg/config"
	"github.com/bxrne/beacon/daemon/pkg/stats"

	"github.com/tarm/serial"
)

// ReadFromPort reads data from the specified serial port and logs it
func ReadFromPort(cfg *config.Config, telemetryChan chan<- stats.DeviceMetrics) error {
	c := &serial.Config{Name: cfg.Telemetry.ExClientPort, Baud: cfg.Telemetry.ExClientBaud}
	s, err := serial.OpenPort(c)
	if err != nil {
		return fmt.Errorf("failed to open port: %w", err)
	}
	defer s.Close()

	scanner := bufio.NewScanner(s)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "[TELEMETRY]") {
			start := strings.Index(line, "{")
			end := strings.LastIndex(line, "}")
			if start == -1 || end == -1 || start > end {
				log.Println("Invalid telemetry data format")
				continue
			}

			telemetryData := line[start : end+1]
			telemetryData = cleanTelemetryData(telemetryData)
			var metrics stats.DeviceMetrics
			if err := json.Unmarshal([]byte(telemetryData), &metrics); err != nil {
				log.Printf("Failed to unmarshal telemetry data: %v", err)
				continue
			}

			hostname, err := os.Hostname()
			if err != nil {
				log.Printf("Failed to get hostname: %v", err)
				continue
			}
			metrics.Hostname = hostname + "-exclient"
			log.Printf("Telemetry metrics: %s", metrics.String())
			telemetryChan <- metrics
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read from port: %w", err)
	}

	return nil
}

// cleanTelemetryData removes any unwanted characters from the telemetry data
func cleanTelemetryData(data string) string {
	re := regexp.MustCompile(`[^\x20-\x7E]+`)
	data = re.ReplaceAllString(data, "")
	data = strings.ReplaceAll(data, "\\", "")
	return data
}
