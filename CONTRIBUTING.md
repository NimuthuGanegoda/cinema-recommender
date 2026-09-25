# Contributing to Cinema Recommender 🍏

Thank you for your interest in contributing. We hold our engineering standards to the level of precision, elegance, and reliability expected in Apple developer ecosystems.

---

## 🏛️ Guiding Engineering Principles

1. **Single Responsibility Principle (SRP)**: Each file in the repository must serve a single, clearly understood purpose. Do not accumulate disparate concerns in monolithic files.
2. **Strict Regional Policy Compliance**: The scraping and recommendation pipeline is **strictly scoped for cinema circuits outside the Colombo metropolitan district**. Any attempt to register or scrape theaters within Colombo will be rejected by the domain layer (`models.ErrColomboExcluded`).
3. **Authentic Localized Standards**: We strictly support authentic Sri Lankan payment channels (CBSL LankaQR, FriMi, Genie, eZ Cash, mCash, Koko, Mintpay). Fictional or non-existent payment providers must not be introduced.
4. **Physical Counter Pickup Context**: This application is an in-theater recommender and counter savings pass. Food is never delivered online; vouchers are physically presented to cinema cashiers.

---

## 🛠️ Development Workflow

### 1. Prerequisites
- **Go 1.22+**
- Standard terminal shell (macOS/Linux `zsh`/`bash`, Windows `pwsh`)

### 2. Running Unit & Integration Tests
Ensure 100% of test suites pass before submitting any changes:
```bash
go test -v ./...
```

### 3. Compiling the Standalone Binary
```bash
go build -o bin/cinema-recommender.exe ./cmd/api
```

### 4. Running the Dev Server
```bash
./bin/cinema-recommender.exe -server -port :8080
```
- Desktop Web App: `http://localhost:8080`
- Mobile PWA App: `http://localhost:8080/mobile`

---

## 📍 Adding a New Regional Theater

To add an outstation cinema circuit:
1. Open [`pkg/scraper/registry.go`](file:///e:/Github/cinema-recommender/pkg/scraper/registry.go).
2. Append a new `models.Cinema` definition to `defaultCinemaRegistry()`:
```go
{
    ID:                 "JFN-RJA",
    Name:               "Rajah Cinema Jaffna",
    Chain:              "Northern Cinema Network",
    City:               "Jaffna",
    Province:           "Northern Province",
    Address:            "Hospital Road, Jaffna",
    Screens:            3,
    IsOutside:          true,
    Latitude:           9.6615,
    Longitude:          80.0255,
    MapURL:             "https://maps.google.com/?q=9.6615,80.0255",
    HasInHouseFood:     true,
    FoodPlaceName:      "Rajah Refreshment Parlour & Snack Stand",
    FoodPlaceType:      "Northern Concession Foyer",
    NearbyFoodOptions:  []string{"Mangos Vegetarian Restaurant", "Rio Ice Cream"},
    FoodHours:          "10:00 AM - 10:45 PM",
    FoodDeliveryToSeat: false,
}
```
3. Add the concession catalog in [`pkg/scraper/feed_parser.go`](file:///e:/Github/cinema-recommender/pkg/scraper/feed_parser.go) under `scrapeByCity()`.
4. Run `go test -v ./pkg/scraper` to verify coverage.

---

## 🎨 Design & Human Interface Guidelines (HIG)

All user interfaces adhere to the following design system rules:
- **Obsidian Dark Mode**: Deep black background (`#07070a`), frosted glass cards (`rgba(16, 16, 22, 0.94)`), subtle specular borders (`rgba(212, 175, 55, 0.22)`).
- **Scope Cinematic Gold**: Primary accents use Scope Gold (`#d4af37`), gold glow (`rgba(212, 175, 55, 0.35)`), and gold hover states (`#b8860b`).
- **Apple Typography**: Pairing of **Cinzel** for cinematic headings with **Montserrat** and **Outfit** for tactile UI elements.
- **Haptic Feedback**: Active states scale down to `0.92` on touch with subtle spring easing.
