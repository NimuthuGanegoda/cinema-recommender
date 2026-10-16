/* Storage may be disabled or contain a malformed plan; keep the UI usable. */
(function () {
 const memory = new Map();
 let persistent = true;
 function normalize(key,value) {
  if (key !== 'cinebite_cart' || value === null) return value;
  try {
   const items=JSON.parse(value);
   if (!Array.isArray(items)) return '[]';
   return JSON.stringify(items.filter(item=>item&&typeof item==='object'&&typeof item.name==='string').slice(0,100).map(item=>({...item,quantity:Math.min(99,Math.max(1,Math.floor(Number(item.quantity)||1)))})));
  } catch { return '[]'; }
 }
 window.CineBiteStorage = {
  get persistent(){return persistent;},
  getItem(key){try{return normalize(key,localStorage.getItem(key));}catch{persistent=false;return normalize(key,memory.get(key)??null);}},
  setItem(key,value){value=normalize(key,String(value));memory.set(key,value);try{localStorage.setItem(key,value);}catch{persistent=false;}},
  removeItem(key){memory.delete(key);try{localStorage.removeItem(key);}catch{persistent=false;}}
 };
})();
/* Automatic appearance and device adaptation, kept entirely in this browser. */
(function () {
  'use strict';
  const root = document.documentElement;
  const paletteNames = ['gruvbox','everforest','tokyo-night','rose-pine','catppuccin'];
  function paletteForElapsed(elapsedMs) { return paletteNames[Math.floor(Math.max(0, elapsedMs) / 45000) % paletteNames.length]; }
  let startedAt = Date.now();
  try {
    const stored = Number(sessionStorage.getItem('cinebite_palette_started_at'));
    if (stored > 0 && stored <= startedAt) startedAt = stored;
    else sessionStorage.setItem('cinebite_palette_started_at', String(startedAt));
  } catch (error) { /* Private browsing can restrict storage; this page still rotates. */ }
  function browserForUA(ua) {
    if (/Edg\//.test(ua)) return 'edge';
    if (/Firefox\//.test(ua)) return 'firefox';
    if (/Chrome\/|CriOS\//.test(ua)) return 'chrome';
    if (/Safari\//.test(ua)) return 'safari';
    return 'other';
  }
  function deviceFor(ua, touchPoints, width, coarsePointer) {
    const tablet = /iPad|Tablet/.test(ua) || (/Macintosh/.test(ua) && touchPoints > 1) || (/Android/.test(ua) && !/Mobile/.test(ua));
    if (tablet) return 'tablet';
    if (/Mobi|iPhone|iPod/.test(ua)) return 'phone';
    if (coarsePointer && width <= 600) return 'phone';
    if (coarsePointer && touchPoints > 1) return 'tablet';
    return 'desktop';
  }
  let transitionTimer;
  function refreshAppearance() {
    const palette = paletteForElapsed(window.matchMedia('(prefers-reduced-motion: reduce)').matches ? 0 : Date.now() - startedAt);
    if (root.dataset.colourTheme === palette) return;
    if (root.dataset.colourTheme && root.dataset.colourTheme !== palette) {
      root.classList.add('palette-transition');
      clearTimeout(transitionTimer);
      transitionTimer = setTimeout(() => root.classList.remove('palette-transition'), 2100);
    }
    root.dataset.colourTheme = palette;
    root.style.colorScheme = 'dark';
    const colours = {'gruvbox':'#282828','everforest':'#2d353b','tokyo-night':'#1a1b26','rose-pine':'#191724','catppuccin':'#1e1e2e'};
    const meta = document.querySelector('meta[name="theme-color"]');
    if (meta) meta.content = colours[palette];
  }
  function refreshDevice() {
    const ua = navigator.userAgent || '';
    root.dataset.browser = browserForUA(ua);
    root.dataset.device = deviceFor(ua, navigator.maxTouchPoints || 0, window.innerWidth, window.matchMedia('(pointer: coarse)').matches);
    root.dataset.touch = (navigator.maxTouchPoints || 0) > 0 ? 'true' : 'false';
  }
  refreshAppearance();
  refreshDevice();
  window.addEventListener('resize', refreshDevice, {passive:true});
  window.addEventListener('pageshow', refreshAppearance);
  document.addEventListener('visibilitychange', () => {if (!document.hidden) {refreshAppearance();refreshDevice();}});
  setInterval(refreshAppearance, 1000);
  // Pure helpers allow palette rotation and device cases to be checked without changing the clock.
  window.CineBiteEnvironment = Object.freeze({paletteForElapsed, browserForUA, deviceFor});
})();
