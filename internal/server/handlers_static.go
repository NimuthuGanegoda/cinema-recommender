package server

import (
	"net/http"
	"strings"
	"time"
)

// isMobileBrowser checks if the incoming request's User-Agent represents a mobile browser.
func isMobileBrowser(ua string) bool {
	if ua == "" {
		return false
	}
	lower := strings.ToLower(ua)
	mobileTokens := []string{
		"android", "iphone", "ipod", "ipad", "windows phone",
		"blackberry", "mobile", "opera mini", "iemobile", "webos",
	}
	for _, token := range mobileTokens {
		if strings.Contains(lower, token) {
			return true
		}
	}
	return false
}

// Handler: Root serves index.html (Desktop Web UI) or auto-redirects mobile browsers to /mobile
func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Automatic mobile browser detection: redirect to /mobile unless desktop view explicitly requested (?view=desktop)
	if r.URL.Query().Get("view") != "desktop" && isMobileBrowser(r.UserAgent()) {
		target := "/mobile"
		if r.URL.RawQuery != "" {
			target += "?" + r.URL.RawQuery
		}
		http.Redirect(w, r, target, http.StatusFound)
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
