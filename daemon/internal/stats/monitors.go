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
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}
	return vmStat, nil
}

type DiskMon struct{}

func (DiskMon) Usage(path string) (*disk.UsageStat, error) {
	usageStat, err := disk.Usage(path)
	if err != nil {
		return nil, err
	}
	return usageStat, nil
}

type HostMon struct{}

func (HostMon) Uptime() (uint64, error) {
	return host.Uptime()
}
