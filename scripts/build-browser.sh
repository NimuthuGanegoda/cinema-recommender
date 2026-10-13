#!/bin/sh
set -eu
cd "$(dirname "$0")/.."
mkdir -p web/assets/wasm
GOOS=js GOARCH=wasm go build -trimpath -ldflags="-s -w" -o web/assets/wasm/cinebite.wasm ./cmd/browser
runtime="$(go env GOROOT)/lib/wasm/wasm_exec.js"
if [ ! -f "$runtime" ]; then runtime="$(go env GOROOT)/misc/wasm/wasm_exec.js"; fi
cp "$runtime" web/assets/wasm/wasm_exec.js
