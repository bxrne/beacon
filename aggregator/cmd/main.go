package main

import (
	"log"

	"github.com/bxrne/beacon/aggregator/internal/config"
	"github.com/bxrne/beacon/aggregator/internal/logger"
	"github.com/bxrne/beacon/aggregator/pkg/poller"
)

func main() {
	cfg, err := config.Load("config.toml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	log := logger.NewLogger(cfg)
	log.Infof("Starting service %s in %s environment", cfg.Labels.Service, cfg.Labels.Environment)

	pollers := make([]*poller.Poller, 0)
	for i := 0; i < len(cfg.Targets.Hosts); i++ {
		p := poller.NewPoller(cfg.Targets.Hosts[i], cfg.Targets.Ports[i], cfg.Targets.Frequencies[i], cfg)
		pollers = append(pollers, p)
		log.Infof("Created poller for host %s with frequency %d", p.Host, p.Frequency)
	}

	for _, p := range pollers {
		log.Infof("Starting poller for host %s", p.Host)
		go p.Start()
	}

	select {} // Block forever
}
