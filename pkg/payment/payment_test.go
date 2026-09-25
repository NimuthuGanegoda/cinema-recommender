package payment

import (
	"strings"
	"testing"

	"cinema-recommender/internal/models"
	"cinema-recommender/pkg/recommender"
)

func TestSupportedPaymentMethods(t *testing.T) {
	svc := NewService(recommender.NewEngine())
	methods := svc.GetSupportedPaymentMethods()

	if len(methods) < 7 {
		t.Fatalf("expected at least 7 authentic Sri Lankan payment methods, got %d", len(methods))
	}

	foundLankaQR := false
	foundFriMi := false
	foundGenie := false
	foundEzCash := false
	foundKoko := false
	foundMintpay := false
	foundRedopay := false

	for _, m := range methods {
		switch m.ID {
		case models.PaymentMethodLankaQR:
			foundLankaQR = true
		case models.PaymentMethodFriMi:
			foundFriMi = true
		case models.PaymentMethodGenie:
			foundGenie = true
		case models.PaymentMethodEzCash:
			foundEzCash = true
		case models.PaymentMethodKoko:
			foundKoko = true
		case models.PaymentMethodMintpay:
			foundMintpay = true
		case models.PaymentMethodRedopay:
			foundRedopay = true
		}
	}

	if !foundLankaQR || !foundFriMi || !foundGenie || !foundEzCash || !foundKoko || !foundMintpay {
		t.Fatal("expected LankaQR, FriMi, Genie, eZ Cash, Koko, and Mintpay in authentic Sri Lankan payment methods")
	}

	if foundRedopay {
		t.Fatal("Redopay does not exist in Sri Lanka and must not be in supported payment methods list")
	}
}

func TestProcessCheckoutScopePrivilegePass(t *testing.T) {
	svc := NewService(recommender.NewEngine())
	cinema := models.Cinema{ID: "KND-KCC", Name: "KCC Multiplex", City: "Kandy"}

	items := []models.ConcessionItem{
		{ID: "1", Name: "Jumbo Popcorn", Price: 1200.00, Category: models.CategoryPopcorn},
		{ID: "2", Name: "Iced Soda", Price: 1000.00, Category: models.CategoryBeverage},
	}

	req := models.CheckoutRequest{
		CinemaID:      "KND-KCC",
		PaymentMethod: models.PaymentMethodLankaQR,
		PrivilegeTier: models.TierPlatinum,
	}

	res, err := svc.ProcessCheckout(cinema, items, req)
	if err != nil {
		t.Fatalf("unexpected error in checkout: %v", err)
	}

	if res.SubtotalLKR != 2200.00 {
		t.Fatalf("expected subtotal 2200.00, got %.2f", res.SubtotalLKR)
	}
	if res.DiscountLKR <= 0 {
		t.Fatal("expected Scope Privilege discount to be applied")
	}
	if res.Status != "PENDING_LANKAQR_SCAN" {
		t.Fatalf("expected status PENDING_LANKAQR_SCAN, got %s", res.Status)
	}
}

func TestProcessCheckoutLankaQR(t *testing.T) {
	svc := NewService(recommender.NewEngine())
	cinema := models.Cinema{ID: "GLE-QNS", Name: "Queens Cinema", City: "Galle"}

	items := []models.ConcessionItem{
		{ID: "1", Name: "Nachos", Price: 1000.00, Category: models.CategorySnack},
	}

	req := models.CheckoutRequest{
		CinemaID:      "GLE-QNS",
		PaymentMethod: models.PaymentMethodLankaQR,
	}

	res, err := svc.ProcessCheckout(cinema, items, req)
	if err != nil {
		t.Fatalf("unexpected error in checkout: %v", err)
	}

	// 5% CBSL incentive discount on 1000 = 50.00
	if res.DiscountLKR != 50.00 {
		t.Fatalf("expected 50.00 LankaQR discount, got %.2f", res.DiscountLKR)
	}
	if res.FinalPayableLKR != 950.00 {
		t.Fatalf("expected 950.00 final payable, got %.2f", res.FinalPayableLKR)
	}
	if res.Status != "PENDING_LANKAQR_SCAN" {
		t.Fatalf("expected status PENDING_LANKAQR_SCAN, got %s", res.Status)
	}
	if !strings.HasPrefix(res.LankaQREMV, "000201010212") {
		t.Fatalf("expected valid LankaQR EMV header, got: %s", res.LankaQREMV)
	}
}

func TestProcessCheckoutKoko(t *testing.T) {
	svc := NewService(recommender.NewEngine())
	cinema := models.Cinema{ID: "GMP-REG", Name: "Regal Gampaha", City: "Gampaha"}

	items := []models.ConcessionItem{
		{ID: "1", Name: "Family Box", Price: 3000.00, Category: models.CategoryCombo},
	}

	req := models.CheckoutRequest{
		CinemaID:      "GMP-REG",
		PaymentMethod: models.PaymentMethodKoko,
	}

	res, err := svc.ProcessCheckout(cinema, items, req)
	if err != nil {
		t.Fatalf("unexpected error in checkout: %v", err)
	}

	if !strings.Contains(res.InstallmentNote, "3 interest-free monthly installments of LKR 1000.00") {
		t.Fatalf("unexpected installment note: %s", res.InstallmentNote)
	}
}

func TestProcessCheckoutEmptyCartError(t *testing.T) {
	svc := NewService(recommender.NewEngine())
	cinema := models.Cinema{ID: "KND-KCC", Name: "KCC", City: "Kandy"}

	req := models.CheckoutRequest{
		PaymentMethod: models.PaymentMethodLankaQR,
	}

	_, err := svc.ProcessCheckout(cinema, nil, req)
	if err == nil {
		t.Fatal("expected error when checking out with empty cart")
	}
}
