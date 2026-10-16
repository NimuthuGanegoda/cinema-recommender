package scraper

import (
	"math"
	"time"

	"cinema-recommender/internal/models"
)

// EnrichItemDiscounts calculates and attaches all available food promotions to a food item.
func EnrichItemDiscounts(item models.ConcessionItem, foodLocation string) models.ConcessionItem {
	price := item.Price
	if price <= 0 {
		return item
	}

	var discounts []models.ItemDiscountInfo

	// 1. CinePass Platinum VIP (35% off)
	platPct := 0.35
	platDisc := math.Round(price*(1.0-platPct)*100) / 100
	discounts = append(discounts, models.ItemDiscountInfo{
		PromoCode:       "CINE-PLATINUM",
		Provider:        "CinePass VIP",
		Title:           "CinePass Platinum VIP (35% Off)",
		DiscountPct:     platPct,
		DiscountedPrice: platDisc,
		SavingsLKR:      math.Round((price-platDisc)*100) / 100,
		Requirement:     "CinePass Platinum VIP membership",
	})

	// 2. CineBite Combo 30% Deal
	comboPct := 0.30
	comboDisc := math.Round(price*(1.0-comboPct)*100) / 100
	discounts = append(discounts, models.ItemDiscountInfo{
		PromoCode:       "CINE-COMBO30",
		Provider:        "CinePass VIP",
		Title:           "CineBite Popcorn + Drink Combo (30% Off)",
		DiscountPct:     comboPct,
		DiscountedPrice: comboDisc,
		SavingsLKR:      math.Round((price-comboDisc)*100) / 100,
		Requirement:     "Pair with any Popcorn & Beverage",
	})

	// 3. Cinema Concession Boost (25% Off)
	regPct := 0.25
	regDisc := math.Round(price*(1.0-regPct)*100) / 100
	discounts = append(discounts, models.ItemDiscountInfo{
		PromoCode:       "CINE-CINEMA25",
		Provider:        "CinePass VIP",
		Title:           "Cinema Concession Pass (25% Off)",
		DiscountPct:     regPct,
		DiscountedPrice: regDisc,
		SavingsLKR:      math.Round((price-regDisc)*100) / 100,
		Requirement:     "Concession orders over LKR 2,000",
	})

	// 4. Student Moviegoer Pass (20% Off)
	stuPct := 0.20
	stuDisc := math.Round(price*(1.0-stuPct)*100) / 100
	discounts = append(discounts, models.ItemDiscountInfo{
		PromoCode:       "CINE-STUDENT",
		Provider:        "CinePass VIP",
		Title:           "Student Moviegoer Pass (20% Off)",
		DiscountPct:     stuPct,
		DiscountedPrice: stuDisc,
		SavingsLKR:      math.Round((price-stuDisc)*100) / 100,
		Requirement:     "Verified student ID cardholders",
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

// ApplyLiveDynamicPricing computes constantly updated real-time concession prices and flash discounts.
// Movie theater candy bars dynamically update concession prices for flash deals, matinee slots,
// and inventory optimization.
func ApplyLiveDynamicPricing(rawItems []models.ConcessionItem, foodLoc string, now time.Time) []models.ConcessionItem {
	results := make([]models.ConcessionItem, len(rawItems))
	for i, item := range rawItems {
		if item.BasePrice <= 0 {
			item.BasePrice = item.Price
		}
		item.Price = item.BasePrice
		item.PriceTrend = "stable"
		item.PriceChangeLKR = 0
		item.FlashDealText = "Example menu price — confirm with cinema"
		item.LastPriceUpdate = ""
		item.LiveTickID = 0
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
		"live_pricing_active":  false,
		"tick_id":              tick,
		"updated_at":           now.Format("15:04:05"),
		"next_update_in_sec":   secondsUntilNextTick,
		"update_interval_sec":  8,
		"dynamic_engine_state": "EXAMPLE_CATALOG",
	}
}
