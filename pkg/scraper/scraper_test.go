package scraper

import (
	"fmt"
	"sync"
	"testing"

	"cinema-recommender/internal/models"
)

func TestValidateLocationRejectsColomboAndAcceptsRegionalCities(t *testing.T) {
	s := NewRegionalScraper()

	colomboAreas := []string{"Colombo", "colombo-03", "Dehiwala", "Mount Lavinia", "Bambalapitiya", "Rajagiriya"}
	for _, loc := range colomboAreas {
		if err := s.ValidateLocation(loc); err == nil {
			t.Fatalf("expected %s to be rejected under Colombo exclusion policy", loc)
		}
	}

	regionalCities := []string{"Kandy", "Gampaha", "Galle", "Kurunegala", "Negombo", "Jaffna", "Matara", "Anuradhapura", "Ratnapura"}
	for _, city := range regionalCities {
		if err := s.ValidateLocation(city); err != nil {
			t.Fatalf("expected regional city %s to be accepted, got: %v", city, err)
		}
	}
}

func TestScrapeCinemaConcessionsSupportsRegionalCities(t *testing.T) {
	s := NewRegionalScraper()

	for _, city := range []string{"Kandy", "Gampaha", "Galle", "Kurunegala", "Negombo", "Jaffna", "Matara", "Anuradhapura", "Ratnapura"} {
		items, err := s.ScrapeCinemaConcessions(models.Cinema{ID: city + "-TEST", Name: "Cinema " + city, City: city})
		if err != nil {
			t.Fatalf("expected %s to be supported, got %v", city, err)
		}
		if len(items) == 0 {
			t.Fatalf("expected concession items for %s", city)
		}
	}
}

func TestCinemaRegistryQueries(t *testing.T) {
	s := NewRegionalScraper()

	cinemas := s.GetRegisteredCinemas()
	if len(cinemas) < 5 {
		t.Fatalf("expected at least 5 registered cinemas, got %d", len(cinemas))
	}

	cinema, err := s.GetCinemaByID("KND-KCC")
	if err != nil {
		t.Fatalf("failed to get cinema by ID: %v", err)
	}
	if cinema.City != "Kandy" {
		t.Fatalf("expected Kandy, got %s", cinema.City)
	}

	matches, err := s.GetCinemasByCity("Gampaha")
	if err != nil || len(matches) == 0 {
		t.Fatalf("expected cinemas for Gampaha, got err: %v", err)
	}
}

func TestScrapeRawConcessionFeed(t *testing.T) {
	s := NewRegionalScraper()

	htmlFeed := `
		<div class="concession-card" data-name="Truffle Popcorn" data-cat="Popcorn" data-price="1350.00"></div>
		<div class="concession-card" data-name="Mango Fizz Drink" data-cat="Beverage" data-price="620.00"></div>
	`
	items, err := s.ScrapeRawConcessionFeed("FEED-01", htmlFeed)
	if err != nil {
		t.Fatalf("unexpected error parsing feed: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}
	if items[0].Price != 1350.00 || items[1].Category != models.CategoryBeverage {
		t.Fatalf("parsed values do not match expected item definitions")
	}
}

func TestScraperConcurrency(t *testing.T) {
	s := NewRegionalScraper()
	var wg sync.WaitGroup

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			cities := []string{"Kandy", "Gampaha", "Galle", "Negombo"}
			city := cities[idx%len(cities)]
			items, err := s.ScrapeCinemaConcessions(models.Cinema{ID: fmt.Sprintf("CONC-%d", idx), City: city})
			if err != nil || len(items) == 0 {
				t.Errorf("concurrent scrape failed for %s: %v", city, err)
			}
		}(i)
	}
	wg.Wait()
}

func TestHaversineAndLocationSorting(t *testing.T) {
	s := NewRegionalScraper()

	// Coordinates near Kandy City Centre: 7.2936, 80.6385
	kandyLat := 7.2940
	kandyLng := 80.6390

	cinemas := s.GetCinemasWithLocationSort(kandyLat, kandyLng, "")
	if len(cinemas) == 0 {
		t.Fatalf("expected non-empty cinemas list")
	}

	// Closest should be Kandy City Centre Multiplex
	if cinemas[0].ID != "KND-KCC" {
		t.Fatalf("expected KND-KCC to be closest cinema, got %s", cinemas[0].ID)
	}

	if cinemas[0].DistanceKM > 2.0 {
		t.Fatalf("expected distance under 2km, got %.2f km", cinemas[0].DistanceKM)
	}

	// Verify food place metadata is populated
	if !cinemas[0].HasInHouseFood || cinemas[0].FoodPlaceName == "" {
		t.Fatalf("expected food place details for KND-KCC")
	}
	if len(cinemas[0].NearbyFoodOptions) == 0 {
		t.Fatalf("expected nearby food options for KND-KCC")
	}
}

func TestEveryFoodHasDiscounts(t *testing.T) {
	s := NewRegionalScraper()

	cinema, err := s.GetCinemaByID("KND-KCC")
	if err != nil {
		t.Fatalf("failed to retrieve cinema: %v", err)
	}

	items, err := s.ScrapeCinemaConcessions(cinema)
	if err != nil {
		t.Fatalf("failed to scrape concessions: %v", err)
	}

	for _, it := range items {
		if !it.HasDiscount {
			t.Errorf("expected food item %s to have HasDiscount=true", it.Name)
		}
		if len(it.ApplicableDiscounts) == 0 {
			t.Errorf("expected food item %s to have applicable discounts", it.Name)
		}
		if it.BestDiscountedPrice >= it.Price {
			t.Errorf("expected best discounted price (%.2f) to be lower than original price (%.2f)", it.BestDiscountedPrice, it.Price)
		}
		if it.FoodLocation == "" {
			t.Errorf("expected food item %s to have a food location specified", it.Name)
		}
	}
}
