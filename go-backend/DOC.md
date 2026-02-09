## Go Backend Service

A production-style Go backend providing REST APIs for managing users and tasks, designed with clean architecture, persistence, caching, middleware, and strong error handling.

### Overview

- This service exposes a REST API that supports:

- User creation and retrieval

- Task creation, querying, and partial updates

- File-based persistence

- Redis-backed caching

- Request logging

- Rate limiting

- Input validation

- Health checks with readiness support

The system is designed to be simple, testable, and extensible, while demonstrating idiomatic Go practices.



### Architecture

The codebase follows a layered architecture with clear separation of concerns:

``` text
go-backend/
├── api/            # Request & response contracts
├── handlers/       # HTTP handlers (controllers)
├── server/         # Server lifecycle, routing
├── middleware/     # request logging, rate limiting
├── storage/        # Persistence layer (file-backed datastore)
│   └── file/
├── cache/          # Redis cache abstraction
├── util/           # Shared helpers (validation, etc.)
├── main.go         # Application entry point
└── docker-compose.yml

```

### Design Principles

Handlers contain HTTP logic only

Storage owns persistence and business data

Middleware handles cross-cutting concerns

API contracts are isolated from handlers

Dependencies are injected explicitly


### API Endpoints

#### Health
GET /health
GET /health?ready=true


Description

/health → liveness check

/health?ready=true → readiness check (verifies datastore availability)

Response

``` json

{
  "status": "ok",
  "message": "Go backend is running"
}
```
Status Codes

200 OK

503 Service Unavailable (readiness failure)


#### Users
GET  /api/users
POST /api/users
GET  /api/users/{id}

POST /api/users

Creates a new user.

Request
``` json
{
  "name": "Alice Johnson",
  "email": "alice@example.com",
  "role": "developer"
}
```

Validation

All fields required

Email format validated

Responses

201 Created

400 Bad Request

500 Internal Server Error


Tasks
GET  /api/tasks
POST /api/tasks
PUT  /api/tasks/{id}

POST /api/tasks

Creates a new task.

Request
``` json
{
  "title": "Implement caching",
  "status": "pending",
  "userId": 1
}
```

Validation

status must be one of: pending, in-progress, completed

userId must reference an existing user

Responses

201 Created

400 Bad Request

500 Internal Server Error


PUT /api/tasks/{id}

Partially updates an existing task.

Request (partial updates supported)

``` json
{
  "status": "completed"
}
```

or

``` json
{
  "title": "Finalize backend",
  "userId": 2
}
```

Behavior

Only provided fields are updated

Non-existent tasks return 404

Responses

200 OK

400 Bad Request

404 Not Found

500 Internal Server Error


### Persistence

The service uses file-based JSON persistence.

Behavior

Data is loaded from disk on startup

All mutations are persisted immediately

Writes are atomic and thread-safe

Missing or empty files are handled gracefully

Guarantees

Data persists across restarts

Concurrent access is safe

Corrupted or empty files do not crash the server


#### Caching

A Redis-backed cache is used to reduce redundant reads.

Cached Endpoints

GET /api/users

GET /api/tasks

Behavior

TTL: 2 minutes

Cache keys include query parameters

Cache is invalidated on:

User creation

Task creation

Task update

Redis is run locally using Docker Compose.

``` bash
docker compose up -d
```

#### Middleware

- Request Logging

All requests are logged with:

HTTP method

Request path

Response status

Duration

Example:

method=GET path=/api/tasks status=200 duration=1.23ms

- Rate Limiting

Per-IP rate limiting

Implemented using Go’s official golang.org/x/time/rate

Excess requests return 429 Too Many Requests


### Error Handling

Appropriate HTTP status codes are returned

Internal errors are logged but not exposed

Validation errors return clear messages

Handlers do not panic under invalid input


### Testing

The project includes unit tests for:

Storage layer (users & tasks)

Partial task updates

File persistence behavior

Run tests with:

``` bash
go test ./...
go test -cover ./...
```

Running the Project
Prerequisites

Go 1.21+

Docker & Docker Compose

Start Services
``` bash
docker compose up -d
go run main.go
```


Server starts on:

http://localhost:8080

### Design Decisions

#### File-Based Storage

Chosen to meet persistence requirements without introducing database complexity.

#### Redis Caching

Used to demonstrate real-world caching patterns and avoid re-implementing concurrency primitives.

#### Middleware-Driven Design

Logging, validation, and rate limiting are implemented as middleware to keep handlers focused.

#### Explicit Dependency Injection

All dependencies (storage, cache) are injected, improving testability and clarity.

#### Scope Control

Advanced features were implemented conservatively to prioritize correctness and maintainability.

### Summary

This backend demonstrates:

Idiomatic Go

Clean architecture

Safe concurrency

Practical middleware usage

Production-ready patterns

The implementation intentionally avoids over-engineering while remaining extensible (adding database and monitoring).