package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/bxrne/beacon/web/pkg/db"
)

func (s *Server) handleDashboardView(w http.ResponseWriter, r *http.Request) {
	page := 1
	if p := r.URL.Query().Get("page"); p != "" {
		if parsedPage, err := strconv.Atoi(p); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	tmpl, err := template.ParseFiles(
		"templates/base.html",
		"templates/partials/device_selection.html",
		"templates/dashboard.html",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.ExecuteTemplate(w, "base", map[string]interface{}{
		"Page": page,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) handleChartsView(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"templates/base.html",
		"templates/partials/device_selection.html",
		"templates/charts.html",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.ExecuteTemplate(w, "base", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Add a new handler for Server-Sent Events
func (s *Server) handleMetricsStream(w http.ResponseWriter, r *http.Request) {
	deviceID := r.URL.Query().Get("deviceID")
	if deviceID == "" {
		http.Error(w, "Missing device ID", http.StatusBadRequest)
		return
	}

	// Set headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	// Stream metrics periodically
	for {
		select {
		case <-r.Context().Done():
			return
		default:
			// Fetch latest metrics for the device
			metricsData, err := s.getLatestMetrics(deviceID)
			if err != nil {
				http.Error(w, "Error fetching metrics", http.StatusInternalServerError)
				return
			}

			// Send metrics as JSON
			data, _ := json.Marshal(metricsData)
			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()

			// Wait before sending the next update
			time.Sleep(5 * time.Second)
		}
	}
}

// Helper function to get latest metrics
func (s *Server) getLatestMetrics(deviceID string) ([]db.Metric, error) {
	var device db.Device
	if err := s.db.First(&device, "name = ?", deviceID).Error; err != nil {
		return nil, err
	}

	var metrics []db.Metric
	if err := s.db.Preload("Type").Preload("Unit").
		Where("device_id = ?", device.ID).
		Order("recorded_at DESC").
		Find(&metrics).Error; err != nil {
		return nil, err
	}

	return metrics, nil
}
