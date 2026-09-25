# System Architecture Specification 🍏🍿

> Engineered with the architectural principles, clarity, and precision expected of Apple system designers.

---

## 🏛️ Executive Summary

The **Cinema Food & Beverage Recommendation Engine** is an intelligent, high-throughput microservice and interactive Progressive Web Application (PWA). It provides real-time pricing telemetry, heuristic knapsack bundle optimization, and authentic Sri Lankan digital payment processing for moviegoers across regional theater circuits in Sri Lanka (strictly targeting circuits outside Colombo metropolitan limits).

---

## 📐 High-Level Topology & Data Flow

```
                                  [ Moviegoers / Cashiers ]
                                             │
                       ┌─────────────────────┴─────────────────────┐
                       │                                           │
             [ Desktop Web Portal ]                       [ Mobile PWA Shell ]
           (Apple HIG Glassmorphism)                    (Touch Haptic Simulator)
                       │                                           │
                       └─────────────────────┬─────────────────────┘
                                             │ HTTP / REST (JSON)
                                             ▼
                     ┌───────────────────────────────────────────────┐
                     │            HTTP Edge Router & Mux             │
                     │          (Go 1.22+ ServeMux Routing)          │
                     └───────┬───────────────┬───────────────┬───────┘
                             │               │               │
               ┌─────────────┘               │               └─────────────┐
               ▼                             ▼                             ▼
     ┌───────────────────┐         ┌───────────────────┐         ┌───────────────────┐
     │  Cinema Handlers  │         │  Recommender Mux  │         │ Checkout Handlers │
     │ (Geo & Registry)  │         │(Cart & Heuristic) │         │ (Sri Lankan Pay)  │
     └─────────┬─────────┘         └─────────┬─────────┘         └─────────┬─────────┘
               │                             │                             │
               ▼                             ▼                             ▼
     ┌───────────────────┐         ┌───────────────────┐         ┌───────────────────┐
     │   `pkg/scraper`   │         │ `pkg/recommender` │         │   `pkg/payment`   │
     │                   │         │                   │         │                   │
     │ ├── scraper.go    │         │ ├── recommender.go│         │ ├── payment.go    │
     │ ├── registry.go   │         │ ├── optimizer.go  │         │ ├── methods.go    │
     │ ├── geo.go        │         │ ├── cart.go       │         │ ├── lankaqr.go    │
     │ ├── pricing.go    │         │ ├── upsell.go     │         │ ├── bnpl.go       │
     │ ├── feed_parser.go│         │ └── promotions.go │         │ └── checkout.go   │
     │ └── ads.go        │         └───────────────────┘         └───────────────────┘
     └───────────────────┘                   │
                                             ▼
                                   ┌───────────────────┐
                                   │ `internal/models` │
                                   │                   │
                                   │ ├── cinema.go     │
                                   │ ├── concession.go │
                                   │ ├── promotion.go  │
                                   │ ├── recommend.go  │
                                   │ ├── payment.go    │
                                   │ ├── ad.go         │
                                   │ └── errors.go     │
                                   └───────────────────┘
```

---

## 📦 Package Architecture & File Responsibilities

Every Go source file in this repository strictly adheres to the **Single Responsibility Principle (SRP)**:

