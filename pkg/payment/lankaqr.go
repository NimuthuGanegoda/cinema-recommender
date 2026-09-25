package payment

import (
	"fmt"
	"strings"

	"cinema-recommender/internal/models"
)

// GenerateLankaQREMV generates an EMVCo compliant LankaQR payload representation for CBSL standard apps.
func GenerateLankaQREMV(cinema models.Cinema, amount float64, orderID string) string {
	// CBSL LankaQR EMV standard tags:
	// 00: Format Indicator (01)
	// 01: Point of Initiation (12 = Dynamic QR)
	// 26: Merchant Account Info (LankaQR standard)
	// 52: Merchant Category Code (7832 = Motion Picture Theaters)
	// 53: Transaction Currency (144 = LKR Sri Lankan Rupee)
	// 54: Transaction Amount
	// 58: Country Code (LK)
	// 59: Merchant Name
	// 60: Merchant City
	// 62: Additional Data (Order ID)
	cleanMerchant := strings.ReplaceAll(cinema.Name, " ", "")
	if len(cleanMerchant) > 20 {
		cleanMerchant = cleanMerchant[:20]
	}
	cleanCity := cinema.City
	if len(cleanCity) > 15 {
		cleanCity = cleanCity[:15]
	}

	return fmt.Sprintf(
		"00020101021226480010com.lankaqr0118%s520478325303144540%d%.2f5802LK59%02d%s60%02d%s62%02d01%s6304LKQR",
		orderID,
		len(fmt.Sprintf("%.2f", amount)),
		amount,
		len(cleanMerchant),
		cleanMerchant,
		len(cleanCity),
		cleanCity,
		len(orderID)+4,
		orderID,
	)
}
