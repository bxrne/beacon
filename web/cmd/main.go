package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bxrne/beacon/web/internal/config"
	"github.com/bxrne/beacon/web/internal/db"
	"github.com/bxrne/beacon/web/internal/logger"
	"github.com/bxrne/beacon/web/internal/server"

	_ "github.com/bxrne/beacon/web/docs" // This line is necessary for go-swagger to find your docs
)

// @title Beacon API
// @version 1.0
// @host beacon-web.fly.dev
// @basePath /api/
// @description Collects device and metric data from clients
func main() {
	cfg, err := config.Load("config.toml")
	if err != nil {
		fmt.Printf("failed to load config: %v\n", err)
		os.Exit(1)
	}

	logger := logger.NewLogger(cfg)
	logger.Infof("Starting Service=%s Environment=%s", cfg.Labels.Service, cfg.Labels.Environment)

	db, err := db.NewDatabase(cfg)
	if err != nil {
		logger.Fatal("failed to connect to database", "error", err)
	}
	logger.Info("Connected to database")

	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	srv := server.New(cfg, logger, db)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := srv.Start(ctx); err != nil {
		logger.Fatal("failed to start server", "error", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	if err := srv.Stop(ctx); err != nil {
		logger.Error("failed to stop server", "error", err)
	}
}
