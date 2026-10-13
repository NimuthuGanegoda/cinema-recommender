package scraper

import (
	"fmt"

	"cinema-recommender/internal/models"
)

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
			Subtitle:        fmt.Sprintf("Complete concession bundle: %s. Compare the snacks in this food bundle.", comboItem.Description),
			BadgeText:       "🎬 CINEBITE COMBO SPECIAL",
			DiscountText:    fmt.Sprintf("Pass Rate LKR %.0f (Save LKR %.0f &bull; 30%% OFF)", comboItem.BestDiscountedPrice, savings),
			PromoCode:       "CINE-COMBO30",
			TargetItemID:    comboItem.ID,
			TargetItemName:  comboItem.Name,
			OriginalPrice:   comboItem.Price,
			DiscountedPrice: comboItem.BestDiscountedPrice,
			SavingsLKR:      savings,
			CallToAction:    "🎟️ Add Combo to Food Plan",
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
			Subtitle:        fmt.Sprintf("Live counter markdown active! %s. Pair with CinePass VIP for up to 35%% off!", flashItem.FlashDealText),
			BadgeText:       "⚡ LIVE FLASH AD",
			DiscountText:    fmt.Sprintf("Live LKR %.0f &bull; Pass LKR %.0f (Total Savings LKR %.0f)", flashItem.Price, flashItem.BestDiscountedPrice, savings),
			PromoCode:       "CINE-PLATINUM",
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
			Subtitle:        fmt.Sprintf("Pair your warm popcorn bucket with %s! Explore example popcorn and drink combo savings.", beverageItem.Description),
			BadgeText:       "🍿 DRINK + SNACK PAIRING",
			DiscountText:    fmt.Sprintf("Starting from LKR %.0f estimated food offer", beverageItem.BestDiscountedPrice),
			PromoCode:       "CINE-COMBO30",
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

	return ads
}
