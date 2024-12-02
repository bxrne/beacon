package stats

import (
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

// INFO: Abstracted for testing
type MemoryMonitor interface {
	VirtualMemory() (*mem.VirtualMemoryStat, error)
}

type DiskMonitor interface {
	Usage(string) (*disk.UsageStat, error)
}

type HostMonitor interface {
	Uptime() (uint64, error)
}

// INFO: Runtime implementations
type MemoryMon struct{}

func (MemoryMon) VirtualMemory() (*mem.VirtualMemoryStat, error) {
	return mem.VirtualMemory()
}

type DiskMon struct{}

func (DiskMon) Usage(path string) (*disk.UsageStat, error) {
	return disk.Usage(path)
}

type HostMon struct{}

func (HostMon) Uptime() (uint64, error) {
	return host.Uptime()
}
