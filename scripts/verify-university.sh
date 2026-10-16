#!/bin/sh
set -eu
cd "$(dirname "$0")/.."
go mod verify
go vet ./...
go test -race -count=1 -cover ./...
./scripts/package-browser.sh
python3 scripts/check-web.py
