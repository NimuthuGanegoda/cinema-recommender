#!/bin/sh
cd "$(dirname "$0")" || exit 1
if command -v go >/dev/null 2>&1; then
  exec go run ./cmd/api -server -port 127.0.0.1:8080
fi
exec ./bin/cinebite -server -port 127.0.0.1:8080
