package stats

import (
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

// INFO: Abstracted for testing
type CPUMonitor interface {
	Percent(time.Duration, bool) ([]float64, error)
}

type MemoryMonitor interface {
	VirtualMemory() (*mem.VirtualMemoryStat, error)
}

type DiskMonitor interface {
	Usage(string) (*disk.UsageStat, error)
}

// INFO: Runtime implementations
type CPUMon struct{}

func (CPUMon) Percent(interval time.Duration, percpu bool) ([]float64, error) {
	return cpu.Percent(interval, percpu)
}

type MemoryMon struct{}

func (MemoryMon) VirtualMemory() (*mem.VirtualMemoryStat, error) {
	return mem.VirtualMemory()
}

type DiskMon struct{}

func (DiskMon) Usage(path string) (*disk.UsageStat, error) {
	return disk.Usage(path)
}
