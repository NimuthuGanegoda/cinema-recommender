package payment

import (
	"cinema-recommender/internal/models"
	"cinema-recommender/pkg/recommender"
)

// Service provides Sri Lankan payment processing, LankaQR EMV generation, and gateway discounts.
type Service struct {
	recEngine *recommender.Engine
	methods   []models.PaymentMethodInfo
}

// NewService constructs an initialized Sri Lankan payment gateway service.
func NewService(recEngine *recommender.Engine) *Service {
	return &Service{
		recEngine: recEngine,
		methods:   defaultPaymentMethods(),
	}
}

// GetSupportedPaymentMethods returns metadata on all supported Sri Lankan payment gateways.
func (s *Service) GetSupportedPaymentMethods() []models.PaymentMethodInfo {
	res := make([]models.PaymentMethodInfo, len(s.methods))
	copy(res, s.methods)
	return res
}
