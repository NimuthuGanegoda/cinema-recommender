package server

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"cinema-recommender/pkg/payment"
	"cinema-recommender/pkg/recommender"
	"cinema-recommender/pkg/scraper"
)

//go:embed ui/*
var uiAssets embed.FS

// Config contains server runtime configurations.
type Config struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// Server encapsulates the HTTP multiplexer, middlewares, and services.
type Server struct {
	cfg        Config
	scraper    *scraper.RegionalScraper
	recEngine  *recommender.Engine
	paymentSvc *payment.Service
	httpServer *http.Server
}

// NewServer initializes a new Server instance with configured timeouts, routing, and middlewares.
func NewServer(cfg Config, sc *scraper.RegionalScraper, recEngine *recommender.Engine) *Server {
	if cfg.Port == "" {
		cfg.Port = ":8080"
	}
	if cfg.ReadTimeout == 0 {
		cfg.ReadTimeout = 10 * time.Second
	}
	if cfg.WriteTimeout == 0 {
		cfg.WriteTimeout = 10 * time.Second
	}
	if cfg.IdleTimeout == 0 {
		cfg.IdleTimeout = 60 * time.Second
	}

	s := &Server{
		cfg:        cfg,
		scraper:    sc,
		recEngine:  recEngine,
		paymentSvc: payment.NewService(recEngine),
	}

	mux := http.NewServeMux()
	s.registerRoutes(mux)

	handler := s.recoveryMiddleware(s.corsMiddleware(s.loggingMiddleware(mux)))

	s.httpServer = &http.Server{
		Addr:         cfg.Port,
		Handler:      handler,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	return s
}

// Start begins listening on the configured TCP address.
func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the HTTP server.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

// GetHandler returns the configured root HTTP handler for testing.
func (s *Server) GetHandler() http.Handler {
	return s.httpServer.Handler
}

// writeJSON serializes data as JSON with appropriate HTTP status.
func (s *Server) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

// handleError formats error responses based on domain errors.
func (s *Server) handleError(w http.ResponseWriter, err error) {
	if errors.Is(err, scraper.ErrColomboExcluded) {
		s.writeJSON(w, http.StatusBadRequest, map[string]string{
			"error":  err.Error(),
			"policy": "Engine is strictly restricted to regional cinema circuits outside Colombo",
		})
		return
	}

	if errors.Is(err, scraper.ErrCinemaNotFound) {
		s.writeJSON(w, http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
		return
	}

	if errors.Is(err, scraper.ErrInvalidInput) {
		s.writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	s.writeJSON(w, http.StatusInternalServerError, map[string]string{
		"error": err.Error(),
	})
}
