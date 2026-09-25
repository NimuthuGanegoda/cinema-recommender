package server

import (
	"net/http"
	"strconv"
)

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
