package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cinema-recommender/internal/models"
	"cinema-recommender/internal/server"
	"cinema-recommender/pkg/recommender"
	"cinema-recommender/pkg/scraper"
)

func main() {
	var (
		serverMode  = flag.Bool("server", false, "Start the HTTP REST API & Web UI server")
		port        = flag.String("port", ":8080", "Port for HTTP server (e.g. :8080)")
		city        = flag.String("city", "Kandy", "Target regional cinema city (outside Colombo)")
		budget      = flag.Float64("budget", 2500.00, "Concession budget in LKR")
		partySize   = flag.Int("party", 2, "Number of moviegoers in the party")
		tier        = flag.String("tier", "Standard", "Redopay cardholder tier (Standard, Student, Silver, Gold, Platinum)")
		promoCode   = flag.String("promo", "", "Specific Redopay promo code override")
		listCinemas = flag.Bool("list-cinemas", false, "List all registered regional out-of-Colombo cinemas")
		listPromos  = flag.Bool("list-promos", false, "List all active Redopay concession promotions")
	)
	flag.Parse()

	printBanner()

	sc := scraper.NewRegionalScraper()
	recEngine := recommender.NewEngine()

	// 1. List Cinemas CLI flag
	if *listCinemas {
		displayCinemas(sc)
		return
	}

	// 2. List Promotions CLI flag
	if *listPromos {
		displayPromotions(recEngine)
		return
	}

	// 3. Server Mode
	if *serverMode {
		runHTTPServer(*port, sc, recEngine)
		return
	}

	// 4. Default Interactive CLI Recommendation Run
	runCLIRecommendation(*city, *budget, *partySize, *tier, *promoCode, sc, recEngine)
}

func printBanner() {
	fmt.Println("================================================================================")
	fmt.Println("🎬 CINEMA FOOD & BEVERAGE RECOMMENDATION ENGINE 🍿🥤")
	fmt.Println("   Regional Outstation Scraping Pipeline & Redopay Discount Optimization Engine")
	fmt.Println("================================================================================")
}

func runCLIRecommendation(
	cityName string,
	budget float64,
	partySize int,
	tierName string,
	promoCode string,
	sc *scraper.RegionalScraper,
	recEngine *recommender.Engine,
) {
	fmt.Printf("\n🔍 Querying Regional Theaters in: %s\n", cityName)

	// Validate location (must be outside Colombo)
	if err := sc.ValidateLocation(cityName); err != nil {
		fmt.Printf("❌ Rejection Policy Triggered: %v\n", err)
		fmt.Println("💡 Reminder: This service strictly targets outstation circuits (Kandy, Gampaha, Galle, Kurunegala, Negombo, etc.)")
		return
	}

	cinemas, err := sc.GetCinemasByCity(cityName)
	if err != nil || len(cinemas) == 0 {
		fmt.Printf("❌ No registered theaters found for city %q: %v\n", cityName, err)
		return
	}

	targetCinema := cinemas[0]
	fmt.Printf("📍 Target Theater: %s (%s, %s)\n", targetCinema.Name, targetCinema.City, targetCinema.Province)

	// Scrape concessions
	items, err := sc.ScrapeCinemaConcessions(targetCinema)
	if err != nil {
		fmt.Printf("❌ Scraping error: %v\n", err)
		return
	}
	fmt.Printf("✅ Scraped %d live concession items.\n", len(items))

	// Run optimization
	req := models.RecommendationRequest{
		CinemaID:    targetCinema.ID,
		City:        targetCinema.City,
		BudgetLKR:   budget,
		PartySize:   partySize,
		RedopayTier: models.RedopayTier(tierName),
		PromoCode:   promoCode,
	}

	result, err := recEngine.RecommendOptimalBundle(targetCinema, items, req)
	if err != nil {
		fmt.Printf("❌ Recommendation failed: %v\n", err)
		return
	}

	fmt.Println("\n--------------------------------------------------------------------------------")
	fmt.Printf("🍿 VALUE-OPTIMIZED CONCESSION BUNDLE (Party of %d | Budget: LKR %.2f)\n", partySize, budget)
	fmt.Println("--------------------------------------------------------------------------------")
	for idx, it := range result.SelectedItems {
		fmt.Printf("   [%d] %-32s | %-10s | LKR %7.2f\n", idx+1, it.Name, it.Category, it.Price)
	}

	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Printf("   Original Menu Subtotal:     LKR %8.2f\n", result.OriginalTotal)
	fmt.Printf("   Redopay Concession Savings: -LKR %8.2f (%s)\n", result.DiscountAmount, result.AppliedPromo)
	fmt.Printf("   👉 Net Payable via Redopay:  LKR %8.2f\n", result.FinalTotal)
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Printf("   💡 %s\n", result.Message)
	if result.NextTierUpsell != "" {
		fmt.Printf("   %s\n", result.NextTierUpsell)
	}
	fmt.Println("--------------------------------------------------------------------------------")

	fmt.Println("\n🚀 Try running as an HTTP Server + Interactive Web UI:")
	fmt.Println("   go run ./cmd/api --server --port :8080")
	fmt.Println("   Then open: http://localhost:8080")
}

