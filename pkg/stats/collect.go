package stats

import (
	"fmt"
	"time"

	"github.com/bxrne/beacon/pkg/config"
)

func Collect(cfg *config.Config, cpu CPUMonitor, memory MemoryMonitor, disk DiskMonitor) (DeviceMetrics, error) {
	metrics := DeviceMetrics{
		DiskUsage: make(map[string]float64),
	}

	if cfg.Monitoring.EnableCPU {
		cpuPercent, err := cpu.Percent(time.Second, false)
		if err != nil {
			return metrics, fmt.Errorf("failed to get CPU usage: %w", err)
		}
		metrics.CPUUsage = cpuPercent[0]
	}

	if cfg.Monitoring.EnableMemory {
		virtualMemory, err := memory.VirtualMemory()
		if err != nil {
			return metrics, fmt.Errorf("failed to get memory usage: %w", err)
		}
		metrics.MemoryUsage = virtualMemory.UsedPercent
	}

	if cfg.Monitoring.EnableDisk {
		for _, path := range cfg.Monitoring.DiskPaths {
			diskUsage, err := disk.Usage(path)
			if err != nil {
				fmt.Printf("Failed to get disk usage for %s: %v\n", path, err)
				continue
			}
			metrics.DiskUsage[path] = diskUsage.UsedPercent
		}
	}

	return metrics, nil
}
