package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bxrne/beacon/api/pkg/db"
	"github.com/bxrne/beacon/api/pkg/metrics"
	"gorm.io/gorm"
)

type errorResponse struct {
	Error string `json:"error"`
}

// handleMetric godoc
// @Summary      Submit metrics
// @Description  Submit metrics for a device
// @Tags         metrics
// @Accept       json
// @Produce      json
// @Param        X-DeviceID  header    string          true  "Device ID"
// @Param        metrics     body      metrics.DeviceMetrics  true  "Metrics data"
// @Success      200         {object}  metrics.DeviceMetrics
// @Failure      400         {object}  errorResponse
// @Failure      500         {object}  errorResponse
// @Router       /metric [post]
func (s *Server) handleMetric(w http.ResponseWriter, r *http.Request) {
	deviceID := r.Header.Get("X-DeviceID")
	if deviceID == "" {
		res := errorResponse{Error: "missing device ID"}
		s.logger.Errorf("handleMetric: %s", res.Error)
		s.respondJSON(w, http.StatusBadRequest, res)
		return
	}

	// Register the device if it doesn't exist
	if err := db.RegisterDevice(s.db, deviceID); err != nil {
		res := errorResponse{Error: "failed to register device"}
		s.logger.Errorf("handleMetric: %s", res.Error)
		s.respondJSON(w, http.StatusInternalServerError, res)
		return
	}

	var deviceMetrics metrics.DeviceMetrics
	if err := json.NewDecoder(r.Body).Decode(&deviceMetrics); err != nil {
		res := errorResponse{Error: "failed to decode request body"}
		s.logger.Errorf("handleMetric: %s", res.Error)
		s.respondJSON(w, http.StatusBadRequest, res)
		return
	}

	for _, metric := range deviceMetrics.Metrics {
		if metric.Type == "" || metric.Unit == "" || metric.Value == "" || metric.RecordedAt == "" {
			res := errorResponse{Error: "all metric fields must be filled"}
			s.logger.Errorf("handleMetric: %s", res.Error)
			s.respondJSON(w, http.StatusBadRequest, res)
			return
		}
	}

	if err := deviceMetrics.Validate(s.db); err != nil {
		res := errorResponse{Error: "invalid metrics: " + err.Error()}
		s.logger.Errorf("handleMetric: %s", res.Error)
		s.respondJSON(w, http.StatusBadRequest, res)
		return
	}

	if err := metrics.PersistMetric(s.db, deviceMetrics, deviceID); err != nil {
		res := errorResponse{Error: "failed to persist metrics"}
		s.logger.Errorf("handleMetric: %s", res.Error, err)
		s.respondJSON(w, http.StatusBadRequest, res)
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
		res := errorResponse{Error: "missing device ID"}
		s.logger.Errorf("handleGetMetric: %s", res.Error)
		s.respondJSON(w, http.StatusBadRequest, res)
		return
	}

	var device db.Device
	if err := s.db.First(&device, "name = ?", deviceID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			res := errorResponse{Error: "device not found"}
			s.logger.Errorf("handleGetMetric: %s", res.Error)
			s.respondJSON(w, http.StatusNotFound, res)
		} else {
			res := errorResponse{Error: "failed to query device"}
			s.logger.Errorf("handleGetMetric: %s", res.Error)
			s.respondJSON(w, http.StatusInternalServerError, res)
		}
		return
	}

	var metrics []db.Metric
	if err := s.db.Preload("Type").Preload("Unit").Where("device_id = ?", device.ID).Find(&metrics).Error; err != nil {
		res := errorResponse{Error: "failed to get metrics"}
		s.logger.Errorf("handleGetMetric: %s", res.Error)
		s.respondJSON(w, http.StatusInternalServerError, res)
		return
	}

	var deviceMetrics []db.Metric
	deviceMetrics = append(deviceMetrics, metrics...)
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

// handleGetDevices godoc
// @Summary List devices
// @Description Find all registered devices
// @Tags devices
// @Produce json
// @Success 200 {object} []string
// @Router /device [get]
func (s *Server) handleGetDevices(w http.ResponseWriter, r *http.Request) {
	var devices []db.Device
	if err := s.db.Find(&devices).Error; err != nil {
		s.logger.Errorf("handleGetDevices: failed to get devices: %s", err)
		s.respondJSON(w, http.StatusInternalServerError, errorResponse{Error: "failed to get devices"})
		return
	}

	var deviceNames []string
	for _, device := range devices {
		deviceNames = append(deviceNames, device.Name)
	}

	s.respondJSON(w, http.StatusOK, deviceNames)
}

func (s *Server) handleGetMetrics(w http.ResponseWriter, r *http.Request) {
	deviceID := r.Header.Get("X-DeviceID")
	if deviceID == "" {
		s.logger.Errorf("handleGetMetrics: missing device ID")
		s.respondJSON(w, http.StatusBadRequest, errorResponse{Error: "missing device ID"})
		return
	}

	var device db.Device
	if err := s.db.First(&device, "name = ?", deviceID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			s.logger.Errorf("handleGetMetrics: device not found")
			s.respondJSON(w, http.StatusNotFound, errorResponse{Error: "device not found"})
		} else {
			s.logger.Errorf("handleGetMetrics: failed to query device: %s", err)
			s.respondJSON(w, http.StatusInternalServerError, errorResponse{Error: "failed to query device"})
		}
		return
	}

	var metrics []db.Metric
	limit := 10                // Default limit per page
	page := 1                  // Default page number
	sort := "recorded_at desc" // Default sort order

	// Parse query parameters
	if p := r.URL.Query().Get("page"); p != "" {
		if parsedPage, err := strconv.Atoi(p); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}
	if s := r.URL.Query().Get("sort"); s != "" {
		sort = s
	}

	offset := (page - 1) * limit

	if err := s.db.Preload("Type").Preload("Unit").Where("device_id = ?", device.ID).Order(sort).Limit(limit).Offset(offset).Find(&metrics).Error; err != nil {
		s.logger.Errorf("handleGetMetrics: failed to get metrics: %s", err)
		s.respondJSON(w, http.StatusInternalServerError, errorResponse{Error: "failed to get metrics"})
		return
	}

	// Get total count for pagination
	var totalRecords int64
	if err := s.db.Model(&db.Metric{}).Where("device_id = ?", device.ID).Count(&totalRecords).Error; err != nil {
		s.logger.Errorf("handleGetMetrics: failed to get total count: %s", err)
		s.respondJSON(w, http.StatusInternalServerError, errorResponse{Error: "failed to get total count"})
		return
	}

	response := map[string]interface{}{
		"metrics":      metrics,
		"totalRecords": totalRecords,
		"currentPage":  page,
		"totalPages":   (totalRecords + int64(limit) - 1) / int64(limit),
	}

	s.respondJSON(w, http.StatusOK, response)
}
