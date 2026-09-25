package payment

import (
	"cinema-recommender/internal/models"
)

// defaultPaymentMethods returns the catalog of authentic supported Sri Lankan payment gateways.
func defaultPaymentMethods() []models.PaymentMethodInfo {
	return []models.PaymentMethodInfo{
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
