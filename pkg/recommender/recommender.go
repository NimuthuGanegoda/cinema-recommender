package recommender

import (
	"fmt"
	"math"
	"sort"
	"strings"

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
				PromoCode:         "REDOPAY-PLATINUM",
				Title:             "Redopay Platinum VIP Snacking",
				Description:       "35% off concession orders over LKR 2,800 for Redopay Platinum cardholders (capped at LKR 3,000)",
				DiscountType:      models.DiscountTypePercentage,
				DiscountPct:       0.35,
				MinSpendLKR:       2800.00,
				MaxDiscountCapLKR: 3000.00,
				EligibleTier:      models.TierPlatinum,
			},
			{
				PromoCode:         "REDOPAY-COMBO30",
				Title:             "Redopay Popcorn + Drink Super Combo",
				Description:       "30% off concession orders over LKR 2,200 containing both Popcorn & Drink (capped at LKR 2,000)",
				DiscountType:      models.DiscountTypeCombo,
				DiscountPct:       0.30,
				MinSpendLKR:       2200.00,
				MaxDiscountCapLKR: 2000.00,
				RequiresCombo:     true,
			},
			{
				PromoCode:         "REDOPAY-CINEMA25",
				Title:             "Redopay Regional 25% Concession Boost",
				Description:       "25% off outstation concession orders over LKR 2,000 paid via Redopay (capped at LKR 1,500)",
				DiscountType:      models.DiscountTypePercentage,
				DiscountPct:       0.25,
				MinSpendLKR:       2000.00,
				MaxDiscountCapLKR: 1500.00,
			},
			{
				PromoCode:         "REDOPAY-STUDENT",
				Title:             "Redopay Student Moviegoer Pass",
				Description:       "20% flat discount on concession orders over LKR 800 for registered student accounts",
				DiscountType:      models.DiscountTypePercentage,
				DiscountPct:       0.20,
				MinSpendLKR:       800.00,
				MaxDiscountCapLKR: 1000.00,
				EligibleTier:      models.TierStudent,
			},
			{
				PromoCode:         "REDOPAY-FIRST500",
				Title:             "Redopay Flat LKR 500 Snacking Voucher",
				Description:       "Flat LKR 500 off any concession order over LKR 1,800 paid via Redopay",
				DiscountType:      models.DiscountTypeFixed,
				FlatDiscountLKR:   500.00,
				MinSpendLKR:       1800.00,
				MaxDiscountCapLKR: 500.00,
			},
			{
				PromoCode:         "REDOPAY-SNACK15",
				Title:             "Redopay Regional 15% Saver",
				Description:       "15% off any regional cinema concession order over LKR 1,000 with Redopay",
				DiscountType:      models.DiscountTypePercentage,
				DiscountPct:       0.15,
				MinSpendLKR:       1000.00,
				MaxDiscountCapLKR: 800.00,
			},
		},
	}
}

// GetAvailablePromotions returns all currently registered Redopay promotional campaigns.
func (e *Engine) GetAvailablePromotions() []models.RedopayPromotion {
	promos := make([]models.RedopayPromotion, len(e.availablePromos))
	copy(promos, e.availablePromos)
	return promos
}

// BestPromotion finds the highest value Redopay promotion applicable to items and subtotal.
func (e *Engine) BestPromotion(subtotal float64) (models.RedopayPromotion, bool) {
	return e.BestPromotionWithContext(subtotal, nil, models.TierStandard, "")
}

