#!/bin/sh
set -eu
cd "$(dirname "$0")"
if [ "$(uname -s)" = Darwin ] && [ "$(uname -m)" = arm64 ] && [ -x bin/cinebite ]; then
  echo 'Open http://127.0.0.1:8080 — press Control-C to stop.'
  exec ./bin/cinebite -server -port 127.0.0.1:8080
fi
if command -v go >/dev/null 2>&1; then
  echo 'Open http://127.0.0.1:8080 — press Control-C to stop.'
  exec go run ./cmd/api -server -port 127.0.0.1:8080
fi
echo 'Install a compatible Go toolchain, or serve dist/browser with a local HTTP server.'
exit 1
