package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"cinema-recommender/internal/models"
	"cinema-recommender/pkg/recommender"
	"cinema-recommender/pkg/scraper"
)

func setupTestServer() *Server {
	sc := scraper.NewRegionalScraper()
	eng := recommender.NewEngine()
	cfg := Config{Port: ":8080"}
	return NewServer(cfg, sc, eng)
}

func TestHealthEndpoint(t *testing.T) {
	srv := setupTestServer()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	w := httptest.NewRecorder()

	srv.GetHandler().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var body map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if body["status"] != "healthy" {
		t.Fatalf("expected healthy, got %v", body["status"])
	}
}

func TestRootEndpointServesUI(t *testing.T) {
	srv := setupTestServer()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	srv.GetHandler().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	if !bytes.Contains(w.Body.Bytes(), []byte("CineBite")) {
		t.Fatal("expected HTML response to contain CineBite branding")
	}
}

func TestMobileBrowserDetectionAndRedirect(t *testing.T) {
	srv := setupTestServer()

	// 1. Mobile browser visiting / should redirect (302) to /mobile
	reqMobile := httptest.NewRequest(http.MethodGet, "/", nil)
	reqMobile.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/604.1")
	wMobile := httptest.NewRecorder()
	srv.GetHandler().ServeHTTP(wMobile, reqMobile)

	if wMobile.Code != http.StatusFound {
		t.Fatalf("expected mobile user-agent to redirect with status 302, got %d", wMobile.Code)
	}
	if loc := wMobile.Header().Get("Location"); loc != "/mobile" {
		t.Fatalf("expected redirect to /mobile, got %s", loc)
	}

	// 2. Android mobile browser visiting / should also redirect
	reqAndroid := httptest.NewRequest(http.MethodGet, "/", nil)
	reqAndroid.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 14; Pixel 8) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36")
	wAndroid := httptest.NewRecorder()
	srv.GetHandler().ServeHTTP(wAndroid, reqAndroid)

	if wAndroid.Code != http.StatusFound || wAndroid.Header().Get("Location") != "/mobile" {
		t.Fatalf("expected Android mobile redirect to /mobile, got status %d, loc: %s", wAndroid.Code, wAndroid.Header().Get("Location"))
	}

	// 3. Mobile browser with ?view=desktop should NOT redirect, serving desktop index.html (200)
	reqOverride := httptest.NewRequest(http.MethodGet, "/?view=desktop", nil)
	reqOverride.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/605.1.15 Mobile/15E148 Safari/604.1")
	wOverride := httptest.NewRecorder()
	srv.GetHandler().ServeHTTP(wOverride, reqOverride)

	if wOverride.Code != http.StatusOK {
		t.Fatalf("expected ?view=desktop to serve status 200, got %d", wOverride.Code)
	}

	// 4. Desktop browser should serve index.html directly (200)
	reqDesktop := httptest.NewRequest(http.MethodGet, "/", nil)
	reqDesktop.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Safari/605.1.15")
	wDesktop := httptest.NewRecorder()
	srv.GetHandler().ServeHTTP(wDesktop, reqDesktop)

	if wDesktop.Code != http.StatusOK {
		t.Fatalf("expected desktop browser to receive status 200, got %d", wDesktop.Code)
	}
}

func TestMobileAndPWAEndpoints(t *testing.T) {
	srv := setupTestServer()

	// 1. Mobile HTML view
	reqMob := httptest.NewRequest(http.MethodGet, "/mobile", nil)
	wMob := httptest.NewRecorder()
	srv.GetHandler().ServeHTTP(wMob, reqMob)
	if wMob.Code != http.StatusOK {
		t.Fatalf("expected /mobile status 200, got %d", wMob.Code)
	}
	if !bytes.Contains(wMob.Body.Bytes(), []byte("CineBite")) {
		t.Fatal("expected /mobile to contain CineBite")
	}

	// 2. Manifest.json
	reqMan := httptest.NewRequest(http.MethodGet, "/manifest.json", nil)
	wMan := httptest.NewRecorder()
	srv.GetHandler().ServeHTTP(wMan, reqMan)
	if wMan.Code != http.StatusOK {
		t.Fatalf("expected /manifest.json status 200, got %d", wMan.Code)
	}

	// 3. Service Worker sw.js
	reqSW := httptest.NewRequest(http.MethodGet, "/sw.js", nil)
	wSW := httptest.NewRecorder()
	srv.GetHandler().ServeHTTP(wSW, reqSW)
	if wSW.Code != http.StatusOK {
		t.Fatalf("expected /sw.js status 200, got %d", wSW.Code)
	}

	// 4. Favicon / App Icon
	reqFav := httptest.NewRequest(http.MethodGet, "/favicon.ico", nil)
	wFav := httptest.NewRecorder()
	srv.GetHandler().ServeHTTP(wFav, reqFav)
	if wFav.Code != http.StatusOK {
		t.Fatalf("expected /favicon.ico status 200, got %d", wFav.Code)
	}
}

