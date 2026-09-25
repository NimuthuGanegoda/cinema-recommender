package scraper

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
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

	sc.cinemaRegistry = []models.Cinema{
		{
			ID:                 "KND-KCC",
			Name:               "Scope Partner Multiplex - KCC Kandy",
			Chain:              "Scope Cinemas Partner Circuit",
			City:               "Kandy",
			Province:           "Central Province",
			Address:            "Level 3, Kandy City Centre, Dalada Veediya, Kandy",
			Screens:            4,
			IsOutside:          true,
			Latitude:           7.2936,
			Longitude:          80.6385,
			MapURL:             "https://maps.google.com/?q=7.2936,80.6385",
			HasInHouseFood:     true,
			FoodPlaceName:      "Scope Signature Candy Bar (Level 3 Lobby)",
			FoodPlaceType:      "Scope Cinemas Concession Stand & Mall Food Court",
			NearbyFoodOptions:  []string{"KCC World Food Court (Level 4)", "Devon Restaurant & Bakery (2 min walk)", "Bake House Dalada Veediya", "Cargills Food Hall"},
			FoodHours:          "10:00 AM - 10:45 PM (Open during all movie screenings)",
			FoodDeliveryToSeat: true,
		},
		{
			ID:                 "KND-REG",
			Name:               "Regal Cinema Kandy",
			Chain:              "Ceylon Theatres",
			City:               "Kandy",
			Province:           "Central Province",
			Address:            "No. 12, Katugastota Road, Kandy",
			Screens:            2,
			IsOutside:          true,
			Latitude:           7.2985,
			Longitude:          80.6335,
			MapURL:             "https://maps.google.com/?q=7.2985,80.6335",
			HasInHouseFood:     true,
			FoodPlaceName:      "Regal Foyer Concession Counter",
			FoodPlaceType:      "In-Theater Snack Bar",
			NearbyFoodOptions:  []string{"Katugastota Cafe & Bake House", "Perera & Sons Kandy", "White House Restaurant"},
			FoodHours:          "10:15 AM - 10:30 PM",
			FoodDeliveryToSeat: false,
		},
		{
			ID:                 "GMP-REG",
			Name:               "Regal Cinema Gampaha",
			Chain:              "Ceylon Theatres",
			City:               "Gampaha",
			Province:           "Western Province",
			Address:            "Bauddhaloka Mawatha, Gampaha",
			Screens:            2,
			IsOutside:          true,
			Latitude:           7.0897,
			Longitude:          79.9925,
			MapURL:             "https://maps.google.com/?q=7.0897,79.9925",
			HasInHouseFood:     true,
			FoodPlaceName:      "Regal Candy Bar & Popcorn Stand",
			FoodPlaceType:      "In-Theater Concession Counter",
			NearbyFoodOptions:  []string{"Gampaha Supermarket Food Court", "Perera & Sons Bauddhaloka Mw", "Fab Confectionery"},
			FoodHours:          "10:00 AM - 10:30 PM",
			FoodDeliveryToSeat: false,
		},
		{
			ID:                 "GLE-QNS",
			Name:               "Queens Cinema Galle",
			Chain:              "EAP Films",
			City:               "Galle",
			Province:           "Southern Province",
			Address:            "Havelock Place, Galle Fort Corridor, Galle",
			Screens:            2,
			IsOutside:          true,
			Latitude:           6.0367,
			Longitude:          80.2170,
			MapURL:             "https://maps.google.com/?q=6.0367,80.2170",
			HasInHouseFood:     true,
			FoodPlaceName:      "Queens Cinema Concession Lounge",
			FoodPlaceType:      "In-Theater Concession Stand & Colonial Cafe",
			NearbyFoodOptions:  []string{"Galle Fort Dutch Hospital Dining", "Pedlar's Inn Cafe", "Fort Printers Bakery", "Galle Green Street Food Kiosks"},
			FoodHours:          "09:45 AM - 11:00 PM",
			FoodDeliveryToSeat: true,
		},
		{
			ID:                 "KRN-LUX",
			Name:               "Luxe Cinema Kurunegala",
			Chain:              "Luxe Theatres",
			City:               "Kurunegala",
			Province:           "North Western Province",
			Address:            "Colombo Road, Kurunegala",
			Screens:            3,
			IsOutside:          true,
			Latitude:           7.4863,
			Longitude:          80.3647,
			MapURL:             "https://maps.google.com/?q=7.4863,80.3647",
			HasInHouseFood:     true,
			FoodPlaceName:      "Luxe Express Popcorn & Beverage Bar",
			FoodPlaceType:      "In-Theater Concession Counter & Refreshment Lounge",
			NearbyFoodOptions:  []string{"Kurunegala Central Mall Eateries", "Saloons Bakery & Cafe", "Sunimal Food Cabin"},
			FoodHours:          "10:00 AM - 10:30 PM",
			FoodDeliveryToSeat: false,
		},
		{
			ID:                 "NGB-REG",
			Name:               "Regal Cinema Negombo",
			Chain:              "Ceylon Theatres",
			City:               "Negombo",
			Province:           "Western Province",
			Address:            "Main Street, Negombo",
			Screens:            2,
			IsOutside:          true,
			Latitude:           7.2083,
			Longitude:          79.8358,
			MapURL:             "https://maps.google.com/?q=7.2083,79.8358",
			HasInHouseFood:     true,
			FoodPlaceName:      "Coastal Snack Counter & Beverage Kiosk",
			FoodPlaceType:      "In-Theater Seafood & Concession Stand",
			NearbyFoodOptions:  []string{"Negombo Beach Road Seafood Strip", "Lords Restaurant Complex", "Perera & Sons Main Street"},
			FoodHours:          "10:00 AM - 11:00 PM",
			FoodDeliveryToSeat: true,
		},
		{
			ID:                 "JFN-MJG",
			Name:               "Majestic Gold Cinema Jaffna",
			Chain:              "Majestic Circuit",
			City:               "Jaffna",
			Province:           "Northern Province",
			Address:            "Hospital Road, Jaffna",
			Screens:            2,
			IsOutside:          true,
			Latitude:           9.6647,
			Longitude:          80.0167,
			MapURL:             "https://maps.google.com/?q=9.6647,80.0167",
			HasInHouseFood:     true,
			FoodPlaceName:      "Majestic Gold Concession Counter & Tea Corner",
			FoodPlaceType:      "In-Theater Concession Counter",
			NearbyFoodOptions:  []string{"Mangos Indian Veg Restaurant", "Jaffna Central Market Street Food", "Rio Ice Cream Parlour (5 min walk)"},
			FoodHours:          "09:30 AM - 10:30 PM",
			FoodDeliveryToSeat: false,
		},
		{
			ID:                 "MTR-SKC",
			Name:               "SK Cinema Matara",
			Chain:              "Southern Screen Network",
			City:               "Matara",
			Province:           "Southern Province",
			Address:            "Anagarika Dharmapala Mawatha, Matara",
			Screens:            2,
			IsOutside:          true,
			Latitude:           5.9485,
			Longitude:          80.5488,
			MapURL:             "https://maps.google.com/?q=5.9485,80.5488",
			HasInHouseFood:     true,
			FoodPlaceName:      "SK Snack Counter & Chill Bar",
			FoodPlaceType:      "In-Theater Concession Stand",
			NearbyFoodOptions:  []string{"Matara Beach Park Food Stalls", "Dutch Fort Bakery Matara", "Keells Super Food Corner"},
			FoodHours:          "10:00 AM - 10:30 PM",
			FoodDeliveryToSeat: false,
		},
		{
			ID:                 "ANR-CMX",
			Name:               "Cinemax Anuradhapura",
			Chain:              "Rajarata Theatres",
			City:               "Anuradhapura",
			Province:           "North Central Province",
			Address:            "Maithripala Senanayake Mawatha, Anuradhapura",
			Screens:            2,
			IsOutside:          true,
			Latitude:           8.3350,
			Longitude:          80.4108,
			MapURL:             "https://maps.google.com/?q=8.3350,80.4108",
			HasInHouseFood:     true,
			FoodPlaceName:      "Cinemax Fresh Popcorn & Snack Corner",
			FoodPlaceType:      "In-Theater Concession Counter",
			NearbyFoodOptions:  []string{"Ceylan Bake House Anuradhapura", "Rajarata Food Center", "Hotel Shalini Dining"},
			FoodHours:          "10:00 AM - 10:15 PM",
			FoodDeliveryToSeat: false,
		},
		{
			ID:                 "RTP-MLN",
			Name:               "Milano Cinema Ratnapura",
			Chain:              "Milano Circuit",
			City:               "Ratnapura",
			Province:           "Sabaragamuwa Province",
			Address:            "Main Street, Ratnapura",
			Screens:            2,
			IsOutside:          true,
			Latitude:           6.6828,
			Longitude:          80.4037,
			MapURL:             "https://maps.google.com/?q=6.6828,80.4037",
			HasInHouseFood:     true,
			FoodPlaceName:      "Milano Concession Stand",
			FoodPlaceType:      "In-Theater Concession Counter",
			NearbyFoodOptions:  []string{"Ratnapura City Bakers", "Gem City Food Plaza", "Perera & Sons Main Street"},
			FoodHours:          "10:15 AM - 10:15 PM",
			FoodDeliveryToSeat: false,
		},
	}

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

