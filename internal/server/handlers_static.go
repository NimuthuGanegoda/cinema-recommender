package server

import (
	"bytes"
	"net/http"
	"time"
)

// handleRoot serves the shared responsive cinema and food website.
func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	s.serveHTMLPage(w, "index.html")
}

func (s *Server) serveHTMLPage(w http.ResponseWriter, pageFile string) {
	content, err := uiAssets.ReadFile("ui/" + pageFile)
	if err != nil {
		http.Error(w, pageFile+" asset not found", http.StatusInternalServerError)
		return
	}

	content = bytes.Replace(content, []byte("</head>"), []byte(`<meta name="cinebite-runtime" content="server"></head>`), 1)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(content)
}

func (s *Server) handleMoviesPage(w http.ResponseWriter, r *http.Request) {
	s.serveHTMLPage(w, "movies.html")
}

func (s *Server) handleCinemasPage(w http.ResponseWriter, r *http.Request) {
	s.serveHTMLPage(w, "cinemas.html")
}

func (s *Server) handleMenuPage(w http.ResponseWriter, r *http.Request) {
	s.serveHTMLPage(w, "menu.html")
}

func (s *Server) handleOptimizerPage(w http.ResponseWriter, r *http.Request) {
	s.serveHTMLPage(w, "optimizer.html")
}

func (s *Server) handleDealsPage(w http.ResponseWriter, r *http.Request) {
	s.serveHTMLPage(w, "deals.html")
}

func (s *Server) handleFoodPlanPage(w http.ResponseWriter, r *http.Request) {
	s.serveHTMLPage(w, "food-plan.html")
}

// handleMobile preserves old bookmarks while using the shared responsive page.
func (s *Server) handleMobile(w http.ResponseWriter, r *http.Request) {
	s.serveHTMLPage(w, "index.html")
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
		"scope":             "sri_lanka",
		"promotion_engine":  "scope_privilege_optimized",
		"timestamp":         time.Now().UTC().Format(time.RFC3339),
	})
}
