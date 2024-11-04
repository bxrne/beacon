package main

import (
	"fmt"
	"os"

	"github.com/bxrne/beacon/pkg/config"
	"github.com/bxrne/beacon/pkg/logger"
	"github.com/bxrne/beacon/pkg/stats"
)

func main() {
	cfg, err := config.Load("config.toml")
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	log := logger.NewLogger(cfg)

	log.Info("starting with",
		"frequency", cfg.FrequencyDuration,
		"environment", cfg.Labels.Environment,
	)

	stats.RunCollector(cfg, log)
}
