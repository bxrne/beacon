package server

import (
	"fmt"
	"net/http"

	"github.com/bxrne/beacon/daemon/internal/config"
	"github.com/bxrne/beacon/daemon/internal/stats"
	"github.com/charmbracelet/log"
)

type HTTPServer struct {
	cfg    *config.Config
	logger *log.Logger
}

func NewHTTPServer(cfg *config.Config, logger *log.Logger) *HTTPServer {
	return &HTTPServer{
		cfg:    cfg,
		logger: logger,
	}
}

func (s *HTTPServer) Start() error {
	http.HandleFunc("/metric", s.handleMetrics)

	addr := fmt.Sprintf(":%d", s.cfg.Server.Port)
	s.logger.Infof("HTTP server listening on port %d", s.cfg.Server.Port)

	return http.ListenAndServe(addr, nil)
}

func (s *HTTPServer) handleMetrics(w http.ResponseWriter, r *http.Request) {
	deviceMetrics, err := stats.CollectMetrics()
	if err != nil {
		s.logger.Errorf("Failed to collect metrics: %v", err)
		http.Error(w, "Failed to collect metrics", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(deviceMetrics.String()))
}
