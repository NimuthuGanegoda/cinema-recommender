package payment

import (
	"fmt"
	"math"
)

// BNPLInstallmentPlan describes a 3-part split payment calculation for Koko or Mintpay.
type BNPLInstallmentPlan struct {
	Provider        string  `json:"provider"`
	Installments    int     `json:"installments"`
	MonthlyAmount   float64 `json:"monthly_amount_lkr"`
	InstallmentNote string  `json:"installment_note"`
	Instructions    string  `json:"instructions"`
}

// CalculateBNPLPlan computes the installment breakdown for Sri Lankan BNPL gateways.
func CalculateBNPLPlan(provider string, finalPayable float64) BNPLInstallmentPlan {
	const installments = 3
	monthly := math.Round((finalPayable/float64(installments))*100) / 100

	switch provider {
	case "Koko":
		return BNPLInstallmentPlan{
			Provider:        "Koko",
			Installments:    installments,
			MonthlyAmount:   monthly,
			InstallmentNote: fmt.Sprintf("Pay in %d interest-free monthly installments of LKR %.2f with Koko", installments, monthly),
			Instructions:    "Koko BNPL schedule created. First installment charged today, remaining 2 installments due over the next 60 days.",
		}
	default: // Mintpay
		return BNPLInstallmentPlan{
			Provider:        "Mintpay",
			Installments:    installments,
			MonthlyAmount:   monthly,
			InstallmentNote: fmt.Sprintf("Pay in %d interest-free debit card installments of LKR %.2f with Mintpay", installments, monthly),
			Instructions:    "Mintpay BNPL installment plan activated. Split across 3 monthly debit card payments with zero interest.",
		}
	}
}