// BestPromotionWithContext evaluates promos matching subtotal, items, cardholder tier, and optional promo code.
func (e *Engine) BestPromotionWithContext(
	subtotal float64,
	items []models.ConcessionItem,
	tier models.RedopayTier,
	forcedPromoCode string,
) (models.RedopayPromotion, bool) {
	var bestPromo models.RedopayPromotion
	var maxDiscount float64
	found := false

	hasPopcorn, hasBeverage := hasPopcornAndBeverage(items)

	cleanCode := strings.ToUpper(strings.TrimSpace(forcedPromoCode))

	for _, promo := range e.availablePromos {
		// If user specified an explicit promo code, only evaluate that code
		if cleanCode != "" && strings.ToUpper(promo.PromoCode) != cleanCode {
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

// CalculateBestDeal evaluates items from a regional cinema and computes the optimized Redopay discount.
// Preserves backwards compatibility with original signature while providing rich savings metrics.
func (e *Engine) CalculateBestDeal(cinema models.Cinema, items []models.ConcessionItem) models.RecommendationResult {
	return e.EvaluateCart(cinema, items, "", models.TierStandard)
}

// EvaluateCart evaluates a cart of concession items against all Redopay promotions.
func (e *Engine) EvaluateCart(
	cinema models.Cinema,
	items []models.ConcessionItem,
	promoCode string,
	tier models.RedopayTier,
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
		result.Message = fmt.Sprintf("Redopay deal applied: %s (Saved LKR %.2f)", promo.Description, discount)
	} else {
		result.DiscountAmount = 0
		result.FinalTotal = subtotal
		result.AppliedPromo = "NONE"
		result.Message = "No Redopay threshold met. Add more items to unlock Redopay savings."
	}

	effectivePct := 0.0
	if subtotal > 0 {
		effectivePct = round((discount / subtotal) * 100)
	}

	result.Savings = models.SavingsBreakdown{
		SubtotalLKR:          subtotal,
		RedopayDiscountLKR:   discount,
		TotalSavingsLKR:      discount,
		EffectiveDiscountPct: effectivePct,
	}

	result.NextTierUpsell = e.calculateUpsellAdvice(subtotal)

	return result
}

// RecommendOptimalBundle generates a value-maximizing concession bundle for a movie party within budget.
func (e *Engine) RecommendOptimalBundle(
	cinema models.Cinema,
	availableItems []models.ConcessionItem,
	req models.RecommendationRequest,
) (models.RecommendationResult, error) {
	if len(availableItems) == 0 {
		return models.RecommendationResult{}, fmt.Errorf("no concession items available for %s", cinema.Name)
	}

	partySize := req.PartySize
	if partySize <= 0 {
		partySize = 2 // default to movie couple
	}

	budget := req.BudgetLKR
	if budget <= 0 {
		// Provide an intelligent default budget based on party size (~LKR 1,500 per person)
		budget = float64(partySize) * 1500.00
	}

	// Categorize items
	var popcorns, beverages, snacks, combos []models.ConcessionItem
	for _, it := range availableItems {
		if !it.InStock {
			continue
		}
		switch it.Category {
		case models.CategoryPopcorn:
			popcorns = append(popcorns, it)
		case models.CategoryBeverage:
			beverages = append(beverages, it)
		case models.CategorySnack:
			snacks = append(snacks, it)
		case models.CategoryCombo:
			combos = append(combos, it)
		}
	}

	// Sort categories by price ascending
	sortByPrice := func(items []models.ConcessionItem) {
		sort.Slice(items, func(i, j int) bool {
			return items[i].Price < items[j].Price
		})
	}
	sortByPrice(popcorns)
	sortByPrice(beverages)
	sortByPrice(snacks)
	sortByPrice(combos)

	var selected []models.ConcessionItem
	currentSubtotal := 0.0

	// Strategy A: If budget allows and suitable combo exists for party, start with combo
	if len(combos) > 0 && budget >= combos[0].Price {
		for i := len(combos) - 1; i >= 0; i-- {
			if combos[i].Price <= budget*0.85 {
				selected = append(selected, combos[i])
				currentSubtotal += combos[i].Price
				break
			}
		}
	}

	// Strategy B: Ensure core snacking balance (Popcorn + Drinks per person + Snack)
	if len(selected) == 0 {
		// 1. Add popcorn (shared or individual)
		if len(popcorns) > 0 {
			bestPopcorn := popcorns[0]
			if partySize >= 2 && len(popcorns) > 1 {
				bestPopcorn = popcorns[len(popcorns)-1] // Jumbo / Large
			}
			selected = append(selected, bestPopcorn)
			currentSubtotal += bestPopcorn.Price
		}

		// 2. Add drinks (aim for 1 drink per person up to 4)
		drinkCount := partySize
		if drinkCount > 4 {
			drinkCount = 4
		}
		if len(beverages) > 0 {
			chosenDrink := beverages[0]
			if len(beverages) > 1 {
				chosenDrink = beverages[1]
			}
			for d := 0; d < drinkCount; d++ {
				if currentSubtotal+chosenDrink.Price <= budget*1.15 {
					selected = append(selected, chosenDrink)
					currentSubtotal += chosenDrink.Price
				}
			}
		}

		// 3. Add snacks if budget allows
		if len(snacks) > 0 {
			for _, s := range snacks {
				if currentSubtotal+s.Price <= budget*1.15 {
					selected = append(selected, s)
					currentSubtotal += s.Price
					if partySize <= 2 {
						break
					}
				}
			}
		}
	}

	// Check if adding one small item would unlock a significantly higher Redopay promotion tier
	selected = e.optimizeForPromotionThreshold(selected, availableItems, currentSubtotal)

	// Evaluate final deal
	res := e.EvaluateCart(cinema, selected, req.PromoCode, req.RedopayTier)
	res.PartySize = partySize

	// Compute value score (coverage of party requirements & savings efficiency)
	res.ValueScore = calculateValueScore(res, partySize)

	return res, nil
}

// optimizeForPromotionThreshold checks if adding a low-cost item bridges the gap to unlock higher savings.
func (e *Engine) optimizeForPromotionThreshold(
	currentItems []models.ConcessionItem,
	availableItems []models.ConcessionItem,
	subtotal float64,
) []models.ConcessionItem {
	thresholds := []float64{1000.00, 1800.00, 2000.00, 2200.00, 2800.00}

	for _, target := range thresholds {
		gap := target - subtotal
		// If gap is between LKR 50 and LKR 400, adding a small drink/snack actually increases net savings!
		if gap > 0 && gap <= 400 {
			var candidate models.ConcessionItem
			candidateFound := false
			for _, item := range availableItems {
				if item.Price >= gap && item.Price <= gap+350 {
					candidate = item
					candidateFound = true
					break
				}
			}
			if candidateFound {
				currentItems = append(currentItems, candidate)
				break
			}
		}
	}

	return currentItems
}

// calculateUpsellAdvice computes user advice on how close they are to unlocking the next promotion tier.
func (e *Engine) calculateUpsellAdvice(subtotal float64) string {
	if subtotal < 1000.00 {
		gap := 1000.00 - subtotal
		return fmt.Sprintf("💡 Add LKR %.0f more to unlock 15%% Redopay savings with REDOPAY-SNACK15!", gap)
	} else if subtotal < 2000.00 {
		gap := 2000.00 - subtotal
		return fmt.Sprintf("🚀 Spend LKR %.0f more to unlock 25%% Redopay savings with REDOPAY-CINEMA25!", gap)
	} else if subtotal < 2200.00 {
		gap := 2200.00 - subtotal
		return fmt.Sprintf("🍿 Combo Deal: Add Popcorn & Drink for LKR %.0f more to get 30%% off via REDOPAY-COMBO30!", gap)
	} else if subtotal < 2800.00 {
		gap := 2800.00 - subtotal
		return fmt.Sprintf("💎 Platinum VIP: Reach LKR 2,800 (LKR %.0f away) to unlock 35%% off with REDOPAY-PLATINUM!", gap)
	}
	return "✨ Maximum Redopay promotional tier unlocked! Enjoy your movie."
}

func calculateDiscount(subtotal float64, promo models.RedopayPromotion) float64 {
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

func round(val float64) float64 {
	return math.Round(val*100) / 100
}
