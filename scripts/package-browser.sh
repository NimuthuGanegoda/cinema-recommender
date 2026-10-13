#!/bin/sh
set -eu
cd "$(dirname "$0")/.."
./scripts/build-browser.sh
mkdir -p dist/browser
cp web/*.html web/manifest.json web/sw.js dist/browser/
cp -R web/assets web/data dist/browser/