func TestListCinemas(t *testing.T) {
	srv := setupTestServer()

	// 1. List all regional cinemas
	req := httptest.NewRequest(http.MethodGet, "/api/v1/cinemas", nil)
	w := httptest.NewRecorder()
	srv.GetHandler().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var cinemas []models.Cinema
	if err := json.Unmarshal(w.Body.Bytes(), &cinemas); err != nil {
		t.Fatalf("failed to parse cinemas: %v", err)
	}
	if len(cinemas) == 0 {
		t.Fatal("expected non-empty list of cinemas")
	}

	// 2. Filter by city: Gampaha
	reqCity := httptest.NewRequest(http.MethodGet, "/api/v1/cinemas?city=Gampaha", nil)
	wCity := httptest.NewRecorder()
	srv.GetHandler().ServeHTTP(wCity, reqCity)

	if wCity.Code != http.StatusOK {
		t.Fatalf("expected status 200 for Gampaha, got %d", wCity.Code)
	}

	// 3. Filter by Colombo -> Returns 3 Colombo Scope Cinemas (CCC, HCM, LBT)
	reqCol := httptest.NewRequest(http.MethodGet, "/api/v1/cinemas?city=Colombo", nil)
	wCol := httptest.NewRecorder()
	srv.GetHandler().ServeHTTP(wCol, reqCol)

	if wCol.Code != http.StatusOK {
		t.Fatalf("expected status 200 for Colombo filter, got %d", wCol.Code)
	}
	var colCinemas []models.Cinema
	if err := json.Unmarshal(wCol.Body.Bytes(), &colCinemas); err != nil {
		t.Fatalf("failed to decode colombo cinemas: %v", err)
	}
	if len(colCinemas) != 3 {
		t.Fatalf("expected 3 Colombo Scope Cinemas, got %d", len(colCinemas))
	}

	// Rejection for unknown / unlisted city -> 404 Not Found
	reqInv := httptest.NewRequest(http.MethodGet, "/api/v1/cinemas?city=UnknownCity", nil)
	wInv := httptest.NewRecorder()
	srv.GetHandler().ServeHTTP(wInv, reqInv)
	if wInv.Code != http.StatusNotFound {
		t.Fatalf("expected status 404 for unknown city, got %d", wInv.Code)
	}

	// 4. Auto-fetch by GPS coordinates (near Kandy: 7.294, 80.639)
	reqGeo := httptest.NewRequest(http.MethodGet, "/api/v1/cinemas?lat=7.2940&lng=80.6390", nil)
	wGeo := httptest.NewRecorder()
	srv.GetHandler().ServeHTTP(wGeo, reqGeo)

	if wGeo.Code != http.StatusOK {
		t.Fatalf("expected status 200 for geo auto-fetch, got %d", wGeo.Code)
	}

	var geoCinemas []models.Cinema
	if err := json.Unmarshal(wGeo.Body.Bytes(), &geoCinemas); err != nil {
		t.Fatalf("failed to decode geo cinemas: %v", err)
	}
	if len(geoCinemas) == 0 || geoCinemas[0].ID != "KND-KCC" {
		t.Fatalf("expected closest cinema to be KND-KCC, got %v", geoCinemas[0].ID)
	}
	if geoCinemas[0].FoodPlaceName == "" || !geoCinemas[0].HasInHouseFood {
		t.Fatalf("expected food place details for %s", geoCinemas[0].Name)
	}
}

