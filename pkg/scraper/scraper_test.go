package scraper

import (
	"testing"

	"cinema-recommender/internal/models"
)

func TestValidateLocationRejectsColomboAndAcceptsRegionalCities(t *testing.T) {
	s := NewRegionalScraper()

	if err := s.ValidateLocation("Colombo"); err == nil {
		t.Fatal("expected Colombo to be rejected")
	}

	for _, city := range []string{"Kandy", "Gampaha", "Galle", "Kurunegala", "Negombo", "Jaffna", "Matara"} {
		if err := s.ValidateLocation(city); err != nil {
			t.Fatalf("expected %s to be accepted, got %v", city, err)
		}
	}
}

func TestScrapeCinemaConcessionsSupportsRegionalCities(t *testing.T) {
	s := NewRegionalScraper()

	for _, city := range []string{"Kandy", "Gampaha", "Galle", "Kurunegala", "Negombo"} {
		items, err := s.ScrapeCinemaConcessions(models.Cinema{ID: city, Name: "Cinema " + city, City: city})
		if err != nil {
			t.Fatalf("expected %s to be supported, got %v", city, err)
		}
		if len(items) == 0 {
			t.Fatalf("expected items for %s", city)
		}
	}
}
