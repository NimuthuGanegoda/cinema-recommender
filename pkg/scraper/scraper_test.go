package scraper

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"cinema-recommender/internal/models"
)

func TestValidateLocationSupportsColomboAndRegionalCities(t *testing.T) {
	s := NewRegionalScraper()

	acceptedCities := []string{"Colombo", "colombo-03", "Kandy", "Gampaha", "Galle", "Kurunegala", "Negombo", "Jaffna", "Matara", "Anuradhapura", "Ratnapura"}
	for _, city := range acceptedCities {
		if err := s.ValidateLocation(city); err != nil {
			t.Fatalf("expected city %s to be accepted, got: %v", city, err)
		}
	}

	invalidCities := []string{"", "Paris", "Sydney", "New York", "Tokyo"}
	for _, city := range invalidCities {
		if err := s.ValidateLocation(city); err == nil {
			t.Fatalf("expected invalid city %q to be rejected", city)
		}
	}
}

func TestScrapeCinemaConcessionsSupportsAllCities(t *testing.T) {
	s := NewRegionalScraper()

	for _, city := range []string{"Colombo", "Kandy", "Gampaha", "Galle", "Kurunegala", "Negombo", "Jaffna", "Matara", "Anuradhapura", "Ratnapura"} {
		items, err := s.ScrapeCinemaConcessions(models.Cinema{ID: city + "-TEST", Name: "Cinema " + city, City: city})
		if err != nil {
			t.Fatalf("expected %s to be supported, got %v", city, err)
		}
		if len(items) == 0 {
			t.Fatalf("expected concession items for %s", city)
		}
	}
}

func TestColomboScopeCinemasRegistry(t *testing.T) {
	s := NewRegionalScraper()

	colomboCinemas, err := s.GetCinemasByCity("Colombo")
	if err != nil {
		t.Fatalf("expected to find Colombo cinemas, got: %v", err)
	}
	if len(colomboCinemas) != 3 {
		t.Fatalf("expected 3 Colombo Scope Cinemas, got %d", len(colomboCinemas))
	}

	expectedIDs := map[string]bool{"CMB-CCC": true, "CMB-HCM": true, "CMB-LBT": true}
	for _, c := range colomboCinemas {
		if !expectedIDs[c.ID] {
			t.Errorf("unexpected cinema ID: %s", c.ID)
		}
		if !c.HasFoodCourtInFront {
			t.Errorf("cinema %s expected to have food court in front", c.ID)
		}
		if !c.HasInHouseFood {
			t.Errorf("cinema %s expected to have in-house food", c.ID)
		}
		if c.Chain != "Scope Cinemas" {
			t.Errorf("expected chain Scope Cinemas, got %s", c.Chain)
		}
	}

	ccc, err := s.GetCinemaByID("CMB-CCC")
	if err != nil {
		t.Fatalf("failed to get CMB-CCC: %v", err)
	}
	if ccc.Name != "Scope Cinemas Multiplex - Colombo City Centre" {
		t.Fatalf("expected Scope CCC name, got %s", ccc.Name)
	}
}

func TestCinemaRegistryQueries(t *testing.T) {
	s := NewRegionalScraper()

	cinemas := s.GetRegisteredCinemas()
	if len(cinemas) < 8 {
		t.Fatalf("expected at least 8 registered cinemas, got %d", len(cinemas))
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

func TestLivePriceUpdatesConstantly(t *testing.T) {
	s := NewRegionalScraper()

	cinema, err := s.GetCinemaByID("KND-KCC")
	if err != nil {
		t.Fatalf("failed to retrieve cinema: %v", err)
	}

	ticker := s.GetLiveConcessionTicker(cinema)
	if ticker["live_pricing_active"] != true {
		t.Fatalf("expected live_pricing_active to be true")
	}

	// Verify price calculation at two different time points (tick 0 vs tick 1)
	rawItems, err := s.scrapeByCity("Kandy", cinema.ID)
	if err != nil {
		t.Fatalf("failed to scrape raw items: %v", err)
	}

	t1 := time.Unix(1700000000, 0)
	t2 := time.Unix(1700000015, 0) // 15 seconds later (next dynamic pricing tick)

	itemsT1 := ApplyLiveDynamicPricing(rawItems, cinema.FoodPlaceName, t1)
	itemsT2 := ApplyLiveDynamicPricing(rawItems, cinema.FoodPlaceName, t2)

	if len(itemsT1) != len(itemsT2) {
		t.Fatalf("expected item counts to match across ticks")
	}

	// Verify dynamic pricing fields are populated
	for _, it := range itemsT1 {
		if it.BasePrice <= 0 || it.Price <= 0 {
			t.Errorf("expected positive base and dynamic prices for %s", it.Name)
		}
		if it.LastPriceUpdate == "" {
			t.Errorf("expected LastPriceUpdate timestamp for %s", it.Name)
		}
		if it.PriceTrend == "" {
			t.Errorf("expected PriceTrend for %s", it.Name)
		}
	}

	// Verify that at least some item's price or tick state changes between tick 1 and tick 2
	var hasDifferentState bool
	for i := range itemsT1 {
		if itemsT1[i].LiveTickID != itemsT2[i].LiveTickID || itemsT1[i].Price != itemsT2[i].Price {
			hasDifferentState = true
			break
		}
	}

	if !hasDifferentState {
		t.Fatalf("expected dynamic pricing engine to advance live tick and update prices")
	}
}

func TestAllCinemasHaveFoodCourtInFrontForCinemaFoodAndDrinksOnly(t *testing.T) {
	s := NewRegionalScraper()

	cinemas := s.GetRegisteredCinemas()
	if len(cinemas) == 0 {
		t.Fatalf("expected registered cinemas list to be non-empty")
	}

	for _, c := range cinemas {
		if !c.HasFoodCourtInFront {
			t.Errorf("cinema %s (%s) must have HasFoodCourtInFront=true", c.Name, c.ID)
		}
		if !c.HasInHouseFood {
			t.Errorf("cinema %s (%s) must have HasInHouseFood=true", c.Name, c.ID)
		}
		if c.FoodCourtName == "" {
			t.Errorf("cinema %s (%s) must have FoodCourtName specified", c.Name, c.ID)
		}
		if c.FoodCourtLocation == "" {
			t.Errorf("cinema %s (%s) must have FoodCourtLocation specified", c.Name, c.ID)
		}
	}

	// Also verify location sorted query preserves strict in-front food court rule
	sortedCinemas := s.GetCinemasWithLocationSort(7.2936, 80.6385, "")
	for _, c := range sortedCinemas {
		if !c.HasFoodCourtInFront || !c.HasInHouseFood {
			t.Errorf("sorted cinema %s must have HasFoodCourtInFront && HasInHouseFood", c.ID)
		}
	}
}
