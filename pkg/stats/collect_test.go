package stats_test

import (
	"testing"

	"github.com/bxrne/beacon/pkg/config"
	"github.com/bxrne/beacon/pkg/stats"
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

	res, err := stats.Collect(cfg)

	require.Nil(t, err)
	require.Zero(t, res.CPUUsage)
	require.Zero(t, res.MemoryUsage)
	for _, v := range res.DiskUsage {
		require.Zero(t, v)
	}
}
