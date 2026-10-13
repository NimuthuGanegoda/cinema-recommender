/* Runs the compiled Go API inside a browser worker; no API server is contacted. */
let failed = false;
self.cinebiteReady = () => self.postMessage({type:'ready'});
self.onmessage = ({data}) => {
  if (failed || typeof self.cinebiteInvoke !== 'function') return;
  try {
    const result = JSON.parse(self.cinebiteInvoke(data.method,data.path,JSON.stringify(data.body ?? {})));
    self.postMessage({id:data.id,...result});
  } catch(error) { self.postMessage({id:data.id,status:500,body:{error:'Could not calculate this food plan. Please try again.'}}); }
};
(async () => {
  try {
    importScripts(new URL('../../assets/wasm/wasm_exec.js',self.location.href).href);
    const go = new Go();
    const response = await fetch(new URL('../../assets/wasm/cinebite.wasm',self.location.href));
    if (!response.ok) throw new Error('The browser engine files are missing from this website.');
    let result;
    try { result = await WebAssembly.instantiateStreaming(response.clone(),go.importObject); }
    catch(error) { result = await WebAssembly.instantiate(await response.arrayBuffer(),go.importObject); }
    go.run(result.instance).catch(error => {failed=true;self.postMessage({type:'error',message:'The browser engine stopped. Reload this page to try again.'});});
  } catch(error) { failed=true;self.postMessage({type:'error',message:error.message || 'The browser engine could not load. Please reload this page.'}); }
})();
