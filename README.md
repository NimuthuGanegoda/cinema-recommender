# CineBite — Cinema and Food Planner

CineBite helps people plan a cinema visit and choose snacks for their budget. It brings cinema locations, movie listings, food menus, and snack suggestions into one application.

**Go is the main application language.** Cinema data, menu processing, promotions, and recommendation calculations are implemented in Go. HTML and CSS provide the interface, while JavaScript connects the interface to the Go engine and handles browser interactions.

## Run locally for evaluation

### Requirements

- Go 1.24 or newer.
- A current browser such as Chrome, Firefox, Safari, or Edge.

Open a terminal in the project directory and run:

```sh
go run ./cmd/api -server -port 127.0.0.1:8080
```

Open **http://127.0.0.1:8080/** in the browser. The Go server serves both the website and its API, so no separate frontend installation is needed. Stop the application with **Ctrl+C** in the terminal.

To build an executable instead:

```sh
go build -o bin/cinebite ./cmd/api
```

On macOS or Linux, run:

```sh
./bin/cinebite -server -port 127.0.0.1:8080
```

On Windows, build and run:

```powershell
go build -o bin/cinebite.exe ./cmd/api
.\bin\cinebite.exe -server -port 127.0.0.1:8080
```

The included `run-local.command` is a convenience launcher for the supplied macOS executable. Building from source is recommended for evaluation on other operating systems.

### Run the tests

```sh
go test ./...
```

The tests cover cinema queries, food menus, movie listings, promotions, budget limits, API responses, and the removal of payment endpoints.

## Current features

- Browse cinema locations, screen details, directions, and food options.
- Search cinemas by name, chain, or city; select results by clicking or using the keyboard.
- Browse movie listings and suggested snack pairings.
- Explore popcorn, drinks, snacks, and sharing combos with matching food photographs.
- Generate a snack recommendation using a budget and party size.
- Compare example food promotions and eligible savings.
- Save a food plan on the current device, adjust quantities, remove items, and print it.
- Use the same responsive interface on phones, tablets, and desktops without a device-mode switch.
- Enjoy five dark colour palettes that rotate every 45 seconds with smooth transitions.

The cinema, movie, menu, and promotion records are **example data for the project**. Refreshing the listings or price display does not establish a live connection to a cinema. Food photographs illustrate the food types and are not official product photographs supplied by the cinemas. The application currently recommends and plans food; it does not submit food orders, reserve pickup slots, issue redeemable vouchers, or process payments.

## Browser version

The same Go logic can also run on a visitor's device through WebAssembly. This version needs only a current browser and ordinary static website files; visitors do not need to install Go or connect to a separate API server.

Build the browser assets using:

```sh
./scripts/build-browser.sh
```

The script produces `cinebite.wasm` and its matching `wasm_exec.js` runtime, and places them in `web/assets/wasm/`. Include those files together with `cinema-browser.js` and `cinema-worker.js` when distributing the browser version. The Go engine runs in a worker so calculations do not block the interface. Asset paths are relative, which allows the site to work inside a subdirectory.

For a local static preview, serve the `web/` directory using any static file server. For example, if Python is available:

```sh
python3 -m http.server 8080 --bind 127.0.0.1 --directory web
```

This server only delivers the files. Cinema data and recommendations execute in the browser. When the Go server is used instead, the application automatically uses its API.

## Suggested improvements and future features

### Official cinema partnerships and advance food ordering

A future version could invite official cinema operators to manage their own menus in CineBite. Visitors could choose their cinema, select food, and send an order to the cinema's concession team before arriving. This would reduce the time spent joining the normal food-ordering queue and give staff time to prepare orders ahead of the show.

The proposed flow would be:

1. The visitor chooses a cinema, film, and pickup time.
2. The visitor selects food from the cinema's official menu and submits the order.
3. Cinema staff accept the order and update its preparation status.
4. The visitor receives a ready-for-pickup notification and a collection reference.
5. The visitor collects the food from a designated pickup counter.

This would shorten queues rather than guarantee that every collection is immediate. Cinemas would need to agree to participate, maintain their listings, and provide staff and collection arrangements. These ordering features are proposed improvements and are not implemented in the current application.

Useful features for this extension include:

- **Cinema partner dashboard:** authorised staff can update prices, food availability, opening hours, and pickup capacity.
- **Order tracking:** show received, accepted, preparing, ready, and collected statuses.
- **Pickup slots:** limit orders per time slot so preparation and collection remain manageable.
- **Stock control:** hide unavailable items and prevent orders that cannot be fulfilled.
- **Order changes and cancellations:** provide clear options before preparation starts.
- **Pickup notifications:** tell visitors when food is ready and where to collect it.
- **Group orders:** allow friends or families to combine their selections in one order.

Real orders would require a shared backend and storage managed by the participating service. The browser-only version can handle recommendations and personal food plans, but cannot deliver orders to cinema staff on its own.

### Other improvements

- Add verified cinema and movie data with visible source and last-updated information.
- Improve recommendations using portion sizes, individual preferences, and dietary requirements.
- Show ingredient and allergen information supplied by cinema partners.
- Add Sinhala and Tamil interface translations alongside English.
- Save favourite cinemas and reusable food plans on the device.
- Improve accessibility with screen-reader testing, clearer labels, and larger text options.
- Add useful offline access for saved plans and recently viewed menus.
- Collect feedback on menus and pickup experiences to guide further improvements.

## Project structure

```text
cmd/                 Go entry points: API server and browser engine
internal/            Go models, server handlers, and tests
pkg/                 Go recommendations and cinema data processing
web/                 One shared website source
  assets/css/        Stylesheets
  assets/js/         Browser interactions and Go worker
  assets/images/     Cinema and food photographs
  assets/wasm/       Compiled Go browser engine and runtime
  data/              Example movie data
  *.html             Website pages
  embed.go           Embed website assets in the native Go server
docs/                API, architecture, contribution notes, photo credits
scripts/             Browser build and packaging scripts
bin/                 Local server executable (generated)
dist/browser/        Standalone browser package (generated)
```

Both runtimes use the same `web/` source. Run `./scripts/package-browser.sh` to prepare a browser-only distribution in `dist/browser/`, without Go source files. Rebuild the native server after editing website files because the executable embeds them.

## Photo credits

Food image sources, authors, and licences are listed in `docs/IMAGE_CREDITS.json` and the website's `web/photo-credits.html` page.
