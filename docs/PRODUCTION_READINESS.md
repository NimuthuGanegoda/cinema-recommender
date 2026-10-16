# CineBite release candidate

This package is hardened for a sample-data cinema and food planner. It is not a verified live cinema service. No deployment was performed.

## Implemented

- Shared Go API, native embedded UI, and WebAssembly engine.
- Removed an unused arbitrary-URL scraping implementation and its third-party dependency tree. Runtime now uses the Go standard library only; the legacy method explicitly reports an unconfigured provider.
- 64 KiB request body limit, 32 KiB HTTP header limit, header/read/write/idle timeouts.
- Strict single-object JSON, unsupported-field rejection, and numeric limits.
- Cart quantities of 1–99, at most 100 lines and 500 total units; unknown and unavailable items are rejected.
- Recommendation budgets up to LKR 1,000,000 and party sizes up to 50. Zero retains the documented default; negative inputs are rejected.
- Unsupported dietary/category filters return an error instead of suggesting they were honored.
- Same-origin access by default and configurable exact cross-origin allowlist.
- 32 concurrent API requests per process, with 503/Retry-After overload responses. This bounds concurrency; it is not a per-client rate limiter.
- Security headers and restricted resource policy; directory listings and non-public files are unavailable.
- Stable example prices; no fabricated eight-second price changes. Example-data labels in pages, API headers, feed metadata, and health responses.
- Removed hourly commits that presented timestamp changes as fetched movie updates.
- Fixed movie ticker shutdown race and prevented stopped engines from restarting.
- Movie filters combine cinema and language; unknown cinema IDs return 404.
- Network/worker request timeouts, bounded saved plans, clear-plan action, keyboard skip links, visible focus, and reduced-motion support.
- Versioned offline cache scoped to the application, no API caching, and no deletion of unrelated caches. Workers wait for existing tabs to close before activating.
- Multi-stage container with non-root runtime, read-only Compose filesystem, dropped capabilities, resource limits, rotated logs, loopback port binding, and health probe.
- CI gates: formatting, dependency verification, vet, race tests, coverage, native build, govulncheck, browser build, asset checks, and container build.

## Launch and operations

Run `go run ./cmd/api -server -port 127.0.0.1:8080` and visit http://127.0.0.1:8080.

For Docker, run `docker compose up --build -d`. The app binds to host loopback only. Put an HTTPS reverse proxy in front of it. Configure `ALLOWED_ORIGINS` with the exact public HTTPS origin, e.g. `ALLOWED_ORIGINS=https://cinema.example.com docker compose up --build -d`. This is required for browser POST requests routed through TLS termination because the backend sees HTTP. Do not trust arbitrary forwarded-origin headers.

Use a supported, patched Go release. The container tracks the Go major tag and Alpine stable tag; for controlled releases pin verified image digests after the container build and security scan pass. The supplied environment has Go 1.27.2; go.mod retains its 1.24 minimum for source compatibility.

At the proxy enforce per-client rate limits, TLS, a request-size limit, and connection limits. Configure HSTS at TLS termination after confirming HTTPS works. Monitor `/health` for availability, 5xx responses, latency, memory, restarts, and logs. Health indicates process availability, not data freshness. SIGTERM drains requests for up to five seconds. Before upgrading, retain the prior image and asset bundle; deploy matched HTML/JS/WASM versions together. Roll back by restoring the previous image/bundle.

The app has no server database or accounts. Plans stay in browser localStorage and are not synchronized or backed up to the server. The clear-plan button removes the saved plan contents. Device loss or browser storage deletion loses plans. Request logs include path, remote address, status, and latency; configure retention and access restrictions for your host.

## Required before claiming a verified live service

1. Obtain licensed/authorized sources for cinema details, menus, promotions, movie metadata and showtimes. Replace example generators with validated provider adapters; include source, last successful fetch, and stale/error states. Do not relabel generated fixtures as live data.
2. Validate image licences and attribution, cinema-brand references, accessibility, and privacy/legal copy for the actual deployment.
3. Run the CI vulnerability scan and container build/scan, then resolve findings. They were not executable locally: network DNS is restricted and Docker is not installed.
4. Perform desktop/mobile browser flows in native and browser-only modes, including first-load/offline/update behavior, keyboard and screen-reader use. Local TCP listeners are denied in this environment, so browser runtime QA could not be completed here.
5. Run load tests against the actual host and proxy, set capacity limits, and confirm monitoring and rollback.
6. Existing pages still contain inline scripts and event handlers; CSP permits unsafe-inline for compatibility. Externalize handlers and adopt hashes/nonces before treating CSP as strong script-injection protection. Current checks are not a complete penetration test.

Payments, reservations, partner order delivery, staff/admin accounts, and allergen guarantees are not implemented. Adding these requires authentication/authorization, durable storage, audit trails, and appropriate external integrations.

## Validation in this environment

See `VALIDATION.md` for the exact completed checks and limits. Do not infer that a CI definition is a completed CI run.
