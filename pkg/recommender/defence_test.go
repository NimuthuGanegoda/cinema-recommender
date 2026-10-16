package recommender

import (
	"cinema-recommender/internal/models"
	"testing"
)

func TestCodesCannotBypassMembership(t *testing.T) {
	e := NewEngine()
	items := []models.ConcessionItem{{ID: "snack", Price: 3000, Category: models.CategorySnack, InStock: true}}
	for _, code := range []string{"CINE-PLATINUM", "SCOPE-PLATINUM", "REDOPAY-PLATINUM", "CINE-STUDENT"} {
		r := e.EvaluateCart(models.Cinema{}, items, code, models.TierStandard)
		if r.DiscountAmount != 0 {
			t.Errorf("%s bypassed membership: %v", code, r)
		}
	}
}

func TestPromotionTopUpMustReduceFinalSpend(t *testing.T) {
	e := NewEngine()
	for _, tc := range []struct {
		name         string
		start, extra float64
		added        bool
	}{
		{"beneficial", 950, 50, true}, {"costly", 600, 400, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			current := []models.ConcessionItem{{ID: "base", Price: tc.start, InStock: true, Category: models.CategorySnack}}
			menu := []models.ConcessionItem{{ID: "extra", Price: tc.extra, InStock: true, Category: models.CategoryBeverage}}
			got := e.optimizeForPromotionThreshold(models.Cinema{}, current, menu, "", models.TierStandard)
			if (len(got) == 2) != tc.added {
				t.Fatalf("unexpected top-up: %v", got)
			}
		})
	}
}

func TestBudgetSweepAndExplanation(t *testing.T) {
	e := NewEngine()
	menu := []models.ConcessionItem{
		{ID: "pop", Price: 950, InStock: true, Category: models.CategoryPopcorn},
		{ID: "drink", Price: 50, InStock: true, Category: models.CategoryBeverage},
		{ID: "snack", Price: 450, InStock: true, Category: models.CategorySnack},
		{ID: "sold", Price: 1, InStock: false, Category: models.CategoryCombo},
	}
	for budget := 1; budget <= 10000; budget += 37 {
		for _, party := range []int{1, 2, 4, 6, 50} {
			for _, tier := range []models.PrivilegeTier{models.TierStandard, models.TierStudent, models.TierPlatinum} {
				r, err := e.RecommendOptimalBundle(models.Cinema{}, menu, models.RecommendationRequest{BudgetLKR: float64(budget), PartySize: party, PrivilegeTier: tier})
				if err != nil {
					t.Fatal(err)
				}
				if r.FinalTotal > float64(budget) || r.FinalTotal < 0 || r.BudgetRemaining < 0 {
					t.Fatalf("budget %d: %+v", budget, r)
				}
				if r.Algorithm == "" || r.Explanation == "" {
					t.Fatal("missing explanation")
				}
				for _, item := range r.SelectedItems {
					if !item.InStock {
						t.Fatal("sold-out recommendation")
					}
				}
			}
		}
	}
}

func BenchmarkRecommendation(b *testing.B) {
	e := NewEngine()
	menu := []models.ConcessionItem{{ID: "pop", Price: 950, InStock: true, Category: models.CategoryPopcorn}, {ID: "drink", Price: 350, InStock: true, Category: models.CategoryBeverage}}
	for i := 0; i < b.N; i++ {
		_, err := e.RecommendOptimalBundle(models.Cinema{}, menu, models.RecommendationRequest{BudgetLKR: 2500, PartySize: 2})
		if err != nil {
			b.Fatal(err)
		}
	}
}
