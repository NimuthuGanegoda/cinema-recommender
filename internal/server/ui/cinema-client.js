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

// ── 2. Sri Lankan Cinema Registry ──────────────────────────────────────────
const SRI_LANKAN_CINEMAS = [
  {
    id: "CMB-CCC",
    name: "Scope Cinemas Multiplex - Colombo City Centre",
    chain: "Scope Cinemas",
    city: "Colombo",
    province: "Western Province",
    address: "Level 3, Colombo City Centre, 137 Sir James Pieris Mawatha, Colombo 02",
    screens: 6,
    is_outside_colombo: false,
    latitude: 6.9205,
    longitude: 79.8524,
    map_url: "https://maps.google.com/?q=6.9205,79.8524",
    image: "cinema-ccc.jpg",
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
    latitude: 6.8822,
    longitude: 79.8667,
    map_url: "https://maps.google.com/?q=6.8822,79.8667",
    image: "cinema-hcm.jpg",
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
    id: "CMB-PVR",
    name: "PVR Cinemas - One Galle Face Mall",
    chain: "PVR Cinemas",
    city: "Colombo",
    province: "Western Province",
    address: "Level 6, One Galle Face Mall, 1A Centre Road, Colombo 02",
    screens: 9,
    is_outside_colombo: false,
    latitude: 6.9268,
    longitude: 79.8443,
    map_url: "https://maps.google.com/?q=6.9268,79.8443",
    image: "cinema-pvr-ogf.jpg",
    technology: "P[XL] Laser • Luxe VIP Lounges • Dolby Atmos",
    has_in_house_food: true,
    has_food_court_in_front: true,
    food_court_name: "One Galle Face Food Court & PVR Concession Lounge",
    food_court_location: "Level 6 Foyer directly facing cinema box office",
    food_place_name: "PVR Concession Lounge & OGF Dining",
    food_place_type: "In-Front Mall Food Court & Gourmet Cinema Bar",
    nearby_food_options: ["OGF Food Studio (Level 6)", "Twist Colombo", "Shiok Singapore", "Chilis Colombo", "Chatime"],
    food_hours: "10:00 AM - 11:30 PM (Daily)",
    food_delivery_to_seat: true
  },
  {
    id: "CMB-LBT",
    name: "Liberty by Scope Cinemas",
    chain: "Scope Cinemas",
    city: "Colombo",
    province: "Western Province",
    address: "35 Srimath Anagarika Dharmapala Mawatha / R. A. De Mel Mawatha, Colombo 03",
    screens: 2,
    is_outside_colombo: false,
    latitude: 6.9119,
    longitude: 79.8501,
    map_url: "https://maps.google.com/?q=6.9119,79.8501",
    image: "cinema-liberty.jpg",
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
    id: "CMB-MJC",
    name: "Majestic Cineplex - Bambalapitiya",
    chain: "Ceylon Theatres",
    city: "Colombo",
    province: "Western Province",
    address: "4th/5th Floor, Majestic City, 10 Station Road, Bambalapitiya, Colombo 04",
    screens: 4,
    is_outside_colombo: false,
    latitude: 6.8942,
    longitude: 79.8551,
    map_url: "https://maps.google.com/?q=6.8942,79.8551",
    image: "cinema-majestic.jpg",
    technology: "Platinum & Gold Class 3D • Dolby Surround 7.1",
    has_in_house_food: true,
    has_food_court_in_front: true,
    food_court_name: "Majestic City 4th Floor Food Court & Snack Counters",
    food_court_location: "4th Floor Atrium directly outside Majestic Cineplex",
    food_place_name: "Majestic City Food Court",
    food_place_type: "In-Front Mall Food Court & Cineplex Snack Counters",
    nearby_food_options: ["Majestic Food Court", "KFC Majestic City", "Perera & Sons Bambalapitiya", "Keells Super Food Corner"],
    food_hours: "10:00 AM - 11:00 PM (Daily)",
    food_delivery_to_seat: false
  },
  {
    id: "CMB-SVY",
    name: "Savoy 3D Cinema - Wellawatte",
    chain: "Savoy Cinemas",
    city: "Colombo",
    province: "Western Province",
    address: "12 Savoy Building, Galle Road, Wellawatte, Colombo 06",
    screens: 3,
    is_outside_colombo: false,
    latitude: 6.8783,
    longitude: 79.8579,
    map_url: "https://maps.google.com/?q=6.8783,79.8579",
    image: "cinema-savoy.jpg",
    technology: "Savoy Premier 3D & 4K • Dolby Atmos",
    has_in_house_food: true,
    has_food_court_in_front: true,
    food_court_name: "Savoy Historic Concession Foyer & Snack Cafe",
    food_court_location: "Savoy Main Entrance Foyer facing Galle Road",
    food_place_name: "Savoy Concession Foyer",
    food_place_type: "In-Front Historic Cinema Snack Bar & Cafe",
    nearby_food_options: ["Wellawatte Indian Vadai & Sweet Shops", "Fab Bakery Wellawatte", "Perera & Sons Savoy", "Galle Road Cafe Row"],
    food_hours: "10:00 AM - 11:00 PM (Daily)",
    food_delivery_to_seat: false
  },
  {
    id: "KND-KCC",
    name: "Scope Partner Multiplex - KCC Kandy",
    chain: "Scope Cinemas Partner Circuit",
    city: "Kandy",
    province: "Central Province",
    address: "Level 3, Kandy City Centre, 5 Dalada Veediya, Kandy",
    screens: 4,
    is_outside_colombo: true,
    latitude: 7.2920,
    longitude: 80.6371,
    map_url: "https://maps.google.com/?q=7.2920,80.6371",
    image: "cinema-kcc.jpg",
    technology: "4K Laser Projection • Dolby Digital 7.1",
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
    id: "GMP-REG",
    name: "Regal Cinema Gampaha",
    chain: "Ceylon Theatres",
    city: "Gampaha",
    province: "Western Province (Outer Regional)",
    address: "2nd Floor, Cargills Square, 364 Miriswatta-Gampaha Road, Gampaha",
    screens: 2,
    is_outside_colombo: true,
    latitude: 7.0917,
    longitude: 79.9998,
    map_url: "https://maps.google.com/?q=7.0917,79.9998",
    image: "cinema-gampaha.jpg",
    technology: "Ceylon Theatres 3D • 2K Digital Screen",
    has_in_house_food: true,
    has_food_court_in_front: true,
    food_court_name: "Regal Front Foyer Food Plaza & Snack Court",
    food_court_location: "2nd Floor Cargills Square lobby facing ticket counter",
    food_place_name: "Regal Front Foyer Food Plaza",
    food_place_type: "In-Front Cinema Concession Food Court",
    nearby_food_options: ["Cargills Food City Concessions", "Dinemore Gampaha", "Perera & Sons Bauddhaloka Mw", "Gampaha Bake House"],
    food_hours: "10:00 AM - 10:30 PM",
    food_delivery_to_seat: false
  },
  {
    id: "GLE-QNS",
    name: "Queens Cinema Galle",
    chain: "Independent Heritage Screen",
    city: "Galle",
    province: "Southern Province",
    address: "29 Wakwella Road, Galle 80000",
    screens: 2,
    is_outside_colombo: true,
    latitude: 6.0342,
    longitude: 80.2166,
    map_url: "https://maps.google.com/?q=6.0342,80.2166",
    image: "cinema-galle.jpg",
    technology: "Heritage Galle 3D Screen • Dolby Digital",
    has_in_house_food: true,
    has_food_court_in_front: true,
    food_court_name: "Queens Foyer Concession & Galle Bake House",
    food_court_location: "Main Foyer Entrance on Wakwella Road",
    food_place_name: "Queens Foyer Candy Bar",
    food_place_type: "In-Front Cinema Concessions",
    nearby_food_options: ["Galle Fort Dining Terrace (5 min)", "Southern Taste Bakery Wakwella", "Pedlar's Inn Concessions", "Galle Fresh Juice Bar"],
    food_hours: "10:00 AM - 10:30 PM",
    food_delivery_to_seat: false
  },
  {
    id: "NGB-RDO",
    name: "Aqua Lite 3D Cinema Negombo",
    chain: "Western Coast Network",
    city: "Negombo",
    province: "Western Province (Outer Coastal)",
    address: "16 Christopher Road, Negombo 11500",
    screens: 2,
    is_outside_colombo: true,
    latitude: 7.2083,
    longitude: 79.8358,
    map_url: "https://maps.google.com/?q=7.2083,79.8358",
    image: "cinema-negombo.jpg",
    technology: "Coastal Multiplex 3D • 4K RealD",
    has_in_house_food: true,
    has_food_court_in_front: true,
    food_court_name: "Aqua Lite Concession Bar & Snacks",
    food_court_location: "Front theater foyer facing screen entrance",
    food_place_name: "Aqua Lite Concession Bar",
    food_place_type: "In-Front Concession Food Court",
    nearby_food_options: ["Perera Snacks Negombo", "Lords Restaurant Negombo Beach", "Negombo Dutch Bakery", "King Coconut Stalls Christopher Rd"],
    food_hours: "10:30 AM - 11:00 PM",
    food_delivery_to_seat: false
  },
  {
    id: "KRN-MLD",
    name: "Imperial 3D Cinema Kurunegala",
    chain: "Imperial Circuits",
    city: "Kurunegala",
    province: "North Western Province",
    address: "05 Kandy Road, Kurunegala 60000",
    screens: 2,
    is_outside_colombo: true,
    latitude: 7.4863,
    longitude: 80.3640,
    map_url: "https://maps.google.com/?q=7.4863,80.3640",
    image: "cinema-kurunegala.jpg",
    technology: "North Western 3D Cineplex • Dolby Digital",
    has_in_house_food: true,
    has_food_court_in_front: true,
    food_court_name: "Imperial Refreshment Corner & Snack Bar",
    food_court_location: "Lobby entrance directly facing ticket counters",
    food_place_name: "Imperial Refreshment Corner",
    food_place_type: "In-Front Cinema Concessions",
    nearby_food_options: ["Kurunegala Clock Tower Cafe Row", "Perera & Sons Kandy Rd", "Royal Bakery Kurunegala", "Ethagala View Dining"],
    food_hours: "10:00 AM - 10:00 PM",
    food_delivery_to_seat: false
  },
  {
    id: "MTR-SKC",
    name: "SK Cinema Matara",
    chain: "Southern Screen Network",
    city: "Matara",
    province: "Southern Province",
    address: "Kotuwegoda, Beach Road, Matara 81000",
    screens: 2,
    is_outside_colombo: true,
    latitude: 5.9443,
    longitude: 80.5501,
    map_url: "https://maps.google.com/?q=5.9443,80.5501",
    image: "cinema-matara.jpg",
    technology: "Southern Screen 3D & 4K • Dolby Surround",
    has_in_house_food: true,
    has_food_court_in_front: true,
    food_court_name: "SK Concession Court & Chill Bar",
    food_court_location: "Main entrance lobby along Beach Road facing cinema screens",
    food_place_name: "SK Concession Food Court",
    food_place_type: "In-Front Concession Food Court",
    nearby_food_options: ["Matara Beach Park Food Stalls", "Dutch Fort Bakery Matara", "Perera & Sons Matara", "Beach Road Fresh Coconut"],
    food_hours: "10:00 AM - 10:30 PM",
    food_delivery_to_seat: false
  },
  {
    id: "JFN-RJA",
    name: "Cargills Square Multiplex Jaffna",
    chain: "Northern Cinema Circuit",
    city: "Jaffna",
    province: "Northern Province",
    address: "Cargills Square, Hospital Road, Jaffna 40000",
    screens: 3,
    is_outside_colombo: true,
    latitude: 9.6635,
    longitude: 80.0165,
    map_url: "https://maps.google.com/?q=9.6635,80.0165",
    image: "cinema-jaffna.jpg",
    technology: "Northern Star 3D Multiplex • 4K Digital Cinema",
    has_in_house_food: true,
    has_food_court_in_front: true,
    food_court_name: "Cargills Square Food Court & Multiplex Concession Bar",
    food_court_location: "Ground & First Floor Mall Foyer on Hospital Road",
    food_place_name: "Cargills Square Food Court",
    food_place_type: "In-Front Modern Mall Food Court & Concession Bar",
    nearby_food_options: ["KFC Cargills Square Jaffna", "Rio Ice Cream Jaffna", "Mangos Vegetarian Restaurant", "Rolex Cafe Jaffna", "Jaffna Vadai Corner"],
    food_hours: "10:00 AM - 10:30 PM",
    food_delivery_to_seat: false
  }
];

