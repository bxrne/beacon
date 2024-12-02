package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type Monitoring struct {
	DiskPaths []string `toml:"disk_paths"`
	Frequency uint     `toml:"frequency"`
}

type Telemetry struct {
	Server       string `toml:"server"`
	ExClientPort string `toml:"ex_client_port"`
	ExClientBaud int    `toml:"ex_client_baud"`
}

type Labels struct {
	Environment string `toml:"environment"`
	Service     string `toml:"service"`
}

type Logging struct {
	Level string `toml:"level"`
}

type Config struct {
	Monitoring Monitoring `toml:"monitoring"`
	Labels     Labels     `toml:"labels"`
	Logging    Logging    `toml:"logging"`
	Telemetry  Telemetry  `toml:"telemetry"`
}

func Load(path string) (*Config, error) {
	config := &Config{}

	if _, err := toml.DecodeFile(path, config); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	return config, nil
}
