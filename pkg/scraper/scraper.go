package scraper

import (
	"errors"
	"fmt"
	"strings"

	"cinema-recommender/internal/models"
)

var (
	// ErrColomboExcluded is returned when a cinema location is within Colombo metropolitan limits.
	ErrColomboExcluded = errors.New("location rejected: scraping engine is strictly scoped for locations outside Colombo")
	// ErrCinemaNotFound is returned when the requested regional cinema cannot be located.
	ErrCinemaNotFound = errors.New("cinema not found in regional registry")
)

// RegionalScraper defines operations for scraping out-of-Colombo cinema concessions.
type RegionalScraper struct {
	// allowedCities caches recognized regional districts/cities outside Colombo
	allowedCities map[string]bool
}

// NewRegionalScraper is an idiomatic Go constructor function.
func NewRegionalScraper() *RegionalScraper {
	return &RegionalScraper{
		allowedCities: map[string]bool{
			"kandy":      true,
			"gampaha":    true,
			"galle":      true,
			"kurunegala": true,
			"negombo":    true,
			"jaffna":     true,
			"matara":     true,
		},
	}
}

// ValidateLocation verifies that the target city is strictly outside Colombo.
func (s *RegionalScraper) ValidateLocation(city string) error {
	cityLower := strings.ToLower(strings.TrimSpace(city))
	if cityLower == "colombo" || strings.HasPrefix(cityLower, "colombo-") {
		return fmt.Errorf("%w: %q is in Colombo metropolitan area", ErrColomboExcluded, city)
	}
	return nil
}

// ScrapeCinemaConcessions fetches available food & beverage items for a regional cinema.
// In Go, functions idiomatically return the result along with an error as the last return value.
func (s *RegionalScraper) ScrapeCinemaConcessions(cinema models.Cinema) ([]models.ConcessionItem, error) {
	// Enforce out-of-Colombo boundary check
	if err := s.ValidateLocation(cinema.City); err != nil {
		return nil, err
	}

	// Mock scraped concession data for regional theaters
	switch strings.ToLower(cinema.City) {
	case "kandy":
		return []models.ConcessionItem{
			{ID: "KND-01", Name: "Jumbo Caramel Popcorn", Category: models.CategoryPopcorn, Price: 1200.00},
			{ID: "KND-02", Name: "Large Iced Mountain Dew", Category: models.CategoryBeverage, Price: 650.00},
			{ID: "KND-03", Name: "Spicy Chicken Hotdog", Category: models.CategorySnack, Price: 950.00},
			{ID: "KND-04", Name: "Hill Country Duo Combo", Category: models.CategoryCombo, Price: 2400.00},
		}, nil

	case "gampaha":
		return []models.ConcessionItem{
			{ID: "GMP-01", Name: "Salted Butter Popcorn (M)", Category: models.CategoryPopcorn, Price: 900.00},
			{ID: "GMP-02", Name: "Cold Milo Float", Category: models.CategoryBeverage, Price: 550.00},
			{ID: "GMP-03", Name: "Crispy Nachos & Cheese", Category: models.CategorySnack, Price: 850.00},
		}, nil

	case "galle":
		return []models.ConcessionItem{
			{ID: "GLE-01", Name: "Southern Cheese Popcorn", Category: models.CategoryPopcorn, Price: 1100.00},
			{ID: "GLE-02", Name: "Fresh Lime & Mint Juice", Category: models.CategoryBeverage, Price: 500.00},
			{ID: "GLE-03", Name: "Sweet Chili Fish & Chips", Category: models.CategorySnack, Price: 1400.00},
		}, nil

	default:
		return nil, fmt.Errorf("%w: no scraper pipeline registered for city %q", ErrCinemaNotFound, cinema.City)
	}
}
