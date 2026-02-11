# Go Backend

A production-style REST API service for managing **users** and **tasks**, built with clean architecture, pluggable persistence (file or SQLite), Redis caching, and standard middleware (logging, rate limiting).

---

## Features

- **REST API** — Users (CRUD-style) and tasks (create, list, partial update)
- **Storage** — File-based JSON or SQLite; switch via `STORAGE_DRIVER`
- **Caching** — Redis with TTL; cache invalidation on writes
- **Middleware** — Request logging (structured), per-IP rate limiting
- **Health** — Liveness and readiness (datastore) checks
- **Validation** — Email format, task status, required fields

---

## Quick Start

**Prerequisites:** Go 1.21+, Docker (optional, for Redis)

```bash
# Dependencies
go mod tidy

# Optional: start Redis for caching
docker compose up -d

# Run (file storage by default, port 8080)
go run main.go
```

Server: **http://localhost:8080**

---

## Configuration

| Variable           | Default       | Description                          |
|-------------------|---------------|--------------------------------------|
| `PORT`            | `8080`        | HTTP server port                     |
| `STORAGE_DRIVER`  | `file`        | Storage backend: `file` or `sqlite`  |

**File storage** uses `data/data.json`.  
**SQLite** uses `data/app.db` (created on first run).

Redis is fixed at `localhost:6380` (match `docker-compose.yml`). Start Redis with:

```bash
docker compose up -d
```

---

## API Overview

### Health

| Method | Endpoint           | Description                    |
|--------|--------------------|--------------------------------|
| GET    | `/health`          | Liveness                       |
| GET    | `/health?ready=true` | Readiness (checks datastore) |

### Users

| Method | Endpoint        | Description        |
|--------|-----------------|--------------------|
| GET    | `/api/users`    | List all users     |
| GET    | `/api/users/{id}` | Get user by ID   |
| POST   | `/api/users`    | Create user        |

**POST body (create user):**

```json
{
  "name": "Alice Johnson",
  "email": "alice@example.com",
  "role": "developer"
}
```

- All fields required; email format validated; duplicate email → 400.

### Tasks

| Method | Endpoint         | Description                    |
|--------|------------------|--------------------------------|
| GET    | `/api/tasks`     | List tasks (optional `?status=`, `?userId=`) |
| POST   | `/api/tasks`     | Create task                    |
| PUT    | `/api/tasks/{id}`| Partial update (title, status, userId) |

**POST body (create task):**

```json
{
  "title": "Implement caching",
  "status": "pending",
  "userId": 1
}
```

- `status`: `pending` \| `in-progress` \| `completed`.  
- `userId` must exist.

**PUT body (partial update):** any subset of `title`, `status`, `userId`.

### Stats

| Method | Endpoint      | Description              |
|--------|---------------|--------------------------|
| GET    | `/api/stats`  | User count and task counts by status |

---

## Project Layout

```
go-backend/
├── api/              # Request/response DTOs
├── handlers/         # HTTP handlers (users, tasks, health, stats)
├── server/           # Server setup and route registration
├── middleware/       # Logging, rate limiting
├── storage/          # DataStore abstraction and implementations
│   ├── file/         # JSON file store
│   ├── db/           # SQLite store
│   └── cache/        # Redis cache
├── util/             # Validation helpers
├── main.go           # Entry point, storage/cache wiring
├── docker-compose.yml
└── data/             # Runtime data (JSON DB or SQLite file)
```

- **Handlers** — HTTP only; delegate to `DataStore` and cache.
- **Storage** — Implements persistence; file and SQLite implement same `DataStore` interface.
- **Middleware** — Cross-cutting: request logging, rate limiting.

---

## Caching

- **Cached:** `GET /api/users`, `GET /api/tasks` (key includes query string).
- **TTL:** 2 minutes.
- **Invalidation:** On any user create, task create, or task update (full cache clear).
- If Redis is down, handlers skip cache and use the datastore.

---

## Middleware

- **Logging** — Every request: method, path, status, duration (structured logs).
- **Rate limiting** — Per-IP; 429 when exceeded (1 req/s, burst 10).

---

## Testing

```bash
go test ./...
go test -cover ./...
```

Tests cover storage (file and DB), task partial updates, and validation.

---

## Build

```bash
go build -o go-backend
./go-backend
```

---

## Design Notes

- **File vs SQLite** — File store keeps the app self-contained; SQLite is for a single binary with no external DB server.
- **Explicit DI** — Store and cache are passed into the server and handlers for testability.
- **Structured logging** — JSON logs with consistent fields for production observability.
