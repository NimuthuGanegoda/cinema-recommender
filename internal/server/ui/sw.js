const CACHE_NAME = 'cinebite-v10-browser-go';
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
  './apple-design-system.css',
  './cinema-polish.css',
  './cinema-environment.js',
  './photo-credits.html',
  './cinema-client.js',
  './cinema-browser.js',
  './cinema-worker.js',
  './wasm_exec.js',
  './cinebite.wasm',
  './movies.json',
  './manifest.json',
  './app-icon.jpg',
  './popcorn.jpg',
  './beverage.jpg',
  './combo.jpg',
  './hotdog.jpg',
  './nachos.jpg',
  './rolls.jpg',
  './food-caramel-popcorn.jpg'
,  './food-butter-popcorn.jpg'
,  './food-chocolate-popcorn.jpg'
,  './food-cheese-sticks.jpg'
,  './food-chicken-strips.jpg'
,  './food-fries.jpg'
,  './food-samosas.jpg'
,  './food-iced-coffee.jpg'
,  './food-iced-tea.jpg'
,  './food-cola.jpg'
,  './food-milo.jpg'
,  './food-lemonade.jpg'

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
