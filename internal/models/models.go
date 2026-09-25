package models

import (
	"fmt"
	"strings"
)

// Cinema represents a movie theater entity with location and food/dining metadata.
type Cinema struct {
	ID                 string   `json:"id"`
	Name               string   `json:"name"`
	Chain              string   `json:"chain,omitempty"`
	City               string   `json:"city"`
	Province           string   `json:"province,omitempty"`
	Address            string   `json:"address,omitempty"`
	Screens            int      `json:"screens,omitempty"`
	IsOutside          bool     `json:"is_outside_colombo"`
	Latitude           float64  `json:"latitude,omitempty"`
	Longitude          float64  `json:"longitude,omitempty"`
	MapURL             string   `json:"map_url,omitempty"`
	HasInHouseFood     bool     `json:"has_in_house_food"`
	FoodPlaceName      string   `json:"food_place_name,omitempty"`
	FoodPlaceType      string   `json:"food_place_type,omitempty"`
	NearbyFoodOptions  []string `json:"nearby_food_options,omitempty"`
	FoodHours          string   `json:"food_hours,omitempty"`
	FoodDeliveryToSeat bool     `json:"food_delivery_to_seat"`
	DistanceKM         float64  `json:"distance_km,omitempty"` // populated dynamically if user coordinates provided
}

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
	Provider        string  `json:"provider"` // "Redopay", "LankaQR", "FriMi", "Genie", "Koko"
	Title           string  `json:"title"`
	DiscountPct     float64 `json:"discount_pct"`
	DiscountedPrice float64 `json:"discounted_price_lkr"`
	SavingsLKR      float64 `json:"savings_lkr"`
	Requirement     string  `json:"requirement"`
}

// ConcessionItem represents a food or beverage product with detailed discount breakdowns.
type ConcessionItem struct {
	ID                  string             `json:"id"`
	CinemaID            string             `json:"cinema_id,omitempty"`
	Name                string             `json:"name"`
	Category            ItemCategory       `json:"category"`
	Size                string             `json:"size,omitempty"`
	Price               float64            `json:"price_lkr"`
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

// RedopayTier defines member loyalty and account tiers.
type RedopayTier string

const (
	TierStandard RedopayTier = "Standard"
	TierSilver   RedopayTier = "Silver"
	TierGold     RedopayTier = "Gold"
	TierPlatinum RedopayTier = "Platinum"
	TierStudent  RedopayTier = "Student"
)

// DiscountType defines how a promotion computes its price reduction.
type DiscountType string

const (
	DiscountTypePercentage DiscountType = "Percentage"
	DiscountTypeFixed      DiscountType = "FixedAmount"
	DiscountTypeCombo      DiscountType = "ComboSpecial"
)

// RedopayPromotion holds promotion metadata for Redopay discounts.
type RedopayPromotion struct {
	PromoCode         string       `json:"promo_code"`
	Title             string       `json:"title"`
	Description       string       `json:"description"`
	DiscountType      DiscountType `json:"discount_type"`
	DiscountPct       float64      `json:"discount_percentage"` // e.g. 0.25 for 25%
	FlatDiscountLKR   float64      `json:"flat_discount_lkr"`   // e.g. 500 for LKR 500 off
	MinSpendLKR       float64      `json:"min_spend_lkr"`
	MaxDiscountCapLKR float64      `json:"max_discount_cap_lkr,omitempty"` // e.g. 1500 max cap
	EligibleTier      RedopayTier  `json:"eligible_tier,omitempty"`
	RequiresCombo     bool         `json:"requires_combo,omitempty"`
}

// SavingsBreakdown details original prices, discounts, and net savings.
type SavingsBreakdown struct {
	SubtotalLKR          float64 `json:"subtotal_lkr"`
	BundleDiscountLKR    float64 `json:"bundle_discount_lkr"`
	RedopayDiscountLKR   float64 `json:"redopay_discount_lkr"`
	TotalSavingsLKR      float64 `json:"total_savings_lkr"`
	EffectiveDiscountPct float64 `json:"effective_discount_pct"`
}

// RecommendationRequest defines input filters for optimal concession generation.
type RecommendationRequest struct {
	CinemaID            string         `json:"cinema_id,omitempty"`
	City                string         `json:"city,omitempty"`
	BudgetLKR           float64        `json:"budget_lkr"`
	PartySize           int            `json:"party_size"`
	PreferredCategories []ItemCategory `json:"preferred_categories,omitempty"`
	DietaryPreferences []string       `json:"dietary_preferences,omitempty"`
	RedopayTier         RedopayTier    `json:"redopay_tier,omitempty"`
	PromoCode           string         `json:"promo_code,omitempty"`
}

// RecommendationResult represents the calculated recommendation with Redopay savings.
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
	CinemaID    string          `json:"cinema_id"`
	City        string          `json:"city,omitempty"`
	Items       []CartItemInput `json:"items"`
	RedopayTier RedopayTier     `json:"redopay_tier,omitempty"`
	PromoCode   string          `json:"promo_code,omitempty"`
}

