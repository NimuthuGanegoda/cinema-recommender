package server

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"cinema-recommender/internal/models"
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

// NewServer initializes a new Server instance.
func NewServer(cfg Config, sc *scraper.RegionalScraper, recEngine *recommender.Engine) *Server {
	if cfg.Port == "" {
		cfg.Port = ":8080"
	}
	if !strings.HasPrefix(cfg.Port, ":") {
		cfg.Port = ":" + cfg.Port
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

// registerRoutes sets up all RESTful endpoints and web UI routing.
func (s *Server) registerRoutes(mux *http.ServeMux) {
	// Web UI & Mobile App
	uiSub, err := fs.Sub(uiAssets, "ui")
	if err == nil {
		mux.Handle("GET /ui/", http.StripPrefix("/ui/", http.FileServer(http.FS(uiSub))))
	}
	mux.HandleFunc("GET /", s.handleRoot)
	mux.HandleFunc("GET /mobile", s.handleMobile)
	mux.HandleFunc("GET /manifest.json", s.handleManifest)
	mux.HandleFunc("GET /sw.js", s.handleServiceWorker)
	mux.HandleFunc("GET /favicon.ico", s.handleFavicon)

	// Health Check
	mux.HandleFunc("GET /health", s.handleHealth)
	mux.HandleFunc("GET /api/v1/health", s.handleHealth)

	// Regional Cinemas
	mux.HandleFunc("GET /api/v1/cinemas", s.handleListCinemas)
	mux.HandleFunc("GET /api/v1/cinemas/{id}", s.handleGetCinema)
	mux.HandleFunc("GET /api/v1/cinemas/{id}/concessions", s.handleGetConcessions)
	mux.HandleFunc("GET /api/v1/cinemas/{id}/concessions/live", s.handleGetLiveConcessions)
	mux.HandleFunc("GET /api/v1/cinemas/{id}/ads", s.handleGetCinemaAds)

	// Scope Privilege Concession Promotions & Authentic Sri Lankan Payment Gateways
	mux.HandleFunc("GET /api/v1/promotions", s.handleListPromotions)
	mux.HandleFunc("GET /api/v1/promotions/privilege", s.handleListPromotions)
	mux.HandleFunc("GET /api/v1/promotions/scope", s.handleListPromotions)
	mux.HandleFunc("GET /api/v1/promotions/redopay", s.handleListPromotions)
	mux.HandleFunc("GET /api/v1/payment-methods", s.handleListPaymentMethods)

	// Optimization, Evaluation & Checkout
	mux.HandleFunc("POST /api/v1/recommend", s.handleRecommend)
	mux.HandleFunc("POST /api/v1/cart/evaluate", s.handleEvaluateCart)
	mux.HandleFunc("POST /api/v1/checkout", s.handleCheckout)
}

// Handler: Root serves index.html (Desktop Web UI)
func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	content, err := uiAssets.ReadFile("ui/index.html")
	if err != nil {
		http.Error(w, "Web UI asset not found", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(content)
}

// Handler: Mobile serves mobile.html (Dedicated Mobile Application)
func (s *Server) handleMobile(w http.ResponseWriter, r *http.Request) {
	content, err := uiAssets.ReadFile("ui/mobile.html")
	if err != nil {
		http.Error(w, "Mobile UI asset not found", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(content)
}

// Handler: Manifest serves manifest.json for PWA installation
func (s *Server) handleManifest(w http.ResponseWriter, r *http.Request) {
	content, err := uiAssets.ReadFile("ui/manifest.json")
	if err != nil {
		http.Error(w, "Manifest not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/manifest+json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(content)
}

// Handler: Service Worker serves sw.js
func (s *Server) handleServiceWorker(w http.ResponseWriter, r *http.Request) {
	content, err := uiAssets.ReadFile("ui/sw.js")
	if err != nil {
		http.Error(w, "Service Worker not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/javascript")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(content)
}

// Handler: Favicon
func (s *Server) handleFavicon(w http.ResponseWriter, r *http.Request) {
	content, err := uiAssets.ReadFile("ui/app-icon.jpg")
	if err != nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(content)
}

// Handler: Health Check
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	cinemas := s.scraper.GetRegisteredCinemas()
	s.writeJSON(w, http.StatusOK, map[string]interface{}{
		"status":            "healthy",
		"service":           "cinema-food-beverage-recommender",
		"version":           "2.0.0",
		"regional_theaters": len(cinemas),
		"scope":             "strictly_outside_colombo",
		"promotion_engine":  "scope_privilege_optimized",
		"timestamp":         time.Now().UTC().Format(time.RFC3339),
	})
}

// Handler: List Regional Cinemas (supports ?city=, ?lat=, and ?lng= for location-based auto-fetch)
func (s *Server) handleListCinemas(w http.ResponseWriter, r *http.Request) {
	cityFilter := r.URL.Query().Get("city")
	latStr := r.URL.Query().Get("lat")
	lngStr := r.URL.Query().Get("lng")

	var lat, lng float64
	if latStr != "" && lngStr != "" {
		var errLat, errLng error
		lat, errLat = strconv.ParseFloat(latStr, 64)
		lng, errLng = strconv.ParseFloat(lngStr, 64)
		if errLat == nil && errLng == nil {
			cinemas := s.scraper.GetCinemasWithLocationSort(lat, lng, cityFilter)
			s.writeJSON(w, http.StatusOK, cinemas)
			return
		}
	}

	if cityFilter != "" {
		cinemas, err := s.scraper.GetCinemasByCity(cityFilter)
		if err != nil {
			s.handleError(w, err)
			return
		}
		s.writeJSON(w, http.StatusOK, cinemas)
		return
	}

	cinemas := s.scraper.GetRegisteredCinemas()
	s.writeJSON(w, http.StatusOK, cinemas)
}

// Handler: Get Cinema by ID
func (s *Server) handleGetCinema(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		s.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "cinema id is required"})
		return
	}

	cinema, err := s.scraper.GetCinemaByID(id)
	if err != nil {
		s.handleError(w, err)
		return
	}

	s.writeJSON(w, http.StatusOK, cinema)
}

// Handler: Get Concessions for a Cinema (with constantly updating live pricing)
func (s *Server) handleGetConcessions(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		s.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "cinema id is required"})
		return
	}

	cinema, err := s.scraper.GetCinemaByID(id)
	if err != nil {
		s.handleError(w, err)
		return
	}

	items, err := s.scraper.ScrapeCinemaConcessions(cinema)
	if err != nil {
		s.handleError(w, err)
		return
	}

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("X-Live-Price-Feed", "active")
	w.Header().Set("X-Prices-Updated-At", time.Now().Format("15:04:05"))
	s.writeJSON(w, http.StatusOK, items)
}

// Handler: Get Live Concessions with Ticker Status & Countdown
func (s *Server) handleGetLiveConcessions(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		s.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "cinema id is required"})
		return
	}

	cinema, err := s.scraper.GetCinemaByID(id)
	if err != nil {
		s.handleError(w, err)
		return
	}

	items, err := s.scraper.ScrapeCinemaConcessions(cinema)
	if err != nil {
		s.handleError(w, err)
		return
	}

	ticker := s.scraper.GetLiveConcessionTicker(cinema)
	ticker["items"] = items
	ticker["cinema_ads"] = s.scraper.GetCinemaAds(cinema, items)

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("X-Live-Price-Feed", "active")
	s.writeJSON(w, http.StatusOK, ticker)
}

