package scraper

import (
	"fmt"
	"strings"

	"cinema-recommender/internal/models"
)

// defaultCinemaRegistry returns the pre-configured outstation cinema theaters across Sri Lanka.
func defaultCinemaRegistry() []models.Cinema {
	return []models.Cinema{
		{
			ID:                  "CMB-CCC",
			Name:                "Scope Cinemas Multiplex - Colombo City Centre",
			Chain:               "Scope Cinemas",
			City:                "Colombo",
			Province:            "Western Province",
			Address:             "Level 3, Colombo City Centre, 137 Sir James Pieris Mawatha, Colombo 02",
			Screens:             6,
			IsOutside:           false,
			Latitude:            6.9175,
			Longitude:           79.8550,
			MapURL:              "https://maps.google.com/?q=6.9175,79.8550",
			HasInHouseFood:      true,
			HasFoodCourtInFront: true,
			FoodCourtName:       "The Food Studio @ CCC & Scope Candy Bar",
			FoodCourtLocation:   "Level 3 Lobby directly in front of cinema entrance",
			FoodPlaceName:       "The Food Studio @ CCC",
			FoodPlaceType:       "In-Front Luxury Mall Food Court & Scope Candy Bar",
			NearbyFoodOptions:   []string{"The Food Studio (Level 3)", "Cargills Food City CCC", "Twist Colombo", "Shiok Singapore Street Food", "McDonald's CCC"},
			FoodHours:           "10:00 AM - 11:30 PM (Open during all movie screenings)",
			FoodDeliveryToSeat:  true,
		},
		{
			ID:                  "CMB-HCM",
			Name:                "Scope Cinemas Multiplex - Havelock City Mall",
			Chain:               "Scope Cinemas",
			City:                "Colombo",
			Province:            "Western Province",
			Address:             "Level 4, Havelock City Mall, 324 Havelock Road, Colombo 05",
			Screens:             6,
			IsOutside:           false,
			Latitude:            6.8778,
			Longitude:           79.8653,
			MapURL:              "https://maps.google.com/?q=6.8778,79.8653",
			HasInHouseFood:      true,
			HasFoodCourtInFront: true,
			FoodCourtName:       "Havelock Food Lounge & Scope IMAX Candy Bar",
			FoodCourtLocation:   "Level 4 Lobby directly in front of IMAX auditorium entrance",
			FoodPlaceName:       "Havelock Food Lounge & Scope IMAX Candy Bar",
			FoodPlaceType:       "In-Front Mall Food Lounge & IMAX Concession Stand",
			NearbyFoodOptions:   []string{"Havelock Food Lounge (Level 4)", "Baskin Robbins HCM", "Subway Havelock City", "Cargills Gourmet Food Hall"},
			FoodHours:           "10:00 AM - 11:30 PM (Open during all movie screenings)",
			FoodDeliveryToSeat:  true,
		},
		{
			ID:                  "CMB-LBT",
			Name:                "Liberty by Scope Cinemas",
			Chain:               "Scope Cinemas",
			City:                "Colombo",
			Province:            "Western Province",
			Address:             "35 R. A. De Mel Mawatha, Colombo 03",
			Screens:             2,
			IsOutside:           false,
			Latitude:            6.9073,
			Longitude:           79.8517,
			MapURL:              "https://maps.google.com/?q=6.9073,79.8517",
			HasInHouseFood:      true,
			HasFoodCourtInFront: true,
			FoodCourtName:       "Liberty Foyer Concession Court & Scope Candy Bar",
			FoodCourtLocation:   "Front main entrance lobby directly before auditorium doors",
			FoodPlaceName:       "Liberty Foyer Concession Court",
			FoodPlaceType:       "In-Front Cinema Concession Food Court",
			NearbyFoodOptions:   []string{"Liberty Plaza Food Court", "Perera & Sons Kollupitiya", "Chit-Chat Coffee Corner", "Mackenzy's"},
			FoodHours:           "10:00 AM - 10:45 PM (Open during all movie screenings)",
			FoodDeliveryToSeat:  true,
		},
		{
			ID:                  "KND-KCC",
			Name:                "Scope Partner Multiplex - KCC Kandy",
			Chain:               "Scope Cinemas Partner Circuit",
			City:                "Kandy",
			Province:            "Central Province",
			Address:             "Level 3, Kandy City Centre, Dalada Veediya, Kandy",
			Screens:             4,
			IsOutside:           true,
			Latitude:            7.2936,
			Longitude:           80.6385,
			MapURL:              "https://maps.google.com/?q=7.2936,80.6385",
			HasInHouseFood:      true,
			HasFoodCourtInFront: true,
			FoodCourtName:       "KCC World Food Court & Scope Candy Bar",
			FoodCourtLocation:   "Level 3 Lobby directly in front of cinema entrance",
			FoodPlaceName:       "KCC World Food Court & Scope Candy Bar",
			FoodPlaceType:       "In-Front Mall Food Court & Cinema Concession Stand",
			NearbyFoodOptions:   []string{"KCC World Food Court (Level 4)", "Devon Restaurant & Bakery (2 min walk)", "Bake House Dalada Veediya", "Cargills Food Hall"},
			FoodHours:           "10:00 AM - 10:45 PM (Open during all movie screenings)",
			FoodDeliveryToSeat:  true,
		},
		{
			ID:                  "KND-REG",
			Name:                "Regal Cinema Kandy",
			Chain:               "Ceylon Theatres",
			City:                "Kandy",
			Province:            "Central Province",
			Address:             "No. 12, Katugastota Road, Kandy",
			Screens:             2,
			IsOutside:           true,
			Latitude:            7.2985,
			Longitude:           80.6335,
			MapURL:              "https://maps.google.com/?q=7.2985,80.6335",
			HasInHouseFood:      true,
			HasFoodCourtInFront: true,
			FoodCourtName:       "Regal Foyer Concession Food Court",
			FoodCourtLocation:   "Directly in front of cinema auditorium entrance",
			FoodPlaceName:       "Regal Foyer Concession Food Court",
			FoodPlaceType:       "In-Front Theater Snack Food Court",
			NearbyFoodOptions:   []string{"Katugastota Cafe & Bake House", "Perera & Sons Kandy", "White House Restaurant"},
			FoodHours:           "10:15 AM - 10:30 PM",
			FoodDeliveryToSeat:  false,
		},
		{
			ID:                  "GMP-REG",
			Name:                "Regal Cinema Gampaha",
			Chain:               "Ceylon Theatres",
			City:                "Gampaha",
			Province:            "Western Province (Outer Regional)",
			Address:             "Bauddhaloka Mawatha, Gampaha",
			Screens:             2,
			IsOutside:           true,
			Latitude:            7.0917,
			Longitude:           79.9998,
			MapURL:              "https://maps.google.com/?q=7.0917,79.9998",
			HasInHouseFood:      true,
			HasFoodCourtInFront: true,
			FoodCourtName:       "Regal Front Foyer Food Plaza & Snack Court",
			FoodCourtLocation:   "Front entrance foyer directly facing ticket counter",
			FoodPlaceName:       "Regal Front Foyer Food Plaza",
			FoodPlaceType:       "In-Front Cinema Concession Food Court",
			NearbyFoodOptions:   []string{"Gampaha Food City Concessions", "Dinemore Gampaha", "Perera & Sons Bauddhaloka Mw"},
			FoodHours:           "10:00 AM - 10:30 PM",
			FoodDeliveryToSeat:  false,
		},
		{
			ID:                  "GLE-QNS",
			Name:                "Queens Cinema Galle",
			Chain:               "Independent Heritage Screen",
			City:                "Galle",
			Province:            "Southern Province",
			Address:             "Wakwella Road, Galle",
			Screens:             2,
			IsOutside:           true,
			Latitude:            6.0535,
			Longitude:           80.2210,
			MapURL:              "https://maps.google.com/?q=6.0535,80.2210",
			HasInHouseFood:      true,
			HasFoodCourtInFront: true,
			FoodCourtName:       "Queens Front Concession Court & Refreshment Bar",
			FoodCourtLocation:   "Front lobby terrace directly before screen entrance",
			FoodPlaceName:       "Queens Front Concession Court",
			FoodPlaceType:       "In-Front Heritage Lobby Food Court",
			NearbyFoodOptions:   []string{"Galle Fort Dining Terrace (5 min drive)", "Pedlar's Inn Concessions", "Southern Taste Bakery Wakwella"},
			FoodHours:           "10:30 AM - 10:15 PM",
			FoodDeliveryToSeat:  false,
		},
		{
			ID:                  "KRN-MLD",
			Name:                "Maldini Cinema Kurunegala",
			Chain:               "Maldini Circuits",
			City:                "Kurunegala",
			Province:            "North Western Province",
			Address:             "Colombo Road, Kurunegala",
			Screens:             2,
			IsOutside:           true,
			Latitude:            7.4863,
			Longitude:           80.3623,
			MapURL:              "https://maps.google.com/?q=7.4863,80.3623",
			HasInHouseFood:      true,
			HasFoodCourtInFront: true,
			FoodCourtName:       "Maldini Concession Food Court & Snack Foyer",
			FoodCourtLocation:   "Cinema lobby entrance directly in front of screens",
			FoodPlaceName:       "Maldini Concession Food Court",
			FoodPlaceType:       "In-Front Theater Snack Food Court",
			NearbyFoodOptions:   []string{"Kurunegala Clock Tower Cafe Row", "Perera & Sons Colombo Rd", "Royal Bakery Kurunegala"},
			FoodHours:           "10:00 AM - 10:30 PM",
			FoodDeliveryToSeat:  false,
		},
		{
			ID:                  "NGB-RDO",
			Name:                "Rodeo Cinema Negombo",
			Chain:               "Western Coast Network",
			City:                "Negombo",
			Province:            "Western Province (Outer Coastal)",
			Address:             "Main Street, Negombo",
			Screens:             2,
			IsOutside:           true,
			Latitude:            7.2083,
			Longitude:           79.8358,
			MapURL:              "https://maps.google.com/?q=7.2083,79.8358",
			HasInHouseFood:      true,
			HasFoodCourtInFront: true,
			FoodCourtName:       "Rodeo Beach Concession Court & Chill Bar",
			FoodCourtLocation:   "Front theater foyer facing screen entrance",
			FoodPlaceName:       "Rodeo Beach Concession Court",
			FoodPlaceType:       "In-Front Concession Food Court",
			NearbyFoodOptions:   []string{"Lords Restaurant Negombo Beach", "Negombo Dutch Bakery", "King Coconut Stalls Main St"},
			FoodHours:           "10:30 AM - 11:00 PM",
			FoodDeliveryToSeat:  false,
		},
		{
			ID:                  "JFN-RJA",
			Name:                "Rajah Cinema Jaffna",
			Chain:               "Northern Cinema Network",
			City:                "Jaffna",
			Province:            "Northern Province",
			Address:             "Hospital Road, Jaffna",
			Screens:             3,
			IsOutside:           true,
			Latitude:            9.6615,
			Longitude:           80.0255,
			MapURL:              "https://maps.google.com/?q=9.6615,80.0255",
			HasInHouseFood:      true,
			HasFoodCourtInFront: true,
			FoodCourtName:       "Rajah Concession Food Court & Refreshment Parlour",
			FoodCourtLocation:   "Front main entrance foyer directly in front of cinema",
			FoodPlaceName:       "Rajah Concession Food Court",
			FoodPlaceType:       "In-Front Northern Concession Food Court",
			NearbyFoodOptions:   []string{"Mangos Vegetarian Restaurant", "Rolex Cafe Jaffna", "Rio Ice Cream Jaffna"},
			FoodHours:           "10:00 AM - 10:45 PM",
			FoodDeliveryToSeat:  false,
		},
		{
			ID:                  "MTR-SKC",
			Name:                "SK Cinema Matara",
			Chain:               "Southern Screen Network",
			City:                "Matara",
			Province:            "Southern Province",
			Address:             "Anagarika Dharmapala Mawatha, Matara",
			Screens:             2,
			IsOutside:           true,
			Latitude:            5.9485,
			Longitude:           80.5488,
			MapURL:              "https://maps.google.com/?q=5.9485,80.5488",
			HasInHouseFood:      true,
			HasFoodCourtInFront: true,
			FoodCourtName:       "SK Concession Court & Chill Bar",
			FoodCourtLocation:   "Main entrance lobby directly in front of cinema screens",
			FoodPlaceName:       "SK Concession Food Court",
			FoodPlaceType:       "In-Front Concession Food Court",
			NearbyFoodOptions:   []string{"Matara Beach Park Food Stalls", "Dutch Fort Bakery Matara", "Keells Super Food Corner"},
			FoodHours:           "10:00 AM - 10:30 PM",
			FoodDeliveryToSeat:  false,
		},
		{
			ID:                  "ANR-CMX",
			Name:                "Cinemax Anuradhapura",
			Chain:               "Rajarata Theatres",
			City:                "Anuradhapura",
			Province:            "North Central Province",
			Address:             "Maithripala Senanayake Mawatha, Anuradhapura",
			Screens:             2,
			IsOutside:           true,
			Latitude:            8.3350,
			Longitude:           80.4108,
			MapURL:              "https://maps.google.com/?q=8.3350,80.4108",
			HasInHouseFood:      true,
			HasFoodCourtInFront: true,
			FoodCourtName:       "Cinemax Concession Food Court & Snack Corner",
			FoodCourtLocation:   "Front lobby directly before auditorium doors",
			FoodPlaceName:       "Cinemax Concession Food Court",
			FoodPlaceType:       "In-Front Concession Food Court",
			NearbyFoodOptions:   []string{"Ceylan Bake House Anuradhapura", "Rajarata Food Center", "Hotel Shalini Dining"},
			FoodHours:           "10:00 AM - 10:15 PM",
			FoodDeliveryToSeat:  false,
		},
		{
			ID:                  "RTP-MLN",
			Name:                "Milano Cinema Ratnapura",
			Chain:               "Milano Circuit",
			City:                "Ratnapura",
			Province:            "Sabaragamuwa Province",
			Address:             "Main Street, Ratnapura",
			Screens:             2,
			IsOutside:           true,
			Latitude:            6.6828,
			Longitude:           80.4037,
			MapURL:              "https://maps.google.com/?q=6.6828,80.4037",
			HasInHouseFood:      true,
			HasFoodCourtInFront: true,
			FoodCourtName:       "Milano Front Concession Court & Snack Bar",
			FoodCourtLocation:   "Main lobby directly in front of cinema screens",
			FoodPlaceName:       "Milano Front Concession Court",
			FoodPlaceType:       "In-Front Concession Food Court",
			NearbyFoodOptions:   []string{"Ratnapura City Bakers", "Gem City Food Plaza", "Perera & Sons Main Street"},
			FoodHours:           "10:15 AM - 10:15 PM",
			FoodDeliveryToSeat:  false,
		},
	}
}

// GetRegisteredCinemas returns all registered cinemas that have a verified food court in front for cinema food and drinks.
func (s *RegionalScraper) GetRegisteredCinemas() []models.Cinema {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var cinemas []models.Cinema
	for _, c := range s.cinemaRegistry {
		// Strict filter: cinema must have an in-front food court for cinema food and drinks only
		if c.HasFoodCourtInFront && c.HasInHouseFood {
			cinemas = append(cinemas, c)
		}
	}
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
