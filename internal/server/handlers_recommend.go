package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

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
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON payload: " + err.Error()})
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
			cinema = models.Cinema{
				ID:        fmt.Sprintf("GEN-%s", strings.ToUpper(req.City)),
				Name:      fmt.Sprintf("Regional Cinema %s", req.City),
				City:      req.City,
				IsOutside: true,
			}
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
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "malformed JSON payload: " + err.Error()})
		return
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
		if original, found := itemMap[cartItem.ItemID]; found {
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
