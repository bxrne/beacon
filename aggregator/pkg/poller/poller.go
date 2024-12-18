package poller

import (
	"net"
	"time"

	"github.com/bxrne/beacon/aggregator/internal/config"
	"github.com/bxrne/beacon/aggregator/internal/logger"
	"github.com/bxrne/beacon/aggregator/pkg/bproto"
	"github.com/charmbracelet/log"
)

// Poller is a service that will send request objects at a frequency to a host
type Poller struct {
	Host      string
	Port      string
	Frequency int
	logger    *log.Logger
}

func NewPoller(host, port string, frequency int, cfg *config.Config) *Poller {
	log := logger.NewLogger(cfg)
	return &Poller{
		Host:      host,
		Port:      port,
		Frequency: frequency,
		logger:    log,
	}
}

// Start begins the polling process
func (p *Poller) Start() {
	for {
		p.sendRequest()
		time.Sleep(time.Duration(p.Frequency) * time.Second)
	}
}

// sendRequest sends to host
func (p *Poller) sendRequest() {
	conn, err := net.Dial("tcp", net.JoinHostPort(p.Host, p.Port))
	if err != nil {
		p.logger.Errorf("Failed to connect to %s:%s: %v\n", p.Host, p.Port, err)
		return
	}
	defer conn.Close()

	// Send GET /metric request
	request := "GET /metric HTTP/1.0\r\n\r\n"
	_, err = conn.Write([]byte(request))
	if err != nil {
		p.logger.Errorf("Failed to send request to %s:%s: %v", p.Host, p.Port, err)
		return
	}

	// Receive response
	response := make([]byte, 1024)
	n, err := conn.Read(response)
	if err != nil {
		p.logger.Errorf("Failed to read response from %s:%s: %v", p.Host, p.Port, err)
		return
	}

	_, err = bproto.ParseResponse(response)
	if err != nil {
		p.logger.Errorf("Failed to parse response from %s:%s: %v", p.Host, p.Port, err)
		return
	}

	p.logger.Debugf("Received %d bytes from %s:%s", n, p.Host, p.Port)
}
