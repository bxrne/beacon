package server

import (
	"encoding/json"
	"net/http"

	"github.com/bxrne/beacon-web/pkg/db"
	"github.com/bxrne/beacon-web/pkg/metrics"
)

// handleMetric godoc
// @Summary      Submit metrics
// @Description  Submit metrics for a device
// @Tags         metrics
// @Accept       json
// @Produce      json
// @Param        X-DeviceID  header    string          true  "Device ID"
// @Param        metrics     body      metrics.DeviceMetrics  true  "Metrics data"
// @Success      200         {object}  metrics.DeviceMetrics
// @Failure      400         {object}  map[string]string
// @Failure      500         {object}  map[string]string
// @Router       /metric [post]
func (s *Server) handleMetric(w http.ResponseWriter, r *http.Request) {
	deviceID := r.Header.Get("X-DeviceID")
	if deviceID == "" {
		s.logger.Errorf("missing device ID")
		s.respondError(w, http.StatusBadRequest, "missing device ID")
		return
	}

	// Register the device if it doesn't exist
	if err := db.RegisterDevice(s.db, deviceID); err != nil {
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

	if err := db.PersistMetric(s.db, deviceMetrics, deviceID); err != nil {
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
// @Param        X-DeviceID  header    string  true  "Device ID"
// @Success      200         {object}  metrics.DeviceMetrics
// @Failure      400         {object}  map[string]string
// @Failure      404         {object}  map[string]string
// @Router       /metric [get]
func (s *Server) handleGetMetric(w http.ResponseWriter, r *http.Request) {
	deviceID := r.Header.Get("X-DeviceID")
	if deviceID == "" {
		s.logger.Errorf("missing device ID")
		s.respondError(w, http.StatusBadRequest, "missing device ID")
		return
	}

	var device db.Device
	if err := s.db.First(&device, "name = ?", deviceID).Error; err != nil {
		s.logger.Errorf("device not found: %v", err)
		s.respondError(w, http.StatusNotFound, "device not found")
		return
	}

	var metrics []db.Metric
	if err := s.db.Where("device_id = ?", device.ID).Find(&metrics).Error; err != nil {
		s.logger.Errorf("failed to get metrics: %v", err)
		s.respondError(w, http.StatusInternalServerError, "failed to get metrics")
		return
	}

	var deviceMetrics []db.Metric
	for _, m := range metrics {
		deviceMetrics = append(deviceMetrics, m)
	}
	s.respondJSON(w, http.StatusOK, deviceMetrics)
}

type healthResponse struct {
	Status string `json:"status"`
}

// handleHealth godoc
// @Summary      Health check
// @Description  Get the health status of the server
// @Tags         health
// @Produce      json
// @Success      200  {object}  healthResponse
// @Router       /health [get]
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	s.respondJSON(w, http.StatusOK, healthResponse{Status: "ok"})
}
