Project 1 — URL Shortener

Functional Requirements

POST /shorten — accepts a long URL, returns a short code (custom alias optional)
GET /:code — redirects to original URL (301/302)
GET /:code/stats — returns click count, created date, last accessed
Expiring links (optional TTL on creation)
Basic user accounts so links can be listed per user (JWT auth)

Non-Functional Requirements

Postgres for persistence, with a proper migration tool (golang-migrate or atlas)
Redis cache for hot redirects (cache-aside pattern) — measure cache hit ratio
Rate limiting per IP/user (token bucket, implement it yourself once — don't just import a middleware blindly)
Structured logging (slog or zap) with request IDs
Graceful shutdown (drain in-flight requests on SIGTERM)
Config via env vars, validated at startup (fail fast if missing)
Dockerfile (multi-stage build, final image <20MB using scratch or distroless)

Testing Checklist

 Unit tests for the short-code generation logic (collision handling)
 Table-driven tests for input validation (malformed URLs, etc.)
 Integration test hitting a real Postgres via testcontainers-go
 Load test: how many redirects/sec before latency degrades? (use k6)
 Race detector clean (go test -race ./...)

## TODO — 2-Week Drill

 1. Fix the data races and the concurrency bugs — mutex on FileDB (snowflake already gives us the injected, thread-safe ID source), `return` after every `WriteHeader`. Must pass `go test -race ./...`.
    - FileDB.Get(): no mutex lock on read of db.data map — note concurrent map access in Go is a fatal runtime throw ("concurrent map read and map write"), not just a -race finding
    - FileDB.Insert(): TOCTOU race — existence check outside mutex (line 72); the method never locks at all
    - FileDB.saveToFile(): no lock during Marshal + WriteFile; concurrent callers corrupt output; WriteFile is non-atomic — a crash mid-write corrupts the file and the next boot fails (write to temp + rename)
    - FileDB.Insert(): double I/O (save then reload) — the reload can also clobber concurrent inserts
    - urls.go:137-138 — missing `return` after WriteHeader, continues to w.Write(); same handler sets Content-Type after WriteHeader (132-133) and misspells it `application/josn` -> DONE 
    - MessageQueue: no mutex on channel/connection; Push/Pull/Close race — amqp091-go v1.13.0 documents channels as not thread-safe ("do not share a single Channel instance between multiple goroutines"); Push runs on every redirect goroutine
    - No test files exist; worse, `go test ./...` fails at the vet step (see item 9) and `-race` can't run on this machine — gcc is 32-bit MinGW ("cc1.exe: 64-bit mode not compiled in"); install a 64-bit toolchain or run tests in WSL/CI
    - Unused `counter uint64` in application struct (main.go:17) — delete it; snowflake Generate() is already mutex-protected
2. Make the redirect path non-blocking and fault-tolerant — queue push goes async (buffered channel + background flusher with retry), so a dead broker costs zero redirect latency.
    - Push's 5s context only checks cancellation before publishing — it does not bound the socket write, so a wedged broker connection can stall a redirect goroutine past 5s (bounded only by heartbeat teardown, tens of seconds)
3. Single module, one messagequeue package — remove the copy-pasted `internals/messageQueue`.
    - The two copies have already diverged (analytics dropped Push's nil-check, added consumer sleeps/comments) — decide: one module with cmd/oltp + cmd/analytics, or go.work + a shared internal/queue that also owns the ClickEvent schema
4. Define the `ClickEvent` schema (shortCode, timestamp, referrer, userAgent, IP); producer pushes the envelope, not a raw urlID (currently it pushes url.ShortKey — and "" on lookup misses). Use manual acks + prefetch + a context-cancellable consumer; document the at-least-once / idempotency decision.
    - Decide producer delivery semantics too: publisher confirms (ch.Confirm + PublishWithDeferredConfirm) or documented fire-and-forget loss
5. Add `Update()` to the data layer and make stats real — then Redis cache-aside becomes obvious, not academic.
    - TotalHits starts at 1 on creation (urls.go:58 — should be 0) and ModifiedAt never changes
    - getUrlDetails is hardcoded to GetUrl(1) — use the path value and return a small DTO, not the raw DB struct
6. HTTP middleware: structured logging with request IDs, recover-from-panic, server timeouts, and a token-bucket rate limiter (implement it yourself — NFR).
    - Also: graceful shutdown on SIGTERM (drain in-flight requests) + handle http.ErrServerClosed, and a /healthz endpoint
7. Config from env vars, validated at startup (fail fast). Then a multi-stage Dockerfile with distroless final image.
    - Include DB_PATH / AMQP_URL / QUEUE_NAME / BASE_URL — today the DB path (D:\Downloads\urlShortner) and amqp://guest:guest@localhost:5672 are hardcoded, and the queue name "url-shortner" is a typo shared with analytics; add .env.example and gitignore .env; godotenv.Load() is CWD-dependent
8. Load test with k6 — measure redirects/sec before latency degrades — and report the numbers here in the readme.
9. Build gates — do this first, it blocks everything else: `go vet` fails on server.go:15 (`string(app.cfg.port)` converts the int to a rune, not digits), which also fails `go test ./...` before any test runs; fix with strconv.Itoa; run `go mod tidy` (snowflake, godotenv and amqp091 are marked `// indirect` but are directly imported).
10. Redirect correctness — the lookup flow is broken today: generated codes that decode to a missing ID return 500; anything else (every custom alias, random paths) falls through, pushes "" to the queue and answers 302 → Location: / instead of 404; implement the short-key lookup (no GetByShortKey in the data layer yet); decide the alias-vs-ID namespace rule (an alias that happens to be valid sqids can resolve to the wrong URL); enforce alias uniqueness; add handler tests for all three cases.
11. Input validation & security: scheme allowlist (http/https only), alias charset/length/reserved-path rules, http.MaxBytesReader on the body, 400 on any JSON decode error (a wrong-typed alias currently sneaks through and creates the URL), configured BASE_URL instead of reflecting r.Host in the response.
12. Cache package: the linked list is broken and unreferenced — PushFront panics on a 1-node list and orphans a node on 2+, Delete panics on the only node and infinite-loops on a missing ID (all reproduced); delete it or rewrite with tests before wiring cache-aside.
13. Postgres + a migration tool (core NFR, currently not in this drill): FileDB is a stopgap; the Database interface is the seam — port it to Postgres (golang-migrate/atlas) and add a testcontainers-go integration test.
14. Spec gaps not yet in this drill: JWT accounts + per-user link listing (URL model needs user_id) and expiring links (expiresAt + check on redirect); document the 301-vs-302 decision.
15. Test suite: table-driven validation tests, sqids/snowflake collision unit tests, handler/mux tests (redirect miss, alias, stats), FileDB tests, plus a CI job that runs `go vet ./...` and `go test -race ./...`.
16. Housekeeping: os.Exit(0) on embed failure (routes.go:16), sqids.New() built per request (hoist into application), debug fetch('http://localhost:4000/bM') shipped in the embedded web/index.html, "Url Shorted" typo + no Location header on 201, identifier typos (RabbmitMQ/messaageQueue/URlModel/DouobleLinkedList), internals/ → internal/, fmt.Println → slog.
