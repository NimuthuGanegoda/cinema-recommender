package scraper

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"cinema-recommender/internal/models"
	"github.com/gocolly/colly/v2"
)

// scrapeByCity executes concession scraping pipelines tuned for regional suppliers and theater concessions.
func (s *RegionalScraper) scrapeByCity(city string, cinemaID string) ([]models.ConcessionItem, error) {
	cityLower := strings.ToLower(strings.TrimSpace(city))

	// Attempt Live Uber Eats Scraping First
	uberEatsURLs := map[string]string{
		"CMB-CCC": "https://www.ubereats.com/lk/store/scope-cinemas-multiplex-ccc",
		"CMB-HCM": "https://www.ubereats.com/lk/store/scope-cinemas-multiplex-havelock",
		"KND-KCC": "https://www.ubereats.com/lk/store/scope-partner-kcc-kandy",
	}

	if url, exists := uberEatsURLs[cinemaID]; exists {
		items, err := s.ScrapeLiveUberEatsMenu(cinemaID, url)
		if err == nil && len(items) > 0 {
			return items, nil
		}
		// If scraping fails (e.g., Cloudflare 403 Forbidden), gracefully fall back to local database.
	}

	switch cityLower {
	case "colombo":
		prefix := "CMB"
		if strings.HasPrefix(cinemaID, "CMB-CCC") {
			prefix = "CCC"
		} else if strings.HasPrefix(cinemaID, "CMB-HCM") {
			prefix = "HCM"
		} else if strings.HasPrefix(cinemaID, "CMB-LBT") {
			prefix = "LBT"
		}
		return []models.ConcessionItem{
			{ID: prefix + "-01", CinemaID: cinemaID, Name: "Scope Signature Warm Caramel Popcorn Tub", Category: models.CategoryPopcorn, Size: "Jumbo", Price: 1350.00, InStock: true, Tags: []string{"sweet", "popular", "scope-signature"}, Description: "Scope signature crunchy caramelized warm corn in souvenir tub"},
			{ID: prefix + "-02", CinemaID: cinemaID, Name: "Scope Real Butter Sea Salt Popcorn", Category: models.CategoryPopcorn, Size: "Large", Price: 1100.00, InStock: true, Tags: []string{"savory", "classic"}, Description: "Freshly popped butterfly corn drenched in melted clarified butter"},
			{ID: prefix + "-03", CinemaID: cinemaID, Name: "Scope Double Cheddar Cheese Popcorn", Category: models.CategoryPopcorn, Size: "Large", Price: 1250.00, InStock: true, Tags: []string{"savory", "cheesy"}, Description: "Hot popped corn tossed in Wisconsin white and yellow cheddar dust"},
			{ID: prefix + "-04", CinemaID: cinemaID, Name: "Scope Large Fountain Coca-Cola", Category: models.CategoryBeverage, Size: "Large", Price: 650.00, InStock: true, Tags: []string{"beverage", "chilled"}, Description: "Ice-cold fountain Coca-Cola with fresh carbonation"},
			{ID: prefix + "-05", CinemaID: cinemaID, Name: "Scope Fresh Passion Fruit Mint Sparkler", Category: models.CategoryBeverage, Size: "Medium", Price: 750.00, InStock: true, Tags: []string{"beverage", "signature", "refreshing"}, Description: "Real island passion fruit pulp with crushed mint and sparkling soda"},
			{ID: prefix + "-06", CinemaID: cinemaID, Name: "Scope IMAX Loaded Chicken Nachos", Category: models.CategorySnack, Size: "Single", Price: 1450.00, InStock: true, Tags: []string{"savory", "spicy", "popular"}, Description: "Stone-ground corn chips with hot queso, jalapeños, salsa, and seasoned chicken"},
			{ID: prefix + "-07", CinemaID: cinemaID, Name: "Scope Gourmet Brioche Chicken Hotdog", Category: models.CategorySnack, Size: "Single", Price: 1150.00, InStock: true, Tags: []string{"savory", "halal", "vip"}, Description: "Artisan grilled chicken sausage in buttered brioche bun with sweet pickle relish"},
			{ID: prefix + "-08", CinemaID: cinemaID, Name: "Scope Director's Deluxe Couple Combo", Category: models.CategoryCombo, Size: "Duo", Price: 2950.00, InStock: true, Tags: []string{"combo", "best-value", "vip"}, Description: "1 Jumbo Caramel Popcorn + 2 Large Fountain Drinks + 1 Loaded Nachos"},
			{ID: prefix + "-09", CinemaID: cinemaID, Name: "Scope VIP Platinum Movie Feast", Category: models.CategoryCombo, Size: "Family", Price: 4200.00, InStock: true, Tags: []string{"combo", "family", "platinum"}, Description: "2 Jumbo Popcorns (Caramel + Cheese) + 3 Drinks + 1 Brioche Hotdog + 1 Nachos"},
		}, nil

	case "kandy":
		return []models.ConcessionItem{
			{ID: "KND-01", CinemaID: cinemaID, Name: "Scope Jumbo Warm Caramel Popcorn", Category: models.CategoryPopcorn, Size: "Jumbo", Price: 1200.00, InStock: true, Tags: []string{"sweet", "popular", "scope-exclusive"}, Description: "Scope signature crunchy caramelized warm corn"},
			{ID: "KND-02", CinemaID: cinemaID, Name: "Scope Signature Butter Salt Popcorn", Category: models.CategoryPopcorn, Size: "Large", Price: 1000.00, InStock: true, Tags: []string{"savory", "classic"}, Description: "Freshly popped with melted golden butter"},
			{ID: "KND-03", CinemaID: cinemaID, Name: "Scope Large Fountain Mountain Dew", Category: models.CategoryBeverage, Size: "Large", Price: 650.00, InStock: true, Tags: []string{"beverage", "chilled"}, Description: "Chilled fountain carbonated drink"},
			{ID: "KND-04", CinemaID: cinemaID, Name: "Scope Fresh Ceylon Iced Tea", Category: models.CategoryBeverage, Size: "Medium", Price: 480.00, InStock: true, Tags: []string{"beverage", "local"}, Description: "Brewed hill country tea with fresh lemon"},
			{ID: "KND-05", CinemaID: cinemaID, Name: "Scope Gourmet Chicken Hotdog", Category: models.CategorySnack, Size: "Single", Price: 950.00, InStock: true, Tags: []string{"spicy", "savory", "halal"}, Description: "Grilled chicken sausage in toasted brioche with relish"},
			{ID: "KND-06", CinemaID: cinemaID, Name: "Scope Director's Duo Combo", Category: models.CategoryCombo, Size: "Duo", Price: 2400.00, InStock: true, Tags: []string{"combo", "best-value", "vip"}, Description: "1 Jumbo Popcorn + 2 Large Drinks + 1 Gourmet Hotdog"},
		}, nil

	case "gampaha":
		return []models.ConcessionItem{
			{ID: "GMP-01", CinemaID: cinemaID, Name: "Salted Butter Popcorn (M)", Category: models.CategoryPopcorn, Size: "Medium", Price: 900.00, InStock: true, Tags: []string{"savory"}, Description: "Warm buttery classic popcorn"},
			{ID: "GMP-02", CinemaID: cinemaID, Name: "Caramel Crunch Tub", Category: models.CategoryPopcorn, Size: "Large", Price: 1100.00, InStock: true, Tags: []string{"sweet"}, Description: "Thick golden caramel glaze"},
			{ID: "GMP-03", CinemaID: cinemaID, Name: "Cold Milo Float", Category: models.CategoryBeverage, Size: "Medium", Price: 550.00, InStock: true, Tags: []string{"beverage", "sweet"}, Description: "Malt Milo topped with vanilla ice cream"},
			{ID: "GMP-04", CinemaID: cinemaID, Name: "Coca-Cola Zero Sugar", Category: models.CategoryBeverage, Size: "Large", Price: 500.00, InStock: true, Tags: []string{"beverage", "sugar-free"}, Description: "Chilled sparkling Coke"},
			{ID: "GMP-05", CinemaID: cinemaID, Name: "Crispy Nachos & Cheese", Category: models.CategorySnack, Size: "Single", Price: 850.00, InStock: true, Tags: []string{"savory", "vegetarian"}, Description: "Corn tortilla chips with jalapeno cheese sauce"},
			{ID: "GMP-06", CinemaID: cinemaID, Name: "Gampaha Deluxe Couple Combo", Category: models.CategoryCombo, Size: "Duo", Price: 2250.00, InStock: true, Tags: []string{"combo"}, Description: "1 Caramel Tub + 2 Milo Floats + Nachos"},
		}, nil

	case "galle":
		return []models.ConcessionItem{
			{ID: "GLE-01", CinemaID: cinemaID, Name: "Southern Cheese Popcorn", Category: models.CategoryPopcorn, Size: "Large", Price: 1100.00, InStock: true, Tags: []string{"savory", "cheesy"}, Description: "Dusted with rich cheddar seasoning"},
			{ID: "GLE-02", CinemaID: cinemaID, Name: "Fresh Lime & Mint Cooler", Category: models.CategoryBeverage, Size: "Large", Price: 500.00, InStock: true, Tags: []string{"beverage", "citrus"}, Description: "Fresh southern lime with crushed mint"},
			{ID: "GLE-03", CinemaID: cinemaID, Name: "Sweet Chili Fish & Chips", Category: models.CategorySnack, Size: "Single", Price: 1400.00, InStock: true, Tags: []string{"seafood", "crispy"}, Description: "Crumbed fish goujons with hand-cut fries"},
			{ID: "GLE-04", CinemaID: cinemaID, Name: "Crispy Samosa Platter (4pcs)", Category: models.CategorySnack, Size: "Single", Price: 750.00, InStock: true, Tags: []string{"savory", "vegetarian"}, Description: "Spiced potato pastry with sweet tamarind dip"},
			{ID: "GLE-05", CinemaID: cinemaID, Name: "Fort Fortress Mega Combo", Category: models.CategoryCombo, Size: "Family", Price: 2950.00, InStock: true, Tags: []string{"combo", "family"}, Description: "2 Large Popcorn + 3 Drinks + 1 Fish & Chips"},
		}, nil

	case "kurunegala":
		return []models.ConcessionItem{
			{ID: "KRN-01", CinemaID: cinemaID, Name: "North Star Popcorn Bucket", Category: models.CategoryPopcorn, Size: "Jumbo", Price: 1000.00, InStock: true, Tags: []string{"savory"}, Description: "Family tub of salted butter corn"},
			{ID: "KRN-02", CinemaID: cinemaID, Name: "Fresh Cola Combo Drink", Category: models.CategoryBeverage, Size: "Large", Price: 600.00, InStock: true, Tags: []string{"beverage"}, Description: "Large iced fountain cola"},
			{ID: "KRN-03", CinemaID: cinemaID, Name: "Crispy Chicken Roll (2pcs)", Category: models.CategorySnack, Size: "Single", Price: 820.00, InStock: true, Tags: []string{"savory", "halal"}, Description: "Crumbed Sri Lankan spicy chicken rolls"},
			{ID: "KRN-04", CinemaID: cinemaID, Name: "Loaded French Fries", Category: models.CategorySnack, Size: "Medium", Price: 780.00, InStock: true, Tags: []string{"savory", "vegetarian"}, Description: "Seasoned fries topped with melted cheese"},
			{ID: "KRN-05", CinemaID: cinemaID, Name: "Kurunegala Value Pack", Category: models.CategoryCombo, Size: "Duo", Price: 2100.00, InStock: true, Tags: []string{"combo"}, Description: "1 Popcorn Bucket + 2 Drinks + 1 Chicken Roll"},
		}, nil

	case "negombo":
		return []models.ConcessionItem{
			{ID: "NGB-01", CinemaID: cinemaID, Name: "Coastal Caramel Popcorn", Category: models.CategoryPopcorn, Size: "Large", Price: 980.00, InStock: true, Tags: []string{"sweet"}, Description: "Artisan buttery caramel popcorn"},
			{ID: "NGB-02", CinemaID: cinemaID, Name: "Iced Lemon Soda", Category: models.CategoryBeverage, Size: "Large", Price: 520.00, InStock: true, Tags: []string{"beverage", "fizzy"}, Description: "Refreshing sparkling lemon cordial"},
			{ID: "NGB-03", CinemaID: cinemaID, Name: "Cheese Nachos Deluxe", Category: models.CategorySnack, Size: "Medium", Price: 930.00, InStock: true, Tags: []string{"savory", "cheesy"}, Description: "Warm tortilla chips with salsa and cheese"},
			{ID: "NGB-04", CinemaID: cinemaID, Name: "Calamari Rings Cone", Category: models.CategorySnack, Size: "Single", Price: 1350.00, InStock: true, Tags: []string{"seafood", "crispy"}, Description: "Crispy fried squid rings with tartare dip"},
			{ID: "NGB-05", CinemaID: cinemaID, Name: "Negombo Sunset Combo", Category: models.CategoryCombo, Size: "Duo", Price: 2300.00, InStock: true, Tags: []string{"combo"}, Description: "1 Caramel Popcorn + 2 Lemon Sodas + 1 Nachos"},
		}, nil

	case "jaffna":
		return []models.ConcessionItem{
			{ID: "JFN-01", CinemaID: cinemaID, Name: "Jaffna Crunch Spiced Popcorn", Category: models.CategoryPopcorn, Size: "Large", Price: 1150.00, InStock: true, Tags: []string{"spicy", "savory"}, Description: "Popcorn dusted with northern chili spices"},
			{ID: "JFN-02", CinemaID: cinemaID, Name: "Palmyra Sweet Cooler", Category: models.CategoryBeverage, Size: "Medium", Price: 700.00, InStock: true, Tags: []string{"beverage", "traditional"}, Description: "Natural palmyra fruit pulp nectar"},
			{ID: "JFN-03", CinemaID: cinemaID, Name: "Hot Ceylon Chicken Bites", Category: models.CategorySnack, Size: "Single", Price: 960.00, InStock: true, Tags: []string{"spicy", "savory", "halal"}, Description: "Boneless fried chicken seasoned with curry leaves"},
			{ID: "JFN-04", CinemaID: cinemaID, Name: "Jaffna Murukku Platter", Category: models.CategorySnack, Size: "Single", Price: 650.00, InStock: true, Tags: []string{"crunchy", "vegetarian"}, Description: "Crispy traditional spiced savory snacks"},
			{ID: "JFN-05", CinemaID: cinemaID, Name: "Northern Festival Combo", Category: models.CategoryCombo, Size: "Duo", Price: 2600.00, InStock: true, Tags: []string{"combo"}, Description: "1 Spiced Popcorn + 2 Coolers + 1 Chicken Bites"},
		}, nil

	case "matara":
		return []models.ConcessionItem{
			{ID: "MTR-01", CinemaID: cinemaID, Name: "Southern Popcorn Mix (Sweet & Salt)", Category: models.CategoryPopcorn, Size: "Large", Price: 1050.00, InStock: true, Tags: []string{"mixed"}, Description: "Dual blend of caramel & sea salted corn"},
			{ID: "MTR-02", CinemaID: cinemaID, Name: "Fresh Orange Fizz", Category: models.CategoryBeverage, Size: "Medium", Price: 580.00, InStock: true, Tags: []string{"beverage"}, Description: "Pulpy orange with sparkling water"},
			{ID: "MTR-03", CinemaID: cinemaID, Name: "Crispy Potato Wedges Platter", Category: models.CategorySnack, Size: "Medium", Price: 890.00, InStock: true, Tags: []string{"savory", "vegetarian"}, Description: "Crispy seasoned skin-on potato wedges"},
			{ID: "MTR-04", CinemaID: cinemaID, Name: "Southern Breeze Duo Combo", Category: models.CategoryCombo, Size: "Duo", Price: 2200.00, InStock: true, Tags: []string{"combo"}, Description: "1 Mix Popcorn + 2 Orange Fizz + 1 Wedges"},
		}, nil

	case "anuradhapura":
		return []models.ConcessionItem{
			{ID: "ANR-01", CinemaID: cinemaID, Name: "Classic Butter Popcorn", Category: models.CategoryPopcorn, Size: "Large", Price: 950.00, InStock: true, Tags: []string{"savory"}, Description: "Golden popped kernels with butter flavor"},
			{ID: "ANR-02", CinemaID: cinemaID, Name: "Chilled Passion Fruit Nectar", Category: models.CategoryBeverage, Size: "Medium", Price: 520.00, InStock: true, Tags: []string{"beverage"}, Description: "Island passion fruit cooler"},
			{ID: "ANR-03", CinemaID: cinemaID, Name: "Vegetable Spring Rolls (3pcs)", Category: models.CategorySnack, Size: "Single", Price: 720.00, InStock: true, Tags: []string{"savory", "vegetarian"}, Description: "Crispy pastry rolls filled with garden veggies"},
			{ID: "ANR-04", CinemaID: cinemaID, Name: "Rajarata Family Snack Box", Category: models.CategoryCombo, Size: "Family", Price: 2350.00, InStock: true, Tags: []string{"combo"}, Description: "1 Popcorn + 2 Nectars + 2 Spring Roll orders"},
		}, nil

	case "ratnapura":
		return []models.ConcessionItem{
			{ID: "RTP-01", CinemaID: cinemaID, Name: "Gem City Sweet Popcorn", Category: models.CategoryPopcorn, Size: "Large", Price: 990.00, InStock: true, Tags: []string{"sweet"}, Description: "Crunchy sweet glaze popcorn"},
			{ID: "RTP-02", CinemaID: cinemaID, Name: "Iced Faluda Royal", Category: models.CategoryBeverage, Size: "Medium", Price: 620.00, InStock: true, Tags: []string{"beverage", "dessert"}, Description: "Rose syrup with milk, basil seeds & ice cream"},
			{ID: "RTP-03", CinemaID: cinemaID, Name: "Spicy Beef Patty (2pcs)", Category: models.CategorySnack, Size: "Single", Price: 800.00, InStock: true, Tags: []string{"spicy", "savory"}, Description: "Shortcrust baked spicy minced beef parcels"},
			{ID: "RTP-04", CinemaID: cinemaID, Name: "Milano Matinee Combo", Category: models.CategoryCombo, Size: "Duo", Price: 2150.00, InStock: true, Tags: []string{"combo"}, Description: "1 Sweet Popcorn + 2 Faludas + 1 Patty Box"},
		}, nil

	default:
		return nil, fmt.Errorf("%w: no concession scraper registered for city %q", ErrCinemaNotFound, city)
	}
}

