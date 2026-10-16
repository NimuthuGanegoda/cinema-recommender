# CineBite — demonstration and technical defence

## Project statement

CineBite helps Sri Lankan moviegoers browse cinema information, explore example snack menus, and create a food plan within a chosen budget. It combines a responsive interface with a shared Go recommendation engine that can run on a native HTTP server or in a browser using WebAssembly.

The engineering contribution is integrating budget enforcement, category-based selection, eligibility-aware promotion evaluation and two execution environments. It is a working software prototype. Commercial partnerships, live ticket inventory, user research and machine-learning accuracy have not been established.

## Five-minute demonstration

| Time | Show | Explain |
| --- | --- | --- |
| 0:00–0:40 | Home and cinema browser | The problem: planning cinema snacks with a limited group budget. Identify example data immediately. |
| 0:40–1:20 | Food menu for one cinema | Prices belong to the selected example menu. Browse before accepting a suggestion. |
| 1:20–2:30 | Snack planner: LKR 2,500, two people, Standard | Explain selected items, subtotal, one promotion, final cost and budget remaining. Read the result explanation. |
| 2:30–3:10 | Change budget and tier | A tier affects eligibility. Entering a higher-tier code does not bypass membership rules. Tier selection is for the demo. |
| 3:10–3:50 | Add suggestions to food plan; adjust quantities | This saves a device-local plan. It does not book, pay or order. |
| 3:50–4:30 | Architecture and tests | Show shared Go logic, automated invariants, and the recorded test output. |
| 4:30–5:00 | Limitations and next step | Explain the heuristic and propose provider integration followed by user evaluation. |

Before presenting: run on the actual laptop, open every page, check browser console and network requests, test small-screen layout, and repeat the full workflow. Keep the native server as the primary demo. Use the browser-only build only after testing its WASM startup in that browser. A successful compilation is not a browser execution test.

## Questions judges may ask

**Is this AI or machine learning?** No. It is an explainable rule-based recommender. It does not learn preferences or predict ratings. A later collaborative-filtering model would require real interaction data and a separate evaluation.

**Does it find the optimal bundle?** No global optimum is claimed. It constructs a category-based candidate and repairs it to satisfy the budget. `RecommendOptimalBundle` is a legacy internal name; the returned algorithm identifies the heuristic.

**Why this approach?** It is inspectable, works without training data, and supports reproducible demonstrations. Its weaknesses include category bias, limited group coverage, and potentially missing better combinations.

**How do you measure correctness?** Budget and availability invariants, exact promotion outcomes, rejected invalid inputs, race checks, and asset validation. These show implementation behavior, not whether users prefer the recommendation.

**Can I use a platinum code as a standard member?** No: code matching and tier eligibility are independent checks. The regression tests cover CINE, SCOPE and REDOPAY aliases. Real account verification is future work.

**Are savings real?** They are arithmetic results against sample campaign rules, not guaranteed cinema offers. No purchase is made.

**Why WebAssembly?** It reuses the Go engine in the browser. The tradeoff is download/startup overhead. JavaScript would also be a reasonable solution; WASM is an architectural choice rather than proof of better performance.

**Where is the database?** There is no database-backed user system. Catalogs are in memory and the food plan uses browser storage. The current scope does not require login. Multi-device persistence would need a backend datastore and account design.

**What does the value score mean?** It is items per person weighted by 25 plus discount fraction weighted by 50, rounded to one decimal. It is neither a percentage nor a validated satisfaction measure. It does not optimize the selection.

**What is original?** Describe the integration and changes you can personally explain. Identify external assets, tools and assisted development according to your university rules. Do not claim original algorithms, field research or partnerships without evidence.

## Reproducible arithmetic example

For an example basket containing popcorn LKR 900 and a drink LKR 600, the subtotal is LKR 1,500. Standard tier with no forced code receives CINE-SNACK15: discount LKR 225, final LKR 1,275. The existing promotion test asserts this exact case. It is not a promise that the heuristic chooses this basket for every budget.

## Evaluation plan for the next iteration

1. Freeze a catalog and create explicit expected baskets for budget, tier, cap and combination boundaries.
2. Compare this heuristic with cheapest-first and an exhaustive-search oracle on tiny bounded menus. Define an objective first: category coverage, drink coverage, then final cost. Do not report superiority until measured.
3. Recruit consenting participants to complete cinema selection, planning and quantity-editing tasks. Record completion rate, errors, time and feedback; report participant count and conditions.
4. Measure native and browser latency separately on named devices, with cold and warm startup. Keep engine microbenchmarks separate from user-visible latency.
5. Integrate an authorized provider feed with source timestamps, availability semantics and failure handling.

## Short presentation outline

1. Problem and intended users: cinema food planning within a budget.
2. Scope and working workflow: discover, inspect, suggest, save a local plan.
3. Architecture: shared Go engine, HTTP and WASM execution paths.
4. Algorithm: category heuristic, eligibility checks, budget repair and explanation.
5. Evidence: exact arithmetic, regressions, budget sweep and recorded checks.
6. Limitations and roadmap: sample data, no optimality proof, no user study, provider integration next.

Bring the marking rubric if available and map each criterion to a screen, source file, test or documented limitation. No alignment with an unseen rubric is claimed.
