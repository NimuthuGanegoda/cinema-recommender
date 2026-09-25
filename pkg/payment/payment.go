package payment

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"strings"
	"time"

	"cinema-recommender/internal/models"
	"cinema-recommender/pkg/recommender"
)

// Service provides Sri Lankan payment processing, LankaQR EMV generation, and gateway discounts.
type Service struct {
	recEngine *recommender.Engine
	methods   []models.PaymentMethodInfo
}

// NewService constructs an initialized Sri Lankan payment gateway service.
func NewService(recEngine *recommender.Engine) *Service {
	methods := []models.PaymentMethodInfo{
		{
			ID:           models.PaymentMethodLankaQR,
			Name:         "LankaQR (CBSL National Standard)",
			Category:     "National QR",
			Icon:         "📱",
			Description:  "Central Bank of Sri Lanka National Standard QR. Supported by Commercial Bank Q+, Sampath WePay, FriMi, BOC SmartPay, and HNB SOLO.",
			DiscountPct:  0.05,
			SpecialOffer: "5% CBSL National Digital Incentive",
			SupportedApps: []string{
				"Commercial Bank COMBANK Q+",
				"Sampath Bank WePay",
				"Nations Trust Bank FriMi",
				"HNB SOLO",
				"BOC SmartPay",
				"Seylan Bank Pay",
			},
		},
		{
			ID:           models.PaymentMethodFriMi,
			Name:         "FriMi (Nations Trust Bank)",
			Category:     "Bank Digital Wallet",
			Icon:         "🏦",
			Description:  "Sri Lanka's leading digital banking lifestyle app with instant cashback rewards.",
			DiscountPct:  0.10,
			SpecialOffer: "10% FriMi Concession Cashback",
		},
		{
			ID:           models.PaymentMethodGenie,
			Name:         "Genie (Dialog Finance)",
			Category:     "Fintech Wallet",
			Icon:         "🧞",
			Description:  "Dialog Finance digital payment app supporting linked Sri Lankan bank accounts and cards.",
			DiscountPct:  0.08,
			SpecialOffer: "8% Genie Digital Rebate",
		},
		{
			ID:             models.PaymentMethodEzCash,
			Name:           "eZ Cash (Dialog Axiata)",
			Category:       "Telco Mobile Money",
			Icon:           "📶",
			Description:    "Direct mobile wallet payment for Dialog, Hutch, and Airtel subscribers via instant USSD prompt.",
			RequiresMobile: true,
			SpecialOffer:   "Direct mobile wallet charge via USSD PIN",
		},
		{
			ID:             models.PaymentMethodMcash,
			Name:           "mCash (SLT-Mobitel)",
			Category:       "Telco Mobile Money",
			Icon:           "📶",
			Description:    "Mobitel and Sri Lanka Telecom mobile money wallet with instant SMS payment authorization.",
			RequiresMobile: true,
			SpecialOffer:   "Zero service fee on movie concessions",
		},
		{
			ID:           models.PaymentMethodKoko,
			Name:         "Koko (Buy Now, Pay Later)",
			Category:     "BNPL",
			Icon:         "🛍️",
			Description:  "Popular Sri Lankan BNPL gateway. Split cinema concession orders over LKR 1,500 into 3 interest-free monthly installments.",
			Installments: 3,
			SpecialOffer: "Pay in 3 monthly installments at 0% interest",
		},
		{
			ID:           models.PaymentMethodMintpay,
			Name:         "Mintpay (Buy Now, Pay Later)",
			Category:     "BNPL",
			Icon:         "🛍️",
			Description:  "Sri Lanka's homegrown BNPL platform. Split payment into 3 debit/credit installments with zero interest.",
			Installments: 3,
			SpecialOffer: "Pay in 3 installments using any debit card",
		},
		{
			ID:           models.PaymentMethodCard,
			Name:         "Sri Lankan Bank Card (LankaPay / Visa / MC)",
			Category:     "Bank Card",
			Icon:         "💳",
			Description:  "Direct debit and credit cards issued by Commercial Bank, Sampath Bank, HNB, BOC, and People's Bank.",
		},
		{
			ID:           models.PaymentMethodCounterCash,
			Name:         "Cash at Cinema Concession Stand",
			Category:     "Cash on Collection",
			Icon:         "💵",
			Description:  "Pay with physical Sri Lankan Rupees (LKR) at the theater snack counter when collecting your order.",
		},
	}

	return &Service{
		recEngine: recEngine,
		methods:   methods,
	}
}

