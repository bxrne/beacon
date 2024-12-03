package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/bxrne/beacon-web/pkg/config"
	"github.com/bxrne/beacon-web/pkg/db"
	"github.com/bxrne/beacon-web/pkg/server"

	_ "github.com/bxrne/beacon-web/docs" // This line is necessary for go-swagger to find your docs
	"github.com/charmbracelet/log"
)

func main() {
	logger := log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
	})

	cfg, err := config.Load("config.toml")
	if err != nil {
		logger.Fatal("failed to load config", "error", err)
	}

	db, err := db.NewDatabase(cfg.Database.DSN)
	if err != nil {
		logger.Fatal("failed to connect to database", "error", err)
	}
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
