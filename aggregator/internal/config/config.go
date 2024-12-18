package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type Telemetry struct {
	Server        string `toml:"server"`
	RetryInterval int    `toml:"retry_interval"`
}

type Labels struct {
	Environment string `toml:"environment"`
	Service     string `toml:"service"`
}

type Logging struct {
	Level string `toml:"level"`
}

type Targets struct {
	Hosts       []string `toml:"hosts"`
	Frequencies []int    `toml:"frequencies"`
	Ports       []string `toml:"ports"`
}

type Config struct {
	Labels    Labels    `toml:"labels"`
	Logging   Logging   `toml:"logging"`
	Telemetry Telemetry `toml:"telemetry"`
	Targets   Targets   `toml:"targets"`
	WebAPI    struct {
		BaseURL string `toml:"base_url"`
		Timeout int    `toml:"timeout"` // in seconds
		Retries int    `toml:"retries"`
	} `toml:"web_api"`
}

func Load(path string) (*Config, error) {
	config := &Config{
		Labels:    Labels{},
		Logging:   Logging{},
		Telemetry: Telemetry{},
		Targets:   Targets{},
	}

	if _, err := toml.DecodeFile(path, config); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	// if missing fields, return an error
	if config.Labels.Environment == "" {
		return nil, fmt.Errorf("missing environment field in config")
	}
	if config.Labels.Service == "" {
		return nil, fmt.Errorf("missing service field in config")
	}
	if config.Logging.Level == "" {
		return nil, fmt.Errorf("missing level field in config")
	}
	if config.Telemetry.Server == "" {
		return nil, fmt.Errorf("missing server field in config")
	}
	if config.Targets.Hosts == nil {
		return nil, fmt.Errorf("missing hosts field in config")
	}
	if config.Telemetry.RetryInterval == 0 {
		return nil, fmt.Errorf("missing retry_interval field in config")
	}
	if config.Targets.Frequencies == nil {
		return nil, fmt.Errorf("missing frequencies field in config")
	}
	if config.Targets.Ports == nil {
		return nil, fmt.Errorf("missing ports field in config")
	}
	if len(config.Targets.Hosts) != len(config.Targets.Frequencies) {
		return nil, fmt.Errorf("hosts and frequencies fields must be equal in length")
	}
	if len(config.Targets.Hosts) != len(config.Targets.Ports) {
		return nil, fmt.Errorf("hosts and ports fields must be equal in length")
	}

	return config, nil
}