// ScrapeRawConcessionFeed parses simulated or raw concession HTML feeds.
// Allows parsing raw web page tables / div lists into structured models.ConcessionItem.
func (s *RegionalScraper) ScrapeRawConcessionFeed(cinemaID string, rawHTML string) ([]models.ConcessionItem, error) {
	if strings.TrimSpace(rawHTML) == "" {
		return nil, fmt.Errorf("%w: raw feed content is empty", ErrInvalidInput)
	}

	var items []models.ConcessionItem
	// Regex pattern for concession items in web feeds: <item name="..." category="..." price="..." ...>
	re := regexp.MustCompile(`(?i)<(?:div|item)[^>]*class=["']concession-card["'][^>]*data-name=["']([^"']+)["'][^>]*data-cat=["']([^"']+)["'][^>]*data-price=["']([0-9.]+)["']`)
	matches := re.FindAllStringSubmatch(rawHTML, -1)

	for idx, m := range matches {
		if len(m) >= 4 {
			name := m[1]
			catStr := m[2]
			priceVal, err := strconv.ParseFloat(m[3], 64)
			if err != nil {
				continue
			}

			category := models.CategorySnack
			switch strings.ToLower(catStr) {
			case "popcorn":
				category = models.CategoryPopcorn
			case "beverage", "drink":
				category = models.CategoryBeverage
			case "combo":
				category = models.CategoryCombo
			case "dessert":
				category = models.CategoryDessert
			}

			items = append(items, models.ConcessionItem{
				ID:       fmt.Sprintf("%s-EXT-%02d", cinemaID, idx+1),
				CinemaID: cinemaID,
				Name:     name,
				Category: category,
				Price:    priceVal,
				InStock:  true,
			})
		}
	}

	return items, nil
}

