# University edition validation — 16 October 2026

## Executed for this edition

| Check | Result |
| --- | --- |
| Go toolchain | go1.27.2, macOS ARM64 |
| `go test -race -count=1 -cover ./...` | Passed |
| Server statement coverage | 77.8% |
| Recommender statement coverage | 89.5% |
| Catalog/scraper statement coverage | 84.1% |
| `go vet ./...` | Passed |
| Browser package rebuild | Passed; Go WASM and matching JS runtime regenerated |
| `python3 scripts/check-web.py` | Passed; packaged references, unique IDs and required assets |
| Promotion membership regressions | Passed; CINE, SCOPE and REDOPAY codes cannot bypass tier checks |
| Promotion top-up regressions | Passed; beneficial top-up accepted, cost-increasing top-up rejected |
| Budget sweep | Passed; 271 budgets × 5 party sizes × 3 tiers = 4,065 cases on a fixed four-item menu |
| Availability and explanation checks | Passed within the budget sweep |

The budget sweep covers budgets 1 through 9,991 in increments of 37, party sizes 1/2/4/6/50 and Standard/Student/Platinum tiers. It checks nonnegative final cost, budget compliance, in-stock selections, algorithm identifier and explanation. It is not an exhaustive proof for every catalog or input.

## Microbenchmark

`go test ./pkg/recommender -run '^$' -bench BenchmarkRecommendation -benchmem -count=1`

Apple M1, darwin/arm64: 254,223 iterations; 4,495 ns/op; 9,500 B/op; 35 allocations/op. This benchmark uses a two-item menu and a two-person LKR 2,500 request. It measures only the Go recommendation function. It excludes HTTP, UI, WASM startup and realistic catalog scale. One run does not establish comparative performance.

Native macOS ARM64 executable rebuilt; CLI smoke run for Kandy, LKR 2,500, two people passed (example final cost LKR 1,750). `go mod verify` passed.

Updated planner inline JavaScript passed Node VM syntax parsing. This is a syntax check, not an interaction test.

## Remaining validation

A native server launch was attempted, but the execution environment rejected TCP binding (`bind: operation not permitted`). Browser interaction, mobile visual QA and full user workflows could not be verified here. Use `Start_Demo.command` on the presentation laptop and rehearse the guide before judging.

The earlier release also left vulnerability scanning, Docker execution, production load testing and external integrations unverified. This edition does not claim those checks passed. CLI entry points and model helpers still have no direct statement coverage. No user study, recommendation preference accuracy or globally optimal bundle result is claimed.

## Reproduce

Run `./scripts/verify-university.sh` with a compatible Go toolchain and Python 3. This runs module verification, vet, race/coverage tests, WASM packaging and asset validation. Results may differ by machine and toolchain. Keep logs from the actual demonstration laptop alongside this file.
