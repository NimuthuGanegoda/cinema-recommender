# Cinema Recommender REST API Reference 📡

> Complete, OpenAPI-aligned HTTP API specification for client integrations.

Base URL: `http://localhost:8080`  
Protocol: `HTTP/1.1`, `HTTP/2`  
Content-Type: `application/json; charset=utf-8`

---

## 📑 Table of Endpoints

| Category | Method | Path | Description |
| :--- | :--- | :--- | :--- |
| **System** | `GET` | `/api/v1/health` | Service health status and metadata |
| **Cinemas** | `GET` | `/api/v1/cinemas` | List regional theaters with optional location sorting |
| | `GET` | `/api/v1/cinemas/{id}` | Retrieve individual theater profile and dining details |
| **Concessions**| `GET` | `/api/v1/cinemas/{id}/concessions` | Retrieve scraped food and beverage menu items |
| | `GET` | `/api/v1/cinemas/{id}/concessions/live` | Stream real-time dynamic pricing feed & ticker |
| | `GET` | `/api/v1/cinemas/{id}/ads` | Retrieve in-theater promotional lobby screen ads |
| **Promotions**| `GET` | `/api/v1/promotions/privilege` | List active Scope Privilege Club discount programs |
| **Intelligence**| `POST` | `/api/v1/recommend` | Compute AI/knapsack value-maximizing concession bundle |
| | `POST` | `/api/v1/cart/evaluate` | Evaluate custom user cart against promotion thresholds |

---

## 1. System Health

### `GET /api/v1/health`
Returns runtime metrics and regional theater availability.

**Response `200 OK`**:
```json
{
  "status": "healthy",
  "service": "cinema-food-beverage-recommender",
  "version": "2.0.0",
  "regional_theaters": 10,
  "scope": "strictly_outside_colombo",
  "promotion_engine": "scope_privilege_optimized",
  "timestamp": "2026-09-26T04:45:00Z"
}
```

---

## 2. Regional Cinemas

### `GET /api/v1/cinemas`
Lists regional theaters outside Colombo. Supports auto-detection by GPS coordinates.

**Query Parameters:**
- `city` (string, optional): Filter by regional city (e.g. `Kandy`, `Galle`, `Kurunegala`). Note: Colombo queries return `400 Bad Request`.
- `lat` (float, optional): User latitude for Haversine distance sorting.
- `lng` (float, optional): User longitude for Haversine distance sorting.

**cURL Example:**
```bash
curl -i "http://localhost:8080/api/v1/cinemas?lat=7.2906&lng=80.6337"
```

---

## 3. Real-Time Concessions & Dynamic Pricing

### `GET /api/v1/cinemas/{id}/concessions/live`
Returns live dynamic pricing updated every 8 seconds, simulating matinee discounts, happy hours, and snack counter rush deals.

**Response `200 OK`**:
```json
{
  "cinema_id": "KND-KCC",
  "cinema_name": "Scope Partner Multiplex - KCC Kandy",
  "dynamic_engine_state": "ACTIVE_MATINEE_FLASH_FEED",
  "live_pricing_active": true,
  "tick_id": 22445890,
  "next_update_in_sec": 6,
  "items": [
    {
      "id": "KND-01",
      "name": "Scope Jumbo Warm Caramel Popcorn",
      "category": "Popcorn",
      "base_price_lkr": 1200,
      "price_lkr": 1150,
      "price_trend": "flash_drop",
      "price_change_lkr": -50,
      "flash_deal_text": "⚡ FLASH CONCESSION DROP (-LKR 50)",
      "applicable_discounts": [
        {
          "promo_code": "SCOPE-PLATINUM",
          "provider": "Scope Privilege",
          "title": "Scope Privilege Platinum VIP (35% Off)",
          "discount_pct": 0.35,
          "discounted_price_lkr": 747.5,
          "savings_lkr": 402.5
        }
      ]
    }
  ]
}
```

---

## 4. AI Knapsack Recommendation

### `POST /api/v1/recommend`
Calculates an optimal snacking combination based on budget, party size, and loyalty tier.

**Request Payload:**
```json
{
  "cinema_id": "KND-KCC",
  "budget_lkr": 2500,
  "party_size": 2,
  "privilege_tier": "Platinum",
  "promo_code": "SCOPE-PLATINUM"
}
```

---



**Request Payload:**
```json
{
  "cinema_id": "KND-KCC",
  "items": [
    { "item_id": "KND-01", "quantity": 1 },
    { "item_id": "KND-03", "quantity": 2 }
  ],
  "privilege_tier": "Gold",
  "customer_phone": "+94 77 123 4567",
  "customer_name": "Amal Perera"
}
```

**Response `200 OK`**:

`GET /food-plan` serves the local food planning page. Cinema, movie, concession, promotion, recommendation, and basket evaluation endpoints remain available.
