package stats_test

import (
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

type MockMemoryMonitor struct {
	VirtualMemoryResult *mem.VirtualMemoryStat
	VirtualMemoryErr    error
}

func (m *MockMemoryMonitor) VirtualMemory() (*mem.VirtualMemoryStat, error) {
	return m.VirtualMemoryResult, m.VirtualMemoryErr
}

type MockDiskMonitor struct {
	UsageResults map[string]*disk.UsageStat
	UsageErrs    map[string]error
}

func (m *MockDiskMonitor) Usage(path string) (*disk.UsageStat, error) {
	if err, exists := m.UsageErrs[path]; exists {
		return nil, err
	}
	return m.UsageResults[path], nil
}

type MockHostMonitor struct {
	UptimeResult uint64
	UptimeErr    error
}

func (m *MockHostMonitor) Uptime() (uint64, error) {
	return m.UptimeResult, m.UptimeErr
}
