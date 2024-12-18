package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/bxrne/beacon/web/internal/db"
	"github.com/bxrne/beacon/web/internal/metrics"
	"gorm.io/gorm"
)

type errorResponse struct {
	Error string `json:"error"`
}

type commandRequest struct {
	Device  string `json:"device"`
	Command string `json:"command"`
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
		s.logger.Errorf("Missing device ID")
		http.Error(w, "Missing device ID", http.StatusBadRequest)
		return
	}

	var deviceMetrics metrics.DeviceMetrics
	if err := json.NewDecoder(r.Body).Decode(&deviceMetrics); err != nil {
		s.logger.Errorf("Failed to decode metrics: %v", err)
		http.Error(w, "Invalid metrics format", http.StatusBadRequest)
		return
	}

	if err := s.persistMetrics(deviceID, deviceMetrics); err != nil {
		s.logger.Errorf("Failed to persist metrics: %v", err)
		http.Error(w, "Failed to persist metrics", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) persistMetrics(deviceID string, deviceMetrics metrics.DeviceMetrics) error {
	for _, metric := range deviceMetrics.Metrics {
		var metricType db.MetricType
		if err := s.db.FirstOrCreate(&metricType, db.MetricType{Name: metric.Type}).Error; err != nil {
			return err
		}

		var unit db.Unit
		if err := s.db.FirstOrCreate(&unit, db.Unit{Name: metric.Unit}).Error; err != nil {
			return err
		}

		var device db.Device
		if err := s.db.FirstOrCreate(&device, db.Device{Name: deviceID}).Error; err != nil {
			return err
		}

		// Clean up the RecordedAt string by removing any control characters
		cleanTime := strings.Map(func(r rune) rune {
			if r >= 32 && r != 127 { // Keep only printable characters
				return r
			}
			return -1
		}, metric.RecordedAt)

		recordedAt, err := time.Parse(time.RFC3339, cleanTime)
		if err != nil {
			return fmt.Errorf("invalid RecordedAt format: %v (raw: %q)", err, metric.RecordedAt)
		}

		dbMetric := db.Metric{
			TypeID:     metricType.ID,
			Value:      metric.Value,
			UnitID:     unit.ID,
			DeviceID:   device.ID,
			RecordedAt: recordedAt,
		}
		if err := s.db.Create(&dbMetric).Error; err != nil {
			return err
		}
	}
	return nil
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

// handleGetMetrics godoc
// @Summary      Get metrics with pagination and filtering
// @Description  Get metrics for a device with pagination and filtering options
// @Tags         metrics
// @Produce      json
// @Param        X-DeviceID  header    string  true  "Device ID"
// @Param        page        query     int     false "Page number"
// @Param        sort        query     string  false "Sort order"
// @Param        type        query     string  false "Metric type"
// @Param        view        query     string  false "View type (charts)"
// @Success      200         {object}  map[string]interface{}
// @Failure      400         {object}  errorResponse
// @Failure      404         {object}  errorResponse
// @Failure      500         {object}  errorResponse
// @Router       /metrics [get]
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
	metricType := r.URL.Query().Get("type")

	offset := (page - 1) * limit

	query := s.db.Preload("Type").Preload("Unit").Where("device_id = ?", device.ID).Order(sort).Limit(limit).Offset(offset)
	var metricTypeIDs []uint
	if metricType != "" {
		if err := s.db.Model(&db.MetricType{}).Select("id").Where("name LIKE ?", "%"+metricType+"%").Scan(&metricTypeIDs).Error; err != nil {
			s.logger.Errorf("handleGetMetrics: failed to get metric type IDs: %s", err)
			s.respondJSON(w, http.StatusInternalServerError, errorResponse{Error: "failed to get metric type IDs"})
			return
		}
		query = query.Where("type_id IN ?", metricTypeIDs)
	}
	if err := query.Find(&metrics).Error; err != nil {
		s.logger.Errorf("handleGetMetrics: failed to get metrics: %s", err)
		s.respondJSON(w, http.StatusInternalServerError, errorResponse{Error: "failed to get metrics"})
		return
	}

	// Check if the request is for charts view
	isChartsView := r.URL.Query().Get("view") == "charts"

	var responseMetrics interface{}
	if isChartsView {
		// For charts view, we want both percent and color metrics
		latestMetrics := make(map[string]db.Metric)
		for _, metric := range metrics {
			if metric.Unit.Name == "percent" || metric.Unit.Name == "color" {
				existing, exists := latestMetrics[metric.Type.Name]
				if !exists || metric.RecordedAt.After(existing.RecordedAt) {
					latestMetrics[metric.Type.Name] = metric
				}
			}
		}

		// Convert map to slice
		var latestMetricsSlice []db.Metric
		for _, metric := range latestMetrics {
			latestMetricsSlice = append(latestMetricsSlice, metric)
		}

		// Ensure we always return an empty array instead of null
		if latestMetricsSlice == nil {
			latestMetricsSlice = []db.Metric{}
		}
		responseMetrics = latestMetricsSlice
	} else {
		responseMetrics = metrics
	}

	// Get total count for pagination
	var totalRecords int64
	countQuery := s.db.Model(&db.Metric{}).Where("device_id = ?", device.ID)
	if metricType != "" {
		countQuery = countQuery.Where("type_id IN ?", metricTypeIDs)
	}
	if err := countQuery.Count(&totalRecords).Error; err != nil {
		s.logger.Errorf("handleGetMetrics: failed to get total count: %s", err)
		s.respondJSON(w, http.StatusInternalServerError, errorResponse{Error: "failed to get total count"})
		return
	}

	response := map[string]interface{}{
		"metrics":      responseMetrics,
		"totalRecords": totalRecords,
		"currentPage":  page,
		"totalPages":   (totalRecords + int64(limit) - 1) / int64(limit),
	}

	s.respondJSON(w, http.StatusOK, response)
}

// handleCommand godoc
// @Summary      Submit command
// @Description  Submit a command for a device
// @Tags         command
// @Accept       json
// @Produce      json
// @Param        commandRequest  body      commandRequest  true  "Command request"
// @Success      200             {object}  map[string]string
// @Failure      400             {object}  errorResponse
// @Failure      404             {object}  errorResponse
// @Failure      500             {object}  errorResponse
// @Router       /command [post]
func (s *Server) handleCommand(w http.ResponseWriter, r *http.Request) {
	var req commandRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid request"})
		return
	}

	// Check if device exists
	var device db.Device
	if err := s.db.First(&device, "name = ?", req.Device).Error; err != nil {
		s.respondJSON(w, http.StatusNotFound, errorResponse{Error: "device not found"})
		return
	}

	// Validate command against command types
	var commandType db.CommandType
	if err := s.db.First(&commandType, "name = ?", req.Command).Error; err != nil {
		s.respondJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid command type"})
		return
	}

	// Create command with status
	command := db.Command{
		Name:     req.Command,
		DeviceID: device.ID,
		Status:   "pending",
	}

	if err := s.db.Create(&command).Error; err != nil {
		s.logger.Errorf("handleCommand: failed to create command: %s", err)
		s.respondJSON(w, http.StatusInternalServerError, errorResponse{Error: "failed to create command"})
		return
	}

	s.respondJSON(w, http.StatusOK, map[string]string{"status": "success", "message": "Command queued successfully"})
}

