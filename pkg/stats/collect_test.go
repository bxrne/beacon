package stats_test

import (
	"errors"
	"testing"

	"github.com/bxrne/beacon/pkg/config"
	"github.com/bxrne/beacon/pkg/stats"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TEST: GIVEN no enabled monitoring options WHEN Collect is called THEN an empty Metric set is returned
func TestCollect_NoEnabledMonitoring(t *testing.T) {
	cfg := &config.Config{
		Monitoring: config.MonitoringConfig{
			EnableCPU:    false,
			EnableMemory: false,
			EnableDisk:   false,
		},
	}
	mockCPU := &MockCPUMonitor{
		PercentResult: []float64{50.5},
	}
	mockMemory := &MockMemoryMonitor{
		VirtualMemoryResult: &mem.VirtualMemoryStat{UsedPercent: 70.3},
	}
	mockDisk := &MockDiskMonitor{
		UsageResults: map[string]*disk.UsageStat{
			"/":     {UsedPercent: 80.1},
			"/home": {UsedPercent: 65.7},
		},
		UsageErrs: map[string]error{
			// simulate an error for an invalid path (if any)
		},
	}

	res, err := stats.Collect(cfg, mockCPU, mockMemory, mockDisk)

	require.Nil(t, err)
	require.Zero(t, res.CPUUsage)
	require.Zero(t, res.MemoryUsage)
	for _, v := range res.DiskUsage {
		require.Zero(t, v)
	}
}

func TestCollect(t *testing.T) {
	// Define configurations for different monitoring scenarios
	cfg := &config.Config{
		Monitoring: config.MonitoringConfig{
			EnableCPU:    true,
			EnableMemory: true,
			EnableDisk:   true,
			DiskPaths:    []string{"/", "/home"},
		},
	}

	// Initialize mock data for CPU, Memory, and Disk usage
	mockCPU := &MockCPUMonitor{
		PercentResult: []float64{50.5},
	}
	mockMemory := &MockMemoryMonitor{
		VirtualMemoryResult: &mem.VirtualMemoryStat{UsedPercent: 70.3},
	}
	mockDisk := &MockDiskMonitor{
		UsageResults: map[string]*disk.UsageStat{
			"/":     {UsedPercent: 80.1},
			"/home": {UsedPercent: 65.7},
		},
		UsageErrs: map[string]error{
			// simulate an error for an invalid path (if any)
		},
	}

	// Run the test
	metrics, err := stats.Collect(cfg, mockCPU, mockMemory, mockDisk)
	assert.NoError(t, err)

	// Verify CPU usage
	assert.Equal(t, 50.5, metrics.CPUUsage)

	// Verify memory usage
	assert.Equal(t, 70.3, metrics.MemoryUsage)

	// Verify disk usage for each path
	assert.Equal(t, 80.1, metrics.DiskUsage["/"])
	assert.Equal(t, 65.7, metrics.DiskUsage["/home"])
}

func TestCollectWithCPUErr(t *testing.T) {
	cfg := &config.Config{
		Monitoring: config.MonitoringConfig{
			EnableCPU:    true,
			EnableMemory: false,
			EnableDisk:   false,
		},
	}

	// Simulate a CPU error
	mockCPU := &MockCPUMonitor{
		PercentErr: errors.New("CPU monitor failed"),
	}
	mockMemory := &MockMemoryMonitor{}
	mockDisk := &MockDiskMonitor{}

	// Run the test
	metrics, err := stats.Collect(cfg, mockCPU, mockMemory, mockDisk)

	// Expect an error due to CPU monitoring failure
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get CPU usage")
	assert.Equal(t, 0.0, metrics.CPUUsage)
}

func TestCollectWithMemoryErr(t *testing.T) {
	cfg := &config.Config{
		Monitoring: config.MonitoringConfig{
			EnableCPU:    false,
			EnableMemory: true,
			EnableDisk:   false,
		},
	}

	// Simulate a memory error
	mockCPU := &MockCPUMonitor{}
	mockMemory := &MockMemoryMonitor{
		VirtualMemoryErr: errors.New("memory monitor failed"),
	}
	mockDisk := &MockDiskMonitor{}

	// Run the test
	metrics, err := stats.Collect(cfg, mockCPU, mockMemory, mockDisk)

	// Expect an error due to memory monitoring failure
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get memory usage")
	assert.Equal(t, 0.0, metrics.MemoryUsage)
}

func TestCollectWithDiskErr(t *testing.T) {
	cfg := &config.Config{
		Monitoring: config.MonitoringConfig{
			EnableCPU:    false,
			EnableMemory: false,
			EnableDisk:   true,
			DiskPaths:    []string{"/", "/invalid"},
		},
	}

	// Simulate a disk error for "/invalid" path
	mockCPU := &MockCPUMonitor{}
	mockMemory := &MockMemoryMonitor{}
	mockDisk := &MockDiskMonitor{
		UsageResults: map[string]*disk.UsageStat{
			"/": {UsedPercent: 80.1},
		},
		UsageErrs: map[string]error{
			"/invalid": errors.New("disk usage retrieval failed"),
		},
	}

	// Run the test
	metrics, err := stats.Collect(cfg, mockCPU, mockMemory, mockDisk)

	// Verify no error in the main function (disk errors are logged, not returned)
	assert.NoError(t, err)

	// Verify disk usage for valid path
	assert.Equal(t, 80.1, metrics.DiskUsage["/"])

	// Verify that the invalid path is not added to the DiskUsage map
	_, exists := metrics.DiskUsage["/invalid"]
	assert.False(t, exists)
}
