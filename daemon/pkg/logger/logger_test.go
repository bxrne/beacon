package logger_test

import (
	"testing"

	"github.com/bxrne/beacon/pkg/config"
	"github.com/bxrne/beacon/pkg/logger"
	"github.com/charmbracelet/log"
	"github.com/stretchr/testify/assert"
)

// TEST: GIVEN a new logger request and no previous logger WHEN NewLogger is called THEN it should return a new logger
func TestNewLogger(t *testing.T) {
	cfg := &config.Config{
		Monitoring: config.Monitoring{},
		Labels:     config.Labels{Service: "beacon"},
		Logging:    config.Logging{Level: "debug"},
	}

	logr := logger.NewLogger(cfg)

	assert.NotNil(t, logr)
	assert.Equal(t, log.DebugLevel, logr.GetLevel())
	assert.Equal(t, "beacon", logr.GetPrefix())
}

// TEST: GIVEN a new logger request and a previous logger WHEN NewLogger is called THEN it should return the previous logger
func TestNewLoggerPrevious(t *testing.T) {
	cfg := &config.Config{
		Monitoring: config.Monitoring{},
		Labels:     config.Labels{Service: "beacon"},
		Logging:    config.Logging{Level: "debug"},
	}

	log := logger.NewLogger(cfg)
	log.Info("test")

	newLog := logger.NewLogger(cfg)

	assert.NotNil(t, newLog)
	assert.Equal(t, log, newLog)
}
