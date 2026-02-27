# 🚛 HeavyPark API

QR-code-based parking management system for heavy load vehicles.

## Tech Stack

- **Go 1.24** — API server
- **chi/v5** — HTTP router
- **pgx/v5** — PostgreSQL driver
- **Supabase** — Auth (JWT) + Storage
- **Stripe/MercadoPago** — Billing (mock for now)
- **zap** — Structured logging

## Quick Start

```bash
# Set environment
export DATABASE_URL="postgres://..."
export PORT=8080

# Run
go run ./cmd/server

# Health check
curl http://localhost:8080/health
```

## Project Structure

See `internal/` for the Go code organized by clean architecture layers:
- `handler/` — HTTP handlers
- `service/` — Business logic
- `repository/` — Database access
- `model/` — Domain entities
- `dto/` — Request/response types
- `middleware/` — Auth, CORS, rate limiting
- `config/` — Environment config
- `router/` — Route registration

## Migrations

SQL files in `migrations/` — applied to Supabase PostgreSQL.
