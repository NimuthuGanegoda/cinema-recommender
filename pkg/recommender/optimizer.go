package recommender

import (
	"fmt"
	"sort"

	"cinema-recommender/internal/models"
)

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

	// Check if adding one small item would unlock a significantly higher promotion tier
	selected = e.optimizeForPromotionThreshold(selected, availableItems, currentSubtotal)

	// Resolve member tier
	tier := req.PrivilegeTier
	if tier == "" {
		tier = req.RedopayTier
	}

	// Evaluate final deal
	res := e.EvaluateCart(cinema, selected, req.PromoCode, tier)
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
