package logger

import (
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
		instance = log.NewWithOptions(os.Stdout, log.Options{
			Level:           cfg.ParsedLogLevel,
			ReportCaller:    true,
			ReportTimestamp: true,
			TimeFormat:      time.RFC3339,
			Prefix:          cfg.Labels.Service,
		})
	})
	return instance
}
