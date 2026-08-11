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

Skills Gained: layered architecture, caching, rate limiting, migrations, containerization basics.


my-backend-app/
├── cmd/                  # Application entry points (e.g., main.go, server.js)
├── src/ (or pkg/)
│   ├── config/           # Environment variables and configuration
│   ├── handlers/         # Presentation Layer (HTTP/REST controllers, gRPC)
│   ├── services/         # Business Logic Layer (Core rules, calculations)
│   ├── repositories/     # Data Access Layer (Database queries, external APIs)
│   ├── models/           # Data structures and entities
│   └── interfaces/       # Explicit interface definitions (if required by language)
├── tests/                # Unit and integration tests
├── README.md
└── package.json / go.mod # Dependency management


📚 Books (The Gold Standards)
"Clean Architecture" by Robert C. Martin (Uncle Bob) – The definitive guide on structuring applications using boundaries and interfaces.
"Domain-Driven Design (DDD)" by Eric Evans – Advanced concept, but great for learning how to model complex backend systems.

🌐 Articles & Concepts to Google
SOLID Principles: Focus heavily on the Dependency Inversion Principle (D), which dictates why we use interfaces.
The Repository Pattern: Read articles specifically breaking down the separation of services and repositories.
Hexagonal Architecture (Ports and Adapters): A widely adopted industry framework built entirely around interfaces.