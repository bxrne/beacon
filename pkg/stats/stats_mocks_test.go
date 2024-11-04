package stats_test

import (
	"time"

	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

type MockCPUMonitor struct {
	PercentResult []float64
	PercentErr    error
}

func (m *MockCPUMonitor) Percent(interval time.Duration, percpu bool) ([]float64, error) {
	return m.PercentResult, m.PercentErr
}

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
