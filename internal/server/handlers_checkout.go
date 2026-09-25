package server

import (
	"encoding/json"
	"net/http"

	"cinema-recommender/internal/models"
)

// Handler: List Supported Sri Lankan Payment Methods
func (s *Server) handleListPaymentMethods(w http.ResponseWriter, r *http.Request) {
	methods := s.paymentSvc.GetSupportedPaymentMethods()
	s.writeJSON(w, http.StatusOK, methods)
}

// Handler: Checkout via Sri Lankan Payment Channels
func (s *Server) handleCheckout(w http.ResponseWriter, r *http.Request) {
	var req models.CheckoutRequest
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

	res, err := s.paymentSvc.ProcessCheckout(cinema, selected, req)
	if err != nil {
		s.writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	s.writeJSON(w, http.StatusOK, res)
}
