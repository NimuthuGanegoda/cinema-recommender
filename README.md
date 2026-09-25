# Cinema Food & Beverage Recommendation Engine 🍿🥤

An intelligent, production-ready backend service and interactive dashboard built in Go designed to scrape, aggregate, and recommend optimal food and beverage deals across regional cinema theaters.

---

## 🎯 Project Overview

The **Cinema Food & Beverage Recommendation Engine** provides moviegoers with real-time, value-optimized recommendations on cinema concessions. By combining automated menu and pricing extraction with an advanced promotion engine, users can find the most cost-effective snacking combinations for their movie experience.

### Key Scope & Design Directives
- **Scope Privilege Promotion Optimization**: The discount and offer evaluation engine is specifically optimized for processing **Scope Privilege Club** promotional campaigns, partner tier discounts (Standard, Student, Silver, Gold, Platinum), combo rebates (`SCOPE-COMBO30`, `SCOPE-CINEMA25`), and automated threshold upsell intelligence.

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
│   │   ├── payment.go              # Sri Lankan payment channels, method metadata & checkout DTOs
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
│       ├── handlers_checkout.go    # Sri Lankan payment gateway listing & checkout handlers
│       ├── handlers_static.go      # Desktop Web UI, Mobile PWA, manifest, sw & health handlers
│       ├── server_test.go          # End-to-end HTTP endpoint, checkout & policy tests
│       └── ui/
│           ├── index.html          # Embedded dark-mode interactive web portal & mobile simulator
│           ├── mobile.html         # Dedicated touch-optimized mobile web app shell (PWA)
│           ├── apple-design-system.css # Apple HIG glassmorphism design tokens & touch animations
│           ├── cinema-client.js    # Clean Promise-based JavaScript client SDK
│           ├── manifest.json       # Web App Manifest for mobile installation
│           └── sw.js               # Service Worker for offline PWA caching
├── pkg/
│   ├── payment/
│   │   ├── payment.go              # Payment service struct, coordinator & method listing
│   │   ├── methods.go              # Sri Lankan payment channels catalog (LankaQR, FriMi, etc.)
│   │   ├── lankaqr.go              # Central Bank of Sri Lanka (CBSL) dynamic EMVCo QR generator
│   │   ├── bnpl.go                 # Koko & Mintpay 3-month interest-free split payment engine
│   │   ├── checkout.go             # Transaction processing & counter pickup voucher issuance
│   │   └── payment_test.go         # CBSL EMV standard, USSD push & BNPL installment tests
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

## 💳 Sri Lankan Payment Channels Supported

Both the Desktop Web Portal and Mobile App Shell support authentic Sri Lankan payment channels (strictly adhering to genuine Sri Lankan financial systems):

1. **LankaQR (CBSL National Standard)**: Generates dynamic EMV-compliant LankaQR payloads compatible with **Commercial Bank COMBANK Q+**, **Sampath WePay**, **Nations Trust FriMi**, **HNB SOLO**, **BOC SmartPay**, and **Seylan Pay**, including a 5% digital incentive.
2. **FriMi (Nations Trust Bank)**: Leading digital banking lifestyle app with instant 10% concession cashback.
3. **Genie (Dialog Finance)**: Digital wallet rebate with 8% instant savings.
4. **eZ Cash (Dialog Axiata)**: Direct mobile wallet charging via instant USSD PIN push notification to Dialog, Hutch, and Airtel numbers.
5. **mCash (SLT-Mobitel)**: Direct carrier billing and mobile wallet authorization for Mobitel subscribers.
6. **Koko (Buy Now, Pay Later)**: Splits cinema concession orders over LKR 1,500 into 3 interest-free monthly installments.
7. **Mintpay (Buy Now, Pay Later)**: Alternative Sri Lankan split-payment installment option.
8. **Sri Lankan Bank Cards**: Visa / MasterCard / LankaPay issued by local banks (Commercial Bank, Sampath Bank, HNB, BOC, People's Bank).
9. **Counter Cash**: Physical cash payment in LKR at the theater concession stand upon collection.

*(Plus **Scope Privilege Club Concession Pass** unlocking up to 35% automated concession tier discounts and combos)*

---

## 🚀 Getting Started

### Prerequisites
- **Go**: Version 1.22 or higher (automatically verified and supported)

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

# List all registered outstation theaters outside Colombo
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
- 🖥️ **Desktop Web Application**: Visit **[http://localhost:8080](http://localhost:8080)** to access the full cinema concessions planner, visual deal optimizer, live concession menu browser, embedded smartphone simulator, and Sri Lankan checkout drawer.
- 📱 **Mobile Application**: Visit **[http://localhost:8080/mobile](http://localhost:8080/mobile)** for the dedicated touch-optimized Mobile App shell with bottom navigation bar (Explore, Optimizer, VIP Pass, Voucher Pass).
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
| `GET` | `/api/v1/payment-methods` | List supported Sri Lankan payment gateways with partner discounts |
| `POST` | `/api/v1/recommend` | Compute AI/heuristic optimal concession bundle for budget & party size |
| `POST` | `/api/v1/cart/evaluate` | Evaluate custom user cart against Scope Privilege promotion thresholds |
| `POST` | `/api/v1/checkout` | Process order and checkout via LankaQR, eZ Cash, mCash, FriMi, or Koko |

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

**List Supported Sri Lankan Payment Methods:**
```bash
curl -s http://localhost:8080/api/v1/payment-methods
```

**Place Concession Order via LankaQR:**
```bash
curl -X POST http://localhost:8080/api/v1/checkout \
  -H "Content-Type: application/json" \
  -d '{
    "cinema_id": "KND-KCC",
    "payment_method": "LANKAQR",
    "customer_name": "Kasun Perera",
    "customer_phone": "0771234567",
    "items": [
      { "item_id": "KND-01", "quantity": 1 },
      { "item_id": "KND-03", "quantity": 2 }
    ]
  }'
```

**Place Concession Order via Koko BNPL (3 Installments):**
```bash
curl -X POST http://localhost:8080/api/v1/checkout \
  -H "Content-Type: application/json" \
  -d '{
    "cinema_id": "KND-KCC",
    "payment_method": "KOKO_BNPL",
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

All packages (`internal/server`, `pkg/payment`, `pkg/recommender`, `pkg/scraper`) include tests covering:
- Strict rejection of Colombo metropolitan locations (`ErrColomboExcluded`)
- Regional cinema registry lookup and concession feed parsing
- Promotion tier discounts, caps, and combo specials (`SCOPE-COMBO30`, `SCOPE-CINEMA25`)
- Upsell intelligence and bundle recommendation heuristics
- Sri Lankan payment gateways (CBSL LankaQR EMV strings, eZ Cash USSD push, Koko 3x installment calculation)
- REST endpoint responses, CORS handling, checkout processing, and panic recovery

---

## 🗺️ Roadmap Status

- [x] Phase 1: Repository Scaffolding & Architecture Blueprint
- [x] Phase 2: Domain Modeling (`internal/models`)
- [x] Phase 3: Out-of-Colombo Scraping Pipeline Implementation (`pkg/scraper`)
- [x] Phase 4: Scope Privilege Promotion & Discount Calculator Engine (`pkg/recommender`)
- [x] Phase 5: REST API Endpoints, Delivery Layer & Embedded Web UI (`internal/server` & `cmd/api`)
- [x] Phase 6: Mobile Application Shell & Progressive Web App (`internal/server/ui/mobile.html`)
- [x] Phase 7: Sri Lankan Payment Gateway Integrations (LankaQR, eZ Cash, mCash, FriMi, Koko BNPL)
