const CACHE_NAME = 'cinebite-v13-university';
const SHELL = ['./','./index.html','./movies.html','./cinemas.html','./menu.html','./optimizer.html','./deals.html','./food-plan.html','./photo-credits.html','./manifest.json','./assets/css/apple-design-system.css','./assets/css/cinema-polish.css','./assets/js/cinema-client.js','./assets/js/cinema-browser.js','./assets/js/cinema-worker.js','./assets/js/cinema-environment.js','./assets/wasm/wasm_exec.js','./assets/wasm/cinebite.wasm'];
self.addEventListener('install', event => {
 event.waitUntil(caches.open(CACHE_NAME).then(cache => cache.addAll(SHELL)));
 // New workers wait for old tabs to close, preventing mixed release assets.
});
self.addEventListener('activate', event => {
 event.waitUntil((async () => {
  await Promise.all((await caches.keys()).filter(key => key.startsWith('cinebite-') && key !== CACHE_NAME).map(key => caches.delete(key)));
  await self.clients.claim();
 })());
});
self.addEventListener('fetch', event => {
 const request = event.request, url = new URL(request.url), scope = new URL(self.registration.scope);
 if (request.method !== 'GET' || url.origin !== scope.origin || !url.pathname.startsWith(scope.pathname) || url.pathname.includes('/api/')) return;
 event.respondWith((async () => {
  const cache = await caches.open(CACHE_NAME);
  try {
   const response = await fetch(request);
   if (response.ok && response.type === 'basic') {
    const cacheable = request.mode === 'navigate' || SHELL.some(path => new URL(path,scope).href === url.href) || url.pathname.includes('/assets/');
    if (cacheable) await cache.put(request,response.clone());
   }
   return response;
  } catch {
   const saved = await cache.match(request);
   if (saved) return saved;
   return new Response('This page is unavailable offline. Reconnect and try again.',{status:503,headers:{'Content-Type':'text/plain; charset=utf-8'}});
  }
 })());
});
