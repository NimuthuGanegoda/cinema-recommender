package recommender

import (
	"fmt"
	"math"

	"cinema-recommender/internal/models"
)

// calculateUpsellAdvice computes user advice on how close they are to unlocking the next promotion tier.
func (e *Engine) calculateUpsellAdvice(subtotal float64) string {
	if subtotal < 1000.00 {
		gap := 1000.00 - subtotal
		return fmt.Sprintf("💡 Add LKR %.0f more to unlock 15%% Scope Privilege savings with SCOPE-SNACK15!", gap)
	} else if subtotal < 2000.00 {
		gap := 2000.00 - subtotal
		return fmt.Sprintf("🚀 Spend LKR %.0f more to unlock 25%% Scope Privilege savings with SCOPE-CINEMA25!", gap)
	} else if subtotal < 2200.00 {
		gap := 2200.00 - subtotal
		return fmt.Sprintf("🍿 Combo Deal: Add Popcorn & Drink for LKR %.0f more to get 30%% off via SCOPE-COMBO30!", gap)
	} else if subtotal < 2800.00 {
		gap := 2800.00 - subtotal
		return fmt.Sprintf("💎 Platinum VIP: Reach LKR 2,800 (LKR %.0f away) to unlock 35%% off with SCOPE-PLATINUM!", gap)
	}
	return "✨ Maximum Scope Privilege promotional tier unlocked! Enjoy your movie."
}

func calculateValueScore(res models.RecommendationResult, partySize int) float64 {
	if res.FinalTotal <= 0 {
		return 0
	}
	// Items per person ratio + discount percentage weight
	itemsPerPerson := float64(len(res.SelectedItems)) / float64(partySize)
	discountBonus := (res.DiscountAmount / (res.OriginalTotal + 1)) * 50
	score := (itemsPerPerson * 25) + discountBonus
	return math.Round(score*10) / 10
}

func hasPopcornAndBeverage(items []models.ConcessionItem) (bool, bool) {
	hasPopcorn := false
	hasBeverage := false
	for _, it := range items {
		if it.Category == models.CategoryPopcorn {
			hasPopcorn = true
		}
		if it.Category == models.CategoryBeverage {
			hasBeverage = true
		}
		if it.Category == models.CategoryCombo {
			hasPopcorn = true
			hasBeverage = true
		}
	}
	return hasPopcorn, hasBeverage
}

func round(val float64) float64 {
	return math.Round(val*100) / 100
}
