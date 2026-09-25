package models

// PrivilegeTier defines cinema loyalty and membership tiers (Scope Privilege Club).
type PrivilegeTier string

const (
	TierStandard PrivilegeTier = "Standard"
	TierSilver   PrivilegeTier = "Silver"
	TierGold     PrivilegeTier = "Gold"
	TierPlatinum PrivilegeTier = "Platinum"
	TierStudent  PrivilegeTier = "Student"
)

// RedopayTier is a backward-compatible alias for PrivilegeTier.
type RedopayTier = PrivilegeTier

// DiscountType defines how a promotion computes its price reduction.
type DiscountType string

const (
	DiscountTypePercentage DiscountType = "Percentage"
	DiscountTypeFixed      DiscountType = "FixedAmount"
	DiscountTypeCombo      DiscountType = "ComboSpecial"
)

// ConcessionPromotion holds promotion metadata for cinema concession discounts and VIP passes.
type ConcessionPromotion struct {
	PromoCode         string        `json:"promo_code"`
	Title             string        `json:"title"`
	Description       string        `json:"description"`
	DiscountType      DiscountType  `json:"discount_type"`
	DiscountPct       float64       `json:"discount_percentage"` // e.g. 0.25 for 25%
	FlatDiscountLKR   float64       `json:"flat_discount_lkr"`   // e.g. 500 for LKR 500 off
	MinSpendLKR       float64       `json:"min_spend_lkr"`
	MaxDiscountCapLKR float64       `json:"max_discount_cap_lkr,omitempty"` // e.g. 1500 max cap
	EligibleTier      PrivilegeTier `json:"eligible_tier,omitempty"`
	RequiresCombo     bool          `json:"requires_combo,omitempty"`
}

// RedopayPromotion is a backward-compatible alias for ConcessionPromotion.
type RedopayPromotion = ConcessionPromotion

// SavingsBreakdown details original prices, discounts, and net savings.
type SavingsBreakdown struct {
	SubtotalLKR          float64 `json:"subtotal_lkr"`
	BundleDiscountLKR    float64 `json:"bundle_discount_lkr"`
	PrivilegeDiscountLKR float64 `json:"privilege_discount_lkr"`
	RedopayDiscountLKR   float64 `json:"redopay_discount_lkr,omitempty"`
	TotalSavingsLKR      float64 `json:"total_savings_lkr"`
	EffectiveDiscountPct float64 `json:"effective_discount_pct"`
}
