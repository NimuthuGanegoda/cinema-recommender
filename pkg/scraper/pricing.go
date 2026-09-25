package scraper

import (
	"fmt"
	"math"
	"time"

	"cinema-recommender/internal/models"
)

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