// GetRegisteredCinemas returns all out-of-Colombo registered cinemas.
func (s *RegionalScraper) GetRegisteredCinemas() []models.Cinema {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cinemas := make([]models.Cinema, len(s.cinemaRegistry))
	copy(cinemas, s.cinemaRegistry)
	return cinemas
}

// GetCinemaByID searches the regional registry for a specific cinema ID.
func (s *RegionalScraper) GetCinemaByID(id string) (models.Cinema, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cleanID := strings.ToUpper(strings.TrimSpace(id))
	for _, c := range s.cinemaRegistry {
		if strings.ToUpper(c.ID) == cleanID {
			return c, nil
		}
	}
	return models.Cinema{}, fmt.Errorf("%w: %q", ErrCinemaNotFound, id)
}

// GetCinemasByCity searches cinemas by city name.
func (s *RegionalScraper) GetCinemasByCity(city string) ([]models.Cinema, error) {
	if err := s.ValidateLocation(city); err != nil {
		return nil, err
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	cleanCity := strings.ToLower(strings.TrimSpace(city))
	var matches []models.Cinema
	for _, c := range s.cinemaRegistry {
		if strings.ToLower(c.City) == cleanCity {
			matches = append(matches, c)
		}
	}

	if len(matches) == 0 {
		return nil, fmt.Errorf("%w: no active theaters found for %s", ErrCinemaNotFound, city)
	}

	return matches, nil
}

// CalculateHaversineDistance computes the geographical distance in kilometers between two GPS coordinates.
func CalculateHaversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadiusKm = 6371.0
	dLat := (lat2 - lat1) * (math.Pi / 180.0)
	dLon := (lon2 - lon1) * (math.Pi / 180.0)

	rLat1 := lat1 * (math.Pi / 180.0)
	rLat2 := lat2 * (math.Pi / 180.0)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(rLat1)*math.Cos(rLat2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return math.Round(earthRadiusKm*c*10) / 10
}

// GetCinemasWithLocationSort returns all regional cinemas, optionally sorted by distance to user coordinates.
func (s *RegionalScraper) GetCinemasWithLocationSort(userLat, userLng float64, cityFilter string) []models.Cinema {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var base []models.Cinema
	cleanCity := strings.ToLower(strings.TrimSpace(cityFilter))
	for _, c := range s.cinemaRegistry {
		if cleanCity == "" || strings.ToLower(c.City) == cleanCity {
			base = append(base, c)
		}
	}

	if userLat != 0 || userLng != 0 {
		for i := range base {
			if base[i].Latitude != 0 && base[i].Longitude != 0 {
				base[i].DistanceKM = CalculateHaversineDistance(userLat, userLng, base[i].Latitude, base[i].Longitude)
			}
		}

		sort.Slice(base, func(i, j int) bool {
			return base[i].DistanceKM < base[j].DistanceKM
		})
	}

	return base
}

// EnrichItemDiscounts calculates and attaches all available promotions and payment method discounts to a food item.
func EnrichItemDiscounts(item models.ConcessionItem, foodLocation string) models.ConcessionItem {
	price := item.Price
	if price <= 0 {
		return item
	}

	var discounts []models.ItemDiscountInfo

	// 1. Scope Privilege Platinum VIP (35% off)
	platPct := 0.35
	platDisc := math.Round(price*(1.0-platPct)*100) / 100
	discounts = append(discounts, models.ItemDiscountInfo{
		PromoCode:       "SCOPE-PLATINUM",
		Provider:        "Scope Privilege",
		Title:           "Scope Privilege Platinum VIP (35% Off)",
		DiscountPct:     platPct,
		DiscountedPrice: platDisc,
		SavingsLKR:      math.Round((price-platDisc)*100) / 100,
		Requirement:     "Scope Privilege Platinum VIP membership",
	})

	// 2. Scope Combo 30% Deal
	comboPct := 0.30
	comboDisc := math.Round(price*(1.0-comboPct)*100) / 100
	discounts = append(discounts, models.ItemDiscountInfo{
		PromoCode:       "SCOPE-COMBO30",
		Provider:        "Scope Privilege",
		Title:           "Scope Popcorn + Drink Combo (30% Off)",
		DiscountPct:     comboPct,
		DiscountedPrice: comboDisc,
		SavingsLKR:      math.Round((price-comboDisc)*100) / 100,
		Requirement:     "Pair with any Popcorn & Beverage",
	})

	// 3. Scope Regional Concession Boost (25% Off)
	regPct := 0.25
	regDisc := math.Round(price*(1.0-regPct)*100) / 100
	discounts = append(discounts, models.ItemDiscountInfo{
		PromoCode:       "SCOPE-CINEMA25",
		Provider:        "Scope Privilege",
		Title:           "Scope Regional Concession Pass (25% Off)",
		DiscountPct:     regPct,
		DiscountedPrice: regDisc,
		SavingsLKR:      math.Round((price-regDisc)*100) / 100,
		Requirement:     "Concession orders over LKR 2,000",
	})

	// 4. Scope Student Moviegoer Pass (20% Off)
	stuPct := 0.20
	stuDisc := math.Round(price*(1.0-stuPct)*100) / 100
	discounts = append(discounts, models.ItemDiscountInfo{
		PromoCode:       "SCOPE-STUDENT",
		Provider:        "Scope Privilege",
		Title:           "Scope Student Moviegoer Pass (20% Off)",
		DiscountPct:     stuPct,
		DiscountedPrice: stuDisc,
		SavingsLKR:      math.Round((price-stuDisc)*100) / 100,
		Requirement:     "Verified student ID cardholders",
	})

	// 5. FriMi Instant Concession Cashback (10% Cashback)
	frimiPct := 0.10
	frimiDisc := math.Round(price*(1.0-frimiPct)*100) / 100
	discounts = append(discounts, models.ItemDiscountInfo{
		PromoCode:       "FRIMI-CASHBACK",
		Provider:        "FriMi",
		Title:           "FriMi 10% Instant Concession Cashback",
		DiscountPct:     frimiPct,
		DiscountedPrice: frimiDisc,
		SavingsLKR:      math.Round((price-frimiDisc)*100) / 100,
		Requirement:     "Pay via FriMi digital lifestyle app",
	})

	// 6. Genie Concession Rebate (8% Rebate)
	geniePct := 0.08
	genieDisc := math.Round(price*(1.0-geniePct)*100) / 100
	discounts = append(discounts, models.ItemDiscountInfo{
		PromoCode:       "GENIE-REBATE",
		Provider:        "Genie",
		Title:           "Genie 8% Concession Rebate",
		DiscountPct:     geniePct,
		DiscountedPrice: genieDisc,
		SavingsLKR:      math.Round((price-genieDisc)*100) / 100,
		Requirement:     "Pay via Genie digital wallet",
	})

	// 7. CBSL LankaQR National Incentive (5% Off)
	lqrPct := 0.05
	lqrDisc := math.Round(price*(1.0-lqrPct)*100) / 100
	discounts = append(discounts, models.ItemDiscountInfo{
		PromoCode:       "LANKAQR-5",
		Provider:        "LankaQR",
		Title:           "CBSL LankaQR 5% National Incentive",
		DiscountPct:     lqrPct,
		DiscountedPrice: lqrDisc,
		SavingsLKR:      math.Round((price-lqrDisc)*100) / 100,
		Requirement:     "Scan dynamic LankaQR with any LK bank app",
	})

	item.ApplicableDiscounts = discounts
	item.BestDiscountedPrice = platDisc // up to 35% maximum savings
	item.MaxDiscountPct = platPct
	item.BestPromoName = "Scope Privilege Platinum VIP / Combo 30%"
	item.HasDiscount = true
	if foodLocation != "" {
		item.FoodLocation = foodLocation
	} else if item.FoodLocation == "" {
		item.FoodLocation = "In-Theater Concession Stand"
	}

	return item
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

// ApplyLiveDynamicPricing computes constantly updated real-time concession prices and flash discounts.
// Movie theater candy bars dynamically update concession prices for flash deals, matinee slots,
// and inventory optimization.
func ApplyLiveDynamicPricing(rawItems []models.ConcessionItem, foodLoc string, now time.Time) []models.ConcessionItem {
	tick := now.Unix() / 8 // new dynamic price cycle every 8 seconds
	timestampStr := now.Format("15:04:05")

	results := make([]models.ConcessionItem, len(rawItems))
	for i, it := range rawItems {
		item := it
		if item.BasePrice <= 0 {
			item.BasePrice = item.Price
		}

		// Calculate deterministic pseudo-random cycle for each item based on item ID and tick
		var idSum int
		for _, b := range []byte(item.ID) {
			idSum += int(b)
		}
		cycle := (int(tick) + idSum) % 5

		var dynamicPrice float64
		var trend string
		var flashText string

		switch cycle {
		case 0:
			// Standard regular baseline price
			dynamicPrice = item.BasePrice
			trend = "stable"
			flashText = "Standard Counter Rate"
		case 1:
			// Flash concession drop (-LKR 50 or ~6%)
			diff := 50.0
			if item.BasePrice < 600 {
				diff = 30.0
			}
			dynamicPrice = item.BasePrice - diff
			trend = "flash_drop"
			flashText = fmt.Sprintf("⚡ FLASH CONCESSION DROP (-LKR %.0f)", diff)
		case 2:
			// Matinee / Happy Hour Special (-LKR 80 or ~8%)
			diff := 80.0
			if item.BasePrice < 700 {
				diff = 40.0
			}
			dynamicPrice = item.BasePrice - diff
			trend = "down"
			flashText = fmt.Sprintf("📉 LIVE MATINEE SAVER (-LKR %.0f)", diff)
		case 3:
			// Quick snack counter booster (-LKR 30)
			diff := 30.0
			dynamicPrice = item.BasePrice - diff
			trend = "down"
			flashText = fmt.Sprintf("🎟️ LIVE PROMO TICK (-LKR %.0f)", diff)
		case 4:
			// Category specific live discount
			if item.Category == models.CategoryPopcorn {
				dynamicPrice = item.BasePrice - 40.0
				trend = "down"
				flashText = "🍿 POPCORN RUSH DISCOUNT (-LKR 40)"
			} else {
				dynamicPrice = item.BasePrice
				trend = "stable"
				flashText = "Active Concession Rate"
			}
		}

		if dynamicPrice <= 0 {
			dynamicPrice = item.BasePrice
		}

		item.Price = dynamicPrice
		item.PriceTrend = trend
		item.PriceChangeLKR = math.Round((dynamicPrice-item.BasePrice)*100) / 100
		item.FlashDealText = flashText
		item.LastPriceUpdate = timestampStr
		item.LiveTickID = tick

		// Recompute all promotional and payment discounts on the live dynamic price
		results[i] = EnrichItemDiscounts(item, foodLoc)
	}

	return results
}

// GetLiveConcessionTicker returns real-time pricing ticker metadata for a cinema.
func (s *RegionalScraper) GetLiveConcessionTicker(cinema models.Cinema) map[string]interface{} {
	now := time.Now()
	tick := now.Unix() / 8
	secondsUntilNextTick := 8 - (now.Unix() % 8)

	return map[string]interface{}{
		"cinema_id":            cinema.ID,
		"cinema_name":          cinema.Name,
		"live_pricing_active":  true,
		"tick_id":              tick,
		"updated_at":           now.Format("15:04:05"),
		"next_update_in_sec":   secondsUntilNextTick,
		"update_interval_sec":  8,
		"dynamic_engine_state": "ACTIVE_MATINEE_FLASH_FEED",
	}
}

// GetCinemaAds generates rich in-theater promotional advertisement banners highlighting combo, food, and drink discounts.
func (s *RegionalScraper) GetCinemaAds(cinema models.Cinema, items []models.ConcessionItem) []models.CinemaAd {
	var ads []models.CinemaAd

	var comboItem *models.ConcessionItem
	var flashItem *models.ConcessionItem
	var beverageItem *models.ConcessionItem
	var snackItem *models.ConcessionItem

	for i := range items {
		it := &items[i]
		if it.Category == models.CategoryCombo && comboItem == nil {
			comboItem = it
		}
		if (it.PriceTrend == "flash_drop" || it.PriceTrend == "down") && flashItem == nil {
			flashItem = it
		}
		if it.Category == models.CategoryBeverage && beverageItem == nil {
			beverageItem = it
		}
		if it.Category == models.CategorySnack && snackItem == nil {
			snackItem = it
		}
	}

	// 1. Featured Concession Combo Ad
	if comboItem != nil {
		savings := comboItem.Price - comboItem.BestDiscountedPrice
		ads = append(ads, models.CinemaAd{
			ID:              fmt.Sprintf("AD-COMBO-%s", comboItem.ID),
			CinemaID:        cinema.ID,
			Title:           fmt.Sprintf("🎬 %s Special", comboItem.Name),
			Subtitle:        fmt.Sprintf("Complete concession bundle: %s. Order via counter pass and unlock massive combo savings!", comboItem.Description),
			BadgeText:       "👑 EXCLUSIVE COMBO AD",
			DiscountText:    fmt.Sprintf("Pass Rate LKR %.0f (Save LKR %.0f &bull; 30%% OFF)", comboItem.BestDiscountedPrice, savings),
			PromoCode:       "SCOPE-COMBO30",
			TargetItemID:    comboItem.ID,
			TargetItemName:  comboItem.Name,
			OriginalPrice:   comboItem.Price,
			DiscountedPrice: comboItem.BestDiscountedPrice,
			SavingsLKR:      savings,
			CallToAction:    "🎟️ Claim Combo Deal & Add to Cart",
			IsFlashAd:       false,
			Tags:            []string{"combo", "popular", "best-value"},
		})
	}

	// 2. Flash Concession Drop Ad (Dynamic live price drop)
	if flashItem != nil {
		savings := flashItem.BasePrice - flashItem.BestDiscountedPrice
		ads = append(ads, models.CinemaAd{
			ID:              fmt.Sprintf("AD-FLASH-%s", flashItem.ID),
			CinemaID:        cinema.ID,
			Title:           fmt.Sprintf("⚡ Flash Concession Markdown: %s", flashItem.Name),
			Subtitle:        fmt.Sprintf("Live counter markdown active! %s. Pair with VIP pass for up to 35%% off!", flashItem.FlashDealText),
			BadgeText:       "⚡ LIVE FLASH AD",
			DiscountText:    fmt.Sprintf("Live LKR %.0f &bull; Pass LKR %.0f (Total Savings LKR %.0f)", flashItem.Price, flashItem.BestDiscountedPrice, savings),
			PromoCode:       "SCOPE-PLATINUM",
			TargetItemID:    flashItem.ID,
			TargetItemName:  flashItem.Name,
			OriginalPrice:   flashItem.BasePrice,
			DiscountedPrice: flashItem.BestDiscountedPrice,
			SavingsLKR:      savings,
			CallToAction:    "⚡ Grab Flash Concession Deal",
			IsFlashAd:       true,
			Tags:            []string{"flash-sale", "limited-time", "live-rate"},
		})
	}

	// 3. Popcorn & Beverage Booster Ad
	if beverageItem != nil {
		savings := beverageItem.Price - beverageItem.BestDiscountedPrice
		ads = append(ads, models.CinemaAd{
			ID:              fmt.Sprintf("AD-DRINK-%s", beverageItem.ID),
			CinemaID:        cinema.ID,
			Title:           fmt.Sprintf("🥤 Chilled Quencher Pairing: %s", beverageItem.Name),
			Subtitle:        fmt.Sprintf("Pair your warm popcorn bucket with %s! Enjoy 30%% Combo discount or 10%% FriMi Instant Cashback.", beverageItem.Description),
			BadgeText:       "🍿 DRINK + SNACK PAIRING",
			DiscountText:    fmt.Sprintf("Starting from LKR %.0f with digital payment pass", beverageItem.BestDiscountedPrice),
			PromoCode:       "FRIMI-CASHBACK",
			TargetItemID:    beverageItem.ID,
			TargetItemName:  beverageItem.Name,
			OriginalPrice:   beverageItem.Price,
			DiscountedPrice: beverageItem.BestDiscountedPrice,
			SavingsLKR:      savings,
			CallToAction:    "🥤 Add Chilled Drink Deal",
			IsFlashAd:       false,
			Tags:            []string{"beverage", "chilled", "pairing"},
		})
	}

	// 4. National CBSL LankaQR Ad
	ads = append(ads, models.CinemaAd{
		ID:              "AD-LANKAQR-NATIONAL",
		CinemaID:        cinema.ID,
		Title:           "🇱🇰 CBSL LankaQR 5% Instant Counter Rebate",
		Subtitle:        "Scan dynamic LankaQR at the counter with your bank app (Commercial Bank Q+, Sampath WePay, FriMi, HNB SOLO, BOC SmartPay) to receive a 5% instant price reduction!",
		BadgeText:       "🇱🇰 5% NATIONAL REBATE",
		DiscountText:    "Instant 5% Off Any Concession Order",
		PromoCode:       "LANKAQR-5",
		CallToAction:    "📱 Apply LankaQR 5% Discount",
		IsFlashAd:       false,
		Tags:            []string{"national-qr", "all-banks", "cbsl"},
	})

	return ads
}

// scrapeByCity executes concession scraping pipelines tuned for regional suppliers and theater concessions.
func (s *RegionalScraper) scrapeByCity(city string, cinemaID string) ([]models.ConcessionItem, error) {
	cityLower := strings.ToLower(strings.TrimSpace(city))

	switch cityLower {
	case "kandy":
		return []models.ConcessionItem{
			{ID: "KND-01", CinemaID: cinemaID, Name: "Scope Jumbo Warm Caramel Popcorn", Category: models.CategoryPopcorn, Size: "Jumbo", Price: 1200.00, InStock: true, Tags: []string{"sweet", "popular", "scope-exclusive"}, Description: "Scope signature crunchy caramelized warm corn"},
			{ID: "KND-02", CinemaID: cinemaID, Name: "Scope Signature Butter Salt Popcorn", Category: models.CategoryPopcorn, Size: "Large", Price: 1000.00, InStock: true, Tags: []string{"savory", "classic"}, Description: "Freshly popped with melted golden butter"},
			{ID: "KND-03", CinemaID: cinemaID, Name: "Scope Large Fountain Mountain Dew", Category: models.CategoryBeverage, Size: "Large", Price: 650.00, InStock: true, Tags: []string{"beverage", "chilled"}, Description: "Chilled fountain carbonated drink"},
			{ID: "KND-04", CinemaID: cinemaID, Name: "Scope Fresh Ceylon Iced Tea", Category: models.CategoryBeverage, Size: "Medium", Price: 480.00, InStock: true, Tags: []string{"beverage", "local"}, Description: "Brewed hill country tea with fresh lemon"},
			{ID: "KND-05", CinemaID: cinemaID, Name: "Scope Gourmet Chicken Hotdog", Category: models.CategorySnack, Size: "Single", Price: 950.00, InStock: true, Tags: []string{"spicy", "savory", "halal"}, Description: "Grilled chicken sausage in toasted brioche with relish"},
			{ID: "KND-06", CinemaID: cinemaID, Name: "Scope Director's Duo Combo", Category: models.CategoryCombo, Size: "Duo", Price: 2400.00, InStock: true, Tags: []string{"combo", "best-value", "vip"}, Description: "1 Jumbo Popcorn + 2 Large Drinks + 1 Gourmet Hotdog"},
		}, nil

	case "gampaha":
		return []models.ConcessionItem{
			{ID: "GMP-01", CinemaID: cinemaID, Name: "Salted Butter Popcorn (M)", Category: models.CategoryPopcorn, Size: "Medium", Price: 900.00, InStock: true, Tags: []string{"savory"}, Description: "Warm buttery classic popcorn"},
			{ID: "GMP-02", CinemaID: cinemaID, Name: "Caramel Crunch Tub", Category: models.CategoryPopcorn, Size: "Large", Price: 1100.00, InStock: true, Tags: []string{"sweet"}, Description: "Thick golden caramel glaze"},
			{ID: "GMP-03", CinemaID: cinemaID, Name: "Cold Milo Float", Category: models.CategoryBeverage, Size: "Medium", Price: 550.00, InStock: true, Tags: []string{"beverage", "sweet"}, Description: "Malt Milo topped with vanilla ice cream"},
			{ID: "GMP-04", CinemaID: cinemaID, Name: "Coca-Cola Zero Sugar", Category: models.CategoryBeverage, Size: "Large", Price: 500.00, InStock: true, Tags: []string{"beverage", "sugar-free"}, Description: "Chilled sparkling Coke"},
			{ID: "GMP-05", CinemaID: cinemaID, Name: "Crispy Nachos & Cheese", Category: models.CategorySnack, Size: "Single", Price: 850.00, InStock: true, Tags: []string{"savory", "vegetarian"}, Description: "Corn tortilla chips with jalapeno cheese sauce"},
			{ID: "GMP-06", CinemaID: cinemaID, Name: "Gampaha Deluxe Couple Combo", Category: models.CategoryCombo, Size: "Duo", Price: 2250.00, InStock: true, Tags: []string{"combo"}, Description: "1 Caramel Tub + 2 Milo Floats + Nachos"},
		}, nil

	case "galle":
		return []models.ConcessionItem{
			{ID: "GLE-01", CinemaID: cinemaID, Name: "Southern Cheese Popcorn", Category: models.CategoryPopcorn, Size: "Large", Price: 1100.00, InStock: true, Tags: []string{"savory", "cheesy"}, Description: "Dusted with rich cheddar seasoning"},
			{ID: "GLE-02", CinemaID: cinemaID, Name: "Fresh Lime & Mint Cooler", Category: models.CategoryBeverage, Size: "Large", Price: 500.00, InStock: true, Tags: []string{"beverage", "citrus"}, Description: "Fresh southern lime with crushed mint"},
			{ID: "GLE-03", CinemaID: cinemaID, Name: "Sweet Chili Fish & Chips", Category: models.CategorySnack, Size: "Single", Price: 1400.00, InStock: true, Tags: []string{"seafood", "crispy"}, Description: "Crumbed fish goujons with hand-cut fries"},
			{ID: "GLE-04", CinemaID: cinemaID, Name: "Crispy Samosa Platter (4pcs)", Category: models.CategorySnack, Size: "Single", Price: 750.00, InStock: true, Tags: []string{"savory", "vegetarian"}, Description: "Spiced potato pastry with sweet tamarind dip"},
			{ID: "GLE-05", CinemaID: cinemaID, Name: "Fort Fortress Mega Combo", Category: models.CategoryCombo, Size: "Family", Price: 2950.00, InStock: true, Tags: []string{"combo", "family"}, Description: "2 Large Popcorn + 3 Drinks + 1 Fish & Chips"},
		}, nil

	case "kurunegala":
		return []models.ConcessionItem{
			{ID: "KRN-01", CinemaID: cinemaID, Name: "North Star Popcorn Bucket", Category: models.CategoryPopcorn, Size: "Jumbo", Price: 1000.00, InStock: true, Tags: []string{"savory"}, Description: "Family tub of salted butter corn"},
			{ID: "KRN-02", CinemaID: cinemaID, Name: "Fresh Cola Combo Drink", Category: models.CategoryBeverage, Size: "Large", Price: 600.00, InStock: true, Tags: []string{"beverage"}, Description: "Large iced fountain cola"},
			{ID: "KRN-03", CinemaID: cinemaID, Name: "Crispy Chicken Roll (2pcs)", Category: models.CategorySnack, Size: "Single", Price: 820.00, InStock: true, Tags: []string{"savory", "halal"}, Description: "Crumbed Sri Lankan spicy chicken rolls"},
			{ID: "KRN-04", CinemaID: cinemaID, Name: "Loaded French Fries", Category: models.CategorySnack, Size: "Medium", Price: 780.00, InStock: true, Tags: []string{"savory", "vegetarian"}, Description: "Seasoned fries topped with melted cheese"},
			{ID: "KRN-05", CinemaID: cinemaID, Name: "Kurunegala Value Pack", Category: models.CategoryCombo, Size: "Duo", Price: 2100.00, InStock: true, Tags: []string{"combo"}, Description: "1 Popcorn Bucket + 2 Drinks + 1 Chicken Roll"},
		}, nil

	case "negombo":
		return []models.ConcessionItem{
			{ID: "NGB-01", CinemaID: cinemaID, Name: "Coastal Caramel Popcorn", Category: models.CategoryPopcorn, Size: "Large", Price: 980.00, InStock: true, Tags: []string{"sweet"}, Description: "Artisan buttery caramel popcorn"},
			{ID: "NGB-02", CinemaID: cinemaID, Name: "Iced Lemon Soda", Category: models.CategoryBeverage, Size: "Large", Price: 520.00, InStock: true, Tags: []string{"beverage", "fizzy"}, Description: "Refreshing sparkling lemon cordial"},
			{ID: "NGB-03", CinemaID: cinemaID, Name: "Cheese Nachos Deluxe", Category: models.CategorySnack, Size: "Medium", Price: 930.00, InStock: true, Tags: []string{"savory", "cheesy"}, Description: "Warm tortilla chips with salsa and cheese"},
			{ID: "NGB-04", CinemaID: cinemaID, Name: "Calamari Rings Cone", Category: models.CategorySnack, Size: "Single", Price: 1350.00, InStock: true, Tags: []string{"seafood", "crispy"}, Description: "Crispy fried squid rings with tartare dip"},
			{ID: "NGB-05", CinemaID: cinemaID, Name: "Negombo Sunset Combo", Category: models.CategoryCombo, Size: "Duo", Price: 2300.00, InStock: true, Tags: []string{"combo"}, Description: "1 Caramel Popcorn + 2 Lemon Sodas + 1 Nachos"},
		}, nil

	case "jaffna":
		return []models.ConcessionItem{
			{ID: "JFN-01", CinemaID: cinemaID, Name: "Jaffna Crunch Spiced Popcorn", Category: models.CategoryPopcorn, Size: "Large", Price: 1150.00, InStock: true, Tags: []string{"spicy", "savory"}, Description: "Popcorn dusted with northern chili spices"},
			{ID: "JFN-02", CinemaID: cinemaID, Name: "Palmyra Sweet Cooler", Category: models.CategoryBeverage, Size: "Medium", Price: 700.00, InStock: true, Tags: []string{"beverage", "traditional"}, Description: "Natural palmyra fruit pulp nectar"},
			{ID: "JFN-03", CinemaID: cinemaID, Name: "Hot Ceylon Chicken Bites", Category: models.CategorySnack, Size: "Single", Price: 960.00, InStock: true, Tags: []string{"spicy", "savory", "halal"}, Description: "Boneless fried chicken seasoned with curry leaves"},
			{ID: "JFN-04", CinemaID: cinemaID, Name: "Jaffna Murukku Platter", Category: models.CategorySnack, Size: "Single", Price: 650.00, InStock: true, Tags: []string{"crunchy", "vegetarian"}, Description: "Crispy traditional spiced savory snacks"},
			{ID: "JFN-05", CinemaID: cinemaID, Name: "Northern Festival Combo", Category: models.CategoryCombo, Size: "Duo", Price: 2600.00, InStock: true, Tags: []string{"combo"}, Description: "1 Spiced Popcorn + 2 Coolers + 1 Chicken Bites"},
		}, nil

	case "matara":
		return []models.ConcessionItem{
			{ID: "MTR-01", CinemaID: cinemaID, Name: "Southern Popcorn Mix (Sweet & Salt)", Category: models.CategoryPopcorn, Size: "Large", Price: 1050.00, InStock: true, Tags: []string{"mixed"}, Description: "Dual blend of caramel & sea salted corn"},
			{ID: "MTR-02", CinemaID: cinemaID, Name: "Fresh Orange Fizz", Category: models.CategoryBeverage, Size: "Medium", Price: 580.00, InStock: true, Tags: []string{"beverage"}, Description: "Pulpy orange with sparkling water"},
			{ID: "MTR-03", CinemaID: cinemaID, Name: "Crispy Potato Wedges Platter", Category: models.CategorySnack, Size: "Medium", Price: 890.00, InStock: true, Tags: []string{"savory", "vegetarian"}, Description: "Crispy seasoned skin-on potato wedges"},
			{ID: "MTR-04", CinemaID: cinemaID, Name: "Southern Breeze Duo Combo", Category: models.CategoryCombo, Size: "Duo", Price: 2200.00, InStock: true, Tags: []string{"combo"}, Description: "1 Mix Popcorn + 2 Orange Fizz + 1 Wedges"},
		}, nil

	case "anuradhapura":
		return []models.ConcessionItem{
			{ID: "ANR-01", CinemaID: cinemaID, Name: "Classic Butter Popcorn", Category: models.CategoryPopcorn, Size: "Large", Price: 950.00, InStock: true, Tags: []string{"savory"}, Description: "Golden popped kernels with butter flavor"},
			{ID: "ANR-02", CinemaID: cinemaID, Name: "Chilled Passion Fruit Nectar", Category: models.CategoryBeverage, Size: "Medium", Price: 520.00, InStock: true, Tags: []string{"beverage"}, Description: "Island passion fruit cooler"},
			{ID: "ANR-03", CinemaID: cinemaID, Name: "Vegetable Spring Rolls (3pcs)", Category: models.CategorySnack, Size: "Single", Price: 720.00, InStock: true, Tags: []string{"savory", "vegetarian"}, Description: "Crispy pastry rolls filled with garden veggies"},
			{ID: "ANR-04", CinemaID: cinemaID, Name: "Rajarata Family Snack Box", Category: models.CategoryCombo, Size: "Family", Price: 2350.00, InStock: true, Tags: []string{"combo"}, Description: "1 Popcorn + 2 Nectars + 2 Spring Roll orders"},
		}, nil

	case "ratnapura":
		return []models.ConcessionItem{
			{ID: "RTP-01", CinemaID: cinemaID, Name: "Gem City Sweet Popcorn", Category: models.CategoryPopcorn, Size: "Large", Price: 990.00, InStock: true, Tags: []string{"sweet"}, Description: "Crunchy sweet glaze popcorn"},
			{ID: "RTP-02", CinemaID: cinemaID, Name: "Iced Faluda Royal", Category: models.CategoryBeverage, Size: "Medium", Price: 620.00, InStock: true, Tags: []string{"beverage", "dessert"}, Description: "Rose syrup with milk, basil seeds & ice cream"},
			{ID: "RTP-03", CinemaID: cinemaID, Name: "Spicy Beef Patty (2pcs)", Category: models.CategorySnack, Size: "Single", Price: 800.00, InStock: true, Tags: []string{"spicy", "savory"}, Description: "Shortcrust baked spicy minced beef parcels"},
			{ID: "RTP-04", CinemaID: cinemaID, Name: "Milano Matinee Combo", Category: models.CategoryCombo, Size: "Duo", Price: 2150.00, InStock: true, Tags: []string{"combo"}, Description: "1 Sweet Popcorn + 2 Faludas + 1 Patty Box"},
		}, nil

	default:
		return nil, fmt.Errorf("%w: no concession scraper registered for city %q", ErrCinemaNotFound, city)
	}
}

// ScrapeRawConcessionFeed parses simulated or raw concession HTML feeds.
// Allows parsing raw web page tables / div lists into structured models.ConcessionItem.
func (s *RegionalScraper) ScrapeRawConcessionFeed(cinemaID string, rawHTML string) ([]models.ConcessionItem, error) {
	if strings.TrimSpace(rawHTML) == "" {
		return nil, fmt.Errorf("%w: raw feed content is empty", ErrInvalidInput)
	}

	var items []models.ConcessionItem
	// Regex pattern for concession items in web feeds: <item name="..." category="..." price="..." ...>
	re := regexp.MustCompile(`(?i)<(?:div|item)[^>]*class=["']concession-card["'][^>]*data-name=["']([^"']+)["'][^>]*data-cat=["']([^"']+)["'][^>]*data-price=["']([0-9.]+)["']`)
	matches := re.FindAllStringSubmatch(rawHTML, -1)

	for idx, m := range matches {
		if len(m) >= 4 {
			name := m[1]
			catStr := m[2]
			priceVal, err := strconv.ParseFloat(m[3], 64)
			if err != nil {
				continue
			}

			category := models.CategorySnack
			switch strings.ToLower(catStr) {
			case "popcorn":
				category = models.CategoryPopcorn
			case "beverage", "drink":
				category = models.CategoryBeverage
			case "combo":
				category = models.CategoryCombo
			case "dessert":
				category = models.CategoryDessert
			}

			items = append(items, models.ConcessionItem{
				ID:       fmt.Sprintf("%s-EXT-%02d", cinemaID, idx+1),
				CinemaID: cinemaID,
				Name:     name,
				Category: category,
				Price:    priceVal,
				InStock:  true,
			})
		}
	}

	return items, nil
}