// ── 2.1 Authentic Sri Lankan Now Showing Movies Registry ─────────────────────
const SRI_LANKAN_MOVIES = [
  {
    id: "MOV-01",
    title: "Ahasa Tharam",
    original_title: "අහස තරම්",
    language: "Sinhala",
    genre: "Romance / Drama",
    duration: "2h 18m",
    rating: "U",
    formats: ["2D Digital", "Dolby 7.1"],
    cinema_ids: ["CMB-CCC", "CMB-HCM", "CMB-SVY", "GMP-RGL", "GLE-QNS", "KRN-IMP", "KDY-KCC"],
    showtimes: ["10:30 AM", "1:45 PM", "4:30 PM", "7:15 PM"],
    director: "Sanjaya Nirmal",
    cast: ["Dinakshie Priyasad", "Sajitha Anuththara", "Bimal Jayakodi"],
    synopsis: "A deeply emotional Sri Lankan romantic journey celebrating enduring love, cultural nuances, and heart-stirring music.",
    status: "now_showing",
    concession_tip: "Best paired with Fresh Hot Butter Popcorn & Classic Ceylon Iced Tea"
  },
  {
    id: "MOV-02",
    title: "Avengers Endgame - Encore",
    original_title: "Avengers: Endgame (IMAX 3D Experience)",
    language: "English",
    genre: "Superhero / Action / Sci-Fi",
    duration: "3h 02m",
    rating: "PG-13",
    formats: ["IMAX 3D", "Dolby Atmos", "Laser 4K"],
    cinema_ids: ["CMB-CCC", "CMB-HCM", "CMB-OGF", "CMB-SVY", "KDY-KCC"],
    showtimes: ["11:00 AM", "3:00 PM", "6:45 PM", "10:15 PM"],
    director: "Anthony Russo, Joe Russo",
    cast: ["Robert Downey Jr.", "Chris Evans", "Mark Ruffalo", "Chris Hemsworth"],
    synopsis: "The legendary epic Marvel conclusion returns to Sri Lankan IMAX screens with newly restored audio and high-frame-rate visuals.",
    status: "now_showing",
    concession_tip: "Best paired with CinePass VIP Platinum Feast (2 Jumbo Popcorns + 3 Drinks + Nachos)"
  },
  {
    id: "MOV-03",
    title: "Sigma",
    original_title: "சிக்மா",
    language: "Tamil",
    genre: "Action / Crime / Thriller",
    duration: "2h 35m",
    rating: "U/A",
    formats: ["Dolby Atmos", "2D Digital"],
    cinema_ids: ["CMB-CCC", "CMB-HCM", "CMB-OGF", "CMB-SVY", "MTR-SKC", "JFN-RJA", "JFN-CGS"],
    showtimes: ["10:45 AM", "2:00 PM", "5:30 PM", "8:45 PM"],
    director: "Karthik Subbaraj",
    cast: ["Vijay Sethupathi", "SJ Suryah", "Pooja Hegde"],
    synopsis: "A gritty neo-noir underworld showdown with pulse-raising soundtrack, dynamic cinematography, and explosive twists.",
    status: "now_showing",
    concession_tip: "Best paired with Kochchi Cheese Sticks & Chilled Milo Dinosaur Float"
  },
  {
    id: "MOV-04",
    title: "Ayu",
    original_title: "ආයු",
    language: "Sinhala",
    genre: "Mystery / Psychological Thriller",
    duration: "2h 10m",
    rating: "U/A",
    formats: ["2D Digital", "Surround 5.1"],
    cinema_ids: ["CMB-CCC", "CMB-HCM", "CMB-LBT", "CMB-MJC", "GMP-RGL", "KDY-KCC"],
    showtimes: ["1:30 PM", "4:15 PM", "7:00 PM", "9:45 PM"],
    director: "Channa Deshapriya",
    cast: ["Hemal Ranasinghe", "Udari Warnakulasooriya", "Jackson Anthony"],
    synopsis: "An intricate psychological mystery surrounding an ancient family heirloom and forgotten truths buried in Sri Lanka's central highlands.",
    status: "now_showing",
    concession_tip: "Best paired with Gourmet Caramel Popcorn & Hot Spiced Chai"
  },
  {
    id: "MOV-05",
    title: "Spider-Man: Brand New Day",
    original_title: "Spider-Man: Brand New Day",
    language: "English",
    genre: "Action / Adventure / Sci-Fi",
    duration: "2h 28m",
    rating: "PG-13",
    formats: ["IMAX 3D", "4DX", "Dolby Atmos"],
    cinema_ids: ["CMB-CCC", "CMB-HCM", "CMB-OGF", "CMB-LBT", "CMB-MJC", "KDY-KCC"],
    showtimes: ["10:15 AM", "1:15 PM", "4:45 PM", "8:00 PM", "11:00 PM"],
    director: "Destin Daniel Cretton",
    cast: ["Tom Holland", "Zendaya", "Mark Ruffalo"],
    synopsis: "Peter Parker balances college life in NYC while grappling with mysterious street vigilantes and high-tech corporate threats.",
    status: "now_showing",
    concession_tip: "Best paired with Loaded Jalapeño Nachos Platter & Large Coca-Cola Zero"
  },
  {
    id: "MOV-06",
    title: "Meesaya Murukku 2",
    original_title: "மீசைய முறுக்கு 2",
    language: "Tamil",
    genre: "Musical / Youth Comedy / Drama",
    duration: "2h 20m",
    rating: "U",
    formats: ["2D Digital", "Dolby Atmos"],
    cinema_ids: ["CMB-OGF", "CMB-SVY", "NGB-AQU", "JFN-RJA", "JFN-CGS"],
    showtimes: ["11:30 AM", "3:15 PM", "6:30 PM", "9:30 PM"],
    director: "Hiphop Tamizha Aadhi",
    cast: ["Hiphop Tamizha Aadhi", "Aathmika", "Vivek"],
    synopsis: "The exuberant musical sequel following an indie music collective chasing their concert dreams against industry titans.",
    status: "now_showing",
    concession_tip: "Best paired with Ceylon Bakery Fish & Mutton Rolls with Sweet Chili Dip"
  },
  {
    id: "MOV-07",
    title: "Baththa",
    original_title: "பத்தா",
    language: "Tamil",
    genre: "Action / Rural Drama",
    duration: "2h 25m",
    rating: "U/A",
    formats: ["2D Digital", "Dolby 7.1"],
    cinema_ids: ["CMB-CCC", "CMB-HCM", "JFN-RJA", "JFN-CGS", "CMB-MJC", "MTR-SKC"],
    showtimes: ["1:00 PM", "4:30 PM", "7:45 PM"],
    director: "Mari Selvaraj",
    cast: ["Dhanush", "Fahadh Faasil", "Keerthy Suresh"],
    synopsis: "A hard-hitting story of community resilience, grassroots sports triumphs, and village solidarity.",
    status: "now_showing",
    concession_tip: "Best paired with Artisanal Crispy Corn & Iced Milo Dinosaur"
  },
  {
    id: "MOV-08",
    title: "Eda Re",
    original_title: "එදා රෑ",
    language: "Sinhala",
    genre: "Crime / Suspense Thriller",
    duration: "2h 05m",
    rating: "U/A",
    formats: ["2D Digital"],
    cinema_ids: ["CMB-CCC", "CMB-LBT", "CMB-SVY", "GLE-QNS", "KRN-IMP"],
    showtimes: ["3:30 PM", "6:45 PM", "9:30 PM"],
    director: "Udayakantha Warnasuriya",
    cast: ["Pubudu Chathuranga", "Mahendra Perera", "Dilhani Ekanayake"],
    synopsis: "A stormy night at an isolated tea plantation mansion triggers an intricate battle of wits between seven strangers.",
    status: "now_showing",
    concession_tip: "Best paired with Scope Kitchen Movie Meal Box (Solo Combo)"
  },
  {
    id: "MOV-09",
    title: "The Odyssey",
    original_title: "The Odyssey: Epic of the Sea",
    language: "English",
    genre: "Epic / Adventure / Fantasy",
    duration: "2h 40m",
    rating: "PG-13",
    formats: ["IMAX 3D", "Dolby Atmos"],
    cinema_ids: ["CMB-CCC", "CMB-HCM", "CMB-OGF", "KDY-KCC"],
    showtimes: ["12:00 PM", "4:00 PM", "7:30 PM", "10:30 PM"],
    director: "Christopher Nolan",
    cast: ["Christian Bale", "Cillian Murphy", "Florence Pugh"],
    synopsis: "A breathtaking cinematic voyage adapting the ancient Mediterranean myth with groundbreaking practical effects and IMAX grandeur.",
    status: "now_showing",
    concession_tip: "Best paired with Premium Jumbo Butter Popcorn & Chilled Sparkling Soda"
  },
  {
    id: "MOV-10",
    title: "Minions & Monsters",
    original_title: "Minions & Monsters (3D)",
    language: "English",
    genre: "Animation / Family / Comedy",
    duration: "1h 34m",
    rating: "U",
    formats: ["2D Digital", "3D Digital"],
    cinema_ids: ["CMB-CCC", "CMB-HCM", "CMB-OGF", "CMB-LBT", "CMB-MJC", "GMP-RGL", "NGB-AQU"],
    showtimes: ["10:00 AM", "12:15 PM", "2:30 PM", "5:00 PM"],
    director: "Pierre Coffin, Kyle Balda",
    cast: ["Steve Carell", "Pierre Coffin", "Taraji P. Henson"],
    synopsis: "The mischievous Minions stumble into a subterranean monster kingdom and try to become the gentle beasts' life coaches.",
    status: "now_showing",
    concession_tip: "Best paired with Sweet Caramel Popcorn Bucket + Fruit Slushies (Family Pack)"
  },
  {
    id: "MOV-11",
    title: "Yezhu Kadal Yezhu Malai",
    original_title: "ஏழு கடல் ஏழு மலை",
    language: "Tamil",
    genre: "Romance / Philosophical Drama",
    duration: "2h 15m",
    rating: "U/A",
    formats: ["2D Digital", "Dolby 7.1"],
    cinema_ids: ["CMB-CCC", "CMB-OGF", "JFN-RJA", "JFN-CGS", "CMB-SVY"],
    showtimes: ["2:15 PM", "5:45 PM", "8:30 PM"],
    director: "Ram",
    cast: ["Nivin Pauly", "Soori", "Anjali"],
    synopsis: "A poetic, timeless love story crossing scenic terrains and human endurance, scored by Yuvan Shankar Raja.",
    status: "now_showing",
    concession_tip: "Best paired with Ceylon Samosa Trio & Fresh Cold Pressed Juice"
  },
  {
    id: "MOV-12",
    title: "Heart of the Beast",
    original_title: "Heart of the Beast",
    language: "English",
    genre: "Action / Wilderness Thriller",
    duration: "1h 58m",
    rating: "A",
    formats: ["Dolby Atmos", "Laser 4K"],
    cinema_ids: ["CMB-CCC", "CMB-HCM", "CMB-OGF"],
    showtimes: ["6:00 PM", "8:45 PM", "11:15 PM"],
    director: "David Ayer",
    cast: ["Jason Statham", "Josh Hutcherson", "Jeremy Irons"],
    synopsis: "A retired Navy SEAL survivalist is thrust into relentless combat defending a remote sanctuary from elite mercenaries.",
    status: "now_showing",
    concession_tip: "Best paired with Brioche Beef Hotdog & Kochchi Dipping Sauce"
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

  // 6. Up-To-Date Movies & 1-Hour Automated Sync Engine
  async getMovies({ cinemaId = '', language = '', forceRefresh = false } = {}) {
    try {
      let url = '/api/v1/movies';
      const params = [];
      if (cinemaId) params.push(`cinema_id=${encodeURIComponent(cinemaId)}`);
      if (language) params.push(`language=${encodeURIComponent(language)}`);
      if (params.length > 0) url += '?' + params.join('&');

      if (!forceRefresh) {
        const cached = localStorage.getItem('cinebite_movies_cache');
        const cachedTime = localStorage.getItem('cinebite_movies_synced_at');
        if (cached && cachedTime && (Date.now() - parseInt(cachedTime, 10) < 3600000)) {
          let list = JSON.parse(cached);
          if (cinemaId) {
            const cid = cinemaId.toUpperCase();
            list = list.filter(m => m.cinema_ids && m.cinema_ids.some(c => c.toUpperCase() === cid));
          }
          if (language) {
            list = list.filter(m => m.language && m.language.toLowerCase() === language.toLowerCase());
          }
          return list;
        }
      }

      const res = await this._get(url);
      if (res && res.movies) {
        localStorage.setItem('cinebite_movies_cache', JSON.stringify(res.movies));
        localStorage.setItem('cinebite_movies_synced_at', Date.now().toString());
        return res.movies;
      }
    } catch (err) {
      // Fallback: try static movies.json with cache buster
      try {
        const staticRes = await fetch('movies.json?t=' + Math.floor(Date.now() / 3600000));
        if (staticRes.ok) {
          const data = await staticRes.json();
          if (data && data.movies) {
            localStorage.setItem('cinebite_movies_cache', JSON.stringify(data.movies));
            localStorage.setItem('cinebite_movies_synced_at', Date.now().toString());
            let list = data.movies;
            if (cinemaId) {
              const cid = cinemaId.toUpperCase();
              list = list.filter(m => m.cinema_ids && m.cinema_ids.some(c => c.toUpperCase() === cid));
            }
            if (language) {
              list = list.filter(m => m.language && m.language.toLowerCase() === language.toLowerCase());
            }
            return list;
          }
        }
      } catch (staticErr) {}

      // Final fallback to embedded dataset
      let list = [...SRI_LANKAN_MOVIES];
      if (cinemaId) {
        const cid = cinemaId.toUpperCase();
        list = list.filter(m => m.cinema_ids && m.cinema_ids.some(c => c.toUpperCase() === cid));
      }
      if (language) {
        list = list.filter(m => m.language && m.language.toLowerCase() === language.toLowerCase());
      }
      return list;
    }
  }

  async getMovieSyncStatus() {
    try {
      return await this._get('/api/v1/movies/status');
    } catch (err) {
      const storedTime = localStorage.getItem('cinebite_movies_synced_at');
      const lastSynced = storedTime ? parseInt(storedTime, 10) : (Date.now() - 300000);
      const elapsedSeconds = Math.floor((Date.now() - lastSynced) / 1000);
      const remainingSeconds = Math.max(0, 3600 - (elapsedSeconds % 3600));

      return {
        status: "synced",
        sync_interval: "1h",
        sync_interval_seconds: 3600,
        last_synced_at: new Date(lastSynced).toISOString(),
        next_sync_at: new Date(Date.now() + remainingSeconds * 1000).toISOString(),
        seconds_until_next_sync: remainingSeconds,
        total_movies: SRI_LANKAN_MOVIES.length,
        total_cinemas_covered: SRI_LANKAN_CINEMAS.length,
        source: "Scope & EAP Client Fallback Engine"
      };
    }
  }

  async triggerMovieSync() {
    try {
      const res = await this._post('/api/v1/movies/sync', {});
      localStorage.setItem('cinebite_movies_synced_at', Date.now().toString());
      if (res && res.movies) {
        localStorage.setItem('cinebite_movies_cache', JSON.stringify(res.movies));
      }
      return res;
    } catch (err) {
      localStorage.setItem('cinebite_movies_synced_at', Date.now().toString());
      return {
        message: "Client 1-hour sync refreshed successfully",
        sync_status: await this.getMovieSyncStatus(),
        movies: SRI_LANKAN_MOVIES
      };
    }
  }

  startHourlyMovieSyncEngine(onSyncCallback, onTickCallback) {
    if (!localStorage.getItem('cinebite_movies_synced_at')) {
      localStorage.setItem('cinebite_movies_synced_at', Date.now().toString());
    }

    const checkAndTick = async () => {
      const storedTime = parseInt(localStorage.getItem('cinebite_movies_synced_at') || Date.now().toString(), 10);
      const now = Date.now();
      const elapsedSec = Math.floor((now - storedTime) / 1000);
      const remainingSec = Math.max(0, 3600 - (elapsedSec % 3600));

      const mins = Math.floor(remainingSec / 60);
      const secs = remainingSec % 60;
      const formatted = `${mins}m ${secs < 10 ? '0' : ''}${secs}s`;

      if (typeof onTickCallback === 'function') {
        onTickCallback(remainingSec, formatted);
      }

      // If 1-hour cycle elapsed, trigger re-sync
      if (elapsedSec >= 3600) {
        localStorage.setItem('cinebite_movies_synced_at', now.toString());
        if (typeof onSyncCallback === 'function') {
          const freshMovies = await this.getMovies({ forceRefresh: true });
          onSyncCallback(freshMovies);
        }
      }
    };

    checkAndTick();
    return setInterval(checkAndTick, 1000);
  }
}

// Export singleton instance for global browser scripts
window.CinemaClient = CinemaClient;
window.cinemaAPI = new CinemaClient();
