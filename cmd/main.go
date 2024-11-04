package main

import (
	"fmt"
	"os"
	"time"

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

	cpuMonitor := stats.CPUMon{}
	memoryMonitor := stats.MemoryMon{}
	diskMonitor := stats.DiskMon{}

	for {
		metrics, err := stats.Collect(cfg, cpuMonitor, memoryMonitor, diskMonitor)
		if err != nil {
			fmt.Printf("Failed to collect metrics: %v\n", err)
			continue
		}

		log.Debug(metrics.String())

		time.Sleep(cfg.FrequencyDuration)
	}
}
