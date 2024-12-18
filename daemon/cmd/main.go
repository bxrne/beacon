package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/bxrne/beacon/daemon/internal/config"
	"github.com/bxrne/beacon/daemon/internal/logger"
)

func main() {
	cfg, err := config.Load("config.toml")
	if err != nil {
		panic(err)
	}

	log := logger.NewLogger(cfg)

	srv, err := NewService(cfg, log)
	if err != nil {
		log.Fatalf("Failed to create service: %v", err)
	}

	// Handle graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		srv.log.Info("Received shutdown signal")
		srv.Shutdown()
		srv.log.Info("Service stopped")
		os.Exit(0)
	}()

	if err := srv.Run(); err != nil {
		log.Fatalf("Failed to start service: %v", err)
	}
}
