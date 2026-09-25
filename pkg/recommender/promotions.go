package recommender

import "cinema-recommender/internal/models"

// defaultConcessionPromotions returns all standard universal CinePass VIP promotional campaigns.
func defaultConcessionPromotions() []models.ConcessionPromotion {
	return []models.ConcessionPromotion{
		{
			PromoCode:         "CINE-PLATINUM",
			Title:             "CinePass Platinum VIP Snacking",
			Description:       "35% off cinema concession orders over LKR 2,800 for CinePass Platinum members (capped at LKR 3,000)",
			DiscountType:      models.DiscountTypePercentage,
			DiscountPct:       0.35,
			MinSpendLKR:       2800.00,
			MaxDiscountCapLKR: 3000.00,
			EligibleTier:      models.TierPlatinum,
		},
		{
			PromoCode:         "CINE-COMBO30",
			Title:             "CineBite Popcorn + Drink Super Combo",
			Description:       "30% off concession orders over LKR 2,200 containing both Popcorn & Drink (capped at LKR 2,000)",
			DiscountType:      models.DiscountTypeCombo,
			DiscountPct:       0.30,
			MinSpendLKR:       2200.00,
			MaxDiscountCapLKR: 2000.00,
			RequiresCombo:     true,
		},
		{
			PromoCode:         "CINE-CINEMA25",
			Title:             "Cinema Concession Pass 25% Boost",
			Description:       "25% off cinema concession orders over LKR 2,000 with CinePass (capped at LKR 1,500)",
			DiscountType:      models.DiscountTypePercentage,
			DiscountPct:       0.25,
			MinSpendLKR:       2000.00,
			MaxDiscountCapLKR: 1500.00,
		},
		{
			PromoCode:         "CINE-STUDENT",
			Title:             "Student Moviegoer Concession Pass",
			Description:       "20% flat discount on concession orders over LKR 800 for registered student accounts",
			DiscountType:      models.DiscountTypePercentage,
			DiscountPct:       0.20,
			MinSpendLKR:       800.00,
			MaxDiscountCapLKR: 1000.00,
			EligibleTier:      models.TierStudent,
		},
		{
			PromoCode:         "CINE-FIRST500",
			Title:             "CinePass LKR 500 Flat Snacking Voucher",
			Description:       "Flat LKR 500 off any concession order over LKR 1,800 with CinePass VIP",
			DiscountType:      models.DiscountTypeFixed,
			FlatDiscountLKR:   500.00,
			MinSpendLKR:       1800.00,
			MaxDiscountCapLKR: 500.00,
		},
		{
			PromoCode:         "CINE-SNACK15",
			Title:             "Cinema Snack Saver 15%",
			Description:       "15% off any cinema concession order over LKR 1,000 with CinePass",
			DiscountType:      models.DiscountTypePercentage,
			DiscountPct:       0.15,
			MinSpendLKR:       1000.00,
			MaxDiscountCapLKR: 800.00,
		},
	}
}
