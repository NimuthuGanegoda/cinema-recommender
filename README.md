# CineBite — University Demonstration Edition

A responsive cinema guide, food menu browser, snack budget recommender, and device-local food planner. Go provides the shared recommendation and catalog engine, with HTML/CSS/JavaScript for the interface. The same engine runs on the native HTTP server and in browser WebAssembly.

**Release candidate with sample data.** Prices, offers and showtimes are examples; confirm real details with cinemas. This package does not reserve seats, place food orders, process payments, or guarantee dietary safety.

## Start here for judging

Read [the demonstration and defence guide](docs/JUDGES_GUIDE.md), [the current architecture](docs/ARCHITECTURE.md), and [university validation evidence](docs/UNIVERSITY_VALIDATION.md). This edition fixes promotion-code membership bypass, checks top-ups against actual final cost, and explains each budget suggestion. The algorithm is a category heuristic; it does not guarantee a globally optimal bundle.

## Run

On an Apple Silicon Mac, run `Start_Demo.command` (or `sh Start_Demo.command` in Terminal), then open http://127.0.0.1:8080. The included executable is macOS ARM64 only; other platforms can use Go or the browser package.

Use a supported patched Go version (module minimum 1.24):

```sh
go run ./cmd/api -server -port 127.0.0.1:8080
```

Open http://127.0.0.1:8080. To build a native executable:

```sh
go build -trimpath -o bin/cinebite ./cmd/api
./bin/cinebite -server -port 127.0.0.1:8080
```

`PORT` sets the listen address when `-port` is omitted. Bare numeric ports are supported. Exact permitted cross-origin browser origins can be set with `ALLOWED_ORIGINS`, comma separated. No wildcard cross-origin access is enabled.

## Browser-only build

```sh
./scripts/package-browser.sh
python3 scripts/check-web.py
python3 -m http.server 8080 --bind 127.0.0.1 --directory dist/browser
```

Serve `dist/browser/` through HTTPS for production; localhost is allowed for development. Visitors do not install Go. Deploy the complete matching WASM/runtime/UI package together. Offline support depends on prior successful caching and browser capabilities.

## Verify

```sh
go mod verify
go vet ./...
go test -race -count=1 ./...
go run golang.org/x/vuln/cmd/govulncheck@latest ./...
./scripts/package-browser.sh
python3 scripts/check-web.py
```

## Container

```sh
docker compose up --build -d
```

The service is available on host loopback at port 8080. Read [the deployment and remaining launch requirements](docs/PRODUCTION_READINESS.md) before exposing it publicly. [Validation results](docs/VALIDATION.md) distinguish completed checks from blocked checks.

## Structure

- `cmd/api`: HTTP server and CLI.
- `cmd/browser`: browser WebAssembly entry point.
- `internal`: models, HTTP handling, validation, security middleware, tests.
- `pkg`: shared catalog and recommendations.
- `web`: single UI source, assets and example data.
- `scripts`: reproducible browser packaging and asset validation.
- `docs`: release notes, API/architecture references, photo credits.
- `.github/workflows`: quality checks and browser distribution workflow.

Photo credits are in `docs/IMAGE_CREDITS.json` and `web/photo-credits.html`. Older API and architecture notes describe the project's history; current runtime security and limits are documented in `docs/PRODUCTION_READINESS.md`.
