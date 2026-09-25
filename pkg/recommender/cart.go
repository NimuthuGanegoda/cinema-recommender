package recommender

import (
	"fmt"

	"cinema-recommender/internal/models"
)

// CalculateBestDeal evaluates items from a regional cinema and computes the optimized Scope Privilege discount.
func (e *Engine) CalculateBestDeal(cinema models.Cinema, items []models.ConcessionItem) models.RecommendationResult {
	return e.EvaluateCart(cinema, items, "", models.TierStandard)
}

// EvaluateCart evaluates a cart of concession items against all Scope Privilege promotions.
func (e *Engine) EvaluateCart(
	cinema models.Cinema,
	items []models.ConcessionItem,
	promoCode string,
	tier models.PrivilegeTier,
) models.RecommendationResult {
	var subtotal float64
	for _, item := range items {
		subtotal += item.Price
	}
	subtotal = round(subtotal)

	result := models.RecommendationResult{
		CinemaName:    cinema.Name,
		Location:      cinema.City,
		SelectedItems: items,
		OriginalTotal: subtotal,
	}

	if len(items) == 0 {
		result.DiscountAmount = 0
		result.FinalTotal = 0
		result.AppliedPromo = "NONE"
		result.Message = "No concession items selected."
		result.Savings = models.SavingsBreakdown{
			SubtotalLKR: 0,
		}
		return result
	}

	promo, eligible := e.BestPromotionWithContext(subtotal, items, tier, promoCode)

	var discount float64
	if eligible {
		discount = round(calculateDiscount(subtotal, promo))
		result.DiscountAmount = discount
		result.FinalTotal = round(subtotal - discount)
		result.AppliedPromo = promo.PromoCode
		result.Message = fmt.Sprintf("Scope Privilege deal applied: %s (Saved LKR %.2f)", promo.Description, discount)
	} else {
		result.DiscountAmount = 0
		result.FinalTotal = subtotal
		result.AppliedPromo = "NONE"
		result.Message = "No Scope Privilege threshold met. Add more items to unlock concession savings."
	}

	effectivePct := 0.0
	if subtotal > 0 {
		effectivePct = round((discount / subtotal) * 100)
	}

	result.Savings = models.SavingsBreakdown{
		SubtotalLKR:          subtotal,
		PrivilegeDiscountLKR: discount,
		RedopayDiscountLKR:   discount,
		TotalSavingsLKR:      discount,
		EffectiveDiscountPct: effectivePct,
	}

	result.NextTierUpsell = e.calculateUpsellAdvice(subtotal)

	return result
}

func calculateDiscount(subtotal float64, promo models.ConcessionPromotion) float64 {
	var discount float64
	switch promo.DiscountType {
	case models.DiscountTypeFixed:
		discount = promo.FlatDiscountLKR
	case models.DiscountTypePercentage, models.DiscountTypeCombo:
		discount = subtotal * promo.DiscountPct
	default:
		discount = subtotal * promo.DiscountPct
	}

	if promo.MaxDiscountCapLKR > 0 && discount > promo.MaxDiscountCapLKR {
		discount = promo.MaxDiscountCapLKR
	}

	return discount
}
