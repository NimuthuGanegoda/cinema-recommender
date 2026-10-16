# Explaining CineBite

## A short introduction you can practise

“CineBite helps users plan cinema snacks within a budget. Users choose a cinema, enter their budget and group size, and receive a suggested food combination. Go handles the cinema data, recommendation rules, and example discounts. HTML and CSS display the website, and JavaScript connects the buttons and forms to Go. Users can save and print their food plan.”

## Learn these five parts first

| Part | Its job | Where to look |
| --- | --- | --- |
| Pages | Show cinema information, menus, and forms | `web/*.html` |
| Styles | Control layout, mobile sizing, and colours | `web/assets/css/` |
| Browser interactions | Read inputs, request results, update the page | `web/assets/js/cinema-client.js` |
| Go handlers | Receive requests and call the correct Go function | `internal/server/` |
| Recommendation rules | Choose food and calculate estimated savings | `pkg/recommender/` |

A **model** is a Go structure that groups related fields. For example, a food item has a name, category, price, and stock status. Models are in `internal/models/`.

## Follow one click: “Find my snack mix”

1. The user selects a cinema, budget, group size, and optional example discount tier or code.
2. The page collects those values. `CinemaClient.recommendBundle` prepares a request.
3. The request reaches Go's `handleRecommend` function in `internal/server/handlers_recommend.go`.
4. The handler finds the selected cinema and obtains its example food menu.
5. `RecommendOptimalBundle` in `pkg/recommender/optimizer.go` selects items.
6. `EvaluateCart` in `pkg/recommender/cart.go` adds their prices and applies an eligible example promotion.
7. Go returns the chosen items, subtotal, discount, final estimate, and explanation.
8. The page displays the result. Adding the snacks to a food plan saves them in this browser.

The values travel as **JSON**, a text format containing named fields. A simplified input looks like this:

```json
{"cinema_id":"CMB-CCC","budget_lkr":2500,"party_size":2}
```

An **API** is the interface the page uses to ask Go for data or calculations. For example, `/api/v1/recommend` means “calculate a recommendation.”

## How the recommendation works

This is a rule-based heuristic: a practical set of selection rules. It is not trained AI, and it does not test every possible combination or guarantee the mathematically best result.

The main rules are:

1. Separate available food into popcorn, drinks, snacks, and combos; exclude out-of-stock items from these category lists.
2. Sort each category by price.
3. Try a combo costing at most 85% of the budget.
4. If no combo is chosen, start with popcorn, then try drinks for the group (up to four), then snacks when the selection rules allow them.
5. Try adding an item near a promotion spending threshold.
6. Calculate eligible savings.
7. If the discounted total exceeds the budget, remove the most expensive selected item and recalculate until the result fits.
8. If nothing remains, look for the cheapest in-stock item that fits. If no item fits, return an empty selection.

The temporary selection can go above the budget before discounts are calculated. The final budget check is what keeps the returned estimate within the requested budget.

```text
Choose candidate snacks
Calculate subtotal and eligible discount
While final cost is over budget:
    Remove the most expensive selected item
    Calculate the cost again
Return the suggestion
```

The promotion-threshold helper currently scans the original menu, rather than the filtered category lists. A useful future improvement is to apply the same stock check there as well. Group size influences selection but does not guarantee that every person gets a drink or that a combo has enough portions.

## Explain the money calculation

```text
Subtotal = sum of selected item prices
Discount = largest eligible example discount
Final estimate = subtotal - discount
```

For a simple arithmetic example, if the subtotal is LKR 2,400 and an eligible 30% discount applies without a lower cap:

```text
Discount = 2,400 × 0.30 = LKR 720
Final estimate = 2,400 - 720 = LKR 1,680
```

This illustrates the calculation; it is not a promise of a real cinema offer. Actual eligibility also depends on the code, tier, minimum spend, required food categories, and discount cap implemented by the promotion rules.

The printed food plan adds saved item prices multiplied by quantity. It does not currently recalculate the recommendation's promotion, so its total can differ from the discounted recommendation estimate.

## Why there are two ways to run Go

### Local server — the easiest version to demonstrate

```text
Browser → local Go server → recommendation rules → browser
```

`cmd/api/main.go` starts the Go application. The server listens on the local machine and serves both the website and API requests. `web/embed.go` includes the website files inside the compiled server executable.

