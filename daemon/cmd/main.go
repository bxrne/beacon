package main

import (
	"fmt"
	"os"

	"github.com/bxrne/beacon/daemon/internal/server"
)

func main() {
	svc, err := NewService("config.toml")
	if err != nil {
		fmt.Printf("Error creating service: %v\n", err)
		os.Exit(1)
	}

	httpServer := server.NewHTTPServer(svc.cfg, svc.log)
	err = httpServer.Start()
	if err != nil {
		svc.log.Errorf("Failed to start HTTP server: %v", err)
		os.Exit(1)
	}

	svc.Run()
}
