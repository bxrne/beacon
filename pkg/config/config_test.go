package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/bxrne/beacon/pkg/config"
	"github.com/charmbracelet/log"
	"github.com/stretchr/testify/assert"
)

// TEST: GIVEN a cfg file WHEN Load is called THEN it should return a valid config
func TestLoad(t *testing.T) {
	configContent := `
frequency = 5
log_level = "debug"

[monitoring]
enable_cpu = true
enable_memory = true
enable_disk = true
disk_paths = ["/", "/dev/disk0"]

[labels]
environment = "production"
service = "beacon"
`
	tmpFile, err := os.CreateTemp("", "config-*.toml")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write([]byte(configContent))
	assert.NoError(t, err)
	assert.NoError(t, tmpFile.Close())

	cfg, err := config.Load(tmpFile.Name())
	assert.NoError(t, err)

	assert.Equal(t, int64(5), cfg.Frequency)
	assert.Equal(t, log.DebugLevel, cfg.ParsedLogLevel)
	assert.Equal(t, true, cfg.Monitoring.EnableCPU)
	assert.Equal(t, true, cfg.Monitoring.EnableMemory)
	assert.Equal(t, true, cfg.Monitoring.EnableDisk)
	assert.Equal(t, []string{"/", "/dev/disk0"}, cfg.Monitoring.DiskPaths)
	assert.Equal(t, "production", cfg.Labels.Environment)
	assert.Equal(t, "beacon", cfg.Labels.Service)
	assert.Equal(t, 5*time.Second, cfg.FrequencyDuration)
}
