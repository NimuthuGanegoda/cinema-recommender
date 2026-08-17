package recommender

import (
	"fmt"
	"math"

	"cinema-recommender/internal/models"
)

// Engine manages recommendation calculations and promotion evaluations.
type Engine struct {
	availablePromos []models.RedopayPromotion
}

// NewEngine constructs a recommendation engine with standard Redopay promotional campaigns.
func NewEngine() *Engine {
	return &Engine{
		availablePromos: []models.RedopayPromotion{
			{
				PromoCode:   "REDOPAY-CINEMA25",
				Description: "25% off concession orders over LKR 2,000 paid via Redopay",
				DiscountPct: 0.25,
				MinSpendLKR: 2000.00,
			},
			{
				PromoCode:   "REDOPAY-SNACK15",
				Description: "15% off any regional cinema snack order with Redopay",
				DiscountPct: 0.15,
				MinSpendLKR: 1000.00,
			},
		},
	}
}

// BestPromotion finds the highest value Redopay promotion applicable to the subtotal.
func (e *Engine) BestPromotion(subtotal float64) (models.RedopayPromotion, bool) {
	var bestPromo models.RedopayPromotion
	var maxDiscount float64
	found := false

	for _, promo := range e.availablePromos {
		if subtotal >= promo.MinSpendLKR {
			discount := subtotal * promo.DiscountPct
			if discount > maxDiscount {
				maxDiscount = discount
				bestPromo = promo
				found = true
			}
		}
	}

	return bestPromo, found
}

// CalculateBestDeal evaluates items from a regional cinema and computes the optimized Redopay discount.
func (e *Engine) CalculateBestDeal(cinema models.Cinema, items []models.ConcessionItem) models.RecommendationResult {
	var subtotal float64
	for _, item := range items {
		subtotal += item.Price
	}

	result := models.RecommendationResult{
		CinemaName:    cinema.Name,
		Location:      cinema.City,
		SelectedItems: items,
		OriginalTotal: round(subtotal),
	}

	promo, eligible := e.BestPromotion(subtotal)
	if eligible {
		discount := round(subtotal * promo.DiscountPct)
		result.DiscountAmount = discount
		result.FinalTotal = round(subtotal - discount)
		result.AppliedPromo = promo.PromoCode
		result.Message = fmt.Sprintf("Redopay deal applied: %s (Saved LKR %.2f)", promo.Description, discount)
	} else {
		result.DiscountAmount = 0
		result.FinalTotal = round(subtotal)
		result.AppliedPromo = "NONE"
		result.Message = "No Redopay threshold met. Add more items to unlock Redopay savings."
	}

	return result
}

// round is a private helper function (starts with lowercase) for rounding currencies.
func round(val float64) float64 {
	return math.Round(val*100) / 100
}
