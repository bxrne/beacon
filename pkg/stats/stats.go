package stats

import (
	"fmt"
)

type DeviceMetrics struct {
	CPUUsage    float64            `json:"cpu_usage"`
	MemoryUsage float64            `json:"memory_usage"`
	DiskUsage   map[string]float64 `json:"disk_usage"`
}

func (d *DeviceMetrics) String() string {
	diskUsage := ""
	for path, usage := range d.DiskUsage {
		diskUsage += fmt.Sprintf("(%s) %.2f%% ", path, usage)
	}

	return fmt.Sprintf("CPU: %.2f%% | Memory: %.2f%% | Disk: %s", d.CPUUsage, d.MemoryUsage, diskUsage)
}