```sh
go run ./cmd/api -server -port 127.0.0.1:8080
```

Open `http://127.0.0.1:8080/`.

### Browser version

```text
Page → browser worker → compiled Go engine → page
```

**WebAssembly** lets compiled Go run inside a browser. `cmd/browser/main.go` exposes the same Go handlers to JavaScript without starting a listening server. `cinema-worker.js` loads that engine in a worker: a separate browser execution context that keeps calculations away from the page's main interaction thread.

`cinema-browser.js` sends messages to the worker and matches each reply to its request. `cinema-client.js` automatically chooses the local API when the page has the server marker; otherwise it uses the browser engine.

The website still needs a way to deliver its static files. Visitors do not need their own Go installation or a separate application backend. The browser engine and the local server share the same Go recommendation code.

## Data and storage

- Cinema and movie records, menus, prices, and promotions are examples supplied by the project. They are not verified live cinema feeds.
- Some example prices change using time-based rules. Refreshing them does not contact a cinema's till or stock system.
- The Go menu cache keeps recent menu results in memory for ten minutes. Restarting the Go engine clears that memory.
- The food plan uses `localStorage`, the browser's storage for this website. It survives a page reload, but it does not sync between devices or different website addresses. Clearing browser data can remove it.
- There is no shared order database, account system, payment processing, or food-order delivery to cinema staff.
- The package named `scraper` contains data-loading and parsing code. Its name alone does not mean the demonstration is scraping live official websites.

## Appearance, search, and printing

- Cinema search matches loaded cinema names, cities, and chains; it does not search the internet.
- Responsive CSS changes the layout as screen width changes. Browser and device checks provide additional hints; no Desktop/App switch is required.
- A browser timer chooses one of five dark palettes every 45 seconds. CSS transitions fade between colours. This is elapsed browsing time, not the time of day.
- Printing uses the browser's `window.print()` function. Print-specific CSS hides navigation and buttons, removes screen decorations and overflow, and shows readable quantities and totals on paper. It also supports the browser's Save to PDF option.

## A short demonstration

1. Start the local Go server.
2. Search for a cinema and select it.
3. Open the food menu and explain its example data.
4. Open Snack Planner and enter a budget and group size.
5. Generate a result and explain the subtotal, discount, and final estimate.
6. Add the snacks to the food plan and change a quantity.
7. Reload to show that the plan is stored on this device.
8. Open Print food plan to show the paper layout. Explain that this saved-plan total does not apply the recommendation promotion.

## Questions judges may ask

**Why Go?** Go handles the application's data structures and calculation rules. The same core code can run as a local server or be compiled for the browser.

**Why JavaScript if Go is the main language?** JavaScript handles browser events and connects the interface to Go. Go performs the recommendation and promotion calculations; HTML and CSS provide presentation.

**Is it AI?** No. It uses written selection rules and price comparisons, which can be explained step by step.

**Where is the database?** The current demonstration uses project-supplied data, an in-memory Go cache, and browser storage for personal food plans. It has no shared database.

**Can users really order food?** Not yet. Official cinema partnerships, staff dashboards, order tracking, and pickup slots are future improvements described in the README. Real ordering would need participating cinemas and shared backend storage.

**How did you test it?** Go tests check server responses, cinema and movie queries, calculations, and budget behaviour. Manual browser checks cover search, menus, recommendations, saved plans, responsive layouts, and print preview. Run `go test ./...` to inspect the automated tests.

**What are the limitations?** The data is illustrative, recommendations are heuristic, portions are not fully modelled, saved plans stay on one browser, and actual ordering is not implemented.

**Did you use help building it?** Describe your actual process honestly, including any AI assistance, and follow your assignment's rules. Practise explaining the functions you present rather than claiming authorship of work you cannot explain.

## Suggested reading order

1. `web/optimizer.html` — the form and result display.
2. `web/assets/js/cinema-client.js` — how the form requests Go results.
3. `internal/server/handlers_recommend.go` — the request's entry point in Go.
4. `pkg/recommender/optimizer.go` — food selection and budget check.
5. `pkg/recommender/cart.go` and `promotions.go` — totals and discounts.
6. `web/food-plan.html` — saving, quantities, and printing.

Learn the local-server flow first. Then learn the browser worker as a second way to call the same Go code.
