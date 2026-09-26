package models

// Cinema represents a movie theater entity with location and food/dining metadata.
type Cinema struct {
	ID                 string   `json:"id"`
	Name               string   `json:"name"`
	Chain              string   `json:"chain,omitempty"`
	City               string   `json:"city"`
	Province           string   `json:"province,omitempty"`
	Address            string   `json:"address,omitempty"`
	Screens            int      `json:"screens,omitempty"`
	IsOutside          bool     `json:"is_outside_colombo"`
	Latitude           float64  `json:"latitude,omitempty"`
	Longitude          float64  `json:"longitude,omitempty"`
	MapURL             string   `json:"map_url,omitempty"`
	HasInHouseFood      bool     `json:"has_in_house_food"`
	HasFoodCourtInFront bool     `json:"has_food_court_in_front"`
	FoodCourtName       string   `json:"food_court_name,omitempty"`
	FoodCourtLocation   string   `json:"food_court_location,omitempty"`
	FoodPlaceName       string   `json:"food_place_name,omitempty"`
	FoodPlaceType       string   `json:"food_place_type,omitempty"`
	NearbyFoodOptions   []string `json:"nearby_food_options,omitempty"`
	FoodHours           string   `json:"food_hours,omitempty"`
	FoodDeliveryToSeat  bool     `json:"food_delivery_to_seat"`
	DistanceKM          float64  `json:"distance_km,omitempty"` // populated dynamically if user coordinates provided
}

// GeoCoord represents GPS coordinates for cinema location services and Haversine distance calculations.
type GeoCoord struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