// Handler: Get Cinema Promotional Ads for Food, Drinks & Combos
func (s *Server) handleGetCinemaAds(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		s.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "cinema id is required"})
		return
	}

	cinema, err := s.scraper.GetCinemaByID(id)
	if err != nil {
		s.handleError(w, err)
		return
	}

	items, err := s.scraper.ScrapeCinemaConcessions(cinema)
	if err != nil {
		s.handleError(w, err)
		return
	}

	ads := s.scraper.GetCinemaAds(cinema, items)
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	s.writeJSON(w, http.StatusOK, ads)
}

// Handler: List Available Scope Privilege Promotions
func (s *Server) handleListPromotions(w http.ResponseWriter, r *http.Request) {
	promos := s.recEngine.GetAvailablePromotions()
	s.writeJSON(w, http.StatusOK, promos)
}

// Handler: Recommend Optimal Bundle
func (s *Server) handleRecommend(w http.ResponseWriter, r *http.Request) {
	var req models.RecommendationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON payload: " + err.Error()})
		return
	}

	// Resolve cinema
	var cinema models.Cinema
	var err error

	if req.CinemaID != "" {
		cinema, err = s.scraper.GetCinemaByID(req.CinemaID)
		if err != nil {
			s.handleError(w, err)
			return
		}
	} else if req.City != "" {
		if err := s.scraper.ValidateLocation(req.City); err != nil {
			s.handleError(w, err)
			return
		}
		cinemas, err := s.scraper.GetCinemasByCity(req.City)
		if err != nil || len(cinemas) == 0 {
			cinema = models.Cinema{
				ID:        fmt.Sprintf("GEN-%s", strings.ToUpper(req.City)),
				Name:      fmt.Sprintf("Regional Cinema %s", req.City),
				City:      req.City,
				IsOutside: true,
			}
		} else {
			cinema = cinemas[0]
		}
	} else {
		// Default to first regional cinema (KCC Multiplex)
		cinemas := s.scraper.GetRegisteredCinemas()
		if len(cinemas) > 0 {
			cinema = cinemas[0]
		}
	}

	// Scrape available concessions
	items, err := s.scraper.ScrapeCinemaConcessions(cinema)
	if err != nil {
		s.handleError(w, err)
		return
	}

	// Compute optimal bundle
	result, err := s.recEngine.RecommendOptimalBundle(cinema, items, req)
	if err != nil {
		s.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	s.writeJSON(w, http.StatusOK, result)
}

