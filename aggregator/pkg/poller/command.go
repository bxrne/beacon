package poller

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
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

type Device struct {
	Name string `json:"name"`
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
	// hosts is each the target and ports matched
	var hosts []string
	hosts_num := len(p.cfg.Targets.Hosts)
	for i := 0; i < hosts_num; i++ {
		hosts = append(hosts, fmt.Sprintf("%s:%s", p.cfg.Targets.Hosts[i], p.cfg.Targets.Ports[i]))
	}

	for _, host := range hosts {
		// Create a new request with the X-DeviceID header
		req, err := http.NewRequest("GET", p.cfg.Telemetry.Server+"/api/command", nil)
		if err != nil {
			p.logger.Error("failed to create request", "error", err)
			continue
		}
		req.Header.Set("X-DeviceID", host)

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
		if err := json.NewDecoder(resp.Body).Decode(&commands); err != nil {
			p.logger.Error("failed to decode commands", "error", err, "host", host)
			continue
		}

		// Process each command
		for _, cmd := range commands {
			// Skip empty commands
			if cmd.Command == "" {
				p.logger.Warn("skipping empty command", "device", cmd.Device)
				continue
			}

			// Only process commands for the current host
			targetHost := cmd.Device
			if targetHost == host {
				p.logger.Infof("processing command", "command", cmd.Command, "host", host)
				if err := p.sendCommand(host, cmd.Command); err != nil {
					p.logger.Error("failed to send command", "error", err, "host", host)
					// Update command status to "failed"
					if err := p.updateCommandStatus(host, cmd.Command, "failed"); err != nil {
						p.logger.Error("failed to update command status", "error", err, "host", host)
					}
				} else {
					p.logger.Infof("successfully sent command", "command", cmd.Command, "host", host)
					// Update command status to "completed"
					if err := p.updateCommandStatus(host, cmd.Command, "completed"); err != nil {
						p.logger.Error("failed to update command status", "error", err, "host", host)
					}
				}
			}
		}
	}
}

func (p *CommandPoller) sendCommand(host, command string) error {
	// Connect to device
	// remove :port from host
	hostname := strings.Split(host, ":")[0]
	port := strings.Split(host, ":")[1]
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", hostname, port))
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer conn.Close()

	// Create JSON payload
	payload := struct {
		Command string `json:"command"`
	}{
		Command: command,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal command: %w", err)
	}

	// Construct HTTP POST request with JSON command
	request := fmt.Sprintf("POST /cmd HTTP/1.0\r\nContent-Type: application/json\r\nContent-Length: %d\r\n\r\n%s",
		len(jsonData), jsonData)

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
	if !strings.HasPrefix(resp, "HTTP/1.1 200 OK") && !strings.HasPrefix(resp, "HTTP/1.0 200 OK") {
		return fmt.Errorf("unexpected response: %s", resp)
	}

	return nil
}

func (p *CommandPoller) updateCommandStatus(device, command, status string) error {
	// Create JSON payload
	payload := struct {
		Device  string `json:"device"`
		Command string `json:"command"`
		Status  string `json:"status"`
	}{
		Device:  device,
		Command: command,
		Status:  status,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal command status: %w", err)
	}

	// Send request to update command status
	req, err := http.NewRequest("POST", p.cfg.Telemetry.Server+"/api/command/status", strings.NewReader(string(jsonData)))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to update command status: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update command status, status code: %d", resp.StatusCode)
	}

	return nil
}