func TestGetConcessions(t *testing.T) {
	srv := setupTestServer()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/cinemas/KND-KCC/concessions", nil)
	w := httptest.NewRecorder()
	srv.GetHandler().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var items []models.ConcessionItem
	if err := json.Unmarshal(w.Body.Bytes(), &items); err != nil {
		t.Fatalf("failed to parse concessions: %v", err)
	}
	if len(items) == 0 {
		t.Fatal("expected concession items for KCC Multiplex")
	}

	// Verify all concession food items have discount metadata and food location attached
	for _, it := range items {
		if !it.HasDiscount || len(it.ApplicableDiscounts) == 0 {
			t.Fatalf("expected item %s to have discount breakdowns", it.Name)
		}
		if it.BestDiscountedPrice >= it.Price {
			t.Fatalf("expected discounted price to be lower than base price for %s", it.Name)
		}
		if it.FoodLocation == "" {
			t.Fatalf("expected food location on item %s", it.Name)
		}
	}
}

func TestListPromotions(t *testing.T) {
	srv := setupTestServer()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/promotions/privilege", nil)
	w := httptest.NewRecorder()
	srv.GetHandler().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var promos []models.ConcessionPromotion
	if err := json.Unmarshal(w.Body.Bytes(), &promos); err != nil {
		t.Fatalf("failed to parse promotions: %v", err)
	}
	if len(promos) == 0 {
		t.Fatal("expected registered Scope Privilege promotions")
	}
}

func TestRecommendEndpoint(t *testing.T) {
	srv := setupTestServer()

	// Successful regional recommendation
	recPayload := models.RecommendationRequest{
		CinemaID:      "KND-KCC",
		BudgetLKR:     2500,
		PartySize:     2,
		PrivilegeTier: models.TierStandard,
	}
	body, _ := json.Marshal(recPayload)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/recommend", bytes.NewReader(body))
	w := httptest.NewRecorder()

	srv.GetHandler().ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d (body: %s)", w.Code, w.Body.String())
	}

	var res models.RecommendationResult
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to parse recommendation result: %v", err)
	}
	if len(res.SelectedItems) == 0 {
		t.Fatal("expected selected items in recommendation")
	}

	// Colombo recommendation is now supported for Scope Cinemas
	colomboPayload := models.RecommendationRequest{
		City:      "Colombo",
		BudgetLKR: 2000,
		PartySize: 2,
	}
	bodyCol, _ := json.Marshal(colomboPayload)
	reqCol := httptest.NewRequest(http.MethodPost, "/api/v1/recommend", bytes.NewReader(bodyCol))
	wCol := httptest.NewRecorder()

	srv.GetHandler().ServeHTTP(wCol, reqCol)
	if wCol.Code != http.StatusOK {
		t.Fatalf("expected status 200 for Colombo recommendation, got %d", wCol.Code)
	}

	// Rejection for unknown / unlisted city -> 404 Not Found
	invPayload := models.RecommendationRequest{
		City:      "UnknownCity",
		BudgetLKR: 2000,
		PartySize: 2,
	}
	bodyInv, _ := json.Marshal(invPayload)
	reqInv := httptest.NewRequest(http.MethodPost, "/api/v1/recommend", bytes.NewReader(bodyInv))
	wInv := httptest.NewRecorder()

	srv.GetHandler().ServeHTTP(wInv, reqInv)
	if wInv.Code != http.StatusNotFound {
		t.Fatalf("expected status 404 for invalid city recommendation, got %d", wInv.Code)
	}
}

func TestCartEvaluateEndpoint(t *testing.T) {
	srv := setupTestServer()

	cartReq := models.CartEvaluationRequest{
		CinemaID: "KND-KCC",
		Items: []models.CartItemInput{
			{ItemID: "KND-01", Quantity: 2}, // Jumbo Popcorn 1200 * 2 = 2400 -> qualifies for SCOPE-CINEMA25 (25%)
		},
		PrivilegeTier: models.TierStandard,
	}
	body, _ := json.Marshal(cartReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/cart/evaluate", bytes.NewReader(body))
	w := httptest.NewRecorder()

	srv.GetHandler().ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var res models.RecommendationResult
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to parse cart evaluation: %v", err)
	}

	if res.AppliedPromo != "CINE-CINEMA25" {
		t.Fatalf("expected CINE-CINEMA25, got %s", res.AppliedPromo)
	}
	expectedDiscount := res.OriginalTotal * 0.25
	if res.DiscountAmount != expectedDiscount {
		t.Fatalf("expected discount 25%% of %.2f (LKR %.2f), got LKR %.2f", res.OriginalTotal, expectedDiscount, res.DiscountAmount)
	}
}

