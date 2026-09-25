package models

// CinemaAd represents an in-theater lobby screen promotional advertisement promoting combos, food, and drink discounts.
type CinemaAd struct {
	ID              string   `json:"id"`
	CinemaID        string   `json:"cinema_id,omitempty"`
	Title           string   `json:"title"`
	Subtitle        string   `json:"subtitle"`
	BadgeText       string   `json:"badge_text"` // e.g. "HOT COMBO DEAL", "LOBBY EXCLUSIVE", "FLASH MATINEE AD"
	DiscountText    string   `json:"discount_text"` // e.g. "Save LKR 720 (30% OFF)", "Combo Price LKR 1,680"
	PromoCode       string   `json:"promo_code,omitempty"`
	TargetItemID    string   `json:"target_item_id,omitempty"`
	TargetItemName  string   `json:"target_item_name,omitempty"`
	OriginalPrice   float64  `json:"original_price_lkr,omitempty"`
	DiscountedPrice float64  `json:"discounted_price_lkr,omitempty"`
	SavingsLKR      float64  `json:"savings_lkr,omitempty"`
	CallToAction    string   `json:"call_to_action"` // e.g. "Claim Combo Deal", "Add Discounted Snack"
	IsFlashAd       bool     `json:"is_flash_ad"`
	Tags            []string `json:"tags,omitempty"`
}
