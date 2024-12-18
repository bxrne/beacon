package server

import (
	"html/template"
	"net/http"
	"strconv"
)

// handleDashboardView handles the dashboard view page
// @Summary Show dashboard view
// @Description Get dashboard view page
// @Tags dashboard
// @Accept  json
// @Produce  html
// @Param page query int false "Page number"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {string} string "Internal Server Error"
// @Router /dashboard [get]
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

// handleChartsView handles the charts view page
// @Summary Show charts view
// @Description Get charts view page
// @Tags charts
// @Accept  json
// @Produce  html
// @Success 200 {string} string "OK"
// @Failure 500 {string} string "Internal Server Error"
// @Router /charts [get]
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

func (s *Server) handleCommandPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"templates/base.html",
		"templates/partials/device_selection.html",
		"templates/command.html",
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
