package models

// ItemCategory defines the category of a concession item.
type ItemCategory string

const (
	CategoryPopcorn  ItemCategory = "Popcorn"
	CategoryBeverage ItemCategory = "Beverage"
	CategorySnack    ItemCategory = "Snack"
	CategoryCombo    ItemCategory = "Combo"
	CategoryDessert  ItemCategory = "Dessert"
)

// ItemDiscountInfo describes a specific promotional or digital payment discount applicable to a food item.
type ItemDiscountInfo struct {
	PromoCode       string  `json:"promo_code,omitempty"`
	Provider        string  `json:"provider"` // "Scope Privilege", "LankaQR", "FriMi", "Genie", "Koko"
	Title           string  `json:"title"`
	DiscountPct     float64 `json:"discount_pct"`
	DiscountedPrice float64 `json:"discounted_price_lkr"`
	SavingsLKR      float64 `json:"savings_lkr"`
	Requirement     string  `json:"requirement"`
}

// ConcessionItem represents a food or beverage product with detailed discount breakdowns and dynamic live pricing.
type ConcessionItem struct {
	ID                  string             `json:"id"`
	CinemaID            string             `json:"cinema_id,omitempty"`
	Name                string             `json:"name"`
	Category            ItemCategory       `json:"category"`
	Size                string             `json:"size,omitempty"`
	BasePrice           float64            `json:"base_price_lkr,omitempty"` // standard baseline menu price
	Price               float64            `json:"price_lkr"`                // current live dynamic price
	PriceTrend          string             `json:"price_trend,omitempty"`    // "down", "up", "flash_drop", "stable"
	PriceChangeLKR      float64            `json:"price_change_lkr,omitempty"`
	FlashDealText       string             `json:"flash_deal_text,omitempty"`
	LastPriceUpdate     string             `json:"last_price_update,omitempty"`
	LiveTickID          int64              `json:"live_tick_id,omitempty"`
	InStock             bool               `json:"in_stock"`
	Tags                []string           `json:"tags,omitempty"`
	Description         string             `json:"description,omitempty"`
	FoodLocation        string             `json:"food_location,omitempty"`
	ApplicableDiscounts []ItemDiscountInfo `json:"applicable_discounts,omitempty"`
	BestDiscountedPrice float64            `json:"best_discounted_price_lkr,omitempty"`
	MaxDiscountPct      float64            `json:"max_discount_pct,omitempty"`
	BestPromoName       string             `json:"best_promo_name,omitempty"`
	HasDiscount         bool               `json:"has_discount"`
}
