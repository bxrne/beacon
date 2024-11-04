package stats_test

import (
	"testing"

	"github.com/bxrne/beacon/pkg/stats"
	"github.com/stretchr/testify/assert"
)

func TestDeviceMetricsString(t *testing.T) {
	metrics := stats.DeviceMetrics{
		CPUUsage:    10.0,
		MemoryUsage: 20.0,
		DiskUsage:   map[string]float64{"/": 30.0},
	}

	expected := "CPU: 10.00% | Memory: 20.00% | Disk: (/) 30.00% "
	assert.Equal(t, expected, metrics.String())
}
