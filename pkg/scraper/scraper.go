package scraper

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"cinema-recommender/internal/models"
)

var (
	// ErrColomboExcluded is returned when a cinema location is within Colombo metropolitan limits.
	ErrColomboExcluded = errors.New("location rejected: scraping engine is strictly scoped for regional theaters outside Colombo")
	// ErrCinemaNotFound is returned when the requested regional cinema cannot be located.
	ErrCinemaNotFound = errors.New("cinema not found in regional registry")
	// ErrInvalidInput is returned when an input parameter is missing or malformed.
	ErrInvalidInput = errors.New("invalid cinema or location parameters")
)

// RegionalScraper defines operations for scraping out-of-Colombo cinema concessions.
type RegionalScraper struct {
	mu             sync.RWMutex
	allowedCities  map[string]bool
	cinemaRegistry []models.Cinema
	menuCache      map[string]cachedMenu
	cacheTTL       time.Duration
}

type cachedMenu struct {
	items     []models.ConcessionItem
	timestamp time.Time
}

// NewRegionalScraper creates an initialized regional scraper with outstation cinema registries.
func NewRegionalScraper() *RegionalScraper {
	sc := &RegionalScraper{
		allowedCities: map[string]bool{
			"kandy":        true,
			"gampaha":      true,
			"galle":        true,
			"kurunegala":   true,
			"negombo":      true,
			"jaffna":       true,
			"matara":       true,
			"anuradhapura": true,
			"ratnapura":    true,
		},
		menuCache: make(map[string]cachedMenu),
		cacheTTL:  10 * time.Minute,
	}

	sc.cinemaRegistry = defaultCinemaRegistry()
	return sc
}

// ValidateLocation verifies that the target city is strictly outside Colombo.
func (s *RegionalScraper) ValidateLocation(city string) error {
	cityClean := strings.ToLower(strings.TrimSpace(city))
	if cityClean == "" {
		return fmt.Errorf("%w: city name cannot be empty", ErrInvalidInput)
	}

	// Strictly exclude Colombo metropolitan limits
	colomboPrefixes := []string{"colombo", "colombo-", "dehiwala", "mount lavinia", "kollupitiya", "bambalapitiya", "rajagiriya"}
	for _, p := range colomboPrefixes {
		if cityClean == p || strings.HasPrefix(cityClean, p) {
			return fmt.Errorf("%w: %q is within the Colombo metropolitan exclusion zone", ErrColomboExcluded, city)
		}
	}

	if !s.allowedCities[cityClean] {
		return fmt.Errorf("%w: no scraper pipeline registered for regional city %q", ErrCinemaNotFound, city)
	}

	return nil
}

// ScrapeCinemaConcessions extracts current concession inventory, prices, and attached discounts for a regional cinema.
func (s *RegionalScraper) ScrapeCinemaConcessions(cinema models.Cinema) ([]models.ConcessionItem, error) {
	if err := s.ValidateLocation(cinema.City); err != nil {
		return nil, err
	}

	cacheKey := fmt.Sprintf("%s:%s", strings.ToLower(cinema.City), strings.ToLower(cinema.ID))
	s.mu.RLock()
	cached, found := s.menuCache[cacheKey]
	s.mu.RUnlock()

	var rawItems []models.ConcessionItem
	if found && time.Since(cached.timestamp) < s.cacheTTL {
		rawItems = cached.items
	} else {
		items, err := s.scrapeByCity(cinema.City, cinema.ID)
		if err != nil {
			return nil, err
		}

		s.mu.Lock()
		s.menuCache[cacheKey] = cachedMenu{
			items:     items,
			timestamp: time.Now(),
		}
		s.mu.Unlock()
		rawItems = items
	}

	foodLoc := cinema.FoodPlaceName
	if foodLoc == "" {
		foodLoc = "In-House Concession Stand"
	}

	// Apply constantly updated real-time dynamic pricing
	liveItems := ApplyLiveDynamicPricing(rawItems, foodLoc, time.Now())
	return liveItems, nil
}
