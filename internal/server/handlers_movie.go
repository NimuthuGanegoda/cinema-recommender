package server

import (
	"net/http"
	"strings"

	"cinema-recommender/internal/models"
)

// handleListMovies returns all movies, optionally filtered by cinema_id or language, with a 1-hour cache header.
func (s *Server) handleListMovies(w http.ResponseWriter, r *http.Request) {
	cinemaID := strings.TrimSpace(r.URL.Query().Get("cinema_id"))
	language := strings.TrimSpace(r.URL.Query().Get("language"))

	var movies []models.Movie
	if cinemaID != "" {
		movies = s.movieEngine.GetMoviesByCinema(cinemaID)
	} else if language != "" {
		movies = s.movieEngine.GetMoviesByLanguage(language)
	} else {
		movies = s.movieEngine.GetAllMovies()
	}

	syncStatus := s.movieEngine.GetSyncStatus()

	// 1-hour browser & CDN caching header aligned with hourly sync engine
	w.Header().Set("Cache-Control", "public, max-age=3600")

	s.writeJSON(w, http.StatusOK, models.MovieListResponse{
		Status:     "success",
		Count:      len(movies),
		CinemaID:   cinemaID,
		Language:   language,
		SyncStatus: syncStatus,
		Movies:     movies,
	})
}

// handleMovieSyncStatus returns real-time operational metrics for the 1-hour background recurring sync.
func (s *Server) handleMovieSyncStatus(w http.ResponseWriter, r *http.Request) {
	status := s.movieEngine.GetSyncStatus()
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	s.writeJSON(w, http.StatusOK, status)
}

// handleTriggerMovieSync forces an immediate sync cycle and returns fresh movie telemetry.
func (s *Server) handleTriggerMovieSync(w http.ResponseWriter, r *http.Request) {
	status := s.movieEngine.ForceSync()
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	s.writeJSON(w, http.StatusOK, map[string]interface{}{
		"message":     "1-hour movie catalog sync triggered successfully",
		"sync_status": status,
		"movies":      s.movieEngine.GetAllMovies(),
	})
}

// handleGetCinemaMovies returns movies showing at a specific cinema path param.
func (s *Server) handleGetCinemaMovies(w http.ResponseWriter, r *http.Request) {
	cinemaID := strings.TrimPrefix(r.URL.Path, "/api/v1/cinemas/")
	cinemaID = strings.TrimSuffix(cinemaID, "/movies")
	cinemaID = strings.TrimSpace(cinemaID)

	if cinemaID == "" {
		s.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Cinema ID is required"})
		return
	}

	movies := s.movieEngine.GetMoviesByCinema(cinemaID)
	syncStatus := s.movieEngine.GetSyncStatus()

	w.Header().Set("Cache-Control", "public, max-age=3600")
	s.writeJSON(w, http.StatusOK, models.MovieListResponse{
		Status:     "success",
		Count:      len(movies),
		CinemaID:   cinemaID,
		SyncStatus: syncStatus,
		Movies:     movies,
	})
}