### 1. `internal/models` (Domain Entities & Value Objects)
| File | Responsibility |
| :--- | :--- |
| [`cinema.go`](file:///e:/Github/cinema-recommender/internal/models/cinema.go) | Movie theater entity, regional location metadata, coordinates, and nearby food spots. |
| [`concession.go`](file:///e:/Github/cinema-recommender/internal/models/concession.go) | Concession food/beverage catalog items, categories, price trends, and discount breakdowns. |
| [`promotion.go`](file:///e:/Github/cinema-recommender/internal/models/promotion.go) | Scope Privilege Club loyalty tiers, promotion rules, and savings calculations. |
| [`recommendation.go`](file:///e:/Github/cinema-recommender/internal/models/recommendation.go) | Recommendation request/result DTOs and cart evaluation payloads. |
| [`payment.go`](file:///e:/Github/cinema-recommender/internal/models/payment.go) | Authentic Sri Lankan payment methods, gateway schemas, and checkout transaction results. |
| [`ad.go`](file:///e:/Github/cinema-recommender/internal/models/ad.go) | Cinema lobby screen promotional advertisements and combo deal models. |
| [`errors.go`](file:///e:/Github/cinema-recommender/internal/models/errors.go) | Sentinel domain errors (`ErrColomboExcluded`, `ErrCinemaNotFound`, `ErrEmptyCart`, etc.). |
| [`models.go`](file:///e:/Github/cinema-recommender/internal/models/models.go) | Package documentation and architectural index. |

### 2. `pkg/scraper` (Outstation Ingestion & Telemetry)
| File | Responsibility |
| :--- | :--- |
| [`scraper.go`](file:///e:/Github/cinema-recommender/pkg/scraper/scraper.go) | Scraper lifecycle, thread-safe memory caching, and Colombo exclusion policy enforcement. |
| [`registry.go`](file:///e:/Github/cinema-recommender/pkg/scraper/registry.go) | Central registry of out-of-Colombo theater circuits (Kandy, Galle, Gampaha, Kurunegala, etc.). |
| [`geo.go`](file:///e:/Github/cinema-recommender/pkg/scraper/geo.go) | Haversine trigonometric formula for GPS distance calculation and proximity sorting. |
| [`pricing.go`](file:///e:/Github/cinema-recommender/pkg/scraper/pricing.go) | Real-time dynamic pricing simulation (8-second tick cycle) and multi-channel discount attachment. |
| [`feed_parser.go`](file:///e:/Github/cinema-recommender/pkg/scraper/feed_parser.go) | Regional concession feeds and raw HTML table/div scraping parser. |
| [`ads.go`](file:///e:/Github/cinema-recommender/pkg/scraper/ads.go) | Dynamic lobby screen promotional advertisements generator for combo deals. |

### 3. `pkg/recommender` (Knapsack Heuristic & Promotion Engine)
| File | Responsibility |
| :--- | :--- |
| [`recommender.go`](file:///e:/Github/cinema-recommender/pkg/recommender/recommender.go) | Engine coordinator, promotion retrieval, and context-aware campaign matching. |
| [`optimizer.go`](file:///e:/Github/cinema-recommender/pkg/recommender/optimizer.go) | Bounded knapsack heuristic algorithm maximizing snack value while respecting party budget. |
| [`cart.go`](file:///e:/Github/cinema-recommender/pkg/recommender/cart.go) | Basket discount evaluator and multi-tier promotion calculator. |
| [`upsell.go`](file:///e:/Github/cinema-recommender/pkg/recommender/upsell.go) | Threshold-gap intelligence advising users how to reach the next discount milestone. |
| [`promotions.go`](file:///e:/Github/cinema-recommender/pkg/recommender/promotions.go) | Scope Privilege Club campaign specifications (`SCOPE-PLATINUM`, `SCOPE-COMBO30`, etc.). |

### 4. `pkg/payment` (Authentic Sri Lankan Payment Gateway)
| File | Responsibility |
| :--- | :--- |
| [`payment.go`](file:///e:/Github/cinema-recommender/pkg/payment/payment.go) | Payment service coordinator and method catalog retrieval. |
| [`methods.go`](file:///e:/Github/cinema-recommender/pkg/payment/methods.go) | Sri Lankan payment channels specifications (LankaQR, FriMi, Genie, eZ Cash, mCash, Koko, Mintpay). |
| [`lankaqr.go`](file:///e:/Github/cinema-recommender/pkg/payment/lankaqr.go) | Central Bank of Sri Lanka (CBSL) EMVCo standard dynamic LankaQR payload generator. |
| [`bnpl.go`](file:///e:/Github/cinema-recommender/pkg/payment/bnpl.go) | Buy Now Pay Later (Koko & Mintpay) 3-month interest-free split payment calculator. |
| [`checkout.go`](file:///e:/Github/cinema-recommender/pkg/payment/checkout.go) | Transaction orchestration and in-person candy bar counter voucher issuance. |

### 5. `internal/server` (HTTP REST Delivery & Middleware)
| File | Responsibility |
| :--- | :--- |
| [`server.go`](file:///e:/Github/cinema-recommender/internal/server/server.go) | Server struct, timeout configuration, graceful shutdown, and JSON response helpers. |
| [`routes.go`](file:///e:/Github/cinema-recommender/internal/server/routes.go) | Go 1.22+ pattern-matched route registration. |
| [`middleware.go`](file:///e:/Github/cinema-recommender/internal/server/middleware.go) | Latency logging, CORS headers, and panic recovery middleware. |
| [`handlers_cinema.go`](file:///e:/Github/cinema-recommender/internal/server/handlers_cinema.go) | Cinema listing and geolocation lookup endpoints. |
| [`handlers_concession.go`](file:///e:/Github/cinema-recommender/internal/server/handlers_concession.go) | Concessions, live dynamic rate feeds, and lobby screen ads endpoints. |
| [`handlers_recommend.go`](file:///e:/Github/cinema-recommender/internal/server/handlers_recommend.go) | AI knapsack recommendation and cart evaluation endpoints. |
| [`handlers_checkout.go`](file:///e:/Github/cinema-recommender/internal/server/handlers_checkout.go) | Sri Lankan gateway payment methods and checkout endpoints. |
| [`handlers_static.go`](file:///e:/Github/cinema-recommender/internal/server/handlers_static.go) | Web UI, Mobile PWA, manifest, service worker, and health check endpoints. |

---

## 🧮 Algorithmic Specifications

### 1. Bounded Knapsack Heuristic Optimization
Given:
- Cinema concession items $I = \{i_1, i_2, \dots, i_n\}$ with prices $P(i)$ and categories $C(i)$
- Party size $S \ge 1$
- User budget $B > 0$
- Loyalty tier $T \in \{\text{Standard}, \text{Student}, \text{Silver}, \text{Gold}, \text{Platinum}\}$

The optimizer selects a balanced subset $I^* \subseteq I$ satisfying:
1. **Core Snacking Balance**: $\sum_{i \in I^*} \mathbb{I}_{C(i) = \text{Popcorn}} \ge 1$ and $\sum_{i \in I^*} \mathbb{I}_{C(i) = \text{Beverage}} \ge \min(S, 4)$.
2. **Threshold Bridging**: If $\text{gap} = \text{Target} - \sum P(i) \le 400\text{ LKR}$, an incremental item is added if the resulting discount $\Delta D > P(\text{item})$, maximizing net customer savings.

### 2. Haversine Spherical Distance
Distance between coordinates $(\phi_1, \lambda_1)$ and $(\phi_2, \lambda_2)$ is evaluated via:
$$d = 2R \cdot \arcsin\left(\sqrt{\sin^2\left(\frac{\Delta \phi}{2}\right) + \cos(\phi_1)\cos(\phi_2)\sin^2\left(\frac{\Delta \lambda}{2}\right)}\right)$$
where $R = 6371.0\text{ km}$.

### 3. CBSL LankaQR Dynamic EMV Specification
Dynamic QR codes generated by [`lankaqr.go`](file:///e:/Github/cinema-recommender/pkg/payment/lankaqr.go) adhere to EMVCo Tag-Length-Value (TLV) standards mandated by the Central Bank of Sri Lanka:
- `Tag 00`: Format Indicator (`01`)
- `Tag 01`: Point of Initiation (`12` = Dynamic QR)
- `Tag 26`: National Merchant Account ID (`com.lankaqr`)
- `Tag 52`: Merchant Category Code (`7832` = Motion Picture Theaters)
- `Tag 53`: Transaction Currency (`144` = Sri Lankan Rupee LKR)
- `Tag 54`: Transaction Amount
- `Tag 58`: Country Code (`LK`)
- `Tag 62`: Reference Data (Order ID)
