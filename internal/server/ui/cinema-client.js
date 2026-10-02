/**
 * 🍏 CineBite Cinema Recommender Client SDK & Universal Fallback Engine
 * Engineered for hybrid runtime: communicates with Go REST API when available,
 * and transparently executes client-side simulation on static hosts (e.g. GitHub Pages) or offline.
 */

// ── 1. Geolocation Math: Haversine Spherical Distance ──────────────────────────
function calculateHaversineDistance(lat1, lon1, lat2, lon2) {
  if (!lat1 || !lon1 || !lat2 || !lon2) return 9999;
  const R = 6371.0; // Earth's radius in kilometers
  const dLat = (lat2 - lat1) * (Math.PI / 180.0);
  const dLon = (lon2 - lon1) * (Math.PI / 180.0);
  const rLat1 = lat1 * (Math.PI / 180.0);
  const rLat2 = lat2 * (Math.PI / 180.0);

  const a = Math.sin(dLat / 2) * Math.sin(dLat / 2) +
            Math.sin(dLon / 2) * Math.sin(dLon / 2) * Math.cos(rLat1) * Math.cos(rLat2);
  const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));
  return Math.round(R * c * 10) / 10;
}

// ── 2. Sri Lankanconst SRI_LANKAN_CINEMAS = [
  {
    id: "CMB-CCC",
    name: "Scope Cinemas Multiplex - Colombo City Centre",
    chain: "Scope Cinemas",
    city: "Colombo",
    province: "Western Province",
    address: "Level 3, Colombo City Centre, 137 Sir James Pieris Mawatha, Colombo 02",
    screens: 6,
    is_outside_colombo: false,
    latitude: 6.9175,
    longitude: 79.8550,
    map_url: "https://maps.google.com/?q=6.9175,79.8550",
    image: "cinema-multiplex.jpg",
    technology: "IMAX Laser • Dolby Atmos 7.1",
    has_in_house_food: true,
    has_food_court_in_front: true,
    food_court_name: "The Food Studio @ CCC & Scope Candy Bar",
    food_court_location: "Level 3 Lobby directly in front of cinema entrance",
    food_place_name: "The Food Studio @ CCC",
    food_place_type: "In-Front Luxury Mall Food Court & Scope Candy Bar",
    nearby_food_options: ["The Food Studio (Level 3)", "Cargills Food City CCC", "Twist Colombo", "Shiok Singapore", "McDonald's CCC"],
    food_hours: "10:00 AM - 11:30 PM (Open during all movie screenings)",
    food_delivery_to_seat: true
  },
  {
    id: "CMB-HCM",
    name: "Scope Cinemas Multiplex - Havelock City Mall",
    chain: "Scope Cinemas",
    city: "Colombo",
    province: "Western Province",
    address: "Level 4, Havelock City Mall, 324 Havelock Road, Colombo 05",
    screens: 6,
    is_outside_colombo: false,
    latitude: 6.8778,
    longitude: 79.8653,
    map_url: "https://maps.google.com/?q=6.8778,79.8653",
    image: "cinema-imax.jpg",
    technology: "Sri Lanka's Largest IMAX • 4K Laser",
    has_in_house_food: true,
    has_food_court_in_front: true,
    food_court_name: "Havelock Food Lounge & Scope IMAX Candy Bar",
    food_court_location: "Level 4 Lobby directly in front of IMAX auditorium entrance",
    food_place_name: "Havelock Food Lounge & Scope IMAX Candy Bar",
    food_place_type: "In-Front Mall Food Lounge & IMAX Concession Stand",
    nearby_food_options: ["Havelock Food Lounge (Level 4)", "Baskin Robbins HCM", "Subway Havelock City", "Cargills Gourmet Food Hall"],
    food_hours: "10:00 AM - 11:30 PM (Open during all movie screenings)",
    food_delivery_to_seat: true
  },
  {
    id: "CMB-LBT",
    name: "Liberty by Scope Cinemas",
    chain: "Scope Cinemas",
    city: "Colombo",
    province: "Western Province",
    address: "35 R. A. De Mel Mawatha, Colombo 03",
    screens: 2,
    is_outside_colombo: false,
    latitude: 6.9073,
    longitude: 79.8517,
    map_url: "https://maps.google.com/?q=6.9073,79.8517",
    image: "cinema-heritage.jpg",
    technology: "Scope VIP Gold Class • Dolby 7.1",
    has_in_house_food: true,
    has_food_court_in_front: true,
    food_court_name: "Liberty Foyer Concession Court & Scope Candy Bar",
    food_court_location: "Front main entrance lobby directly before auditorium doors",
    food_place_name: "Liberty Foyer Concession Court",
    food_place_type: "In-Front Cinema Concession Food Court",
    nearby_food_options: ["Liberty Plaza Food Court", "Perera & Sons Kollupitiya", "Chit-Chat Coffee Corner", "Mackenzy's"],
    food_hours: "10:00 AM - 10:45 PM (Open during all movie screenings)",
    food_delivery_to_seat: true
  },
  {
    id: "KND-KCC",
    name: "Scope Partner Multiplex - KCC Kandy",
    chain: "Scope Cinemas Partner Circuit",
    city: "Kandy",
    province: "Central Province",
    address: "Level 3, Kandy City Centre, Dalada Veediya, Kandy",
    screens: 4,
    is_outside_colombo: true,
    latitude: 7.2936,
    longitude: 80.6385,
    map_url: "https://maps.google.com/?q=7.2936,80.6385",
    image: "cinema-multiplex.jpg",
    technology: "4K Laser Projection • Dolby Digital",
    has_in_house_food: true,
    has_food_court_in_front: true,
    food_court_name: "KCC World Food Court & Scope Candy Bar",
    food_court_location: "Level 3 Lobby directly in front of cinema entrance",
    food_place_name: "KCC World Food Court & Scope Candy Bar",
    food_place_type: "In-Front Mall Food Court & Cinema Concession Stand",
    nearby_food_options: ["KCC World Food Court (Level 4)", "Devon Restaurant & Bakery", "Bake House Dalada Veediya", "Cargills Food Hall"],
    food_hours: "10:00 AM - 10:45 PM (Open during all movie screenings)",
    food_delivery_to_seat: true
  },
  {
    id: "KND-REG",
    name: "Regal Cinema Kandy",
    chain: "Ceylon Theatres",
    city: "Kandy",
    province: "Central Province",
    address: "No. 12, Katugastota Road, Kandy",
    screens: 2,
    is_outside_colombo: true,
    latitude: 7.2985,
    longitude: 80.6335,
    map_url: "https://maps.google.com/?q=7.2985,80.6335",
    image: "cinema-heritage.jpg",
    technology: "Ceylon Theatres Classic • Dolby Surround",
    has_in_house_food: true,
    has_food_court_in_front: true,
    food_court_name: "Regal Foyer Concession Food Court",
    food_court_location: "Directly in front of cinema auditorium entrance",
    food_place_name: "Regal Foyer Concession Food Court",
    food_place_type: "In-Front Theater Snack Food Court",
    nearby_food_options: ["Katugastota Cafe & Bake House", "Perera & Sons Kandy", "White House Restaurant"],
    food_hours: "10:15 AM - 10:30 PM",
    food_delivery_to_seat: false
  },
  {
    id: "GMP-REG",
    name: "Regal Cinema Gampaha",
    chain: "Ceylon Theatres",
    city: "Gampaha",
    province: "Western Province (Outer Regional)",
    address: "Bauddhaloka Mawatha, Gampaha",
    screens: 2,
    is_outside_colombo: true,
    latitude: 7.0917,
    longitude: 79.9998,
    map_url: "https://maps.google.com/?q=7.0917,79.9998",
    image: "cinema-heritage.jpg",
    technology: "Ceylon Theatres • Digital Screen",
    has_in_house_food: true,
    has_food_court_in_front: true,
    food_court_name: "Regal Front Foyer Food Plaza & Snack Court",
    food_court_location: "Front entrance foyer directly facing ticket counter",
    food_place_name: "Regal Front Foyer Food Plaza",
    food_place_type: "In-Front Cinema Concession Food Court",
    nearby_food_options: ["Gampaha Food City Concessions", "Dinemore Gampaha", "Perera & Sons Bauddhaloka Mw"],
    food_hours: "10:00 AM - 10:30 PM",
    food_delivery_to_seat: false
  },
  {
    id: "GLE-QNS",
    name: "Queens Cinema Galle",
    chain: "Independent Heritage Screen",
    city: "Galle",
    province: "Southern Province",
    address: "Wakwella Road, Galle Fort Promenade",
    screens: 2,
    is_outside_colombo: true,
    latitude: 6.0535,
    longitude: 80.2210,
    map_url: "https://maps.google.com/?q=6.0535,80.2210",
    image: "cinema-heritage.jpg",
    technology: "Heritage Galle Screen • Dolby Digital",
    has_in_house_food: true,
    has_food_court_in_front: true,
    food_court_name: "Queens Foyer Candy Bar",
    food_court_location: "Main Foyer Entrance",
    food_place_name: "Queens Foyer Candy Bar",
    food_place_type: "In-Front Cinema Concessions",
    nearby_food_options: ["Galle Fort Dining Terrace", "Pedlar's Inn", "Southern Taste Bakery"],
    food_hours: "10:00 AM - 10:30 PM",
    food_delivery_to_seat: false
  },
  {
    id: "NGB-RDO",
    name: "Rodeo Cinema Negombo",
    chain: "Western Coast Network",
    city: "Negombo",
    province: "Western Province (Outer Coastal)",
    address: "Main Street, Negombo",
    screens: 2,
    is_outside_colombo: true,
    latitude: 7.2083,
    longitude: 79.8358,
    map_url: "https://maps.google.com/?q=7.2083,79.8358",
    image: "cinema-multiplex.jpg",
    technology: "Coastal Multiplex • RealD 3D",
    has_in_house_food: true,
    has_food_court_in_front: true,
    food_court_name: "Rodeo Beach Concession Court & Chill Bar",
    food_court_location: "Front theater foyer facing screen entrance",
    food_place_name: "Rodeo Beach Concession Court",
    food_place_type: "In-Front Concession Food Court",
    nearby_food_options: ["Lords Restaurant Negombo Beach", "Negombo Dutch Bakery", "King Coconut Stalls"],
    food_hours: "10:30 AM - 11:00 PM",
    food_delivery_to_seat: false
  },
  {
    id: "KRN-MLD",
    name: "Maldini Cinema Kurunegala",
    chain: "Maldini Circuits",
    city: "Kurunegala",
    province: "North Western Province",
    address: "Colombo Road, Kurunegala",
    screens: 2,
    is_outside_colombo: true,
    latitude: 7.4863,
    longitude: 80.3623,
    map_url: "https://maps.google.com/?q=7.4863,80.3623",
    image: "cinema-multiplex.jpg",
    technology: "North Western Cineplex • Dolby",
    has_in_house_food: true,
    has_food_court_in_front: true,
    food_court_name: "North Star Concession Lobby",
    food_court_location: "Lobby entrance",
    food_place_name: "North Star Concession Lobby",
    food_place_type: "In-Front Cinema Concessions",
    nearby_food_options: ["Clock Tower Cafe", "Perera & Sons Colombo Rd", "Royal Bakery"],
    food_hours: "10:00 AM - 10:00 PM",
    food_delivery_to_seat: false
  },
  {
    id: "MTR-SKC",
    name: "SK Cinema Matara",
    chain: "Southern Screen Network",
    city: "Matara",
    province: "Southern Province",
    address: "Anagarika Dharmapala Mawatha, Matara",
    screens: 2,
    is_outside_colombo: true,
    latitude: 5.9485,
    longitude: 80.5488,
    map_url: "https://maps.google.com/?q=5.9485,80.5488",
    image: "cinema-heritage.jpg",
    technology: "Southern Screen • Dolby Surround",
    has_in_house_food: true,
    has_food_court_in_front: true,
    food_court_name: "SK Concession Court & Chill Bar",
    food_court_location: "Main entrance lobby directly in front of cinema screens",
    food_place_name: "SK Concession Food Court",
    food_place_type: "In-Front Concession Food Court",
    nearby_food_options: ["Dutch Fort Cafe Matara", "Perera & Sons Matara", "Beach Road Fresh Juice"],
    food_hours: "10:30 AM - 10:30 PM",
    food_delivery_to_seat: false
  },
  {
    id: "JFN-RJA",
    name: "Rajah Cinema Jaffna",
    chain: "Northern Cinema Network",
    city: "Jaffna",
    province: "Northern Province",
    address: "Hospital Road, Jaffna",
    screens: 3,
    is_outside_colombo: true,
    latitude: 9.6615,
    longitude: 80.0255,
    map_url: "https://maps.google.com/?q=9.6615,80.0255",
    image: "cinema-heritage.jpg",
    technology: "Northern Star • 4K Digital Cinema",
    has_in_house_food: true,
    has_food_court_in_front: true,
    food_court_name: "Rajah Snack Point",
    food_court_location: "Main Entrance",
    food_place_name: "Rajah Snack Point",
    food_place_type: "In-Front Concession Bar",
    nearby_food_options: ["Mangos Vegetarian", "Rolex Cafe Jaffna", "Rio Ice Cream"],
    food_hours: "10:00 AM - 10:00 PM",
    food_delivery_to_seat: false
  }
];

