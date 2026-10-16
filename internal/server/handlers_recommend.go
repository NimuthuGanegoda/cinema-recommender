package server

import (
	"net/http"

	"cinema-recommender/internal/models"
)

// Handler: List Available Scope Privilege Promotions
func (s *Server) handleListPromotions(w http.ResponseWriter, r *http.Request) {
	promos := s.recEngine.GetAvailablePromotions()
	s.writeJSON(w, http.StatusOK, promos)
}

// Handler: Recommend Optimal Bundle
func (s *Server) handleRecommend(w http.ResponseWriter, r *http.Request) {
	var req models.RecommendationRequest
	if !s.decodeJSON(w, r, &req) {
		return
	}
	if req.BudgetLKR < 0 || req.BudgetLKR > 1000000 || req.PartySize < 0 || req.PartySize > 50 {
		s.writeJSON(w, 400, map[string]string{"error": "budget must be 0–1000000 LKR and party size 0–50; zero uses defaults"})
		return
	}
	if len(req.PreferredCategories) > 0 || len(req.DietaryPreferences) > 0 {
		s.writeJSON(w, 400, map[string]string{"error": "category and dietary filtering are not supported; do not use suggestions for allergy decisions"})
		return
	}
	if !validTier(req.PrivilegeTier) || !validTier(req.RedopayTier) {
		s.writeJSON(w, 400, map[string]string{"error": "invalid membership tier"})
		return
	}

	// Resolve cinema
	var cinema models.Cinema
	var err error

	if req.CinemaID != "" {
		cinema, err = s.scraper.GetCinemaByID(req.CinemaID)
		if err != nil {
			s.handleError(w, err)
			return
		}
	} else if req.City != "" {
		if err := s.scraper.ValidateLocation(req.City); err != nil {
			s.handleError(w, err)
			return
		}
		cinemas, err := s.scraper.GetCinemasByCity(req.City)
		if err != nil || len(cinemas) == 0 {
			s.writeJSON(w, 404, map[string]string{"error": "no registered cinema for this city"})
			return
		} else {
			cinema = cinemas[0]
		}
	} else {
		// Default to first regional cinema (KCC Multiplex)
		cinemas := s.scraper.GetRegisteredCinemas()
		if len(cinemas) > 0 {
			cinema = cinemas[0]
		}
	}

	// Scrape available concessions
	items, err := s.scraper.ScrapeCinemaConcessions(cinema)
	if err != nil {
		s.handleError(w, err)
		return
	}

	// Compute optimal bundle
	result, err := s.recEngine.RecommendOptimalBundle(cinema, items, req)
	if err != nil {
		s.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	s.writeJSON(w, http.StatusOK, result)
}

// Handler: Evaluate Custom User Cart
func (s *Server) handleEvaluateCart(w http.ResponseWriter, r *http.Request) {
	var req models.CartEvaluationRequest
	if !s.decodeJSON(w, r, &req) {
		return
	}
	if len(req.Items) > 100 {
		s.writeJSON(w, 400, map[string]string{"error": "cart may contain at most 100 lines"})
		return
	}
	if !validTier(req.PrivilegeTier) || !validTier(req.RedopayTier) {
		s.writeJSON(w, 400, map[string]string{"error": "invalid membership tier"})
		return
	}
	total := 0
	for _, item := range req.Items {
		total += item.Quantity
		if item.Quantity < 1 || item.Quantity > 99 || total > 500 {
			s.writeJSON(w, 400, map[string]string{"error": "quantities must be 1–99 with at most 500 items in total"})
			return
		}
	}

	if req.CinemaID == "" {
		req.CinemaID = "KND-KCC"
	}

	cinema, err := s.scraper.GetCinemaByID(req.CinemaID)
	if err != nil {
		s.handleError(w, err)
		return
	}

	availableItems, err := s.scraper.ScrapeCinemaConcessions(cinema)
	if err != nil {
		s.handleError(w, err)
		return
	}

	// Match cart items with quantities
	var selected []models.ConcessionItem
	itemMap := make(map[string]models.ConcessionItem)
	for _, it := range availableItems {
		itemMap[it.ID] = it
	}

	for _, cartItem := range req.Items {
		original, found := itemMap[cartItem.ItemID]
		if !found || !original.InStock {
			s.writeJSON(w, 400, map[string]string{"error": "unknown or unavailable cart item: " + cartItem.ItemID})
			return
		}
		if found {
			qty := cartItem.Quantity
			if qty <= 0 {
				qty = 1
			}
			for q := 0; q < qty; q++ {
				selected = append(selected, original)
			}
		}
	}

	tier := req.PrivilegeTier
	if tier == "" {
		tier = req.RedopayTier
	}
	result := s.recEngine.EvaluateCart(cinema, selected, req.PromoCode, tier)
	s.writeJSON(w, http.StatusOK, result)
}