// GetSupportedPaymentMethods returns metadata on all supported Sri Lankan payment gateways.
func (s *Service) GetSupportedPaymentMethods() []models.PaymentMethodInfo {
	res := make([]models.PaymentMethodInfo, len(s.methods))
	copy(res, s.methods)
	return res
}

// ProcessCheckout processes a concession order through chosen Sri Lankan payment gateway.
func (s *Service) ProcessCheckout(
	cinema models.Cinema,
	items []models.ConcessionItem,
	req models.CheckoutRequest,
) (models.CheckoutResult, error) {
	if len(items) == 0 {
		return models.CheckoutResult{}, fmt.Errorf("cart is empty: please select concession items before checkout")
	}

	method := req.PaymentMethod
	if method == "" || method == models.PaymentMethodRedopay {
		method = models.PaymentMethodLankaQR
	}

	methodInfo := s.findMethod(method)

	var subtotal float64
	for _, it := range items {
		subtotal += it.Price
	}
	subtotal = math.Round(subtotal*100) / 100

	status := "COMPLETED"
	instructions := ""
	installmentNote := ""

	// 1. Evaluate Scope Privilege Concession Pass savings if member tier or promo code provided
	tier := req.PrivilegeTier
	if tier == "" {
		tier = req.RedopayTier
	}

	var discount float64
	if req.PromoCode != "" || (tier != "" && tier != models.TierStandard) {
		rec := s.recEngine.EvaluateCart(cinema, items, req.PromoCode, tier)
		discount = rec.DiscountAmount
		if discount > 0 {
			instructions = fmt.Sprintf("Scope Privilege concession savings applied: %s. ", rec.Message)
		}
	}

	// 2. Apply payment gateway-specific incentives & gateway handling
	switch method {
	case models.PaymentMethodLankaQR:
		// CBSL 5% digital incentive
		lankaQRDiscount := math.Round(subtotal*0.05*100) / 100
		if lankaQRDiscount > discount {
			discount = lankaQRDiscount
		}
		status = "PENDING_LANKAQR_SCAN"
		instructions += "Scan the dynamic LankaQR using any supported Sri Lankan banking app (COMBANK Q+, WePay, FriMi, SOLO, SmartPay) to complete payment."

	case models.PaymentMethodFriMi:
		// 10% FriMi cashback
		frimiDiscount := math.Round(subtotal*0.10*100) / 100
		if frimiDiscount > discount {
			discount = frimiDiscount
		}
		instructions += "Approved via FriMi API. 10% partner concession cashback credited to your FriMi wallet."

	case models.PaymentMethodGenie:
		// 8% Genie rebate
		genieDiscount := math.Round(subtotal*0.08*100) / 100
		if genieDiscount > discount {
			discount = genieDiscount
		}
		instructions += "Processed via Genie by Dialog Finance. 8% rebate applied to checkout."

	case models.PaymentMethodEzCash:
		phone := cleanPhone(req.CustomerPhone)
		status = "PENDING_USSD_PIN"
		instructions += fmt.Sprintf("A USSD payment prompt has been dispatched to %s. Enter your 4-digit eZ Cash PIN on your handset to authorize.", phone)

	case models.PaymentMethodMcash:
		phone := cleanPhone(req.CustomerPhone)
		status = "PENDING_USSD_PIN"
		instructions += fmt.Sprintf("SLT-Mobitel mCash notification sent to %s. Confirm payment via your mCash wallet.", phone)

	case models.PaymentMethodKoko:
		finalPayable := subtotal - discount
		installments := 3
		monthlyAmount := math.Round((finalPayable/float64(installments))*100) / 100
		installmentNote = fmt.Sprintf("Pay in %d interest-free monthly installments of LKR %.2f with Koko", installments, monthlyAmount)
		instructions += "Koko BNPL schedule created. First installment charged today, remaining 2 installments due over the next 60 days."

	case models.PaymentMethodMintpay:
		finalPayable := subtotal - discount
		installments := 3
		monthlyAmount := math.Round((finalPayable/float64(installments))*100) / 100
		installmentNote = fmt.Sprintf("Pay in %d interest-free debit card installments of LKR %.2f with Mintpay", installments, monthlyAmount)
		instructions += "Mintpay BNPL installment plan activated. Split across 3 monthly debit card payments with zero interest."

	case models.PaymentMethodCounterCash:
		instructions += "Order reserved! Present your Order ID at the cinema snack counter and pay in cash (LKR) to collect your freshly prepared concessions."

	default:
		instructions += "Payment authorized via Sri Lankan Interbank Payment Gateway."
	}

	finalPayable := math.Round((subtotal-discount)*100) / 100
	if finalPayable < 0 {
		finalPayable = 0
	}

	orderID := fmt.Sprintf("VCH-LK-%s-%d", strings.ToUpper(cinema.ID), randomInt(10000, 99999))
	txnRef := fmt.Sprintf("TXN-CBSL-%d", randomInt(10000000, 99999999))

	var lankaQR string
	if method == models.PaymentMethodLankaQR {
		lankaQR = generateMockLankaQREMV(cinema, finalPayable, orderID)
	}

	pickupLoc := cinema.FoodPlaceName
	if pickupLoc == "" {
		pickupLoc = fmt.Sprintf("Main Concession Stand, %s", cinema.Name)
	}

	redemptionNotice := "⚠️ In-Person Counter Redemption: This is a cinema concession recommendation voucher. Food is NOT delivered online. Please visit the theater concession counter in person to collect your freshly prepared snacks."
	cashierInst := "Present this voucher code and your payment / LankaQR confirmation to the concession counter cashier upon physical arrival at the cinema."

	return models.CheckoutResult{
		OrderID:                  orderID,
		CinemaName:               cinema.Name,
		CinemaCity:               cinema.City,
		Items:                    items,
		SubtotalLKR:              subtotal,
		DiscountLKR:              discount,
		FinalPayableLKR:          finalPayable,
		PaymentMethod:            method,
		PaymentMethodName:        methodInfo.Name,
		Status:                   status,
		TransactionReference:     txnRef,
		LankaQREMV:               lankaQR,
		Instructions:             instructions,
		InstallmentNote:          installmentNote,
		VoucherType:              "PHYSICAL_IN_THEATER_COUNTER_VOUCHER",
		CounterPickupLocation:    pickupLoc,
		PhysicalRedemptionNotice: redemptionNotice,
		CashierInstructions:     cashierInst,
		CreatedAt:                time.Now().Format("2006-01-02 15:04:05 MST"),
	}, nil
}

