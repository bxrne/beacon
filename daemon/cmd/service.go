package main

import (
	"context"

	"github.com/bxrne/beacon/daemon/internal/config"
	"github.com/bxrne/beacon/daemon/internal/server"
	"github.com/charmbracelet/log"
)

type Service struct {
	cfg    *config.Config
	log    *log.Logger
	server *server.HTTPServer
}

func NewService(cfg *config.Config, log *log.Logger) (*Service, error) {
	srv := server.NewHTTPServer(cfg, log)

	return &Service{
		cfg:    cfg,
		log:    log,
		server: srv,
	}, nil
}

func (s *Service) Run() error {
	s.log.Infof("Service initialized (%s)", s.cfg.Labels.Environment)
	return s.server.Start()
}

func (s *Service) Shutdown() {
	s.log.Info("Shutting down service...")
	if err := s.server.Shutdown(context.Background()); err != nil {
		s.log.Errorf("Error shutting down server: %v", err)
	}
}