// ── 3. Concession Items Catalog by City / Cinema ──────────────────────────────
const CONCESSION_CATALOG = {
  default: [
    { id: "P01", category: "Popcorn", name: "Scope Warm Caramel Popcorn Tub", size: "Jumbo Souvenir Tub", price: 1350, original_price: 1600, image: "popcorn.jpg", in_stock: true, price_trend: "down", description: "Scope signature crunchy caramelized warm corn in commemorative cinema souvenir tub" },
    { id: "P02", category: "Popcorn", name: "Kochchi & White Cheddar Popcorn", size: "Large Bucket", price: 1250, original_price: 1450, image: "popcorn.jpg", in_stock: true, price_trend: "down", description: "Hot butterfly corn tossed with fiery Sri Lankan green kochchi & Wisconsin cheddar dust" },
    { id: "P03", category: "Popcorn", name: "Classic Clarified Butter Salted Popcorn", size: "Large Bucket", price: 1100, original_price: 1300, image: "popcorn.jpg", in_stock: true, price_trend: "steady", description: "Freshly popped butterfly corn drenched in golden clarified butter and pure sea salt" },
    { id: "P04", category: "Popcorn", name: "Scope Double Chocolate Glazed Popcorn", size: "Large Bucket", price: 1250, original_price: 1450, image: "popcorn.jpg", in_stock: true, price_trend: "down", description: "Crunchy popped corn generously drizzled with rich Belgian milk and dark chocolate glaze" },
    { id: "B01", category: "Beverage", name: "Chilled Milo Dinosaur Float", size: "Tall Glass", price: 550, original_price: 650, image: "beverage.jpg", in_stock: true, price_trend: "down", description: "Chilled Ceylon chocolate malt Milo topped with vanilla ice cream and heaped cocoa powder" },
    { id: "B02", category: "Beverage", name: "Elephant House Ginger Beer (EGB) Float", size: "Regular 450ml", price: 480, original_price: 580, image: "beverage.jpg", in_stock: true, price_trend: "steady", description: "Authentic spicy Ceylon ginger beer with crushed ice and creamy vanilla ice cream float" },
    { id: "B03", category: "Beverage", name: "Scope Fresh Passion Fruit Mint Sparkler", size: "Medium 500ml", price: 750, original_price: 900, image: "beverage.jpg", in_stock: true, price_trend: "down", description: "Real island passion fruit pulp with crushed garden mint and sparkling soda" },
    { id: "B04", category: "Beverage", name: "Scope Large Fountain Coca-Cola", size: "Large 650ml", price: 650, original_price: 800, image: "beverage.jpg", in_stock: true, price_trend: "steady", description: "Ice-cold fountain Coca-Cola with fresh carbonation and lemon wedge" },
    { id: "B05", category: "Beverage", name: "Scope Creamy Iced Coffee Float", size: "Tall Glass", price: 680, original_price: 800, image: "beverage.jpg", in_stock: true, price_trend: "steady", description: "Rich brewed Ceylon coffee with sweetened milk, dark cocoa dust, and vanilla float" },
    { id: "S01", category: "Snack", name: "Scope Kitchen Kochchi Cheese Sticks (4 pcs)", size: "Hot Snack", price: 950, original_price: 1100, image: "rolls.jpg", in_stock: true, price_trend: "down", description: "Golden crumbed molten mozzarella sticks spiked with spicy green kochchi chili dip" },
    { id: "S02", category: "Snack", name: "Scope Kitchen Crispy Chicken Strips (4 pcs)", size: "Hot Snack", price: 1150, original_price: 1350, image: "rolls.jpg", in_stock: true, price_trend: "steady", description: "Golden crumbed tender chicken breast strips served with house garlic mayonnaise dip" },
    { id: "S03", category: "Snack", name: "Loaded Steak & Queso Fries", size: "Sharing Platter", price: 1450, original_price: 1700, image: "nachos.jpg", in_stock: true, price_trend: "down", description: "Crispy skin-on seasoned fries smothered in hot jalapeño queso and tender steak strips" },
    { id: "S04", category: "Snack", name: "IMAX Loaded Chicken Nachos", size: "Sharing Platter", price: 1450, original_price: 1700, image: "nachos.jpg", in_stock: true, price_trend: "down", description: "Stone-ground corn chips with hot queso, jalapeños, salsa, and seasoned chicken" },
    { id: "S05", category: "Snack", name: "Gourmet Brioche Chicken Hotdog", size: "Single", price: 1150, original_price: 1350, image: "hotdog.jpg", in_stock: true, price_trend: "steady", description: "Artisan grilled chicken sausage in buttered brioche bun with sweet pickle relish" },
    { id: "S06", category: "Snack", name: "Spicy Beef Chili Dog with Jalapeños", size: "Single", price: 1350, original_price: 1550, image: "hotdog.jpg", in_stock: true, price_trend: "down", description: "Grilled beef frankfurter loaded with spiced meat chili, melted cheddar, and jalapeño slices" },
    { id: "S07", category: "Snack", name: "Crispy Ceylon Chicken Rolls (2 pcs)", size: "Hot Snack", price: 820, original_price: 950, image: "rolls.jpg", in_stock: true, price_trend: "down", description: "Crispy golden crumbed Sri Lankan bakery-style spicy chicken rolls with chili dip" },
    { id: "S08", category: "Snack", name: "Crispy Vegetable Samosa Platter (4 pcs)", size: "Sharing Plate", price: 750, original_price: 900, image: "rolls.jpg", in_stock: true, price_trend: "steady", description: "Golden spiced potato and pea pastry triangles served with sweet tamarind chutney" },
    { id: "C01", category: "Combo", name: "Scope Director's Deluxe Couple Combo", size: "Duo Feast", price: 2950, original_price: 3600, image: "combo.jpg", in_stock: true, price_trend: "down", description: "1 Jumbo Caramel Popcorn + 2 Large Fountain Drinks + 1 Loaded Nachos Platter" },
    { id: "C02", category: "Combo", name: "CinePass VIP Platinum Movie Feast", size: "Family (3-4 Pax)", price: 4200, original_price: 5200, image: "combo.jpg", in_stock: true, price_trend: "down", description: "2 Jumbo Popcorns (Caramel + Cheese) + 3 Drinks + 1 Brioche Hotdog + 1 Nachos Platter" },
    { id: "C03", category: "Combo", name: "Scope Kitchen Movie Meal Box", size: "Solo Combo", price: 2250, original_price: 2700, image: "combo.jpg", in_stock: true, price_trend: "down", description: "1 Large Caramel Popcorn + 1 Kochchi Cheese Sticks (4 pcs) + 1 Chilled Milo Dinosaur Float" },
    { id: "C04", category: "Combo", name: "Solo Movie Snack Saver", size: "Single Pack", price: 1650, original_price: 2100, image: "combo.jpg", in_stock: true, price_trend: "down", description: "1 Large Clarified Butter Popcorn + 1 Medium Fountain Drink + 2 Ceylon Bakery Rolls" }
  ]
};

