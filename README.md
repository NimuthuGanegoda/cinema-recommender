# CineBite 🍿 — Universal Sri Lanka Cinema Concession & Snack Recommender

An intelligent, prototype cinema concession recommendation service and Apple HIG-crafted interactive dashboard engineered in Go to scrape, aggregate, and recommend optimal candy bar deals across all Sri Lankan cinema theaters.

---

## 🎯 Project Overview

**CineBite** provides moviegoers with real-time, value-optimized recommendations on cinema concessions. Built for any theater in Sri Lanka (Regal Kandy, Queens Galle, Rodeo Negombo, Rajah Jaffna, SK Matara, Cinemax Anuradhapura, Milano Ratnapura, Scope partner circuits, and heritage independent screens), moviegoers can find the most cost-effective snacking combinations, live prices, and instant candy bar discounts with **CinePass VIP**.

### Key Design & Universal Directives
- **Universal Cinema Coverage**: Designed for ANY theater circuit across Sri Lanka, with automated GPS nearby theater detection and in-person snack guide.
- **Apple Human Interface Guidelines (HIG)**: Crafted with OLED true black canvas, Cupertino system grays (`#1c1c1e`), Apple System Blue (`#0071e3`), Apple Mint Green (`#30d158`), SF Pro typography, Dynamic Island live indicator, and local food planning summary.
- **CinePass VIP Promotion Engine**: Evaluates discount tiers (Standard, Student, Silver, Gold, Platinum), combo rebates (`CINE-COMBO30`, `CINE-CINEMA25`, `CINE-PLATINUM`), and threshold upsell advice.

---

## 🏗️ Architecture & Project Layout

The repository follows standard Go project layout conventions (`cmd/`, `pkg/`, `internal/`):

```
cinema-recommender/
├── ARCHITECTURE.md                 # 🍏 Apple-grade system architecture & design patterns specification
├── API.md                          # Complete REST API reference, JSON schemas & cURL examples
├── CONTRIBUTING.md                 # Developer guide, code conventions & extension instructions
├── Makefile                        # Developer automation (make build, make test, make run)
├── Dockerfile                      # Multi-stage minimal container build
├── docker-compose.yml              # Local container deployment
├── cmd/
│   └── api/
│       └── main.go                 # CLI commands, server bootstrap & graceful shutdown
├── internal/
│   ├── models/
│   │   ├── cinema.go               # Theaters, regional circuits, food spots & GPS coordinates
│   │   ├── concession.go           # Concession catalog items, categories, price trends & discounts
│   │   ├── promotion.go            # Scope Privilege Club loyalty tiers, discount types & savings
│   │   ├── recommendation.go       # Knapsack request/response DTOs & basket evaluation schemas
│   │   ├── ad.go                   # Cinema promotional lobby screen advertisements
│   │   ├── errors.go               # Domain sentinel errors (ErrColomboExcluded, etc.)
│   │   └── models.go               # Package overview and architectural documentation
│   └── server/
│       ├── server.go               # Server struct, configuration, lifecycle & JSON helpers
│       ├── routes.go               # Pattern-matched HTTP route registration (Go 1.22+)
│       ├── middleware.go           # Latency logger, CORS headers & panic recovery
│       ├── handlers_cinema.go      # Regional cinema listing & geolocation lookup handlers
│       ├── handlers_concession.go  # Concessions, dynamic live rate ticker & lobby ads handlers
│       ├── handlers_recommend.go   # AI knapsack recommendation & cart evaluation handlers
│       ├── handlers_static.go      # Desktop Web UI, Mobile PWA, manifest, sw & health handlers
│       └── ui/
│           ├── index.html          # Embedded dark-mode interactive web portal & mobile simulator
│           ├── mobile.html         # Dedicated touch-optimized mobile web app shell (PWA)
│           ├── apple-design-system.css # Apple HIG glassmorphism design tokens & touch animations
│           ├── cinema-client.js    # Clean Promise-based JavaScript client SDK
│           ├── manifest.json       # Web App Manifest for mobile installation
│           └── sw.js               # Service Worker for offline PWA caching
├── pkg/
│   ├── recommender/
│   │   ├── recommender.go          # Recommendation engine coordinator & promotion matcher
│   │   ├── optimizer.go            # Bounded knapsack heuristic algorithm for food & beverage
│   │   ├── cart.go                 # Basket evaluation & multi-tier promotion calculator
│   │   ├── upsell.go               # Threshold-gap advice & party value score heuristics
│   │   ├── promotions.go           # Scope Privilege Club campaigns registry (SCOPE-PLATINUM, etc.)
│   │   └── recommender_test.go     # Comprehensive promotion tier & optimization unit tests
│   └── scraper/
│       ├── scraper.go              # Scraper coordinator, thread-safe memory cache & validation
│       ├── registry.go             # Outstation regional theater circuit database & queries
│       ├── geo.go                  # Haversine spherical distance calculation & proximity sorting
│       ├── pricing.go              # Real-time dynamic pricing simulator (8-sec cycle) & discounts
│       ├── feed_parser.go          # Regional supplier concession catalog & raw HTML feed parser
│       ├── ads.go                  # Cinema lobby screen promotional advertisements generator
│       └── scraper_test.go         # Policy exclusion, scraping and concurrency tests
├── bin/
│   └── cinema-recommender.exe      # Compiled self-contained executable binary
├── go.mod                          # Go module definition (cinema-recommender)
└── README.md                       # Project documentation & architectural specification
```

