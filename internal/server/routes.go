package server

import (
	"io/fs"
	"net/http"
)

// registerRoutes sets up all RESTful endpoints and web UI routing.
func (s *Server) registerRoutes(mux *http.ServeMux) {
	// Web UI & Mobile App
	uiSub, err := fs.Sub(uiAssets, "ui")
	if err == nil {
		fileServer := http.FileServer(http.FS(uiSub))
		mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/" {
				s.handleRoot(w, r)
				return
			}
			fileServer.ServeHTTP(w, r)
		})
	}
	mux.HandleFunc("GET /mobile", s.handleMobile)
	mux.HandleFunc("GET /mobile.html", s.handleMobile)
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
	mux.HandleFunc("GET /api/v1/cinemas/{id}/movies", s.handleGetCinemaMovies)

	// Live Movies & 1-Hour Automated Sync Engine
	mux.HandleFunc("GET /api/v1/movies", s.handleListMovies)
	mux.HandleFunc("GET /api/v1/movies/status", s.handleMovieSyncStatus)
	mux.HandleFunc("POST /api/v1/movies/sync", s.handleTriggerMovieSync)

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
