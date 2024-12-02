package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/bxrne/beacon-web/pkg/db"
	"github.com/bxrne/beacon-web/pkg/metrics"
)

// handleMetric godoc
// @Summary      Submit metrics
// @Description  Submit metrics for a device
// @Tags         metrics
// @Accept       json
// @Produce      json
// @Param        X-Hostname  header    string          true  "Device hostname"
// @Param        metrics     body      metrics.DeviceMetrics  true  "Metrics data"
// @Success      200         {object}  metrics.DeviceMetrics
// @Failure      400         {object}  map[string]string
// @Failure      500         {object}  map[string]string
// @Router       /metric [post]
func (s *Server) handleMetric(w http.ResponseWriter, r *http.Request) {
	hostname := r.Header.Get("X-Hostname")
	if hostname == "" {
		s.logger.Errorf("missing hostname")
		s.respondError(w, http.StatusBadRequest, "missing hostname")
		return
	}

	// Register the device if it doesn't exist
	if err := db.RegisterDevice(s.db, hostname); err != nil {
		s.logger.Errorf("failed to register device: %v", err)
		s.respondError(w, http.StatusInternalServerError, "failed to register device")
		return
	}

	var deviceMetrics metrics.DeviceMetrics
	if err := json.NewDecoder(r.Body).Decode(&deviceMetrics); err != nil {
		s.logger.Errorf("failed to decode request payload: %v", err)
		s.respondError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	if err := deviceMetrics.Validate(s.db); err != nil {
		s.logger.Errorf("invalid metrics: %v", err)
		s.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	s.metricsCache.SetMetrics(hostname, deviceMetrics)
	if err := db.PersistMetric(s.db, deviceMetrics, hostname); err != nil {
		s.logger.Errorf("failed to persist metrics: %v", err)
		s.respondError(w, http.StatusInternalServerError, "failed to persist metrics")
		return
	}

	s.respondJSON(w, http.StatusOK, deviceMetrics)
}

// handleGetMetric godoc
// @Summary      Get metrics
// @Description  Get metrics for a device
// @Tags         metrics
// @Produce      json
// @Param        X-Hostname  header    string  true  "Device hostname"
// @Success      200         {object}  metrics.DeviceMetrics
// @Failure      400         {object}  map[string]string
// @Failure      404         {object}  map[string]string
// @Router       /metric [get]
func (s *Server) handleGetMetric(w http.ResponseWriter, r *http.Request) {
	hostname := r.Header.Get("X-Hostname")
	if hostname == "" {
		s.logger.Errorf("missing hostname")
		s.respondError(w, http.StatusBadRequest, "missing hostname")
		return
	}

	deviceMetrics, exists := s.metricsCache.GetMetrics(hostname)
	if !exists {
		s.respondError(w, http.StatusNotFound, "no metrics found for device")
		return
	}

	s.respondJSON(w, http.StatusOK, deviceMetrics)
}

// handleHealth godoc
// @Summary      Health check
// @Description  Get the health status of the server
// @Tags         health
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /health [get]
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status": "OK",
		"time":   time.Now().Format(time.RFC3339),
	}

	s.respondJSON(w, http.StatusOK, response)
}