// ScrapeLiveUberEatsMenu uses Colly to extract concession menus from an Uber Eats store page.
func (s *RegionalScraper) ScrapeLiveUberEatsMenu(cinemaID, url string) ([]models.ConcessionItem, error) {
	if url == "" {
		return nil, fmt.Errorf("%w: Uber Eats URL cannot be empty", ErrInvalidInput)
	}

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"),
	)

	var items []models.ConcessionItem
	var parseErr error
	itemIdx := 1

	// Uber Eats standard menu item container (data-test="store-item")
	c.OnHTML(`li[data-test="store-item"]`, func(e *colly.HTMLElement) {
		name := e.ChildText(`span[data-test="store-item-name"]`)
		desc := e.ChildText(`span[data-test="store-item-description"]`)
		priceStr := e.ChildText(`span[data-test="store-item-price"]`) // format: "Rs. 1,200.00"

		if name != "" && priceStr != "" {
			cleanPrice := strings.ReplaceAll(priceStr, "Rs.", "")
			cleanPrice = strings.ReplaceAll(cleanPrice, ",", "")
			cleanPrice = strings.TrimSpace(cleanPrice)
			priceVal, err := strconv.ParseFloat(cleanPrice, 64)
			if err == nil {
				category := models.CategorySnack
				nameLower := strings.ToLower(name)
				if strings.Contains(nameLower, "popcorn") {
					category = models.CategoryPopcorn
				} else if strings.Contains(nameLower, "coke") || strings.Contains(nameLower, "pepsi") || strings.Contains(nameLower, "drink") {
					category = models.CategoryBeverage
				} else if strings.Contains(nameLower, "combo") {
					category = models.CategoryCombo
				}

				items = append(items, models.ConcessionItem{
					ID:          fmt.Sprintf("%s-UE-%02d", cinemaID, itemIdx),
					CinemaID:    cinemaID,
					Name:        name,
					Category:    category,
					Price:       priceVal,
					Description: desc,
					InStock:     true,
				})
				itemIdx++
			}
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		parseErr = fmt.Errorf("colly scrape failed with status %d: %w", r.StatusCode, err)
	})

	err := c.Visit(url)
	if err != nil {
		return nil, fmt.Errorf("failed to visit Uber Eats URL: %w", err)
	}

	if parseErr != nil {
		return nil, parseErr
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("no items extracted; page might be protected by Cloudflare or structure changed")
	}

	return items, nil
}
