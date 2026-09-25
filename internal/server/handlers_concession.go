package server

import (
	"net/http"
	"time"
)

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
