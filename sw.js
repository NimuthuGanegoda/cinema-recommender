const CACHE_NAME = 'cinemasnack-v3';
const STATIC_ASSETS = [
  './',
  './index.html',
  './movies.html',
  './cinemas.html',
  './menu.html',
  './optimizer.html',
  './deals.html',
  './checkout.html',
  './mobile.html',
  './apple-design-system.css',
  './cinema-client.js',
  './movies.json',
  './manifest.json',
  './app-icon.jpg'
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
    fetch(event.request).catch(() => caches.match(event.request))
  );
});
