package recommender

import (
	"strings"

	"cinema-recommender/internal/models"
)

// Engine manages recommendation calculations and promotion evaluations.
type Engine struct {
	availablePromos []models.ConcessionPromotion
}

// NewEngine constructs a recommendation engine with standard Scope Privilege promotional campaigns.
func NewEngine() *Engine {
	return &Engine{
		availablePromos: defaultConcessionPromotions(),
	}
}

// GetAvailablePromotions returns all currently registered Scope Privilege promotional campaigns.
func (e *Engine) GetAvailablePromotions() []models.ConcessionPromotion {
	promos := make([]models.ConcessionPromotion, len(e.availablePromos))
	copy(promos, e.availablePromos)
	return promos
}

// BestPromotion finds the highest value Scope Privilege promotion applicable to items and subtotal.
func (e *Engine) BestPromotion(subtotal float64) (models.ConcessionPromotion, bool) {
	return e.BestPromotionWithContext(subtotal, nil, models.TierStandard, "")
}

// BestPromotionWithContext evaluates promos matching subtotal, items, cardholder tier, and optional promo code.
func (e *Engine) BestPromotionWithContext(
	subtotal float64,
	items []models.ConcessionItem,
	tier models.PrivilegeTier,
	forcedPromoCode string,
) (models.ConcessionPromotion, bool) {
	var bestPromo models.ConcessionPromotion
	var maxDiscount float64
	found := false

	hasPopcorn, hasBeverage := hasPopcornAndBeverage(items)

	cleanCode := strings.ToUpper(strings.TrimSpace(forcedPromoCode))
	aliasCode := strings.Replace(cleanCode, "REDOPAY-", "SCOPE-", 1)

	for _, promo := range e.availablePromos {
		// If user specified an explicit promo code, only evaluate that code (supports SCOPE and legacy REDOPAY prefix)
		if cleanCode != "" && strings.ToUpper(promo.PromoCode) != cleanCode && strings.ToUpper(promo.PromoCode) != aliasCode {
			continue
		}

		// Tier eligibility check
		if promo.EligibleTier != "" && promo.EligibleTier != tier && cleanCode == "" {
			continue
		}

		// Item combo requirement check
		if promo.RequiresCombo && items != nil && (!hasPopcorn || !hasBeverage) {
			continue
		}

		// Minimum spend check
		if subtotal >= promo.MinSpendLKR {
			discount := calculateDiscount(subtotal, promo)
			if discount > maxDiscount {
				maxDiscount = discount
				bestPromo = promo
				found = true
			}
		}
	}

	return bestPromo, found
}
