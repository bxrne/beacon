package stats

import (
	"fmt"
	"time"

	"github.com/bxrne/beacon/pkg/config"
	"github.com/charmbracelet/log"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

func RunCollector(cfg *config.Config, log *log.Logger) {
	for {
		metrics, err := collectMetrics(cfg)
		if err != nil {
			fmt.Printf("Failed to collect metrics: %v\n", err)
			continue
		}

		log.Debug(metrics.String())

		time.Sleep(cfg.FrequencyDuration)
	}
}

func collectMetrics(cfg *config.Config) (DeviceMetrics, error) {
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
		virtualMemory, err := mem.VirtualMemory()
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
