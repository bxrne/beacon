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
	ReadTimeout     int      `toml:"read_timeout"`
	WriteTimeout    int      `toml:"write_timeout"`
	IdleTimeout     int      `toml:"idle_timeout"`
	CacheTTL        int      `toml:"cache_ttl"`
}

type Database struct {
	DSN string `toml:"dsn"`
}

type Metrics struct {
	Types    []string `toml:"types"`
	Units    []string `toml:"units"`
	Commands []string `toml:"commands"`
}

type CommandType struct {
	Name string `toml:"name"`
}

type Config struct {
	Labels       Labels        `toml:"labels"`
	Logging      Logging       `toml:"logging"`
	Server       Server        `toml:"server"`
	Database     Database      `toml:"database"`
	Metrics      Metrics       `toml:"metrics"`
	CommandTypes []CommandType `toml:"command_types"`
}

func Load(path string) (*Config, error) {
	var cfg Config
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	return &cfg, nil
}
