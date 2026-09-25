/**
 * 🍏 Cinema Recommender Client SDK
 * Lightweight, Promise-based API layer engineered for web & mobile clients.
 */

class CinemaClient {
  constructor(baseUrl = '') {
    this.baseUrl = baseUrl;
  }

  async _get(path) {
    const res = await fetch(`${this.baseUrl}${path}`);
    if (!res.ok) {
      const err = await res.json().catch(() => ({ error: res.statusText }));
      throw new Error(err.error || `HTTP ${res.status}`);
    }
    return res.json();
  }

  async _post(path, data) {
    const res = await fetch(`${this.baseUrl}${path}`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });
    if (!res.ok) {
      const err = await res.json().catch(() => ({ error: res.statusText }));
      throw new Error(err.error || `HTTP ${res.status}`);
    }
    return res.json();
  }

  // 1. Cinema Registry
  async getCinemas({ lat = 0, lng = 0, city = '' } = {}) {
    const params = new URLSearchParams();
    if (city) params.append('city', city);
    if (lat && lng) {
      params.append('lat', lat);
      params.append('lng', lng);
    }
    const query = params.toString() ? `?${params.toString()}` : '';
    return this._get(`/api/v1/cinemas${query}`);
  }

  async getCinema(cinemaId) {
    return this._get(`/api/v1/cinemas/${encodeURIComponent(cinemaId)}`);
  }

  // 2. Concessions & Live Ticker
  async getLiveConcessions(cinemaId) {
    return this._get(`/api/v1/cinemas/${encodeURIComponent(cinemaId)}/concessions/live`);
  }

  async getCinemaAds(cinemaId) {
    return this._get(`/api/v1/cinemas/${encodeURIComponent(cinemaId)}/ads`);
  }

  // 3. Loyalty & Promotions
  async getPromotions() {
    return this._get('/api/v1/promotions/privilege');
  }

  async getPaymentMethods() {
    return this._get('/api/v1/payment-methods');
  }

  // 4. Recommendation & Cart
  async recommendBundle({ cinemaId, budget, partySize, privilegeTier, promoCode }) {
    return this._post('/api/v1/recommend', {
      cinema_id: cinemaId,
      budget_lkr: budget,
      party_size: partySize,
      privilege_tier: privilegeTier,
      promo_code: promoCode,
    });
  }

  async evaluateCart({ cinemaId, items, privilegeTier, promoCode }) {
    return this._post('/api/v1/cart/evaluate', {
      cinema_id: cinemaId,
      items,
      privilege_tier: privilegeTier,
      promo_code: promoCode,
    });
  }

  // 5. Checkout
  async checkout({ cinemaId, items, paymentMethod, privilegeTier, promoCode, customerPhone, customerName }) {
    return this._post('/api/v1/checkout', {
      cinema_id: cinemaId,
      items,
      payment_method: paymentMethod,
      privilege_tier: privilegeTier,
      promo_code: promoCode,
      customer_phone: customerPhone,
      customer_name: customerName,
    });
  }
}

// Export singleton instance for global browser scripts
window.cinemaAPI = new CinemaClient();
