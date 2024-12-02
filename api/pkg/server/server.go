package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	_ "github.com/bxrne/beacon-web/docs" // This line is necessary for go-swagger to find your docs
	"github.com/bxrne/beacon-web/pkg/config"
	"github.com/bxrne/beacon-web/pkg/metrics"
	"github.com/charmbracelet/log"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Server struct {
	router       *mux.Router
	srv          *http.Server
	logger       *log.Logger
	cfg          *config.Config
	db           *sql.DB
	metricsCache *metrics.MetricsCache
}

func New(cfg *config.Config, logger *log.Logger, db *sql.DB) *Server {
	s := &Server{
		router:       mux.NewRouter(),
		logger:       logger,
		cfg:          cfg,
		db:           db,
		metricsCache: metrics.NewMetricsCache(cfg),
	}

	s.setupRoutes()
	return s
}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lrw := newLoggingResponseWriter(w)

		next.ServeHTTP(lrw, r)

		s.logger.Info("request completed",
			"method", r.Method,
			"path", r.URL.Path,
			"status", lrw.statusCode,
			"duration_ns", time.Since(start).Nanoseconds(),
			"remote_addr", r.RemoteAddr,
			"hostname", r.Header.Get("X-Hostname"),
		)
	})
}

func (s *Server) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		s.logger.Error("failed to encode response", "error", err)
	}
}

func (s *Server) respondError(w http.ResponseWriter, code int, message string) {
	s.respondJSON(w, code, map[string]string{"error": message})
}

func (s *Server) Start(ctx context.Context) error {
	cors := handlers.CORS(
		handlers.AllowedOrigins(s.cfg.Server.AllowedOrigins),
		handlers.AllowedMethods([]string{"GET", "POST"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	s.srv = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.cfg.Server.Port),
		Handler: cors(s.router),
	}

	go func() {
		s.logger.Info("starting server", "port", s.cfg.Server.Port)
		if err := s.srv.ListenAndServe(); err != http.ErrServerClosed {
			s.logger.Error("server error", "error", err)
		}
	}()

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("stopping server")

	shutdownCtx, cancel := context.WithTimeout(ctx, time.Duration(s.cfg.Server.ShutdownTimeout)*time.Second)
	defer cancel()

	if err := s.srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	return nil
}

func (s *Server) setupRoutes() {
	s.router.Use(s.loggingMiddleware)
	s.router.HandleFunc("/health", s.handleHealth).Methods(http.MethodGet)
	s.router.HandleFunc("/metric", s.handleMetric).Methods(http.MethodPost)
	s.router.HandleFunc("/metric", s.handleGetMetric).Methods(http.MethodGet)
	s.router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
}