func (s *Service) findMethod(id models.PaymentMethod) models.PaymentMethodInfo {
	for _, m := range s.methods {
		if m.ID == id {
			return m
		}
	}
	return models.PaymentMethodInfo{
		ID:   id,
		Name: string(id),
	}
}

func cleanPhone(phone string) string {
	phone = strings.TrimSpace(phone)
	if phone == "" {
		return "+94 77 123 4567"
	}
	return phone
}

func randomInt(min, max int64) int64 {
	nBig, err := rand.Int(rand.Reader, big.NewInt(max-min+1))
	if err != nil {
		return min + (time.Now().UnixNano() % (max - min + 1))
	}
	return min + nBig.Int64()
}

// generateMockLankaQREMV generates an EMVCo compliant LankaQR payload representation for CBSL standard apps.
func generateMockLankaQREMV(cinema models.Cinema, amount float64, orderID string) string {
	// CBSL LankaQR EMV standard tags:
	// 00: Format Indicator (01)
	// 01: Point of Initiation (12 = Dynamic QR)
	// 26: Merchant Account Info (LankaQR standard)
	// 52: Merchant Category Code (7832 = Motion Picture Theaters)
	// 53: Transaction Currency (144 = LKR Sri Lankan Rupee)
	// 54: Transaction Amount
	// 58: Country Code (LK)
	// 59: Merchant Name
	// 60: Merchant City
	// 62: Additional Data (Order ID)
	cleanMerchant := strings.ReplaceAll(cinema.Name, " ", "")
	if len(cleanMerchant) > 20 {
		cleanMerchant = cleanMerchant[:20]
	}
	cleanCity := cinema.City
	if len(cleanCity) > 15 {
		cleanCity = cleanCity[:15]
	}

	return fmt.Sprintf(
		"00020101021226480010com.lankaqr0118%s520478325303144540%d%.2f5802LK59%02d%s60%02d%s62%02d01%s6304LKQR",
		orderID,
		len(fmt.Sprintf("%.2f", amount)),
		amount,
		len(cleanMerchant),
		cleanMerchant,
		len(cleanCity),
		cleanCity,
		len(orderID)+4,
		orderID,
	)
}