// Handler: Evaluate Custom User Cart
func (s *Server) handleEvaluateCart(w http.ResponseWriter, r *http.Request) {
	var req models.CartEvaluationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "malformed JSON payload: " + err.Error()})
		return
	}

	if req.CinemaID == "" {
		req.CinemaID = "KND-KCC"
	}

	cinema, err := s.scraper.GetCinemaByID(req.CinemaID)
	if err != nil {
		s.handleError(w, err)
		return
	}

	availableItems, err := s.scraper.ScrapeCinemaConcessions(cinema)
	if err != nil {
		s.handleError(w, err)
		return
	}

	// Match cart items with quantities
	var selected []models.ConcessionItem
	itemMap := make(map[string]models.ConcessionItem)
	for _, it := range availableItems {
		itemMap[it.ID] = it
	}

	for _, cartItem := range req.Items {
		if original, found := itemMap[cartItem.ItemID]; found {
			qty := cartItem.Quantity
			if qty <= 0 {
				qty = 1
			}
			for q := 0; q < qty; q++ {
				selected = append(selected, original)
			}
		}
	}

	tier := req.PrivilegeTier
	if tier == "" {
		tier = req.RedopayTier
	}
	result := s.recEngine.EvaluateCart(cinema, selected, req.PromoCode, tier)
	s.writeJSON(w, http.StatusOK, result)
}

// Handler: List Supported Sri Lankan Payment Methods
func (s *Server) handleListPaymentMethods(w http.ResponseWriter, r *http.Request) {
	methods := s.paymentSvc.GetSupportedPaymentMethods()
	s.writeJSON(w, http.StatusOK, methods)
}

// Handler: Checkout via Sri Lankan Payment Channels
func (s *Server) handleCheckout(w http.ResponseWriter, r *http.Request) {
	var req models.CheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "malformed JSON payload: " + err.Error()})
		return
	}

	if req.CinemaID == "" {
		req.CinemaID = "KND-KCC"
	}

	cinema, err := s.scraper.GetCinemaByID(req.CinemaID)
	if err != nil {
		s.handleError(w, err)
		return
	}

	availableItems, err := s.scraper.ScrapeCinemaConcessions(cinema)
	if err != nil {
		s.handleError(w, err)
		return
	}

	var selected []models.ConcessionItem
	itemMap := make(map[string]models.ConcessionItem)
	for _, it := range availableItems {
		itemMap[it.ID] = it
	}

	for _, cartItem := range req.Items {
		if original, found := itemMap[cartItem.ItemID]; found {
			qty := cartItem.Quantity
			if qty <= 0 {
				qty = 1
			}
			for q := 0; q < qty; q++ {
				selected = append(selected, original)
			}
		}
	}

	res, err := s.paymentSvc.ProcessCheckout(cinema, selected, req)
	if err != nil {
		s.writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	s.writeJSON(w, http.StatusOK, res)
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

// writeJSON serializes data as JSON with appropriate HTTP status.
func (s *Server) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

// Middlewares
func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := &statusResponseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(ww, r)
		log.Printf("[%s] %s %s - %d (%v)", r.Method, r.URL.Path, r.RemoteAddr, ww.status, time.Since(start))
	})
}

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("PANIC recovered in %s: %v", r.URL.Path, rec)
				s.writeJSON(w, http.StatusInternalServerError, map[string]string{
					"error": "internal server panic recovered",
				})
			}
		}()
		next.ServeHTTP(w, r)
	})
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

type statusResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}
