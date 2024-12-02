package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type Labels struct {
	Environment string `toml:"environment"`
	Service     string `toml:"service"`
}

type Logging struct {
	Level string `toml:"level"`
}

type Server struct {
	Port            int      `toml:"port"`
	ShutdownTimeout int      `toml:"shutdown_timeout"`
	AllowedOrigins  []string `toml:"allowed_origins"`
	CacheTTL        int      `toml:"cache_ttl"`
}

type Config struct {
	Labels  Labels  `toml:"labels"`
	Logging Logging `toml:"logging"`
	Server  Server  `toml:"server"`
}

func Load(path string) (*Config, error) {
	config := &Config{}

	if _, err := toml.DecodeFile(path, config); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	return config, nil
}
