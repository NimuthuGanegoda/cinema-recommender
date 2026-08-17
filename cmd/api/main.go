package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"cinema-recommender/internal/models"
	"cinema-recommender/pkg/recommender"
	"cinema-recommender/pkg/scraper"
)

func main() {
	fmt.Println("===============================================================")
	fmt.Println("⚡ CINEMA FOOD & BEVERAGE RECOMMENDATION ENGINE 🍿🥤")
	fmt.Println("   Regional Outstation Scraping & Redopay Optimization Engine")
	fmt.Println("===============================================================\n")

	// 1. Initialize core engines
	sc := scraper.NewRegionalScraper()
	recEngine := recommender.NewEngine()

	// 2. Define test cinemas
	testCinemas := []models.Cinema{
		{ID: "C-01", Name: "Kandy City Centre Cinema", City: "Kandy", IsOutside: true},
		{ID: "C-02", Name: "Regal Cinema Gampaha", City: "Gampaha", IsOutside: true},
		{ID: "C-03", Name: "Liberty by Scope (Colombo)", City: "Colombo", IsOutside: false},
	}

	// 3. Process each cinema through the scraper & recommender
	for _, cinema := range testCinemas {
		fmt.Printf("🔍 Processing: %s (%s)...\n", cinema.Name, cinema.City)

		// Call Scraper
		items, err := sc.ScrapeCinemaConcessions(cinema)
		if err != nil {
			fmt.Printf("   ❌ Skipped / Error: %v\n\n", err)
			continue
		}

		fmt.Printf("   ✅ Scraped %d concession items successfully.\n", len(items))

		// Call Recommender with Redopay optimization
		deal := recEngine.CalculateBestDeal(cinema, items)
		fmt.Println("   🏷️  " + deal.String())
		fmt.Printf("   💡 %s\n\n", deal.Message)
	}

	// If argument "--server" is provided, start HTTP API server
	if len(os.Args) > 1 && os.Args[1] == "--server" {
		startHTTPServer(sc, recEngine)
	} else {
		fmt.Println("💡 Tip: Run with `go run ./cmd/api --server` to start the HTTP REST API server on :8080")
	}
}

// startHTTPServer demonstrates idiomatic Go standard library HTTP routing
func startHTTPServer(sc *scraper.RegionalScraper, recEngine *recommender.Engine) {
	port := ":8080"
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy", "service": "cinema-recommender"})
	})

	http.HandleFunc("/api/recommend", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		city := r.URL.Query().Get("city")
		if city == "" {
			city = "Kandy"
		}

		cinema := models.Cinema{ID: "C-HTTP", Name: fmt.Sprintf("Cinema %s", city), City: city}
		items, err := sc.ScrapeCinemaConcessions(cinema)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		deal := recEngine.CalculateBestDeal(cinema, items)
		json.NewEncoder(w).Encode(deal)
	})

	fmt.Printf("🌐 HTTP API listening on http://localhost%s (Press Ctrl+C to stop)\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
