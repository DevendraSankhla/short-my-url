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

1. Fix the data races and the concurrency bugs — atomic counter or injected ID source, mutex on FileDB, `return` after every `WriteHeader`. Must pass `go test -race ./...`.
2. Make the redirect path non-blocking and fault-tolerant — queue push goes async (buffered channel + background flusher with retry), so a dead broker costs zero redirect latency.
3. Single module, one messagequeue package — remove the copy-pasted `internals/messageQueue`.
4. Define the `ClickEvent` schema (shortCode, timestamp, referrer, userAgent, IP); producer pushes the envelope, not a raw urlID. Use manual acks + prefetch + a context-cancellable consumer; document the at-least-once / idempotency decision.
5. Add `Update()` to the data layer and make stats real — then Redis cache-aside becomes obvious, not academic.
6. HTTP middleware: structured logging with request IDs, recover-from-panic, server timeouts, and a token-bucket rate limiter (implement it yourself — NFR).
7. Config from env vars, validated at startup (fail fast). Then a multi-stage Dockerfile with distroless final image.
8. Load test with k6 — measure redirects/sec before latency degrades — and report the numbers here in the readme.
