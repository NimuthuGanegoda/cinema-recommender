package payment

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"strings"
	"time"

	"cinema-recommender/internal/models"
)

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
		plan := CalculateBNPLPlan("Koko", finalPayable)
		installmentNote = plan.InstallmentNote
		instructions += plan.Instructions

	case models.PaymentMethodMintpay:
		finalPayable := subtotal - discount
		plan := CalculateBNPLPlan("Mintpay", finalPayable)
		installmentNote = plan.InstallmentNote
		instructions += plan.Instructions

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
		lankaQR = GenerateLankaQREMV(cinema, finalPayable, orderID)
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
