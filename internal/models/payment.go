package models

// PaymentMethod defines authentic supported Sri Lankan payment gateways and digital channels.
type PaymentMethod string

const (
	PaymentMethodLankaQR     PaymentMethod = "LANKAQR"
	PaymentMethodFriMi       PaymentMethod = "FRIMI"
	PaymentMethodGenie       PaymentMethod = "GENIE"
	PaymentMethodEzCash      PaymentMethod = "EZCASH"
	PaymentMethodMcash       PaymentMethod = "MCASH"
	PaymentMethodKoko        PaymentMethod = "KOKO_BNPL"
	PaymentMethodMintpay     PaymentMethod = "MINTPAY_BNPL"
	PaymentMethodCard        PaymentMethod = "LK_BANK_CARD"
	PaymentMethodCounterCash PaymentMethod = "COUNTER_CASH"

	// Deprecated: Redopay does not exist in Sri Lanka; alias retained for legacy tests only.
	PaymentMethodRedopay PaymentMethod = "REDOPAY"
)

// PaymentMethodInfo provides metadata, partner discounts, and instructions for Sri Lankan payment methods.
type PaymentMethodInfo struct {
	ID             PaymentMethod `json:"id"`
	Name           string        `json:"name"`
	Category       string        `json:"category"` // "National QR", "Bank Lifestyle App", "Fintech Wallet", "Telco Mobile Money", "BNPL", "Bank Card", "Cash"
	Icon           string        `json:"icon"`
	Description    string        `json:"description"`
	DiscountPct    float64       `json:"discount_pct,omitempty"`
	SpecialOffer   string        `json:"special_offer,omitempty"`
	Installments   int           `json:"installments,omitempty"`
	SupportedApps  []string      `json:"supported_apps,omitempty"`
	RequiresMobile bool          `json:"requires_mobile,omitempty"`
}

// CheckoutRequest contains details for placing a concession order via authentic Sri Lankan payment channels.
type CheckoutRequest struct {
	CinemaID      string          `json:"cinema_id"`
	Items         []CartItemInput `json:"items"`
	PaymentMethod PaymentMethod   `json:"payment_method"`
	PrivilegeTier PrivilegeTier   `json:"privilege_tier,omitempty"`
	RedopayTier   RedopayTier     `json:"redopay_tier,omitempty"`
	PromoCode     string          `json:"promo_code,omitempty"`
	CustomerPhone string          `json:"customer_phone,omitempty"`
	CustomerName  string          `json:"customer_name,omitempty"`
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