// ── 4. Main SDK Class ─────────────────────────────────────────────────────────
class CinemaClient {
  constructor(baseUrl = '') {
    this.baseUrl = baseUrl;
  }

  async _get(path) {
    try {
      const res = await fetch(`${this.baseUrl}${path}`);
      if (!res.ok) {
        throw new Error(`HTTP ${res.status}`);
      }
      return await res.json();
    } catch (err) {
      throw err;
    }
  }

  async _post(path, data) {
    try {
      const res = await fetch(`${this.baseUrl}${path}`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data),
      });
      if (!res.ok) {
        throw new Error(`HTTP ${res.status}`);
      }
      return await res.json();
    } catch (err) {
      throw err;
    }
  }

  // 1. Cinema Registry (with automatic Haversine proximity sort)
  async getCinemas({ lat = 0, lng = 0, city = '' } = {}) {
    try {
      const params = new URLSearchParams();
      if (city) params.append('city', city);
      if (lat && lng) {
        params.append('lat', lat);
        params.append('lng', lng);
      }
      const query = params.toString() ? `?${params.toString()}` : '';
      const data = await this._get(`/api/v1/cinemas${query}`);
      if (Array.isArray(data) && data.length > 0) {
        return data.map(c => {
          const fallback = SRI_LANKAN_CINEMAS.find(x => x.id === c.id) || {};
          return {
            ...fallback,
            ...c,
            image: c.image || fallback.image || (c.name.toLowerCase().includes('imax') ? 'cinema-imax.jpg' : (c.screens >= 4 ? 'cinema-multiplex.jpg' : 'cinema-heritage.jpg')),
            technology: c.technology || fallback.technology || (c.screens >= 4 ? `${c.screens} Screens • Dolby Atmos` : `${c.screens} Screens • Dolby 7.1`)
          };
        });
      }
    } catch (err) {
      // Gracefully fall back to client-side database
    }

    // Client-side fallback engine:
    let list = JSON.parse(JSON.stringify(SRI_LANKAN_CINEMAS));
    if (city) {
      const clean = city.toLowerCase().trim();
      list = list.filter(c => c.city.toLowerCase().includes(clean));
    }

    if (lat && lng) {
      list.forEach(c => {
        c.distance_km = calculateHaversineDistance(lat, lng, c.latitude, c.longitude);
      });
      list.sort((a, b) => (a.distance_km || 9999) - (b.distance_km || 9999));
    }

    return list;
  }

  async getCinema(cinemaId) {
    try {
      const c = await this._get(`/api/v1/cinemas/${encodeURIComponent(cinemaId)}`);
      const fallback = SRI_LANKAN_CINEMAS.find(x => x.id === c.id) || {};
      return {
        ...fallback,
        ...c,
        image: c.image || fallback.image || 'cinema-multiplex.jpg',
        technology: c.technology || fallback.technology || `${c.screens} Screens • Dolby Atmos`
      };
    } catch (err) {
      const found = SRI_LANKAN_CINEMAS.find(c => c.id === cinemaId);
      if (found) return found;
      return SRI_LANKAN_CINEMAS[0];
    }
  }

  // 2. Concessions & Live Ticker
  async getLiveConcessions(cinemaId) {
    try {
      const res = await this._get(`/api/v1/cinemas/${encodeURIComponent(cinemaId)}/concessions/live`);
      if (res && res.items) {
        res.items = res.items.map(item => ({
          ...item,
          price: item.price_lkr !== undefined ? item.price_lkr : (item.price || 0),
          original_price: item.base_price_lkr !== undefined ? item.base_price_lkr : (item.original_price || item.price_lkr || item.price || 0),
          image: item.image || (
            item.category === 'Popcorn' ? 'popcorn.jpg' :
            item.category === 'Beverage' ? 'beverage.jpg' :
            item.category === 'Combo' ? 'combo.jpg' :
            (item.name && item.name.toLowerCase().includes('dog') ? 'hotdog.jpg' :
             item.name && (item.name.toLowerCase().includes('roll') || item.name.toLowerCase().includes('samosa')) ? 'rolls.jpg' : 'nachos.jpg')
          )
        }));
      }
      return res;
    } catch (err) {
      const cinema = SRI_LANKAN_CINEMAS.find(c => c.id === cinemaId) || SRI_LANKAN_CINEMAS[0];
      return {
        cinema_id: cinema.id,
        cinema_name: cinema.name,
        currency: "LKR",
        items: CONCESSION_CATALOG.default.map(i => ({ ...i, cinema_id: cinema.id }))
      };
    }
  }

  async getCinemaAds(cinemaId) {
    try {
      return await this._get(`/api/v1/cinemas/${encodeURIComponent(cinemaId)}/ads`);
    } catch (err) {
      return [
        {
          id: "AD-01",
          cinema_id: cinemaId,
          title: "Combo Flash Deal",
          subtitle: "Save 30% on any duo or family combo this weekend",
          discount_text: "30% OFF",
          promo_code: "CINE-COMBO30",
          is_flash_ad: true,
          badge_text: "FLASH DEAL"
        },
        {
          id: "AD-02",
          cinema_id: cinemaId,
          title: "CinePass Silver+ Perk",
          subtitle: "Unlock 25% on all fountain beverages with Silver tier",
          discount_text: "25% OFF",
          promo_code: "CINE-CINEMA25",
          is_flash_ad: false,
          badge_text: "VIP PERK"
        },
        {
          id: "AD-03",
          cinema_id: cinemaId,
          title: "Student Special",
          subtitle: "Show student ID for 15% off full concession order",
          discount_text: "15% OFF",
          promo_code: "CINE-STUDENT",
          is_flash_ad: false,
          badge_text: "STUDENT"
        }
      ];
    }
  }

  // 3. Loyalty & Promotions
  async getPromotions() {
    try {
      return await this._get('/api/v1/promotions/privilege');
    } catch (err) {
      return [
        { promo_code: "CINE-COMBO30", discount_pct: 30, category_restriction: "Combo", description: "30% off any combo" },
        { promo_code: "CINE-CINEMA25", discount_pct: 25, category_restriction: "Beverage", description: "25% off beverages" },
        { promo_code: "CINE-STUDENT", discount_pct: 15, category_restriction: "", description: "15% off total order" },
        { promo_code: "CINE-PLATINUM", discount_pct: 35, category_restriction: "", description: "35% off VIP feast" }
      ];
    }
  }

  async getPaymentMethods() {
    try {
      return await this._get('/api/v1/payment-methods');
    } catch (err) {
      return [
        { id: "LANKAQR", name: "LankaQR", description: "CBSL National EMV Standard (Q+, WePay, FriMi, SOLO, SmartPay) - 5% Instant Rebate", is_digital: true },
        { id: "FRIMI", name: "FriMi (NTB)", description: "Instant 10% Concession Cashback", is_digital: true },
        { id: "GENIE", name: "Genie (Dialog)", description: "8% Digital Wallet Rebate", is_digital: true },
        { id: "EZCASH", name: "eZ Cash", description: "Direct USSD PIN Push for Dialog / Hutch / Airtel", is_digital: true },
        { id: "MCASH", name: "mCash", description: "Direct Carrier & Mobile Wallet for SLT-Mobitel", is_digital: true },
        { id: "KOKO", name: "Koko (BNPL)", description: "3 Interest-Free Monthly Installments", is_digital: true },
        { id: "CARD", name: "Visa / MasterCard / LankaPay", description: "Local Sri Lankan Bank Cards", is_digital: true },
        { id: "CASH", name: "Cash at Counter", description: "Pay in LKR at Theater Concession Stand", is_digital: false }
      ];
    }
  }

  // 4. Recommendation & Cart Optimizer (Client-side Knapsack Heuristic)
  async recommendBundle({ cinemaId, budget, partySize = 2, privilegeTier = 'Standard', promoCode = '' }) {
    try {
      return await this._post('/api/v1/recommend', {
        cinema_id: cinemaId,
        budget_lkr: budget,
        party_size: partySize,
        privilege_tier: privilegeTier,
        promo_code: promoCode,
      });
    } catch (err) {
      // Run algorithmic knapsack bundle optimization client-side
      const menu = CONCESSION_CATALOG.default;
      const budgetLimit = budget > 0 ? budget : partySize * 1500;

      const combos = menu.filter(i => i.category === 'Combo');
      const popcorns = menu.filter(i => i.category === 'Popcorn');
      const beverages = menu.filter(i => i.category === 'Beverage');
      const snacks = menu.filter(i => i.category === 'Snack');

      let selected = [];
      let subtotal = 0;

      // Strategy A: If combo fits within 85% of budget, pick best combo
      if (combos.length > 0 && budgetLimit >= combos[0].price) {
        for (let i = combos.length - 1; i >= 0; i--) {
          if (combos[i].price <= budgetLimit * 0.85) {
            selected.push(combos[i]);
            subtotal += combos[i].price;
            break;
          }
        }
      }

      // Strategy B: If no combo or remaining budget, add core items
      if (selected.length === 0) {
        // Add Popcorn
        if (popcorns.length > 0) {
          const pop = (partySize >= 2 && popcorns.length > 1) ? popcorns[0] : popcorns[1];
          selected.push(pop);
          subtotal += pop.price;
        }
        // Add Drinks for party
        for (let p = 0; p < partySize; p++) {
          const drink = beverages[p % beverages.length];
          if (subtotal + drink.price <= budgetLimit) {
            selected.push(drink);
            subtotal += drink.price;
          }
        }
        // Add Snack if budget permits
        for (const s of snacks) {
          if (subtotal + s.price <= budgetLimit) {
            selected.push(s);
            subtotal += s.price;
            break;
          }
        }
      } else {
        // Combo selected: top up with extra drinks or snacks if budget remains
        for (const b of beverages) {
          if (subtotal + b.price <= budgetLimit) {
            selected.push(b);
            subtotal += b.price;
            break;
          }
        }
      }

      // Apply Tier Discount (Standard: 0%, Silver: 15%, Gold: 20%, Platinum: 25%)
      let discountRate = 0.0;
      let appliedPromo = `${privilegeTier} Member`;
      if (privilegeTier === 'Silver') discountRate = 0.15;
      else if (privilegeTier === 'Gold') discountRate = 0.20;
      else if (privilegeTier === 'Platinum') discountRate = 0.25;

      // Promo Code Overrides
      const code = (promoCode || '').toUpperCase().trim();
      if (code === 'CINE-COMBO30') {
        discountRate = Math.max(discountRate, 0.30);
        appliedPromo = 'CINE-COMBO30 (30% Combo Flash)';
      } else if (code === 'CINE-CINEMA25') {
        discountRate = Math.max(discountRate, 0.25);
        appliedPromo = 'CINE-CINEMA25 (25% VIP Beverage)';
      } else if (code === 'CINE-STUDENT') {
        discountRate = Math.max(discountRate, 0.15);
        appliedPromo = 'CINE-STUDENT (15% Student Discount)';
      } else if (code === 'CINE-PLATINUM') {
        discountRate = Math.max(discountRate, 0.35);
        appliedPromo = 'CINE-PLATINUM (35% VIP Feast)';
      }

      const discountAmount = Math.round(subtotal * discountRate);
      const finalTotal = Math.max(0, subtotal - discountAmount);

      return {
        cinema_id: cinemaId,
        cinema_name: (SRI_LANKAN_CINEMAS.find(c => c.id === cinemaId) || {}).name || "Selected Cinema",
        budget_lkr: budgetLimit,
        party_size: partySize,
        privilege_tier: privilegeTier,
        applied_promo: appliedPromo,
        original_total: subtotal,
        final_total: finalTotal,
        discount_amount: discountAmount,
        items: selected,
        upsell_advice: `Add LKR ${Math.round(budgetLimit * 0.15 + 200)} more to upgrade drinks to Jumbo Caramel Popcorn!`
      };
    }
  }

  async evaluateCart({ cinemaId, items, privilegeTier = 'Standard', promoCode = '' }) {
    try {
      return await this._post('/api/v1/cart/evaluate', {
        cinema_id: cinemaId,
        items,
        privilege_tier: privilegeTier,
        promo_code: promoCode,
      });
    } catch (err) {
      const subtotal = items.reduce((s, i) => s + (i.price || 0), 0);
      let rate = 0.10;
      if (privilegeTier === 'Silver') rate = 0.15;
      else if (privilegeTier === 'Gold') rate = 0.20;
      else if (privilegeTier === 'Platinum') rate = 0.25;

      const discount = Math.round(subtotal * rate);
      return {
        cinema_id: cinemaId,
        items,
        original_total: subtotal,
        discount_amount: discount,
        final_total: subtotal - discount,
        applied_promo: promoCode || `${privilegeTier} Privilege Club`,
        privilege_tier: privilegeTier
      };
    }
  }

  // 5. Checkout
  async checkout({ cinemaId, items, paymentMethod, privilegeTier, promoCode, customerPhone, customerName }) {
    try {
      return await this._post('/api/v1/checkout', {
        cinema_id: cinemaId,
        items,
        payment_method: paymentMethod,
        privilege_tier: privilegeTier,
        promo_code: promoCode,
        customer_phone: customerPhone,
        customer_name: customerName,
      });
    } catch (err) {
      const subtotal = items.reduce((s, i) => s + (i.price || 0), 0);
      const discount = Math.round(subtotal * 0.15);
      const finalTotal = subtotal - discount;
      const serial = `CB-${cinemaId || 'PASS'}-${Math.floor(10000 + Math.random() * 90000)}-2026`;

      return {
        status: "confirmed",
        transaction_id: "TXN-" + Math.floor(100000 + Math.random() * 900000),
        voucher_code: serial,
        pickup_counter: "Concession Express Counter 1",
        estimated_ready: "5-8 minutes",
        lanka_qr_payload: `00020101021226580009LK.CBSL.QR0112${Math.floor(100000000000 + Math.random() * 900000000000)}520458145303144540${finalTotal.toFixed(2)}5802LK5910CINEBITE6007COLOMBO6304ABCD`,
        message: "Order placed successfully! Present digital pass at cinema concession counter."
      };
    }
  }
}

// Export singleton instance for global browser scripts
window.CinemaClient = CinemaClient;
window.cinemaAPI = new CinemaClient();
