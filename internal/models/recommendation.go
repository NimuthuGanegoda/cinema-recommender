package models

import (
	"fmt"
	"strings"
)

// RecommendationRequest defines input filters for optimal concession generation.
type RecommendationRequest struct {
	CinemaID            string         `json:"cinema_id,omitempty"`
	City                string         `json:"city,omitempty"`
	BudgetLKR           float64        `json:"budget_lkr"`
	PartySize           int            `json:"party_size"`
	PreferredCategories []ItemCategory `json:"preferred_categories,omitempty"`
	DietaryPreferences []string       `json:"dietary_preferences,omitempty"`
	PrivilegeTier      PrivilegeTier  `json:"privilege_tier,omitempty"`
	RedopayTier        RedopayTier    `json:"redopay_tier,omitempty"`
	PromoCode          string         `json:"promo_code,omitempty"`
}

// RecommendationResult represents the calculated recommendation with Scope Privilege and promotional savings.
type RecommendationResult struct {
	CinemaName     string           `json:"cinema_name"`
	Location       string           `json:"location"`
	PartySize      int              `json:"party_size,omitempty"`
	SelectedItems  []ConcessionItem `json:"selected_items"`
	OriginalTotal  float64          `json:"original_total_lkr"`
	DiscountAmount float64          `json:"discount_amount_lkr"`
	FinalTotal     float64          `json:"final_total_lkr"`
	Savings        SavingsBreakdown `json:"savings_breakdown"`
	AppliedPromo   string           `json:"applied_promo"`
	Message        string           `json:"message"`
	NextTierUpsell string           `json:"next_tier_upsell,omitempty"`
	ValueScore     float64          `json:"value_score,omitempty"`
}

// String provides a formatted summary of the recommendation.
func (r RecommendationResult) String() string {
	var itemNames []string
	for _, it := range r.SelectedItems {
		itemNames = append(itemNames, fmt.Sprintf("%s (LKR %.0f)", it.Name, it.Price))
	}
	itemsStr := strings.Join(itemNames, ", ")
	if itemsStr == "" {
		itemsStr = "No items"
	}

	return fmt.Sprintf(
		"🎬 [%s - %s]\n   Items (%d): %s\n   Original: LKR %.2f | Discount: LKR %.2f | Final: LKR %.2f (%s)",
		r.CinemaName, r.Location, len(r.SelectedItems), itemsStr, r.OriginalTotal, r.DiscountAmount, r.FinalTotal, r.AppliedPromo,
	)
}

// CartItemInput represents an item entry in a user's concession order.
type CartItemInput struct {
	ItemID   string `json:"item_id"`
	Quantity int    `json:"quantity"`
}

// CartEvaluationRequest is submitted to evaluate a user-defined basket.
type CartEvaluationRequest struct {
	CinemaID      string          `json:"cinema_id"`
	City          string          `json:"city,omitempty"`
	Items         []CartItemInput `json:"items"`
	PrivilegeTier PrivilegeTier   `json:"privilege_tier,omitempty"`
	RedopayTier   RedopayTier     `json:"redopay_tier,omitempty"`
	PromoCode     string          `json:"promo_code,omitempty"`
}
