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
	if cityLower == "" {
		return fmt.Errorf("%w: city name is required", ErrCinemaNotFound)
	}
	if cityLower == "colombo" || strings.HasPrefix(cityLower, "colombo-") {
		return fmt.Errorf("%w: %q is in Colombo metropolitan area", ErrColomboExcluded, city)
	}
	if !s.allowedCities[cityLower] {
		return fmt.Errorf("%w: no scraper pipeline registered for city %q", ErrCinemaNotFound, city)
	}
	return nil
}

// ScrapeCinemaConcessions fetches available food & beverage items for a regional cinema.
func (s *RegionalScraper) ScrapeCinemaConcessions(cinema models.Cinema) ([]models.ConcessionItem, error) {
	if err := s.ValidateLocation(cinema.City); err != nil {
		return nil, err
	}

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

	case "kurunegala":
		return []models.ConcessionItem{
			{ID: "KRN-01", Name: "North Star Popcorn Bucket", Category: models.CategoryPopcorn, Price: 1000.00},
			{ID: "KRN-02", Name: "Fresh Cola Combo", Category: models.CategoryBeverage, Price: 600.00},
			{ID: "KRN-03", Name: "Crispy Chicken Roll", Category: models.CategorySnack, Price: 820.00},
		}, nil

	case "negombo":
		return []models.ConcessionItem{
			{ID: "NGB-01", Name: "Coastal Caramel Popcorn", Category: models.CategoryPopcorn, Price: 980.00},
			{ID: "NGB-02", Name: "Iced Lemon Soda", Category: models.CategoryBeverage, Price: 520.00},
			{ID: "NGB-03", Name: "Cheese Nachos Deluxe", Category: models.CategorySnack, Price: 930.00},
		}, nil

	case "jaffna":
		return []models.ConcessionItem{
			{ID: "JFN-01", Name: "Jaffna Crunch Popcorn", Category: models.CategoryPopcorn, Price: 1150.00},
			{ID: "JFN-02", Name: "Coconut Cooler", Category: models.CategoryBeverage, Price: 700.00},
			{ID: "JFN-03", Name: "Hot Ceylon Chicken Bites", Category: models.CategorySnack, Price: 960.00},
		}, nil

	case "matara":
		return []models.ConcessionItem{
			{ID: "MTR-01", Name: "Southern Popcorn Mix", Category: models.CategoryPopcorn, Price: 1050.00},
			{ID: "MTR-02", Name: "Fresh Orange Fizz", Category: models.CategoryBeverage, Price: 580.00},
			{ID: "MTR-03", Name: "Crispy Wedges Platter", Category: models.CategorySnack, Price: 890.00},
		}, nil

	default:
		return nil, fmt.Errorf("%w: no scraper pipeline registered for city %q", ErrCinemaNotFound, cinema.City)
	}
}
