package scraper

import (
	"math"
	"sort"
	"strings"

	"cinema-recommender/internal/models"
)

// CalculateHaversineDistance computes the geographical distance in kilometers between two GPS coordinates
// using the spherical law of cosines / Haversine formula.
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