// handleGetCommands godoc
// @Summary      Get pending commands
// @Description  Get pending commands for a device
// @Tags         command
// @Produce      json
// @Param        X-DeviceID  header    string  true  "Device ID"
// @Success      200         {object}  []metrics.CommandResponse
// @Failure      400         {object}  errorResponse
// @Failure      500         {object}  errorResponse
// @Router       /commands [get]
func (s *Server) handleGetCommands(w http.ResponseWriter, r *http.Request) {
	deviceID := r.Header.Get("X-DeviceID")
	if deviceID == "" {
		s.logger.Error("missing device ID")
		s.respondJSON(w, http.StatusBadRequest, errorResponse{Error: "missing device ID"})
		return
	}

	// Get pending commands for the device
	var commands []db.Command
	if err := s.db.Preload("Device").Where("device_id IN (SELECT id FROM devices WHERE name = ?) AND status = ?",
		deviceID, "pending").Find(&commands).Error; err != nil {
		s.logger.Error("failed to get commands", "error", err)
		s.respondJSON(w, http.StatusInternalServerError, errorResponse{Error: "failed to get commands"})
		return
	}

	// Convert to response format
	var response []metrics.CommandResponse
	for _, cmd := range commands {
		response = append(response, metrics.CommandResponse{
			Device:  deviceID,
			Command: cmd.Name,
		})
	}

	// If no commands found, return empty array instead of null
	if response == nil {
		response = []metrics.CommandResponse{}
	}

	s.respondJSON(w, http.StatusOK, response)
}

// handleCommandStatus godoc
// @Summary      Update command status
// @Description  Update the status of a command for a device
// @Tags         command
// @Accept       json
// @Produce      json
// @Param        commandStatusRequest  body      metrics.CommandStatusRequest  true  "Command status request"
// @Success      200                   {object}  map[string]string
// @Failure      400                   {object}  errorResponse
// @Failure      500                   {object}  errorResponse
// @Router       /command/status [post]
func (s *Server) handleCommandStatus(w http.ResponseWriter, r *http.Request) {
	var req metrics.CommandStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid request"})
		return
	}

	// Update command status
	result := s.db.Model(&db.Command{}).
		Where("device_id IN (SELECT id FROM devices WHERE name = ?) AND name = ? AND status = ?",
			req.Device, req.Command, "pending").
		Update("status", req.Status)

	if result.Error != nil {
		s.logger.Error("failed to update command status", "error", result.Error)
		s.respondJSON(w, http.StatusInternalServerError, errorResponse{Error: "failed to update command status"})
		return
	}

	s.respondJSON(w, http.StatusOK, map[string]string{"status": "success"})
}
