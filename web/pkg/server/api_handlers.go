package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/bxrne/beacon/web/pkg/db"
	"github.com/bxrne/beacon/web/pkg/metrics"
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
		// Filter for percent-based metrics and keep only the latest for each type
		latestMetrics := make(map[string]db.Metric)
		for _, metric := range metrics {
			if metric.Unit.Name == "percent" {
				latestMetrics[metric.Type.Name] = metric
			}
		}

		// Convert map to slice
		var latestMetricsSlice []db.Metric
		for _, metric := range latestMetrics {
			latestMetricsSlice = append(latestMetricsSlice, metric)
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
