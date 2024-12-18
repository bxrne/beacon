package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/bxrne/beacon/daemon/internal/config"
	"github.com/bxrne/beacon/daemon/internal/stats"
	"github.com/charmbracelet/log"
)

type HTTPServer struct {
	cfg    *config.Config
	logger *log.Logger
	server *http.Server
}

func NewHTTPServer(cfg *config.Config, logger *log.Logger) *HTTPServer {
	return &HTTPServer{
		cfg:    cfg,
		logger: logger,
	}
}

func (s *HTTPServer) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/metric", s.handleMetrics)
	mux.HandleFunc("/cmd", s.handleCommand)

	addr := fmt.Sprintf(":%d", s.cfg.Server.Port)
	s.logger.Infof("HTTP server listening on port %d", s.cfg.Server.Port)

	s.server = &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	return s.server.ListenAndServe()
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	if s.server != nil {
		return s.server.Shutdown(ctx)
	}
	return nil
}

func (s *HTTPServer) handleCommand(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.logger.Errorf("Failed to read command body: %v", err)
		http.Error(w, "Failed to read command", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var cmd struct {
		Command string `json:"command"`
		Value   int    `json:"value,omitempty"`
	}
	if err := json.Unmarshal(body, &cmd); err != nil {
		s.logger.Errorf("Failed to parse command: %v", err)
		http.Error(w, "Invalid command format", http.StatusBadRequest)
		return
	}

	// Handle different commands
	switch cmd.Command {
	case "notify":
		if err := stats.SendNotification("Beacon Alert", "This is a test notification", s.logger); err != nil {
			s.logger.Error("failed to send notification", "error", err)
			http.Error(w, "Failed to execute command", http.StatusInternalServerError)
			return
		}
	case "brightness":
		if err := stats.SetScreenBrightness(cmd.Value, s.logger); err != nil {
			s.logger.Error("failed to set brightness", "error", err)
			http.Error(w, "Failed to execute command", http.StatusInternalServerError)
			return
		}
	default:
		s.logger.Warn("unknown command received", "command", cmd.Command)
		http.Error(w, "Unknown command", http.StatusBadRequest)
		return
	}

	// Add proper HTTP headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Command executed successfully",
	})
}

func (s *HTTPServer) handleMetrics(w http.ResponseWriter, r *http.Request) {
	deviceMetrics, err := stats.CollectMetrics()
	if err != nil {
		s.logger.Error("failed to collect metrics", "error", err)
		http.Error(w, "Failed to collect metrics", http.StatusInternalServerError)
		return
	}

	s.logger.Debug("metrics collected successfully")

	payload := deviceMetrics.String()
	payloadLength := len(payload)

	// Construct message according to protocol: STX + LEN + PAYLOAD + ETX
	message := make([]byte, payloadLength+3)
	message[0] = 0x02 // STX
	message[1] = byte(payloadLength)
	copy(message[2:], []byte(payload))
	message[payloadLength+2] = 0x03 // ETX

	w.Write(message)
}