func displayCinemas(sc *scraper.RegionalScraper) {
	cinemas := sc.GetRegisteredCinemas()
	fmt.Printf("\n🏛️  Registered Out-of-Colombo Regional Theaters (%d total):\n", len(cinemas))
	fmt.Printf("%-10s | %-32s | %-14s | %-22s\n", "ID", "Cinema Name", "City", "Province")
	fmt.Println("--------------------------------------------------------------------------------")
	for _, c := range cinemas {
		fmt.Printf("%-10s | %-32s | %-14s | %-22s\n", c.ID, c.Name, c.City, c.Province)
	}
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println("ℹ️  Colombo metropolitan theaters are deliberately excluded as per system scope.")
}

func displayPromotions(recEngine *recommender.Engine) {
	promos := recEngine.GetAvailablePromotions()
	fmt.Printf("\n💳 Active Redopay Concession Promotions (%d campaigns):\n", len(promos))
	fmt.Printf("%-18s | %-8s | %-14s | %s\n", "Promo Code", "Type", "Min Spend", "Description")
	fmt.Println("--------------------------------------------------------------------------------")
	for _, p := range promos {
		minSpendStr := fmt.Sprintf("LKR %.0f", p.MinSpendLKR)
		fmt.Printf("%-18s | %-8s | %-14s | %s\n", p.PromoCode, p.DiscountType, minSpendStr, p.Description)
	}
	fmt.Println("--------------------------------------------------------------------------------")
}

func runHTTPServer(port string, sc *scraper.RegionalScraper, recEngine *recommender.Engine) {
	cfg := server.Config{
		Port:         port,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	srv := server.NewServer(cfg, sc, recEngine)

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		fmt.Printf("\n🌐 HTTP REST API & Web UI running at http://localhost%s\n", port)
		fmt.Println("   Endpoints:")
		fmt.Println("   - Dashboard UI:       GET  http://localhost" + port + "/")
		fmt.Println("   - Health Check:       GET  http://localhost" + port + "/api/v1/health")
		fmt.Println("   - List Cinemas:       GET  http://localhost" + port + "/api/v1/cinemas")
		fmt.Println("   - Scrape Concessions: GET  http://localhost" + port + "/api/v1/cinemas/{id}/concessions")
		fmt.Println("   - Redopay Promos:     GET  http://localhost" + port + "/api/v1/promotions/redopay")
		fmt.Println("   - AI Recommendation:  POST http://localhost" + port + "/api/v1/recommend")
		fmt.Println("   - Cart Evaluation:    POST http://localhost" + port + "/api/v1/cart/evaluate")
		fmt.Println("\n(Press Ctrl+C to gracefully shutdown)")

		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ HTTP server error: %v", err)
		}
	}()

	<-stopChan
	fmt.Println("\n🛑 Shutting down HTTP server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("❌ Server forced to shutdown: %v", err)
	}

	fmt.Println("✅ Server exited cleanly.")
}
