const CACHE_NAME = 'cinebite-v12-print-layout';
const STATIC_ASSETS = [
  './',
  './index.html',
  './movies.html',
  './cinemas.html',
  './menu.html',
  './optimizer.html',
  './deals.html',
  './food-plan.html',
  './mobile.html',
  './assets/css/apple-design-system.css',
  './assets/css/cinema-polish.css',
  './assets/js/cinema-environment.js',
  './photo-credits.html',
  './assets/js/cinema-client.js',
  './assets/js/cinema-browser.js',
  './assets/js/cinema-worker.js',
  './assets/wasm/wasm_exec.js',
  './assets/wasm/cinebite.wasm',
  './data/movies.json',
  './manifest.json',
  './assets/images/app-icon.jpg',
  './assets/images/popcorn.jpg',
  './assets/images/beverage.jpg',
  './assets/images/combo.jpg',
  './assets/images/hotdog.jpg',
  './assets/images/nachos.jpg',
  './assets/images/rolls.jpg',
  './assets/images/food-caramel-popcorn.jpg'
,  './assets/images/food-butter-popcorn.jpg'
,  './assets/images/food-chocolate-popcorn.jpg'
,  './assets/images/food-cheese-sticks.jpg'
,  './assets/images/food-chicken-strips.jpg'
,  './assets/images/food-fries.jpg'
,  './assets/images/food-samosas.jpg'
,  './assets/images/food-iced-coffee.jpg'
,  './assets/images/food-iced-tea.jpg'
,  './assets/images/food-cola.jpg'
,  './assets/images/food-milo.jpg'
,  './assets/images/food-lemonade.jpg'

];

self.addEventListener('install', (event) => {
  event.waitUntil(
    caches.open(CACHE_NAME).then((cache) => cache.addAll(STATIC_ASSETS))
  );
  self.skipWaiting();
});

self.addEventListener('activate', (event) => {
  event.waitUntil(
    caches.keys().then((keys) =>
      Promise.all(
        keys.map((key) => {
          if (key !== CACHE_NAME) return caches.delete(key);
        })
      )
    )
  );
  self.clients.claim();
});

self.addEventListener('fetch', (event) => {
  if (event.request.method !== 'GET') return;
  event.respondWith(
    fetch(event.request).catch(() => caches.match(event.request, {ignoreSearch:true}))
  );
});
