package poller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/bxrne/beacon/web/pkg/db"
	"github.com/charmbracelet/log"
	"gorm.io/gorm"
)

type CommandPoller struct {
	db         *gorm.DB
	logger     *log.Logger
	client     *http.Client
	pollTicker *time.Ticker
	stopChan   chan struct{}
}

func NewCommandPoller(db *gorm.DB, logger *log.Logger) *CommandPoller {
	return &CommandPoller{
		db:         db,
		logger:     logger,
		client:     &http.Client{Timeout: 5 * time.Second},
		pollTicker: time.NewTicker(5 * time.Second),
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
	var commands []db.Command
	if err := p.db.Preload("Device").Where("status = ?", "pending").Find(&commands).Error; err != nil {
		p.logger.Error("failed to query pending commands", "error", err)
		return
	}

	for _, cmd := range commands {
		go p.executeCommand(cmd)
	}
}

func (p *CommandPoller) executeCommand(cmd db.Command) {
	payload := map[string]string{
		"command": cmd.Name,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		p.updateCommandStatus(cmd, "failed", fmt.Sprintf("failed to marshal command: %v", err))
		return
	}

	url := fmt.Sprintf("http://%s:%d/cmd", cmd.Device.Name)
	resp, err := p.client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		p.updateCommandStatus(cmd, "failed", fmt.Sprintf("failed to send command: %v", err))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		p.updateCommandStatus(cmd, "failed", fmt.Sprintf("command failed with status: %d", resp.StatusCode))
		return
	}

	now := time.Now()
	p.updateCommandStatus(cmd, "sent", "")
	p.db.Model(&cmd).Update("sent_at", &now)
}

func (p *CommandPoller) updateCommandStatus(cmd db.Command, status, errorMsg string) {
	updates := map[string]interface{}{
		"status": status,
	}
	if errorMsg != "" {
		updates["error_msg"] = errorMsg
	}
	p.db.Model(&cmd).Updates(updates)
}
