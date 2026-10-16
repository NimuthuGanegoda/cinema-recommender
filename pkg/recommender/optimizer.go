package recommender

import (
	"fmt"
	"sort"

	"cinema-recommender/internal/models"
)

// RecommendOptimalBundle builds a budget-constrained category heuristic.
// The legacy function name does not imply a globally optimal solution.
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

	// Check if adding one small item would unlock a significantly higher promotion tier
	tier := req.PrivilegeTier
	if tier == "" {
		tier = req.RedopayTier
	}
	selected = e.optimizeForPromotionThreshold(cinema, selected, availableItems, req.PromoCode, tier)

	// Resolve member tier

	// Evaluate final deal
	res := e.EvaluateCart(cinema, selected, req.PromoCode, tier)
	// A promotion threshold must never make the final recommendation exceed the budget.
	for res.FinalTotal > budget && len(selected) > 0 {
		mostExpensive := 0
		for i := range selected {
			if selected[i].Price > selected[mostExpensive].Price {
				mostExpensive = i
			}
		}
		selected = append(selected[:mostExpensive], selected[mostExpensive+1:]...)
		res = e.EvaluateCart(cinema, selected, req.PromoCode, tier)
	}
	if len(selected) == 0 {
		// A small budget may still afford a drink or snack when popcorn does not fit.
		var cheapest *models.ConcessionItem
		for i := range availableItems {
			item := &availableItems[i]
			if item.InStock && item.Price > 0 && item.Price <= budget &&
				(cheapest == nil || item.Price < cheapest.Price) {
				cheapest = item
			}
		}
		if cheapest != nil {
			res = e.EvaluateCart(cinema, []models.ConcessionItem{*cheapest}, req.PromoCode, tier)
		}
	}
	res.PartySize = partySize
	res.Algorithm = "category-heuristic-v2"
	res.BudgetRemaining = round(budget - res.FinalTotal)
	res.Explanation = fmt.Sprintf("Category-based suggestion for %d people. Final cost LKR %.2f is within the LKR %.2f food budget; LKR %.2f remains. Promotions are example rules. This heuristic does not guarantee the best possible bundle or one drink per person.", partySize, res.FinalTotal, budget, res.BudgetRemaining)

	// Compute value score (coverage of party requirements & savings efficiency)
	res.ValueScore = calculateValueScore(res, partySize)

	return res, nil
}

// optimizeForPromotionThreshold adds an item only when the actual eligible
// promotion makes the entire basket cheaper. Thresholds come from the rules.
func (e *Engine) optimizeForPromotionThreshold(cinema models.Cinema, currentItems, availableItems []models.ConcessionItem, code string, tier models.PrivilegeTier) []models.ConcessionItem {
	baseline := e.EvaluateCart(cinema, currentItems, code, tier)
	bestTotal := baseline.FinalTotal
	var best *models.ConcessionItem
	for i := range availableItems {
		item := &availableItems[i]
		if !item.InStock || item.Price <= 0 {
			continue
		}
		candidate := append(append([]models.ConcessionItem(nil), currentItems...), *item)
		result := e.EvaluateCart(cinema, candidate, code, tier)
		if result.FinalTotal < bestTotal {
			bestTotal = result.FinalTotal
			best = item
		}
	}
	if best != nil {
		return append(currentItems, *best)
	}
	return currentItems
}
