package main

import (
	"fmt"
	"time"

	"github.com/bxrne/beacon/daemon/pkg/config"
	"github.com/bxrne/beacon/daemon/pkg/logger"
	"github.com/bxrne/beacon/daemon/pkg/stats"
	"github.com/charmbracelet/log"
)

type Service struct {
	cfg           *config.Config
	log           *log.Logger
	hostMonitor   stats.HostMon
	memoryMonitor stats.MemoryMon
	diskMonitor   stats.DiskMon
}

func NewService(cfgPath string) (*Service, error) {
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return nil, fmt.Errorf("error loading config: %w", err)
	}

	log := logger.NewLogger(cfg)
	log.Infof("Service initialized (%s)", cfg.Labels.Environment)

	hostMonitor := stats.HostMon{}
	memoryMonitor := stats.MemoryMon{}
	diskMonitor := stats.DiskMon{}

	return &Service{
		cfg:           cfg,
		log:           log,
		hostMonitor:   hostMonitor,
		memoryMonitor: memoryMonitor,
		diskMonitor:   diskMonitor,
	}, nil
}

func (s *Service) Run() {
	ticker := time.NewTicker(time.Duration(s.cfg.Monitoring.Frequency) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		metrics, err := stats.Collect(s.cfg, s.hostMonitor, s.memoryMonitor, s.diskMonitor)
		if err != nil {
			s.log.Error("Failed to collect metrics", "error", err)
			continue
		}

		s.log.Debug(metrics.String())

		if err := stats.Send(s.cfg, metrics); err != nil {
			s.log.Error("Failed to send metrics", "error", err)
		}
	}
}
