package models

import "fmt"

// Cinema represents a movie theater entity.
type Cinema struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	City      string `json:"city"`
	IsOutside bool   `json:"is_outside_colombo"`
}

// ItemCategory defines the category of a concession item.
type ItemCategory string

const (
	CategoryPopcorn  ItemCategory = "Popcorn"
	CategoryBeverage ItemCategory = "Beverage"
	CategorySnack    ItemCategory = "Snack"
	CategoryCombo    ItemCategory = "Combo"
)

// ConcessionItem represents a food or beverage product.
type ConcessionItem struct {
	ID       string       `json:"id"`
	Name     string       `json:"name"`
	Category ItemCategory `json:"category"`
	Price    float64      `json:"price_lkr"`
}

// RedopayPromotion holds promotion metadata for Redopay discounts.
type RedopayPromotion struct {
	PromoCode   string  `json:"promo_code"`
	Description string  `json:"description"`
	DiscountPct float64 `json:"discount_percentage"` // e.g. 0.20 for 20%
	MinSpendLKR float64 `json:"min_spend_lkr"`
}

// RecommendationResult represents the calculated recommendation with Redopay savings.
type RecommendationResult struct {
	CinemaName     string           `json:"cinema_name"`
	Location       string           `json:"location"`
	SelectedItems  []ConcessionItem `json:"selected_items"`
	OriginalTotal  float64          `json:"original_total_lkr"`
	DiscountAmount float64          `json:"discount_amount_lkr"`
	FinalTotal     float64          `json:"final_total_lkr"`
	AppliedPromo   string           `json:"applied_promo"`
	Message        string           `json:"message"`
}

// String provides a formatted summary of the recommendation.
func (r RecommendationResult) String() string {
	return fmt.Sprintf(
		"🎬 [%s - %s]\n   Items: %d | Original: LKR %.2f | Discount: LKR %.2f | Final: LKR %.2f (%s)",
		r.CinemaName, r.Location, len(r.SelectedItems), r.OriginalTotal, r.DiscountAmount, r.FinalTotal, r.AppliedPromo,
	)
}
