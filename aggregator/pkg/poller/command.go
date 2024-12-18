package poller

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/bxrne/beacon/aggregator/internal/config"
	"github.com/charmbracelet/log"
)

type CommandPoller struct {
	logger     *log.Logger
	client     *http.Client
	cfg        *config.Config
	pollTicker *time.Ticker
	stopChan   chan struct{}
}

type Command struct {
	Device  string `json:"device"`
	Command string `json:"command"`
}

func NewCommandPoller(cfg *config.Config, logger *log.Logger) *CommandPoller {
	return &CommandPoller{
		logger:     logger,
		client:     &http.Client{Timeout: 5 * time.Second},
		pollTicker: time.NewTicker(5 * time.Second),
		cfg:        cfg,
		stopChan:   make(chan struct{}),
	}
}

func (p *CommandPoller) Start() {
	go func() {
		for {
			select {
			case <-p.pollTicker.C:
				p.pollCommands()
			case <-p.stopChan:
				p.pollTicker.Stop()
				return
			}
		}
	}()
}

func (p *CommandPoller) Stop() {
	close(p.stopChan)
}

func (p *CommandPoller) pollCommands() {
	// Iterate over all configured hosts
	for _, host := range p.cfg.Targets.Hosts {
		deviceID := host // Adjust based on actual host structure

		// Create a new request with the X-DeviceID header
		req, err := http.NewRequest("GET", p.cfg.WebAPI.BaseURL+"/api/command", nil)
		if err != nil {
			p.logger.Error("failed to create request", "error", err)
			continue
		}
		req.Header.Set("X-DeviceID", deviceID)

		resp, err := p.client.Do(req)
		if err != nil {
			p.logger.Error("failed to get commands", "host", host, "error", err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			p.logger.Error("failed to get commands", "status", resp.StatusCode, "host", host)
			continue
		}

		var commands []Command
		p.logger.Info("received commands", "host", host)
		if err := json.NewDecoder(resp.Body).Decode(&commands); err != nil {
			p.logger.Error("failed to decode commands", "error", err, "host", host)
			continue
		}
		p.logger.Info("processing commands", "cmds", commands)

		// Process each command
		for _, cmd := range commands {
			hostAddr, port, err := net.SplitHostPort(cmd.Device)
			if err != nil {
				p.logger.Error("invalid device address", "device", cmd.Device, "error", err)
				continue
			}

			if err := p.sendCommand(hostAddr, port, cmd.Command); err != nil {
				p.logger.Error("failed to send command", "device", cmd.Device, "error", err)
				continue
			}
		}
	}
}

func (p *CommandPoller) sendCommand(host, port, command string) error {
	// Connect to device
	conn, err := net.Dial("tcp", net.JoinHostPort(host, port))
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer conn.Close()

	// Construct HTTP POST request with command
	request := fmt.Sprintf("POST /cmd HTTP/1.0\r\nContent-Length: %d\r\n\r\n%s",
		len(command), command)

	// Send request
	if _, err = conn.Write([]byte(request)); err != nil {
		return fmt.Errorf("failed to send command: %w", err)
	}

	// Read response
	response := make([]byte, 1024)
	n, err := conn.Read(response)
	if err != nil {

		return fmt.Errorf("failed to read response: %w", err)
	}

	// Check response status
	resp := string(response[:n])
	if resp[:15] != "HTTP/1.1 200 OK" {
		return fmt.Errorf("unexpected response: %s", resp)
	}

	return nil
}
