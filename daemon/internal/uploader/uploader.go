package uploader

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/bxrne/beacon/daemon/internal/config"
	"github.com/bxrne/beacon/daemon/internal/stats"
)

type Uploader struct {
	cfg           *config.Config
	queue         []stats.DeviceMetrics
	mu            sync.Mutex
	retryInterval time.Duration
}

func NewUploader(cfg *config.Config) (*Uploader, error) {
	if cfg.Telemetry.RetryInterval <= 0 {
		return nil, fmt.Errorf("invalid retry interval: %d", cfg.Telemetry.RetryInterval)
	}

	return &Uploader{
		cfg:           cfg,
		queue:         make([]stats.DeviceMetrics, 0),
		retryInterval: time.Duration(cfg.Telemetry.RetryInterval) * time.Second,
	}, nil
}

func (u *Uploader) AddToQueue(metrics stats.DeviceMetrics) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.queue = append(u.queue, metrics)
}

func (u *Uploader) Start() {
	ticker := time.NewTicker(u.retryInterval)
	defer ticker.Stop()

	for range ticker.C {
		u.processQueue()
	}
}

func (u *Uploader) processQueue() {
	u.mu.Lock()
	defer u.mu.Unlock()

	for len(u.queue) > 0 {
		metrics := u.queue[0]
		if err := u.send(metrics); err != nil {
			fmt.Printf("Failed to send metrics: %v\n", err)
			time.Sleep(u.retryInterval) // Retry after the specified interval
			continue
		}
		u.queue = u.queue[1:]
	}
}

func (u *Uploader) send(metrics stats.DeviceMetrics) error {
	url := u.cfg.Telemetry.Server + "/metric"
	data, err := json.Marshal(metrics)
	if err != nil {
		return fmt.Errorf("failed to marshal metrics: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-DeviceID", u.cfg.Labels.Service)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-OK response: %s", resp.Status)
	}

	return nil
}
