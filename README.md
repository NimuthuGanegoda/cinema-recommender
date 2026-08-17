# Cinema Food & Beverage Recommendation Engine 🍿🥤

An intelligent backend service built in Go designed to scrape, aggregate, and recommend optimal food and beverage deals across regional cinema theaters.

---

## 🎯 Project Overview

The **Cinema Food & Beverage Recommendation Engine** provides moviegoers with real-time, value-optimized recommendations on cinema concessions. By combining automated menu and pricing extraction with an advanced promotion engine, users can find the most cost-effective snacking combinations for their movie experience.

### Key Scope & Design Directives
- **Regional Theater Coverage**: The scraping engine strictly targets cinema locations **outside of Colombo** (suburban, provincial, and outstation theater circuits). Concession offerings and price points within Colombo metropolitan limits are excluded by design.
- **Redopay Promotion Optimization**: The discount and offer evaluation engine is specifically optimized for processing **Redopay** promotional campaigns, partner tier discounts, bundle rebates, and checkout savings.

---

## 🏗️ Architecture & Project Layout

The repository follows standard Go project layout conventions (`cmd/`, `pkg/`, `internal/`):

```
cinema-recommender/
├── cmd/
│   └── api/
│       └── main.go           # Application entrypoint & HTTP server bootstrap
├── pkg/
│   ├── scraper/
│   │   └── scraper.go        # Out-of-Colombo cinema web scraping engine
│   └── recommender/
│       └── recommender.go    # Recommendation engine & Redopay discount calculator
├── internal/
│   └── models/
│       └── models.go         # Core data structures (Cinema, Item, Deal, Promotion)
├── .gitignore                # Go-specific ignore rules
├── go.mod                    # Go module definition (cinema-recommender)
└── README.md                 # Project documentation & architectural specification
```

### Module Responsibilities

| Package | Purpose |
| :--- | :--- |
| `cmd/api` | Server initialization, route registration, middleware configuration, and lifecycle management. |
| `pkg/scraper` | Resilient web scraping pipelines targeting cinema concession menus and ticket-combo pricing outside Colombo. |
| `pkg/recommender` | Recommendation heuristics and discount logic optimized for Redopay promotional schemes and bundling. |
| `internal/models` | Domain models, request/response DTOs, and promotion schemas shared across packages. |

---

## 🚀 Getting Started

### Prerequisites
- **Go**: Version 1.22 or higher
- **Git**: For version control

### Installation & Build

1. Clone the repository:
   ```bash
   git clone <repo-url> cinema-recommender
   cd cinema-recommender
   ```

2. Download module dependencies:
   ```bash
   go mod download
   ```

3. Build the API service:
   ```bash
   go build -o bin/api ./cmd/api
   ```

4. Run the service:
   ```bash
   ./bin/api
   ```

---

## 🗺️ Roadmap

- [x] Phase 1: Repository Scaffolding & Architecture Blueprint
- [ ] Phase 2: Domain Modeling (`internal/models`)
- [ ] Phase 3: Out-of-Colombo Scraping Pipeline Implementation (`pkg/scraper`)
- [ ] Phase 4: Redopay Promotion & Discount Calculator Engine (`pkg/recommender`)
- [ ] Phase 5: REST / gRPC API Endpoints & Delivery Layer (`cmd/api`)
