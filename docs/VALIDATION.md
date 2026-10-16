# Validation — 16 October 2026

## Passed

- `go test -race -count=1 ./...` — all packages passed, including new invalid-input, cart-abuse, origin/security-header, static-routing, and concurrent movie lifecycle cases.
- `go vet ./...` — passed.
- `go mod tidy` and `go mod verify` — passed after removal of unused third-party dependencies.
- Native macOS/ARM64 executable rebuilt at `bin/cinebite`.
- Go WebAssembly engine and matching JavaScript runtime rebuilt; complete browser package generated at `dist/browser`.
- `python3 scripts/check-web.py` — packaged HTML file references, unique element IDs, and required runtime assets passed.
- Node VM syntax parsing — 12 JavaScript assets/inline blocks passed.
- Isolated browser-storage behavior — malformed cart, quantity clamp, and disabled-storage memory fallback passed.

A coverage run before unused-scraper removal measured 77.8% of server statements, 87.8% of recommendation statements, and 72.6% of scraper statements. Those percentages are not a security or functional-completeness certification. CLI entry point and models have no direct tests.

## Blocked / not completed

- Browser interaction and visual/mobile QA: local server startup failed because sandbox policy denies TCP listeners. No interactive browser pass is claimed.
- Executing compiled WASM in the isolated Node runtime: the embedder prohibits WebAssembly code generation. The WASM build passed; browser execution still needs verification outside this environment.
- `govulncheck`: downloading its module failed due to DNS/network restrictions. CI is configured to run it. No clean vulnerability scan is claimed.
- Docker/container execution and build: Docker is not installed here. CI is configured to build the image. No successful container run is claimed.
- Production load, real cinema/provider integration, HTTPS proxy, monitoring, and external deployment were not exercised.

The release remains a sample-data release candidate until the applicable checks in PRODUCTION_READINESS.md are completed.
