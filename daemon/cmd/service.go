package main

import (
	"fmt"
	"net"

	"github.com/bxrne/beacon/daemon/internal/config"
	"github.com/bxrne/beacon/daemon/internal/logger"
	"github.com/bxrne/beacon/daemon/internal/stats"
	"github.com/charmbracelet/log"
)

const (
	startByte = 0x02
	endByte   = 0x03
)

type Service struct {
	cfg           *config.Config
	log           *log.Logger
	hostMonitor   stats.HostMon
	memoryMonitor stats.MemoryMon
	diskMonitor   stats.DiskMon
}

func NewService(cfgPath string) (*Service, error) {
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return nil, fmt.Errorf("error loading config: %w", err)
	}

	log := logger.NewLogger(cfg)
	log.Infof("Service initialized (%s)", cfg.Labels.Environment)

	hostMonitor := stats.HostMon{}
	memoryMonitor := stats.MemoryMon{}
	diskMonitor := stats.DiskMon{}

	return &Service{
		cfg:           cfg,
		log:           log,
		hostMonitor:   hostMonitor,
		memoryMonitor: memoryMonitor,
		diskMonitor:   diskMonitor,
	}, nil
}

func (s *Service) Run() {
	listenAddr := fmt.Sprintf(":%d", s.cfg.Server.Port)
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		s.log.Fatalf("Failed to listen on %s: %v", listenAddr, err)
	}
	defer listener.Close()

	s.log.Infof("Listening on %s", listenAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			s.log.Errorf("Failed to accept connection: %v", err)
			continue
		}

		go s.handleConnection(conn)
	}
}

func (s *Service) handleConnection(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		s.log.Errorf("Failed to read from connection: %v", err)
		return
	}

	request := string(buffer[:n])
	request = request[:len(request)-2] // Remove the trailing \r\n
	s.log.Debugf("Received request: %s", request)

	response, err := s.processRequest()
	if err != nil {
		s.log.Errorf("Failed to process request: %v", err)
		return
	}

	_, err = conn.Write(response)
	if err != nil {
		s.log.Errorf("Failed to send response: %v", err)
	}
}

func (s *Service) processRequest() ([]byte, error) {
	metrics, err := stats.Collect(s.cfg, s.hostMonitor, s.memoryMonitor, s.diskMonitor)
	if err != nil {
		return nil, fmt.Errorf("failed to collect metrics: %w", err)
	}

	payload := metrics.String()
	payloadLength := len(payload)
	message := make([]byte, payloadLength+3)
	message[0] = startByte
	message[1] = byte(payloadLength)
	copy(message[2:], []byte(payload))
	message[payloadLength+2] = endByte

	return message, nil
}
