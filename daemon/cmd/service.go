package main

import (
	"fmt"
	"time"

	"github.com/bxrne/beacon/daemon/internal/config"
	"github.com/bxrne/beacon/daemon/internal/logger"
	"github.com/bxrne/beacon/daemon/internal/stats"
	"github.com/bxrne/beacon/daemon/internal/uploader"
	"github.com/charmbracelet/log"
)

type Service struct {
	cfg           *config.Config
	log           *log.Logger
	hostMonitor   stats.HostMon
	memoryMonitor stats.MemoryMon
	diskMonitor   stats.DiskMon
	uploader      *uploader.Uploader
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

	uploader, err := uploader.NewUploader(cfg)
	if err != nil {
		return nil, fmt.Errorf("error creating uploader: %w", err)
	}

	return &Service{
		cfg:           cfg,
		log:           log,
		hostMonitor:   hostMonitor,
		memoryMonitor: memoryMonitor,
		diskMonitor:   diskMonitor,
		uploader:      uploader,
	}, nil
}

func (s *Service) Run() {
	go s.uploader.Start()

	ticker := time.NewTicker(time.Duration(s.cfg.Monitoring.Frequency) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		metrics, err := stats.Collect(s.cfg, s.hostMonitor, s.memoryMonitor, s.diskMonitor)
		if err != nil {
			s.log.Error("Failed to collect metrics", "error", err)
			continue
		}

		s.log.Debug(metrics.String())

		s.uploader.AddToQueue(metrics)
	}
}
