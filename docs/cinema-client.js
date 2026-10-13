function foodImageFor(item) {
  const name=(item.name||'').toLowerCase();
  if(name.includes('combo')||name.includes('feast')||name.includes('meal box')||name.includes('saver')) return 'combo.jpg';
  if(name.includes('popcorn')) {if(name.includes('chocolate'))return 'food-chocolate-popcorn.jpg';if(name.includes('caramel'))return 'food-caramel-popcorn.jpg';return 'food-butter-popcorn.jpg';}
  if(name.includes('cheese sticks')||name.includes('mozzarella'))return 'food-cheese-sticks.jpg';
  if(name.includes('chicken strips')||name.includes('tender'))return 'food-chicken-strips.jpg';
  if(name.includes('fries'))return 'food-fries.jpg';
  if(name.includes('samosa'))return 'food-samosas.jpg';
  if(name.includes('milo'))return 'food-milo.jpg';
  if(name.includes('coffee'))return 'food-iced-coffee.jpg';
  if(name.includes('tea'))return 'food-iced-tea.jpg';
  if(name.includes('coca-cola')||name.includes('cola')||name.includes('coke'))return 'food-cola.jpg';
  if(name.includes('sparkler')||name.includes('lemon')||name.includes('juice'))return 'food-lemonade.jpg';
  if(name.includes('hotdog')||name.includes('hot dog')||name.includes('chili dog'))return 'hotdog.jpg';
  if(name.includes('nachos'))return 'nachos.jpg';
  if(name.includes('roll')||name.includes('patty'))return 'rolls.jpg';
  return item.image || ({Popcorn:'food-butter-popcorn.jpg',Beverage:'beverage.jpg',Combo:'combo.jpg',Snack:'nachos.jpg'}[item.category]) || 'popcorn.jpg';
}


const SRI_LANKAN_CINEMAS = []; // No duplicated browser registry: Go is the source of cinema data.
function calculateHaversineDistance(lat1,lon1,lat2,lon2) {
  const r=Math.PI/180; const a=Math.sin((lat2-lat1)*r/2)**2+Math.cos(lat1*r)*Math.cos(lat2*r)*Math.sin((lon2-lon1)*r/2)**2;
  return Math.round(6371*2*Math.atan2(Math.sqrt(a),Math.sqrt(1-a))*10)/10;
}
class CinemaClient {
 constructor(baseUrl=''){this.baseUrl=baseUrl;this.apiRequired=true;this.nextMovieSync=Date.now()+3600000;}
 async request(path,body){
  const method=body===undefined?'GET':'POST';
  if(document.querySelector('meta[name="cinebite-runtime"]')?.content==='server'){
    const response=await fetch(this.baseUrl+path,{method,headers:body===undefined?{}:{'Content-Type':'application/json'},body:body===undefined?undefined:JSON.stringify(body),cache:'no-store'});
    const result=await response.json().catch(()=>null);
    if(!response.ok)throw new Error(result?.error||'Could not complete this request.');
    return result;
  }
  const response=await window.CineBiteBrowser.request(method,path,body);
  if(response.status<200||response.status>=300)throw new Error(response.body?.error||'Could not complete this request.');
  return response.body;
 }
 async getCinemas({lat=0,lng=0,city=''}={}){const q=new URLSearchParams();if(city)q.set('city',city);if(lat||lng){q.set('lat',lat);q.set('lng',lng);}return this.request('/api/v1/cinemas'+(q.size?'?'+q:''));}
 async getCinema(id){return this.request('/api/v1/cinemas/'+encodeURIComponent(id));}
 normalizeItem(item){return {...item,price:item.price_lkr??item.price??0,original_price:item.base_price_lkr??item.original_price??item.price_lkr??0,image:foodImageFor(item)};}
 async getLiveConcessions(id){const data=await this.request('/api/v1/cinemas/'+encodeURIComponent(id)+'/concessions/live');return {...data,items:(data.items||[]).map(item=>this.normalizeItem(item))};}
 async getCinemaAds(id){return this.request('/api/v1/cinemas/'+encodeURIComponent(id)+'/ads');}
 async getPromotions(){return this.request('/api/v1/promotions');}
 normalizeRecommendation(data){return {...data,items:(data.selected_items||[]).map(item=>this.normalizeItem(item)),original_total:data.original_total_lkr||0,discount_amount:data.discount_amount_lkr||0,final_total:data.final_total_lkr||0,savings_lkr:data.discount_amount_lkr||0,upsell_advice:data.next_tier_upsell||''};}
 async recommendBundle({cinemaId,budget,budgetLKR,budget_lkr,partySize=2,privilegeTier='Standard',promoCode=''}){return this.normalizeRecommendation(await this.request('/api/v1/recommend',{cinema_id:cinemaId,budget_lkr:budget??budgetLKR??budget_lkr??2500,party_size:partySize,privilege_tier:privilegeTier,promo_code:promoCode}));}
 async evaluateCart({cinemaId,items,privilegeTier='Standard',promoCode=''}){return this.normalizeRecommendation(await this.request('/api/v1/cart/evaluate',{cinema_id:cinemaId,items:items.map(item=>({item_id:item.item_id||item.id,quantity:item.quantity||1})),privilege_tier:privilegeTier,promo_code:promoCode}));}
 async getMovies({cinemaId='',language=''}={}){const q=new URLSearchParams();if(cinemaId)q.set('cinema_id',cinemaId);if(language)q.set('language',language);const data=await this.request('/api/v1/movies'+(q.size?'?'+q:''));if(data.sync_status?.next_sync_at)this.nextMovieSync=new Date(data.sync_status.next_sync_at).getTime();return data.movies||[];}
 async getMovieSyncStatus(){return this.request('/api/v1/movies/status');}
 async triggerMovieSync(){const data=await this.request('/api/v1/movies/sync',{});if(data.sync_status?.next_sync_at)this.nextMovieSync=new Date(data.sync_status.next_sync_at).getTime();return data;}
 startHourlyMovieSyncEngine(onSync,onTick){let refreshing=false;const tick=async()=>{const remaining=Math.max(0,Math.ceil((this.nextMovieSync-Date.now())/1000));if(onTick)onTick(remaining,Math.floor(remaining/60)+'m '+String(remaining%60).padStart(2,'0')+'s');if(!remaining&&!refreshing){refreshing=true;try{const movies=await this.getMovies();if(onSync)onSync(movies);}catch(error){console.error(error);}finally{refreshing=false;}}};tick();return setInterval(tick,1000);}
}
if(typeof window!=='undefined'){window.CinemaClient=CinemaClient;window.cinemaAPI=new CinemaClient();}
if(typeof module!=='undefined')module.exports={CinemaClient,calculateHaversineDistance,SRI_LANKAN_CINEMAS};
