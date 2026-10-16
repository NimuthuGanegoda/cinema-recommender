/* Start Go in the visitor's browser, using relative asset paths. */
(function(){
  const assetBase = new URL('.',document.currentScript.src);
  let worker, ready, serial=0;
  const pending = new Map();
  function start(){
    if(ready) return ready;
    ready=new Promise((resolve,reject)=>{
      if(!window.Worker || !window.WebAssembly){reject(new Error('Please use a current browser to open CineBite.'));return;}
      const status=document.createElement('div');status.className='engine-loading';status.setAttribute('role','status');status.textContent='Getting your cinema guide ready…';
      const show=()=>{document.body.append(status);};if(document.body)show();else document.addEventListener('DOMContentLoaded',show,{once:true});
      worker=new Worker(new URL('cinema-worker.js',assetBase));
      const timer=setTimeout(()=>{status.textContent='The cinema guide is taking longer to load. Please reload to try again.';reject(new Error('The cinema guide could not load. Please reload this page.'));worker.terminate();},90000);
      const fail=message=>{clearTimeout(timer);status.textContent=message;reject(new Error(message));pending.forEach(task=>{clearTimeout(task.timer);task.reject(new Error(message));});pending.clear();};
      worker.onerror=()=>fail('The browser engine could not load. Please reload this page.');
      worker.onmessage=({data})=>{
        if(data.type==='ready'){clearTimeout(timer);status.remove();resolve();return;}
        if(data.type==='error'){fail(data.message);return;}
        const task=pending.get(data.id);if(!task)return;pending.delete(data.id);clearTimeout(task.timer);task.resolve(data);
      };
    });return ready;
  }
  window.CineBiteBrowser={async request(method,path,body){await start();return new Promise((resolve,reject)=>{const id=++serial;const timer=setTimeout(()=>{pending.delete(id);reject(new Error('The request timed out. Please reload and try again.'));},15000);pending.set(id,{resolve,reject,timer});worker.postMessage({id,method,path,body});});}};
})();
