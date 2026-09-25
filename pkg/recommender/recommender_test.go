package recommender

import (
	"testing"

	"cinema-recommender/internal/models"
)

func TestRecommenderPromotions(t *testing.T) {
	eng := NewEngine()

	// Under 1000 spend -> No promo
	itemsLow := []models.ConcessionItem{
		{ID: "1", Name: "Popcorn", Category: models.CategoryPopcorn, Price: 600.00},
	}
	resLow := eng.CalculateBestDeal(models.Cinema{ID: "C1", Name: "Test", City: "Kandy"}, itemsLow)
	if resLow.AppliedPromo != "NONE" || resLow.DiscountAmount != 0 {
		t.Fatalf("expected no promo below 1000 spend, got %s (LKR %.2f)", resLow.AppliedPromo, resLow.DiscountAmount)
	}

	// 1000 to 1999 -> REDOPAY-SNACK15 (15%)
	itemsMed := []models.ConcessionItem{
		{ID: "1", Name: "Popcorn", Category: models.CategoryPopcorn, Price: 900.00},
		{ID: "2", Name: "Drink", Category: models.CategoryBeverage, Price: 600.00},
	}
	resMed := eng.CalculateBestDeal(models.Cinema{ID: "C1", Name: "Test", City: "Kandy"}, itemsMed)
	if resMed.AppliedPromo != "REDOPAY-SNACK15" {
		t.Fatalf("expected REDOPAY-SNACK15 for LKR 1500 spend, got %s", resMed.AppliedPromo)
	}
	expectedDiscount := 1500.00 * 0.15
	if resMed.DiscountAmount != expectedDiscount {
		t.Fatalf("expected discount LKR %.2f, got LKR %.2f", expectedDiscount, resMed.DiscountAmount)
	}

	// Over 2000 without combo -> REDOPAY-CINEMA25 (25%)
	itemsHigh := []models.ConcessionItem{
		{ID: "1", Name: "Snack 1", Category: models.CategorySnack, Price: 1200.00},
		{ID: "2", Name: "Snack 2", Category: models.CategorySnack, Price: 900.00},
	}
	resHigh := eng.CalculateBestDeal(models.Cinema{ID: "C1", Name: "Test", City: "Gampaha"}, itemsHigh)
	if resHigh.AppliedPromo != "REDOPAY-CINEMA25" {
		t.Fatalf("expected REDOPAY-CINEMA25 for non-combo >2000 spend, got %s", resHigh.AppliedPromo)
	}

	// Over 2200 with Combo (Popcorn + Drink) -> REDOPAY-COMBO30 (30%)
	itemsCombo := []models.ConcessionItem{
		{ID: "1", Name: "Popcorn Large", Category: models.CategoryPopcorn, Price: 1400.00},
		{ID: "2", Name: "Iced Soda", Category: models.CategoryBeverage, Price: 900.00},
	}
	resCombo := eng.CalculateBestDeal(models.Cinema{ID: "C1", Name: "Test", City: "Galle"}, itemsCombo)
	if resCombo.AppliedPromo != "REDOPAY-COMBO30" {
		t.Fatalf("expected REDOPAY-COMBO30 for combo >2200 spend, got %s", resCombo.AppliedPromo)
	}
	expectedComboDiscount := 2300.00 * 0.30
	if resCombo.DiscountAmount != expectedComboDiscount {
		t.Fatalf("expected discount LKR %.2f, got LKR %.2f", expectedComboDiscount, resCombo.DiscountAmount)
	}
}

func TestPlatinumAndStudentTiers(t *testing.T) {
	eng := NewEngine()
	cinema := models.Cinema{ID: "C1", Name: "Luxe", City: "Kurunegala"}

	// Student Tier: 20% on >= 800
	studentItems := []models.ConcessionItem{
		{ID: "1", Name: "Popcorn", Category: models.CategoryPopcorn, Price: 850.00},
	}
	resStudent := eng.EvaluateCart(cinema, studentItems, "", models.TierStudent)
	if resStudent.AppliedPromo != "REDOPAY-STUDENT" {
		t.Fatalf("expected REDOPAY-STUDENT, got %s", resStudent.AppliedPromo)
	}

	// Platinum Tier: 35% on >= 2800
	platItems := []models.ConcessionItem{
		{ID: "1", Name: "Platter", Category: models.CategorySnack, Price: 3000.00},
	}
	resPlat := eng.EvaluateCart(cinema, platItems, "", models.TierPlatinum)
	if resPlat.AppliedPromo != "REDOPAY-PLATINUM" {
		t.Fatalf("expected REDOPAY-PLATINUM, got %s", resPlat.AppliedPromo)
	}
}

func TestRecommendOptimalBundle(t *testing.T) {
	eng := NewEngine()
	cinema := models.Cinema{ID: "KND-01", Name: "Kandy Cinema", City: "Kandy"}

	menu := []models.ConcessionItem{
		{ID: "1", Name: "Popcorn Caramel", Category: models.CategoryPopcorn, Price: 1200.00, InStock: true},
		{ID: "2", Name: "Mountain Dew", Category: models.CategoryBeverage, Price: 650.00, InStock: true},
		{ID: "3", Name: "Hotdog", Category: models.CategorySnack, Price: 950.00, InStock: true},
	}

	req := models.RecommendationRequest{
		BudgetLKR: 3000.00,
		PartySize: 2,
	}

	result, err := eng.RecommendOptimalBundle(cinema, menu, req)
	if err != nil {
		t.Fatalf("failed to recommend bundle: %v", err)
	}

	if len(result.SelectedItems) == 0 {
		t.Fatal("expected selected items in recommended bundle")
	}
	if result.FinalTotal > req.BudgetLKR*1.20 {
		t.Fatalf("recommended total LKR %.2f exceeded budget constraint significantly", result.FinalTotal)
	}
	if result.ValueScore <= 0 {
		t.Fatal("expected positive value score")
	}
}