// PaymentMethod defines supported Sri Lankan payment gateways and digital channels.
type PaymentMethod string

const (
	PaymentMethodRedopay    PaymentMethod = "REDOPAY"
	PaymentMethodLankaQR    PaymentMethod = "LANKAQR"
	PaymentMethodEzCash     PaymentMethod = "EZCASH"
	PaymentMethodMcash      PaymentMethod = "MCASH"
	PaymentMethodFriMi      PaymentMethod = "FRIMI"
	PaymentMethodGenie      PaymentMethod = "GENIE"
	PaymentMethodKoko       PaymentMethod = "KOKO_BNPL"
	PaymentMethodCard       PaymentMethod = "LK_BANK_CARD"
	PaymentMethodCounterCash PaymentMethod = "COUNTER_CASH"
)

// PaymentMethodInfo provides metadata, partner discounts, and instructions for Sri Lankan payment methods.
type PaymentMethodInfo struct {
	ID             PaymentMethod `json:"id"`
	Name           string        `json:"name"`
	Category       string        `json:"category"` // "Digital Partner", "National QR", "Telco Wallet", "Bank Wallet", "BNPL", "Card"
	Icon           string        `json:"icon"`
	Description    string        `json:"description"`
	DiscountPct    float64       `json:"discount_pct,omitempty"`
	SpecialOffer   string        `json:"special_offer,omitempty"`
	Installments   int           `json:"installments,omitempty"`
	SupportedApps  []string      `json:"supported_apps,omitempty"`
	RequiresMobile bool          `json:"requires_mobile,omitempty"`
}

// CheckoutRequest contains details for placing a concession order via Sri Lankan payment channels.
type CheckoutRequest struct {
	CinemaID       string          `json:"cinema_id"`
	Items          []CartItemInput `json:"items"`
	PaymentMethod  PaymentMethod   `json:"payment_method"`
	RedopayTier    RedopayTier     `json:"redopay_tier,omitempty"`
	PromoCode      string          `json:"promo_code,omitempty"`
	CustomerPhone  string          `json:"customer_phone,omitempty"`
	CustomerName   string          `json:"customer_name,omitempty"`
}

// CheckoutResult details the order confirmation, Sri Lankan payment status, and transaction details.
type CheckoutResult struct {
	OrderID                  string           `json:"order_id"`
	CinemaName               string           `json:"cinema_name"`
	CinemaCity               string           `json:"cinema_city"`
	Items                    []ConcessionItem `json:"items"`
	SubtotalLKR              float64          `json:"subtotal_lkr"`
	DiscountLKR              float64          `json:"discount_lkr"`
	FinalPayableLKR          float64          `json:"final_payable_lkr"`
	PaymentMethod            PaymentMethod    `json:"payment_method"`
	PaymentMethodName        string           `json:"payment_method_name"`
	Status                   string           `json:"status"` // "COMPLETED", "PENDING_LANKAQR_SCAN", "PENDING_USSD_PIN"
	TransactionReference     string           `json:"transaction_reference"`
	LankaQREMV               string           `json:"lanka_qr_emv,omitempty"`
	Instructions             string           `json:"instructions"`
	InstallmentNote          string           `json:"installment_note,omitempty"`
	VoucherType              string           `json:"voucher_type"` // "PHYSICAL_COUNTER_REDEMPTION"
	CounterPickupLocation    string           `json:"counter_pickup_location"`
	PhysicalRedemptionNotice string           `json:"physical_redemption_notice"`
	CashierInstructions     string           `json:"cashier_instructions"`
	CreatedAt                string           `json:"created_at"`
}
