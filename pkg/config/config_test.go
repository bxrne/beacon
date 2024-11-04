package config_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/bxrne/beacon/pkg/config"
	"github.com/charmbracelet/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name           string
		configContent  string
		expectedConfig *config.Config
		expectedError  string
	}{
		{ // TEST: GIVEN a valid (full) config file WHEN Load is called THEN no error is returned AND the config is loaded correctly
			name: "valid config with all fields",
			configContent: `
[monitoring]
enable_cpu = true
enable_memory = true
enable_disk = true
disk_paths = ["/path1", "/path2"]
frequency = 60

[labels]
environment = "production"
service = "test-service"
frequency = 60

[logging]
level = "debug"
file = "app.log"
`,
			expectedConfig: &config.Config{
				Monitoring: config.MonitoringConfig{
					EnableCPU:    true,
					EnableMemory: true,
					EnableDisk:   true,
					DiskPaths:    []string{"/path1", "/path2"},
				},
				Labels: config.Labels{
					Environment: "production",
					Service:     "test-service",
					Frequency:   60,
				},
				Logging: config.Logging{
					Level: "debug",
				},
				FrequencyDuration: 60 * time.Second,
				ParsedLogLevel:    log.DebugLevel,
			},
		},
		{ // TEST: GIVEN a valid config file with minimal fields WHEN Load is called THEN no error is returned AND the config is loaded correctly
			name: "valid config with minimal fields",
			configContent: `
[monitoring]
enable_cpu = false
enable_memory = false
enable_disk = false
frequency = 30

[labels]
environment = "dev"
service = "minimal"
frequency = 30

[logging]
level = "info"
`,
			expectedConfig: &config.Config{
				Monitoring: config.MonitoringConfig{
					EnableCPU:    false,
					EnableMemory: false,
					EnableDisk:   false,
				},
				Labels: config.Labels{
					Environment: "dev",
					Service:     "minimal",
					Frequency:   30,
				},
				Logging: config.Logging{
					Level: "info",
				},
				FrequencyDuration: 30 * time.Second,
				ParsedLogLevel:    log.InfoLevel, // Default log level if not explicitly set
			},
		},
		{ // TEST: GIVEN a config file with all log levels WHEN Load is called THEN no error is returned AND the config is loaded correctly
			name: "config with all log levels",
			configContent: `
[monitoring]
enable_cpu = true
frequency = 10

[labels]
environment = "test"
service = "logger"
frequency = 10

[logging]
level = "error"
file = "error.log"
`,
			expectedConfig: &config.Config{
				Monitoring: config.MonitoringConfig{
					EnableCPU: true,
				},
				Labels: config.Labels{
					Environment: "test",
					Service:     "logger",
					Frequency:   10,
				},
				Logging: config.Logging{
					Level: "error",
				},
				FrequencyDuration: 10 * time.Second,
				ParsedLogLevel:    log.ErrorLevel,
			},
		},
		{
			name: "invalid TOML syntax",
			configContent: `
[monitoring
enable_cpu = true
`,
			expectedError: "failed to decode config file:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			configPath := filepath.Join(tmpDir, "config.toml")
			err := os.WriteFile(configPath, []byte(tt.configContent), 0644)
			require.NoError(t, err, "Failed to write test config file")

			config, err := config.Load(configPath)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, config)

			assert.Equal(t, tt.expectedConfig.Monitoring.EnableCPU, config.Monitoring.EnableCPU)
			assert.Equal(t, tt.expectedConfig.Monitoring.EnableMemory, config.Monitoring.EnableMemory)
			assert.Equal(t, tt.expectedConfig.Monitoring.EnableDisk, config.Monitoring.EnableDisk)
			assert.Equal(t, tt.expectedConfig.Monitoring.DiskPaths, config.Monitoring.DiskPaths)

			assert.Equal(t, tt.expectedConfig.Labels.Environment, config.Labels.Environment)
			assert.Equal(t, tt.expectedConfig.Labels.Service, config.Labels.Service)
			assert.Equal(t, tt.expectedConfig.Labels.Frequency, config.Labels.Frequency)

			assert.Equal(t, tt.expectedConfig.Logging.Level, config.Logging.Level)

			assert.Equal(t, tt.expectedConfig.FrequencyDuration, config.FrequencyDuration)
			assert.Equal(t, tt.expectedConfig.ParsedLogLevel, config.ParsedLogLevel)
		})
	}
}

func TestLoadNonExistentFile(t *testing.T) {
	_, err := config.Load("non_existent_file.toml")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to decode config file")
}
