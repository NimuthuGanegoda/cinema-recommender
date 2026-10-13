package recommender

import (
	"testing"

	"cinema-recommender/internal/models"
)

func TestRecommendationRespectsSmallAndLargeBudgets(t *testing.T) {
	engine := NewEngine()
	menu := []models.ConcessionItem{
		{ID: "popcorn", Category: models.CategoryPopcorn, Price: 1200, InStock: true},
		{ID: "drink", Category: models.CategoryBeverage, Price: 350, InStock: true},
		{ID: "snack", Category: models.CategorySnack, Price: 750, InStock: true},
	}
	for _, budget := range []float64{100, 400, 900, 1500, 2500, 4000} {
		for _, party := range []int{1, 4} {
			result, err := engine.RecommendOptimalBundle(models.Cinema{ID: "test"}, menu,
				models.RecommendationRequest{BudgetLKR: budget, PartySize: party})
			if err != nil {
				t.Fatal(err)
			}
			if result.FinalTotal > budget {
				t.Fatalf("budget %.0f, party %d: recommendation costs %.0f", budget, party, result.FinalTotal)
			}
			if budget >= 350 && len(result.SelectedItems) == 0 {
				t.Fatalf("budget %.0f should allow at least the drink", budget)
			}
		}
	}
}