func TestCORSOptionsPreflight(t *testing.T) {
	srv := setupTestServer()

	req := httptest.NewRequest(http.MethodOptions, "/api/v1/cinemas", nil)
	w := httptest.NewRecorder()
	srv.GetHandler().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 on OPTIONS preflight, got %d", w.Code)
	}
	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Fatal("expected Access-Control-Allow-Origin header")
	}
}

func TestListPaymentMethodsEndpoint(t *testing.T) {
	srv := setupTestServer()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/payment-methods", nil)
	w := httptest.NewRecorder()
	srv.GetHandler().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var methods []models.PaymentMethodInfo
	if err := json.Unmarshal(w.Body.Bytes(), &methods); err != nil {
		t.Fatalf("failed to decode payment methods: %v", err)
	}
	if len(methods) < 7 {
		t.Fatalf("expected at least 7 payment methods, got %d", len(methods))
	}
}

func TestCheckoutEndpoint(t *testing.T) {
	srv := setupTestServer()

	// 1. Checkout via LankaQR
	qrReq := models.CheckoutRequest{
		CinemaID:      "KND-KCC",
		PaymentMethod: models.PaymentMethodLankaQR,
		Items: []models.CartItemInput{
			{ItemID: "KND-01", Quantity: 1}, // Jumbo Popcorn 1200
		},
	}
	body, _ := json.Marshal(qrReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/checkout", bytes.NewReader(body))
	w := httptest.NewRecorder()

	srv.GetHandler().ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 on checkout, got %d (body: %s)", w.Code, w.Body.String())
	}

	var res models.CheckoutResult
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to parse checkout result: %v", err)
	}

	if res.PaymentMethod != models.PaymentMethodLankaQR {
		t.Fatalf("expected LankaQR payment method, got %s", res.PaymentMethod)
	}
	if res.Status != "PENDING_LANKAQR_SCAN" {
		t.Fatalf("expected PENDING_LANKAQR_SCAN status, got %s", res.Status)
	}
	if res.LankaQREMV == "" {
		t.Fatal("expected LankaQR EMV payload to be populated")
	}
}

func TestLiveConcessionsEndpoint(t *testing.T) {
	srv := setupTestServer()

	// 1. Check standard concessions headers for live updates
	reqStd := httptest.NewRequest(http.MethodGet, "/api/v1/cinemas/KND-KCC/concessions", nil)
	wStd := httptest.NewRecorder()
	srv.GetHandler().ServeHTTP(wStd, reqStd)
	if wStd.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", wStd.Code)
	}
	if wStd.Header().Get("X-Live-Price-Feed") != "active" {
		t.Fatalf("expected X-Live-Price-Feed: active header")
	}

	// 2. Check /live endpoint
	reqLive := httptest.NewRequest(http.MethodGet, "/api/v1/cinemas/KND-KCC/concessions/live", nil)
	wLive := httptest.NewRecorder()
	srv.GetHandler().ServeHTTP(wLive, reqLive)
	if wLive.Code != http.StatusOK {
		t.Fatalf("expected status 200 on /live, got %d", wLive.Code)
	}

	var liveData map[string]interface{}
	if err := json.Unmarshal(wLive.Body.Bytes(), &liveData); err != nil {
		t.Fatalf("failed to decode /live response: %v", err)
	}

	if liveData["live_pricing_active"] != true {
		t.Fatalf("expected live_pricing_active=true")
	}
	if liveData["items"] == nil {
		t.Fatalf("expected items array in live response")
	}
	if liveData["cinema_ads"] == nil {
		t.Fatalf("expected cinema_ads in live response")
	}
}

func TestCinemaAdsEndpoint(t *testing.T) {
	srv := setupTestServer()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/cinemas/KND-KCC/ads", nil)
	w := httptest.NewRecorder()
	srv.GetHandler().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 on /ads, got %d", w.Code)
	}

	var ads []models.CinemaAd
	if err := json.Unmarshal(w.Body.Bytes(), &ads); err != nil {
		t.Fatalf("failed to decode ads: %v", err)
	}

	if len(ads) == 0 {
		t.Fatal("expected promotional cinema ads for KCC Multiplex")
	}

	var hasComboAd bool
	for _, ad := range ads {
		if ad.TargetItemID != "" && ad.PromoCode != "" {
			hasComboAd = true
		}
		if ad.Title == "" || ad.DiscountText == "" || ad.CallToAction == "" {
			t.Errorf("ad missing essential fields: %+v", ad)
		}
	}

	if !hasComboAd {
		t.Fatal("expected at least one actionable combo/discount ad with target item and promo code")
	}
}