### 🧩 Module & Architecture Specifications
For in-depth architectural blueprints and developer documentation:
- 🍏 **[ARCHITECTURE.md](file:///e:/Github/cinema-recommender/ARCHITECTURE.md)**: System topology, design patterns, algorithmic proofs, and package responsibilities.
- 📡 **[API.md](file:///e:/Github/cinema-recommender/API.md)**: REST API reference, request/response JSON schemas, and cURL commands.
- 🤝 **[CONTRIBUTING.md](file:///e:/Github/cinema-recommender/CONTRIBUTING.md)**: Engineering guidelines, adding theaters, and test standards.

---



8. **Sri Lankan Bank Cards**: Visa / MasterCard / LankaPay issued by local banks (Commercial Bank, Sampath Bank, HNB, BOC, People's Bank).

*(Plus **Scope Privilege Club Concession Pass** unlocking up to 35% automated concession tier discounts and combos)*

---

## 🚀 Getting Started

### Prerequisites
- **Go**: Version 1.24 or higher (automatically verified and supported)

### Building the Application

Compile the standalone binary into `bin/`:
```bash
go build -o bin/cinema-recommender.exe ./cmd/api
```

---

## 💻 Running the Application

### 1. Interactive CLI Mode

Quickly query optimal bundles right from your terminal:

```bash
# Calculate deal for a party of 2 in Kandy with LKR 2,500 budget
./bin/cinema-recommender.exe -city Kandy -budget 2500 -party 2

# List all registered theaters
./bin/cinema-recommender.exe -list-cinemas

# List all active Scope Privilege promotion campaigns
./bin/cinema-recommender.exe -list-promos

# Calculate with specific Scope Privilege tier and promo code
./bin/cinema-recommender.exe -city Galle -budget 3500 -party 3 -tier Platinum -promo SCOPE-PLATINUM
```

### 2. HTTP REST API, Web Application & Mobile Application Server Mode

Start the server on port `:8080` (or any custom port with `-port`):
```bash
./bin/cinema-recommender.exe -server -port :8080
```
Then explore both application experiences:
- 📱 **Mobile Application**: Visit **[http://localhost:8080/mobile](http://localhost:8080/mobile)** for the dedicated touch-optimized Mobile App shell with bottom navigation bar (Explore, Optimizer, VIP Pass, Food Plan).
- 📲 **Installable PWA**: Open `http://localhost:8080/mobile` on any mobile browser (Safari iOS / Chrome Android) and tap **"Add to Home Screen"** to install it as an app on your phone!
- 🧪 **Live Mobile Simulator**: From the desktop Web UI, click the **"📱 Mobile App Simulator"** button in the top navigation bar to test the live mobile app inside a virtual smartphone mockup directly on your screen!

---

## 📡 REST API Reference

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/` | Desktop Web Application with embedded Mobile Simulator |
| `GET` | `/mobile` | Dedicated Mobile Application (touch-friendly PWA) |
| `GET` | `/manifest.json` | Web App Manifest for mobile installation |
| `GET` | `/sw.js` | PWA Service Worker for offline caching |
| `GET` | `/api/v1/health` | Service health status and metadata |
| `GET` | `/api/v1/cinemas` | List regional out-of-Colombo theaters (supports `?city=`, `?lat=&lng=` for distance auto-fetch) |
| `GET` | `/api/v1/cinemas/{id}` | Retrieve specific theater details with GPS coords, maps, & in-house/nearby food locations |
| `GET` | `/api/v1/cinemas/{id}/concessions` | Scrape and return live concession items with prices, serving location, & discount breakdowns |
| `GET` | `/api/v1/promotions/privilege` | List all active Scope Privilege discount programs and rules (aliases: `/api/v1/promotions/scope`, `/api/v1/promotions`) |
| `POST` | `/api/v1/recommend` | Compute AI/heuristic optimal concession bundle for budget & party size |
| `POST` | `/api/v1/cart/evaluate` | Evaluate custom user cart against Scope Privilege promotion thresholds |

### Example cURL Queries

**Auto-Fetch Cinemas by GPS Location:**
```bash
curl -s "http://localhost:8080/api/v1/cinemas?lat=7.2936&lng=80.6385"
```

**Get Live Concessions with Every Food Item Discount Breakdown:**
```bash
curl -s "http://localhost:8080/api/v1/cinemas/KND-KCC/concessions"
```

**Health Check:**
```bash
curl -i http://localhost:8080/api/v1/health
```

```bash
```

```bash
  -H "Content-Type: application/json" \
  -d '{
    "cinema_id": "KND-KCC",
    "customer_name": "Kasun Perera",
    "customer_phone": "0771234567",
    "items": [
      { "item_id": "KND-01", "quantity": 1 },
      { "item_id": "KND-03", "quantity": 2 }
    ]
  }'
```

```bash
  -H "Content-Type: application/json" \
  -d '{
    "cinema_id": "KND-KCC",
    "customer_name": "Kasun Perera",
    "customer_phone": "0771234567",
    "items": [
      { "item_id": "KND-06", "quantity": 1 }
    ]
  }'
```

---

## 🧪 Running Tests

Execute the full suite of unit, integration, concurrency, and policy tests:

```bash
go test -v ./...
```

- Regional cinema registry lookup and concession feed parsing
- Promotion tier discounts, caps, and combo specials (`SCOPE-COMBO30`, `SCOPE-CINEMA25`)
- Upsell intelligence and bundle recommendation heuristics

---

## 🗺️ Roadmap Status

- [x] Phase 1: Repository Scaffolding & Architecture Blueprint
- [x] Phase 2: Domain Modeling (`internal/models`)
- [x] Phase 3: Out-of-Colombo Scraping Pipeline Implementation (`pkg/scraper`)
- [x] Phase 4: Scope Privilege Promotion & Discount Calculator Engine (`pkg/recommender`)
- [x] Phase 5: REST API Endpoints, Delivery Layer & Embedded Web UI (`internal/server` & `cmd/api`)
- [x] Phase 6: Mobile Application Shell & Progressive Web App (`internal/server/ui/mobile.html`)


## Local food planning scope

CineBite helps users browse cinemas and movies, explore cinema food menus, compare food offers, and build snack recommendations. The food plan is saved in the browser on the current device and can be printed. Listed prices and promotions are estimates or examples; confirm them with the cinema.

Run with Go 1.24 or later. The Go server embeds the pages in `internal/server/ui/`; `docs/` and the root contain equivalent static pages. No publishing or remote deployment is required.

## Run the Go assignment locally

Go is the primary application language. `cmd/api` runs the HTTP server; `internal/server` serves the website and API; `pkg/recommender` calculates snack bundles and discounts; `pkg/scraper` holds cinema, movie, and food catalog logic. HTML and CSS provide the browser layout. JavaScript renders Go API responses and handles browser interactions; it does not calculate recommendations.

With Go 1.24 or newer, run `go run ./cmd/api -server -port 127.0.0.1:8080` and open http://127.0.0.1:8080/. Run `go test ./...` to check the Go application. Cinema, movie, and food catalog data are local examples for the assignment. No payment or remote publishing feature is included.

The shared layout adapts to phones, tablets, and desktops. Five dark Omarchy-inspired colour palettes rotate every 45 seconds with a 1.6-second fade. Palette timing and device detection stay in the browser; no Desktop/App switch is needed.

Photo sources are in `IMAGE_CREDITS.json` and `photo-credits.html`.

A compiled local macOS version is included in `bin/cinebite`. Double-click `run-local.command` to start it; no global Go installation is needed to run that binary. Stop it with Control-C in its terminal.

The mobile layout wraps planner choices and action buttons, keeps touch targets usable, and lazily loads offscreen food and cinema photographs. No light palette or device-mode switch is included.
