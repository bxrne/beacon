package logger

import (
	"os"
	"sync"
	"time"

	"github.com/bxrne/beacon/api/pkg/config"
	"github.com/charmbracelet/log"
)

var (
	instance *log.Logger
	once     sync.Once
)

func parseLogLevel(level string) log.Level {
	switch level {
	case "debug":
		return log.DebugLevel
	case "info":
		return log.InfoLevel
	case "warn":
		return log.WarnLevel
	case "error":
		return log.ErrorLevel
	case "fatal":
		return log.FatalLevel
	default:
		return log.InfoLevel
	}
}

func NewLogger(cfg *config.Config) *log.Logger {
	once.Do(func() {
		logLevel := parseLogLevel(cfg.Logging.Level)

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
