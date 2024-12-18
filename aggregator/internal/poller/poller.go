package poller

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/bxrne/beacon/aggregator/internal/config"
	"github.com/charmbracelet/log"
)

type Poller struct {
	cfg    *config.Config
	logger *log.Logger
}

func NewPoller(cfg *config.Config, logger *log.Logger) *Poller {
	return &Poller{
		cfg:    cfg,
		logger: logger,
	}
}

func (p *Poller) Start() {
	for _, host := range p.cfg.Targets.Hosts {
		go p.pollHost(host)
	}
}

func (p *Poller) pollHost(host string) {
	for {
		conn, err := net.Dial("tcp", host)
		if err != nil {
			p.logger.Errorf("Failed to connect to %s: %v", host, err)
			time.Sleep(time.Duration(p.cfg.Telemetry.RetryInterval) * time.Second)
			continue
		}

		p.logger.Debugf("Connected to %s", host)

		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			p.logger.Errorf("Failed to read from %s: %v", host, err)
			conn.Close()
			time.Sleep(time.Duration(p.cfg.Telemetry.RetryInterval) * time.Second)
			continue
		}

		p.logger.Debugf("Received %d bytes from %s", n, host)

		payload, err := p.extractPayload(buf)
		if err != nil {
			p.logger.Errorf("Failed to parse metrics from %s: %v", host, err)
			conn.Close()
			time.Sleep(time.Duration(p.cfg.Telemetry.RetryInterval) * time.Second)
			continue
		}

		if err := p.sendMetricsToAPI(host, payload); err != nil {
			p.logger.Errorf("Failed to send metrics to API: %v", err)
		}

		conn.Close()
		time.Sleep(time.Duration(p.cfg.Telemetry.RetryInterval) * time.Second)
	}
}

func (p *Poller) extractPayload(data []byte) ([]byte, error) {
	if len(data) < 3 {
		return nil, errors.New("invalid response format: payload too short")
	}

	// Log the raw payload for debugging
	p.logger.Debug("Raw payload", "data", fmt.Sprintf("%x", data))

	// Check for daemon format (STX + LEN + PAYLOAD + ETX)
	if data[0] == 0x02 {
		if data[len(data)-1] != 0x03 {
			return nil, errors.New("invalid response format: missing ETX")
		}
		length := int(data[1])
		if len(data) < length+3 {
			return nil, errors.New("invalid payload length")
		}

		// Extract payload without STX, length byte, and ETX
		payload := make([]byte, length)
		copy(payload, data[2:2+length])
		return payload, nil
	}

	// Check for HTTP response format
	if string(data[:4]) == "HTTP" {
		parts := bytes.Split(data, []byte("\r\n\r\n"))
		if len(parts) != 2 {
			return nil, errors.New("invalid HTTP response format")
		}
		payload := parts[1]

		// Remove the end byte if it exists
		if len(payload) > 0 && payload[len(payload)-1] == 0x03 {
			payload = payload[:len(payload)-1]
		}

		return payload, nil
	}

	// If not a recognized format, return as-is
	return data, nil
}

func (p *Poller) sendMetricsToAPI(host string, payload []byte) error {
	req, err := http.NewRequest("POST", p.cfg.WebAPI.BaseURL+"/metrics", bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-DeviceID", host)

	client := &http.Client{
		Timeout: time.Duration(p.cfg.WebAPI.Timeout) * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		p.logger.Errorf("Unexpected status code: %d, response body: %s", resp.StatusCode, string(body))
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
