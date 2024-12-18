package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"gorm.io/gorm"

	_ "github.com/bxrne/beacon/web/docs" // This line is necessary for go-swagger to find your docs
	"github.com/bxrne/beacon/web/internal/config"
	"github.com/bxrne/beacon/web/pkg/metrics"
	"github.com/charmbracelet/log"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Server struct {
	router       *mux.Router
	srv          *http.Server
	logger       *log.Logger
	cfg          *config.Config
	db           *gorm.DB
	metricsCache *metrics.MetricsCache
}

func New(cfg *config.Config, logger *log.Logger, db *gorm.DB) *Server {
	s := &Server{
		router:       mux.NewRouter(),
		logger:       logger,
		cfg:          cfg,
		db:           db.Session(&gorm.Session{}),
		metricsCache: metrics.NewMetricsCache(cfg),
	}

	s.setupRoutes()
	return s
}

func (s *Server) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		s.logger.Error("failed to encode response", "error", err)
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.srv = &http.Server{
		Addr:         fmt.Sprintf(":%d", s.cfg.Server.Port),
		Handler:      s.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatalf("listen: %s\n", err)
		}
	}()

	s.logger.Infof("Server started on port %d", s.cfg.Server.Port)
	<-ctx.Done()

	s.logger.Info("Server stopping")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), time.Duration(s.cfg.Server.ShutdownTimeout)*time.Second)
	defer cancel()

	if err := s.srv.Shutdown(ctxShutDown); err != nil {
		s.logger.Fatalf("server shutdown failed: %s", err)
	}

	s.logger.Info("Server exited properly")
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.status = code
	rec.ResponseWriter.WriteHeader(code)
}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

		// Log incoming request details
		s.logger.Debugf("REQUEST Method=%s Path=%s Source=%s ",
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
		)

		next.ServeHTTP(rec, r)

		duration := time.Since(start).Nanoseconds()

		// Log response details
		s.logger.Debugf("RESPONSE Method=%s Path=%s Status=%d DurationNS=%d Source=%s",
			r.Method,
			r.URL.Path,
			rec.status,
			duration,
			r.RemoteAddr,
		)
	})
}

func (s *Server) setupRoutes() {
	s.router.Use(s.loggingMiddleware)

	s.router.HandleFunc("/", s.handleDashboardView).Methods(http.MethodGet)
	s.router.HandleFunc("/charts", s.handleChartsView).Methods(http.MethodGet)
	s.router.HandleFunc("/command", s.handleCommandPage).Methods(http.MethodGet)

	s.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	s.router.PathPrefix("/docs/").Handler(httpSwagger.WrapHandler)

	// WARN: Silence favicon warnings
	s.router.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	apiRouter := s.router.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/health", s.handleHealth).Methods(http.MethodGet)
	apiRouter.HandleFunc("/metric", s.handleMetric).Methods(http.MethodPost)
	apiRouter.HandleFunc("/metric", s.handleGetMetric).Methods(http.MethodGet)
	apiRouter.HandleFunc("/device", s.handleGetDevices).Methods(http.MethodGet)
	apiRouter.HandleFunc("/metrics", s.handleGetMetrics).Methods(http.MethodGet)
	apiRouter.HandleFunc("/command", s.handleCommand).Methods(http.MethodPost)
	apiRouter.HandleFunc("/command", s.handleGetCommands).Methods(http.MethodGet)
}
