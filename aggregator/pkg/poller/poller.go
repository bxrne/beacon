package poller

import (
	"fmt"
	"net"
	"strings"
	"time"
)

// Poller is a service that will send request objects at a frequency to a host

// Poller is a service that will send request objects at a frequency to a host
type Poller struct {
	Host      string
	Port      string
	Frequency int
}

func NewPoller(host, port string, frequency int) *Poller {
	return &Poller{
		Host:      host,
		Port:      port,
		Frequency: frequency,
	}
}

// Start begins the polling process
func (p *Poller) Start() {
	for {
		p.sendRequest()
		time.Sleep(time.Duration(p.Frequency) * time.Second)
	}
}

// sendRequest sends the request object to the host
func (p *Poller) sendRequest() {
	conn, err := net.Dial("tcp", net.JoinHostPort(p.Host, p.Port))
	if err != nil {
		fmt.Printf("Failed to connect to %s:%s - %v\n", p.Host, p.Port, err)
		return
	}
	defer conn.Close()

	// Send GET /metric request
	request := "GET /metric HTTP/1.0\r\n\r\n"
	_, err = conn.Write([]byte(request))
	if err != nil {
		fmt.Printf("Failed to send request: %v\n", err)
		return
	}

	// Receive response
	response := make([]byte, 1024)
	n, err := conn.Read(response)
	if err != nil {
		fmt.Printf("Failed to read response: %v\n", err)
		return
	}
	responseStr := string(response[:n])

	// Parse the response
	if idx := strings.Index(responseStr, "\r\n\r\n"); idx != -1 {
		payload := responseStr[idx+4:]
		fmt.Printf("Received payload from %s:%s - %s\n", p.Host, p.Port, payload)
	} else {
		// No headers, assume entire response is payload
		fmt.Printf("Received payload from %s:%s - %s\n", p.Host, p.Port, responseStr)
	}
}

// Stop ends the polling process
func (p *Poller) Stop() {
	// Implement stop logic if needed
}
