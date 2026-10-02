package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"cinema-recommender/internal/models"
	"cinema-recommender/pkg/recommender"
	"cinema-recommender/pkg/scraper"
)

func TestMovieEndpoints(t *testing.T) {
	sc := scraper.NewRegionalScraper()
	recEngine := recommender.NewEngine()
	srv := NewServer(Config{Port: ":9999"}, sc, recEngine)
	handler := srv.GetHandler()

	t.Run("GET /api/v1/movies with 1-hour cache header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/movies", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rr.Code)
		}

		cacheControl := rr.Header().Get("Cache-Control")
		if cacheControl != "public, max-age=3600" {
			t.Errorf("expected 1-hour cache header 'public, max-age=3600', got %q", cacheControl)
		}

		var resp models.MovieListResponse
		if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to decode JSON response: %v", err)
		}

		if resp.Count < 10 {
			t.Errorf("expected at least 10 movies, got %d", resp.Count)
		}
		if resp.SyncStatus.SyncInterval != "1h" {
			t.Errorf("expected sync interval 1h, got %s", resp.SyncStatus.SyncInterval)
		}
	})

	t.Run("GET /api/v1/movies?cinema_id=CMB-CCC", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/movies?cinema_id=CMB-CCC", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rr.Code)
		}

		var resp models.MovieListResponse
		if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to decode JSON response: %v", err)
		}

		if resp.Count == 0 {
			t.Errorf("expected movies for CMB-CCC, got 0")
		}
		for _, m := range resp.Movies {
			found := false
			for _, cid := range m.CinemaIDs {
				if cid == "CMB-CCC" {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("movie %s does not contain cinema CMB-CCC", m.Title)
			}
		}
	})

	t.Run("GET /api/v1/movies?language=Sinhala", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/movies?language=Sinhala", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rr.Code)
		}

		var resp models.MovieListResponse
		if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to decode JSON response: %v", err)
		}

		if resp.Count == 0 {
			t.Errorf("expected Sinhala movies, got 0")
		}
		for _, m := range resp.Movies {
			if m.Language != "Sinhala" {
				t.Errorf("expected Sinhala movie, got %s", m.Language)
			}
		}
	})

	t.Run("GET /api/v1/movies/status", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/movies/status", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rr.Code)
		}

		var status models.MovieSyncStatus
		if err := json.Unmarshal(rr.Body.Bytes(), &status); err != nil {
			t.Fatalf("failed to decode status: %v", err)
		}

		if status.SyncInterval != "1h" {
			t.Errorf("expected 1h sync interval, got %s", status.SyncInterval)
		}
		if status.SecondsUntilNext <= 0 || status.SecondsUntilNext > 3600 {
			t.Errorf("expected seconds until next to be in (0, 3600], got %d", status.SecondsUntilNext)
		}
	})

	t.Run("POST /api/v1/movies/sync triggers immediate refresh", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/movies/sync", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rr.Code)
		}

		var body map[string]interface{}
		if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
			t.Fatalf("failed to decode sync response: %v", err)
		}

		if body["message"] == nil {
			t.Errorf("expected message in response")
		}
	})

	t.Run("GET /api/v1/cinemas/CMB-CCC/movies", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/cinemas/CMB-CCC/movies", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rr.Code)
		}

		var resp models.MovieListResponse
		if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to decode JSON response: %v", err)
		}

		if resp.Count == 0 {
			t.Errorf("expected movies for CMB-CCC, got 0")
		}
	})
}
