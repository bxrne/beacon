package logger

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/bxrne/beacon/pkg/config"
	"github.com/charmbracelet/log"
)

var (
	instance *log.Logger
	once     sync.Once
)

func NewLogger(cfg *config.Config) *log.Logger {
	once.Do(func() {
		logLevel, err := log.ParseLevel(cfg.Logging.Level)
		if err != nil {
			fmt.Printf("Failed to parse log level: %v\n", err)
			os.Exit(1)
		}

		instance = log.NewWithOptions(os.Stdout, log.Options{
			Level:           logLevel,
			ReportCaller:    true,
			ReportTimestamp: true,
			TimeFormat:      time.RFC3339,
			Prefix:          cfg.Labels.Service,
		})
	})
	return instance
}
