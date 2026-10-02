package models

import "time"

// Movie represents a film currently showing or coming soon at Sri Lankan cinemas.
type Movie struct {
	ID            string   `json:"id"`
	Title         string   `json:"title"`
	OriginalTitle string   `json:"original_title,omitempty"`
	Language      string   `json:"language"`       // Sinhala, Tamil, English, Hindi
	Genre         string   `json:"genre"`          // Action, Romance, Drama, Animation, etc.
	Duration      string   `json:"duration"`       // e.g. "2h 18m"
	Rating        string   `json:"rating"`         // e.g. "U", "PG-13", "U/A", "A"
	Formats       []string `json:"formats"`        // ["IMAX 3D", "Dolby Atmos", "Laser 4K", "2D Digital"]
	CinemaIDs     []string `json:"cinema_ids"`     // Cinemas where this movie is playing
	Showtimes     []string `json:"showtimes"`      // e.g. ["10:30 AM", "1:45 PM", "4:30 PM", "7:15 PM"]
	Director      string   `json:"director,omitempty"`
	Cast          []string `json:"cast,omitempty"`
	Synopsis      string   `json:"synopsis"`
	Status        string   `json:"status"`         // "now_showing", "coming_soon"
	ReleaseDate   string   `json:"release_date,omitempty"`
	TrailerURL    string   `json:"trailer_url,omitempty"`
	ConcessionTip string   `json:"concession_tip,omitempty"` // Recommended concession combo for this movie
	LastUpdated   string   `json:"last_updated"`
}

// MovieSyncStatus reports the operational health and countdown of the 1-hour recurring movie sync engine.
type MovieSyncStatus struct {
	Status              string    `json:"status"`                 // "synced", "updating"
	SyncInterval        string    `json:"sync_interval"`          // "1h"
	SyncIntervalSeconds int       `json:"sync_interval_seconds"`  // 3600
	LastSyncedAt        time.Time `json:"last_synced_at"`
	NextSyncAt          time.Time `json:"next_sync_at"`
	SecondsUntilNext    int       `json:"seconds_until_next_sync"`
	TotalMovies         int       `json:"total_movies"`
	TotalCinemasCovered int       `json:"total_cinemas_covered"`
	SyncCount           int       `json:"sync_count"`
	Source              string    `json:"source"`
}

// MovieListResponse represents the JSON envelope for movies query.
type MovieListResponse struct {
	Status     string          `json:"status"`
	Count      int             `json:"count"`
	CinemaID   string          `json:"cinema_id,omitempty"`
	Language   string          `json:"language,omitempty"`
	SyncStatus MovieSyncStatus `json:"sync_status"`
	Movies     []Movie         `json:"movies"`
}
