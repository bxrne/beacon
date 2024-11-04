package config

import (
	"fmt"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/charmbracelet/log"
)

type MonitoringConfig struct {
	EnableCPU    bool     `toml:"enable_cpu"`
	EnableMemory bool     `toml:"enable_memory"`
	EnableDisk   bool     `toml:"enable_disk"`
	DiskPaths    []string `toml:"disk_paths"`
}

type Labels struct {
	Environment string `toml:"environment"`
	Service     string `toml:"service"`
	Frequency   int64  `toml:"frequency"`
}

type Logging struct {
	Level string `toml:"level"`
}

type Config struct {
	Monitoring MonitoringConfig `toml:"monitoring"`
	Labels     Labels           `toml:"labels"`
	Logging    Logging          `toml:"logging"`

	// NOTE: Computed fields
	FrequencyDuration time.Duration
	ParsedLogLevel    log.Level
}

func Load(path string) (*Config, error) {
	config := &Config{}

	if _, err := toml.DecodeFile(path, config); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	// NOTE: Convert frequency to duration
	config.FrequencyDuration = time.Duration(config.
		Labels.Frequency) * time.Second

	switch config.Logging.Level {
	case "debug":
		config.ParsedLogLevel = log.DebugLevel
	case "info":
		config.ParsedLogLevel = log.InfoLevel
	case "warn":
		config.ParsedLogLevel = log.WarnLevel
	case "error":
		config.ParsedLogLevel = log.ErrorLevel
	default:
		config.ParsedLogLevel = log.InfoLevel
	}

	return config, nil
}
