package recommender

import "cinema-recommender/internal/models"

// defaultConcessionPromotions returns all standard Scope Privilege Club promotional campaigns.
func defaultConcessionPromotions() []models.ConcessionPromotion {
	return []models.ConcessionPromotion{
		{
			PromoCode:         "SCOPE-PLATINUM",
			Title:             "Scope Privilege Platinum VIP Snacking",
			Description:       "35% off concession orders over LKR 2,800 for Scope Privilege Platinum members (capped at LKR 3,000)",
			DiscountType:      models.DiscountTypePercentage,
			DiscountPct:       0.35,
			MinSpendLKR:       2800.00,
			MaxDiscountCapLKR: 3000.00,
			EligibleTier:      models.TierPlatinum,
		},
		{
			PromoCode:         "SCOPE-COMBO30",
			Title:             "Scope Popcorn + Drink Super Combo",
			Description:       "30% off concession orders over LKR 2,200 containing both Popcorn & Drink (capped at LKR 2,000)",
			DiscountType:      models.DiscountTypeCombo,
			DiscountPct:       0.30,
			MinSpendLKR:       2200.00,
			MaxDiscountCapLKR: 2000.00,
			RequiresCombo:     true,
		},
		{
			PromoCode:         "SCOPE-CINEMA25",
			Title:             "Scope Regional 25% Concession Boost",
			Description:       "25% off outstation concession orders over LKR 2,000 with Scope Concession Pass (capped at LKR 1,500)",
			DiscountType:      models.DiscountTypePercentage,
			DiscountPct:       0.25,
			MinSpendLKR:       2000.00,
			MaxDiscountCapLKR: 1500.00,
		},
		{
			PromoCode:         "SCOPE-STUDENT",
			Title:             "Scope Student Moviegoer Pass",
			Description:       "20% flat discount on concession orders over LKR 800 for registered student accounts",
			DiscountType:      models.DiscountTypePercentage,
			DiscountPct:       0.20,
			MinSpendLKR:       800.00,
			MaxDiscountCapLKR: 1000.00,
			EligibleTier:      models.TierStudent,
		},
		{
			PromoCode:         "SCOPE-FIRST500",
			Title:             "Scope Privilege Flat LKR 500 Snacking Voucher",
			Description:       "Flat LKR 500 off any concession order over LKR 1,800 with Scope Privilege",
			DiscountType:      models.DiscountTypeFixed,
			FlatDiscountLKR:   500.00,
			MinSpendLKR:       1800.00,
			MaxDiscountCapLKR: 500.00,
		},
		{
			PromoCode:         "SCOPE-SNACK15",
			Title:             "Scope Regional 15% Saver",
			Description:       "15% off any regional cinema concession order over LKR 1,000 with Scope Privilege",
			DiscountType:      models.DiscountTypePercentage,
			DiscountPct:       0.15,
			MinSpendLKR:       1000.00,
			MaxDiscountCapLKR: 800.00,
		},
	}
}
