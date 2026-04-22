# game-score-api

A production-ready Game Score Ranking API built with **Go**, **PostgreSQL**, and **Redis**.
Deployed live on AWS EC2 — [http://13.114.19.198:8080/health](http://13.114.19.198:8080/health)

---

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                        Client                           │
└────────────────────────┬────────────────────────────────┘
                         │ HTTP
┌────────────────────────▼────────────────────────────────┐
│                    AWS EC2 (ap-northeast-1)              │
│  ┌──────────────────────────────────────────────────┐   │
│  │                  Gin HTTP Server                 │   │
│  │  ┌──────────┐  ┌────────────┐  ┌─────────────┐  │   │
│  │  │  Handler │→ │  Service   │→ │ Repository  │  │   │
│  │  │  (HTTP)  │  │ (Business) │  │ (DB Access) │  │   │
│  │  └──────────┘  └────────────┘  └──────┬──────┘  │   │
│  │                                        │         │   │
│  │       ┌────────────────────────────────┘         │   │
│  │       │                                          │   │
│  │  ┌────▼────────┐      ┌───────────────┐          │   │
│  │  │ PostgreSQL  │      │  Redis Cache  │          │   │
│  │  │ (Users /    │      │ (Rankings TTL │          │   │
│  │  │  Scores)    │      │  60s)         │          │   │
│  │  └─────────────┘      └───────────────┘          │   │
│  └──────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────┘
```

**Request flow:**
1. Client sends HTTP request
2. JWT middleware validates Bearer token (auth-required routes)
3. Handler parses request → delegates to Service layer
4. Service applies business logic → calls Repository
5. Repository queries PostgreSQL (write/read) or Redis (ranking cache)
6. Response returned as JSON

---

## Tech Stack

| Layer | Technology | Version | Purpose |
|-------|-----------|---------|---------|
| Language | Go | 1.23 | Application runtime |
| Framework | Gin | v1.9.1 | HTTP router & middleware |
| Database | PostgreSQL | 16 | User & score persistence |
| Cache | Redis | 7 | Ranking cache (TTL: 60s) |
| Auth | golang-jwt | v5.2.1 | JWT token generation & validation |
| DB Driver | pgx | v5.5.5 | PostgreSQL connection pool |
| Redis Client | go-redis | v9.5.1 | Redis client |
| Password | bcrypt | x/crypto | Password hashing |
| Container | Docker / Docker Compose | — | Local dev environment |
| CI/CD | GitHub Actions | — | Auto-deploy to EC2 |
| Deploy | AWS EC2 | t3.micro | Production server |

---

## API Endpoints

| Method | Path | Auth | Description |
|--------|------|:----:|-------------|
| `GET` | `/health` | — | Health check (returns version & status) |
| `POST` | `/api/v1/auth/register` | — | Register a new player account |
| `POST` | `/api/v1/auth/login` | — | Login and receive JWT token |
| `GET` | `/api/v1/rankings` | — | Get top rankings (Redis-cached, TTL 60s) |
| `GET` | `/api/v1/players` | — | List all players |
| `GET` | `/api/v1/players/:id` | — | Get a specific player by UUID |
| `POST` | `/api/v1/scores` | **JWT** | Submit a new score |
| `GET` | `/api/v1/scores/me` | **JWT** | Get your own score history |

> **Auth:** Pass `Authorization: Bearer <token>` header for JWT-required routes.

See full request/response schemas in [`docs/swagger.yaml`](docs/swagger.yaml).

---

## Project Structure

```
game-score-api/
├── cmd/
│   └── server/
│       └── main.go             # Entry point, DI wiring, routing
├── internal/
│   ├── handler/
│   │   ├── auth.go             # POST /register, POST /login
│   │   ├── score.go            # POST /scores, GET /scores/me, GET /rankings
│   │   └── player.go           # GET /players, GET /players/:id
│   ├── service/
│   │   ├── auth_service.go          # Register (bcrypt), Login (JWT issue)
│   │   ├── auth_service_test.go     # Unit tests — Register / Login (5 cases)
│   │   ├── score_service.go         # Score submit, Redis cache logic
│   │   └── score_service_test.go    # Unit tests — PostScore / GetRankings (5 cases)
│   ├── repository/
│   │   ├── interfaces.go       # UserRepositoryInterface, ScoreRepositoryInterface
│   │   ├── user_repo.go        # Create, FindByEmail, FindByID, FindAll
│   │   └── score_repo.go       # Create, FindByUserID, GetRankings, CountAll
│   └── model/
│       ├── user.go             # User struct, RegisterRequest, LoginRequest/Response
│       └── score.go            # Score struct, PostScoreRequest, RankingResponse
├── pkg/
│   ├── auth/
│   │   └── jwt.go              # JWT generate & validate
│   ├── database/
│   │   ├── postgres.go         # PostgreSQL connection pool
│   │   └── redis.go            # Redis connection
│   └── middleware/
│       └── auth.go             # JWT authentication middleware
├── migrations/
│   ├── 001_create_users.sql    # users table with UUID PK
│   └── 002_create_scores.sql   # scores table + performance indexes
├── docs/
│   └── swagger.yaml            # OpenAPI 3.0 spec
├── .github/
│   └── workflows/
│       └── deploy.yml          # GitHub Actions — build & deploy to EC2
├── docker-compose.yml          # PostgreSQL + Redis for local dev
├── Dockerfile                  # Multi-stage build (alpine)
└── .env.example                # Environment variable template
```

---

## Local Development

### Prerequisites

- Go 1.23+
- Docker & Docker Compose
- Git

### 1. Clone & configure

```bash
git clone https://github.com/natuomi/game-score-api.git
cd game-score-api
cp .env.example .env
# Edit .env — change JWT_SECRET at minimum
```

### 2. Start PostgreSQL & Redis

```bash
# Start only the dependencies (PostgreSQL + Redis)
docker compose up -d postgres redis
```

### 3. Run the API server

```bash
go run ./cmd/server
```

The server starts on `http://localhost:8080`.

> **Note:** The Dockerfile-based API container has a known runc path issue on some setups.
> Running `go run ./cmd/server` directly is the recommended local approach.

### 4. Verify it's working

```bash
# Health check
curl http://localhost:8080/health
# Expected: {"status":"ok","version":"1.0.0"}

# Register a player
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","email":"alice@example.com","password":"password123"}'

# Login and capture token
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"alice@example.com","password":"password123"}' | jq -r '.token')

# Submit a score (JWT required)
curl -X POST http://localhost:8080/api/v1/scores \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"score":9500,"game_mode":"classic"}'

# Get rankings (cached after first call)
curl http://localhost:8080/api/v1/rankings
# First call:  "cached": false  (PostgreSQL)
# Second call: "cached": true   (Redis, TTL 60s)
```

### 5. Run DB migrations (first time only)

```bash
# Connect to PostgreSQL and run migration files
psql -h localhost -U postgres -d game_score -f migrations/001_create_users.sql
psql -h localhost -U postgres -d game_score -f migrations/002_create_scores.sql
```

---

## Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `PORT` | HTTP server port | `8080` |
| `GIN_MODE` | Gin mode (`debug` / `release`) | `release` |
| `JWT_SECRET` | Secret key for JWT signing | `change-me-in-production` |
| `DB_HOST` | PostgreSQL host | `localhost` |
| `DB_PORT` | PostgreSQL port | `5432` |
| `DB_USER` | PostgreSQL user | `postgres` |
| `DB_PASSWORD` | PostgreSQL password | `postgres` |
| `DB_NAME` | PostgreSQL database name | `game_score` |
| `REDIS_HOST` | Redis host | `localhost` |
| `REDIS_PORT` | Redis port | `6379` |
| `REDIS_PASSWORD` | Redis password (leave empty if none) | `` |
| `RANKING_CACHE_TTL` | Ranking cache TTL in seconds | `60` |
| `RANKING_LIMIT` | Max entries in rankings | `100` |

Copy `.env.example` to `.env` and fill in your values. **Never commit `.env` to Git.**

---

## CI/CD — GitHub Actions

On every push to `main`, the workflow in `.github/workflows/deploy.yml`:
1. Sets up Go 1.23
2. Cross-compiles a Linux binary (`GOOS=linux GOARCH=amd64`)
3. Transfers the binary to AWS EC2 via SCP
4. Restarts the server process over SSH

**Required GitHub Secrets:**

| Secret | Description |
|--------|-------------|
| `EC2_HOST` | Public IP or hostname of the EC2 instance |
| `EC2_USER` | SSH login user (e.g., `ec2-user`, `ubuntu`) |
| `EC2_SSH_KEY` | Private SSH key (PEM content) for EC2 access |

Set these in **GitHub → Settings → Secrets and variables → Actions**.

---

## API Documentation

Full OpenAPI 3.0 specification: [`docs/swagger.yaml`](docs/swagger.yaml)

You can preview it with:
- [Swagger Editor](https://editor.swagger.io/) — paste the YAML content
- VS Code extension: **OpenAPI (Swagger) Editor**

---

## Running Tests

Unit tests are provided for the service layer (`internal/service/`).  
No database or Redis connection is required — all external dependencies are replaced with in-memory hand-written mocks.

### Run all unit tests

```bash
go test ./internal/service/...
```

### Run with verbose output

```bash
go test ./internal/service/... -v
```

### Run with coverage report

```bash
go test ./internal/service/... -cover
```

### What is tested

| Test file | Covered scenarios |
|-----------|-------------------|
| `auth_service_test.go` | Register success (bcrypt hash verification, repo.Create called) |
| | Register failure (duplicate email) |
| | Login success (JWT token generated, user returned) |
| | Login failure — password mismatch |
| | Login failure — user not found (returns same error as mismatch, prevents enumeration) |
| `score_service_test.go` | PostScore success (repo.Create called, Redis cache invalidated) |
| | PostScore with explicit GameMode |
| | GetRankings cache HIT (Redis hit, DB not called, `cached: true`) |
| | GetRankings cache MISS (DB queried, result stored in Redis, `cached: false`) |
| | GetRankings DB error propagation |

### Design notes

- Mocks are hand-written implementations of `repository.UserRepositoryInterface`, `repository.ScoreRepositoryInterface`, and `service.RedisClient` — no third-party mock libraries required.
- `JWT_SECRET` is set per-test via `t.Setenv` (automatically restored after each test).
- Tests that use `t.Setenv` run sequentially (not parallel) to comply with Go's testing restrictions.

---

## Live Demo

- **Health check:** [http://13.114.19.198:8080/health](http://13.114.19.198:8080/health)
- **Rankings:** [http://13.114.19.198:8080/api/v1/rankings](http://13.114.19.198:8080/api/v1/rankings)

---

## Roadmap

| Phase | Status | Description |
|-------|--------|-------------|
| Phase 1 | Done | Local dev — Go + Gin + PostgreSQL + Redis + JWT |
| Phase 2 | Done | AWS EC2 deploy — binary + systemd/nohup |
| Phase 3 | In Progress | CI/CD — GitHub Actions auto-deploy |
| Phase 4 | Planned | Testing — unit & integration tests |
| Phase 5 | Planned | Kubernetes — EKS migration (optional) |
