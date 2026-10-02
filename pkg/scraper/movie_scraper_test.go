package scraper

import (
	"context"
	"testing"
	"time"
)

func TestMovieSyncEngine_Basics(t *testing.T) {
	engine := NewMovieSyncEngine()

	movies := engine.GetAllMovies()
	if len(movies) < 10 {
		t.Fatalf("expected at least 10 authentic movies, got %d", len(movies))
	}

	// Test cinema filter
	cccMovies := engine.GetMoviesByCinema("CMB-CCC")
	if len(cccMovies) == 0 {
		t.Errorf("expected movies for Colombo City Centre, got 0")
	}

	kccMovies := engine.GetMoviesByCinema("KDY-KCC")
	if len(kccMovies) == 0 {
		t.Errorf("expected movies for KCC Kandy, got 0")
	}

	// Test language filter
	sinhalaMovies := engine.GetMoviesByLanguage("Sinhala")
	if len(sinhalaMovies) == 0 {
		t.Errorf("expected Sinhala movies, got 0")
	}

	tamilMovies := engine.GetMoviesByLanguage("Tamil")
	if len(tamilMovies) == 0 {
		t.Errorf("expected Tamil movies, got 0")
	}

	englishMovies := engine.GetMoviesByLanguage("English")
	if len(englishMovies) == 0 {
		t.Errorf("expected English movies, got 0")
	}

	// Test sync status
	status := engine.GetSyncStatus()
	if status.SyncInterval != "1h" {
		t.Errorf("expected sync interval 1h, got %s", status.SyncInterval)
	}
	if status.SyncIntervalSeconds != 3600 {
		t.Errorf("expected 3600 seconds, got %d", status.SyncIntervalSeconds)
	}
	if status.SecondsUntilNext <= 0 || status.SecondsUntilNext > 3600 {
		t.Errorf("expected seconds until next to be between 1 and 3600, got %d", status.SecondsUntilNext)
	}
	if status.TotalMovies != len(movies) {
		t.Errorf("expected status total movies %d, got %d", len(movies), status.TotalMovies)
	}

	// Test force sync
	initialCount := status.SyncCount
	updatedStatus := engine.ForceSync()
	if updatedStatus.SyncCount != initialCount+1 {
		t.Errorf("expected sync count to increment to %d, got %d", initialCount+1, updatedStatus.SyncCount)
	}
}

func TestMovieSyncEngine_StartStop(t *testing.T) {
	engine := NewMovieSyncEngine()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	engine.Start(ctx)
	// Calling Start again should be idempotent
	engine.Start(ctx)

	time.Sleep(50 * time.Millisecond)
	engine.Stop()
}
