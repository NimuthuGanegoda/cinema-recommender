#!/bin/sh
set -eu
cd "$(dirname "$0")/.."
GOOS=js GOARCH=wasm go build -trimpath -ldflags="-s -w" -o internal/server/ui/cinebite.wasm ./cmd/browser
runtime="$(go env GOROOT)/lib/wasm/wasm_exec.js"
if [ ! -f "$runtime" ]; then runtime="$(go env GOROOT)/misc/wasm/wasm_exec.js"; fi
cp "$runtime" internal/server/ui/wasm_exec.js
for destination in . docs; do
  cp internal/server/ui/cinebite.wasm "$destination/cinebite.wasm"
  cp internal/server/ui/wasm_exec.js "$destination/wasm_exec.js"
done
