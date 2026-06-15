# World Cup 2026 Predictions

A web app where users register, predict match scores for World Cup 2026, and compete on a leaderboard.

- **Backend:** Go 1.23 + Gin, pgx, JWT, bcrypt
- **Frontend:** React 18 + Vite + TypeScript + Tailwind + TanStack Query
- **Database:** PostgreSQL 16
- **Match data:** [football-data.org](https://www.football-data.org/) free tier (optional)
- **Deploy:** Docker Compose

## Project layout

```
worldcup-predictions/
├── backend/
│   ├── cmd/server/main.go      # HTTP server bootstrap
│   ├── internal/
│   │   ├── auth/               # bcrypt + JWT + middleware
│   │   ├── config/             # env loading
│   │   ├── db/                 # pgxpool
│   │   ├── handlers/           # HTTP handlers
│   │   ├── models/             # domain types
│   │   ├── repository/         # SQL queries
│   │   ├── scoring/            # points calculation
│   │   └── sync/               # football-data.org integration
│   ├── migrations/             # SQL migrations (golang-migrate)
│   ├── Dockerfile
│   └── go.mod
├── frontend/
│   ├── src/                    # React app
│   ├── Dockerfile              # nginx-served SPA
│   └── package.json
├── docker-compose.yml
└── .env.example
```

## Quick start (Docker)

```bash
cp .env.example .env
# edit .env — set JWT_SECRET (≥32 chars), strong POSTGRES_PASSWORD
docker compose up --build
```

Then open http://localhost.

## Local dev (without Docker)

**Backend:**
```bash
cd backend
# Iran: set GOPROXY first
export GOPROXY=https://goproxy.cn,direct
export DATABASE_URL=postgres://wcuser:wcpass@localhost:5432/wcdb?sslmode=disable
export JWT_SECRET=$(openssl rand -base64 48)
go run ./cmd/server
```

**Frontend:**
```bash
cd frontend
npm install
npm run dev      # http://localhost:5173, /api proxied to :8080
```

## API endpoints

| Method | Path                          | Auth | Purpose                       |
| ------ | ----------------------------- | ---- | ----------------------------- |
| POST   | /api/v1/auth/register         | —    | Create account                |
| POST   | /api/v1/auth/login            | —    | Get access + refresh tokens   |
| POST   | /api/v1/auth/refresh          | —    | Refresh access token          |
| GET    | /api/v1/me                    | ✓    | Current user                  |
| GET    | /api/v1/matches               | —    | List all fixtures             |
| GET    | /api/v1/predictions           | ✓    | My predictions                |
| PUT    | /api/v1/predictions/:match_id | ✓    | Create/update a prediction    |
| GET    | /api/v1/leaderboard?limit=100 | —    | Ranking                       |
| GET    | /healthz, /readyz             | —    | Liveness/readiness            |

## Scoring rules

Configurable via env (defaults in parens):

- **Exact score** — `POINTS_EXACT` (5)
- **Correct outcome + correct goal difference, wrong exact** — `POINTS_GD` (3)
- **Correct outcome only** — `POINTS_OUTCOME` (1)
- Otherwise 0

Predictions are **locked at kickoff** (server-enforced — client clock not trusted).

## Match data

Two ways to populate fixtures and results:

1. **football-data.org sync** — set `FOOTBALL_DATA_API_KEY` in `.env`. The backend polls every 10 minutes. NOTE: the WC2026 competition code (`WC` in `internal/sync/footballdata.go`) and the team-ID mapping are placeholders — football-data has not published the WC2026 schema at the time of writing. You will need to:
   - Confirm the competition code once published.
   - Implement team-ID mapping in `sync.SyncOnce` (the `TODO` line).
2. **Manual entry** — insert teams + matches directly via SQL or build a small admin panel (`is_admin` flag on `users` + admin endpoints).

## Technical considerations (production checklist)

- **Secrets.** Never commit `.env`. Use Docker secrets or a vault in real production.
- **TLS.** Front the stack with a reverse proxy (Traefik / Caddy / nginx) terminating HTTPS via Let's Encrypt. Add `Strict-Transport-Security`.
- **Rate limiting.** Add `gin-contrib/limiter` or a middleware in front of `/auth/login` and `/auth/register` (e.g. 5/min/IP).
- **Database backups.** `pg_dump` cron + offsite storage.
- **Migrations.** Already wired through `migrate/migrate` container. To add a new one, drop `0002_*.up.sql` and `.down.sql` into `backend/migrations/`.
- **Leaderboard scale.** The current SQL window function is fine up to ~10k users. Beyond that, switch to a Redis sorted set updated at the end of each `ScoreFinishedMatches` pass.
- **Concurrency.** Predictions are upserted with `ON CONFLICT (user_id, match_id)` — safe against double-submit.
- **Timezones.** All times stored as `TIMESTAMPTZ` UTC; rendered in browser locale.
- **Iran network.** Dockerfiles use `goproxy.cn` mirror for Go modules. If blocked, override with `--build-arg GOPROXY=https://mirror-go.runflare.com`. For npm, set `NPM_REGISTRY` build-arg if needed.
- **Observability.** Logs are JSON via `slog`. Add Prometheus by mounting `prometheus/client_golang` in `main.go` and exposing `/metrics`.
- **Anti-cheat.** The lock-out is server-side. Also log IP + user-agent per prediction for auditing.
- **Email verification & password reset** — not in this scaffold; add when going to real users (issue signed tokens, send via SMTP/SES).

## Next steps to ship

1. `go mod tidy` inside `backend/` to populate `go.sum`.
2. `npm install` inside `frontend/` to generate `package-lock.json`.
3. Implement the team-ID mapping in `internal/sync/footballdata.go` once the competition publishes.
4. Add an admin user / admin endpoints to manually create fixtures before launch.
5. Buy a domain, point it at the host, add Traefik + Let's Encrypt.
