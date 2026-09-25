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
