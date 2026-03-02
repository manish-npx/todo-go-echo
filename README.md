# Todo + Blog + User API (Go, Echo, PostgreSQL)

Production-oriented REST API with layered architecture, centralized error handling, structured logging, migrations, Docker setup, and optional ORM support.

## Stack

- Go `1.25.x` (toolchain pinned to patched version)
- Echo v4
- PostgreSQL
- `database/sql` repositories (primary path)
- Optional GORM (`gorm.io/gorm`) bootstrap
- Zap logger (`go.uber.org/zap`)
- `go-playground/validator`
- `golang-migrate` style SQL migrations

## Architecture

Request flow:

`Routes -> Middleware -> Handlers -> Services -> Repositories -> Database`

Design rules used:

- Handlers: HTTP concerns only (bind, validate, status code, response)
- Services: business logic and orchestration
- Repositories: DB queries and persistence details
- Middleware: cross-cutting concerns (auth, logging, error handler)
- DTOs: API response shape

## Folder Structure

```text
todo-go-echo/
  cmd/
    api/
      main.go                     # App bootstrap + DI wiring

  config/
    config.yaml                  # Local config
    config.docker.yaml           # Docker config

  migrations/
    001_create_users.up.sql
    001_create_users.down.sql
    002_create_categories_and_blogs.up.sql
    002_create_categories_and_blogs.down.sql
    003_create_todos.up.sql
    003_create_todos.down.sql

  internal/
    app/
      app.go                     # Single setup place: wiring + server lifecycle

    config/
      config.go                  # YAML config structs + loader

    constants/
      errors.go                  # Shared error messages/codes
      messages.go                # Shared success messages

    database/
      postgres.go                # sql.DB connection (primary)
      gorm.go                    # Optional GORM bootstrap

    dto/
      response.go                # Standard API response wrapper
      error.go                   # DTO types for validation errors

    handlers/
      user_handler.go
      todo_handler.go
      category_handler.go
      blog_handler.go

    logger/
      logger.go                  # Zap singleton logger

    middleware/
      setup.go                   # Recover, CORS, timeout, request logging
      jwt.go                     # JWT auth middleware
      error.go                   # Global HTTP error handler

    models/
      user.go
      todo.go
      category.go
      blog.go

    repository/
      user_repository.go
      todo_repository.go
      category_repository.go
      blog_repository.go

    routes/
      routes.go                  # API route groups and middleware attachment

    service/
      user_service.go
      todo_service.go
      category_service.go
      blog_service.go

    validator/
      validator.go               # Echo validator adapter

  Dockerfile
  docker-compose.yml
```

## Config Strategy

### Should constants be inside `config/`?

Short answer: **No**.

- `config/` should hold runtime/environment configuration (ports, DB host, secrets).
- `constants/` should hold compile-time static values (messages, error codes).

This separation is correct and recommended.

### Environment-based config

- `CONFIG_PATH` env var chooses config file.
- If not set, app uses `config/config.yaml`.

Example:

```bash
CONFIG_PATH=config/config.docker.yaml go run ./cmd/api
```

## Logging

- Zap logger is initialized inside `internal/app/app.go`.
- Request logs are emitted from `middleware/setup.go` with fields:
  - request id
  - method
  - URI
  - status
  - latency

## Centralized Error Handling

- Global handler: `internal/middleware/error.go`
- Wired during app bootstrap (`internal/app/app.go`) via:

```go
e.HTTPErrorHandler = middleware.ErrorHandler
```

This gives one consistent error response shape for uncaught handler errors.

## JWT Authentication

- Middleware: `internal/middleware/jwt.go`
- Protected routes are attached in `internal/routes/routes.go`.

## Database Migrations

The project uses `golang-migrate` naming (`*.up.sql`, `*.down.sql`).

Run locally:

```bash
migrate -path migrations -database "postgres://postgres:password@localhost:5432/todo_app?sslmode=disable" up
```

Rollback one step:

```bash
migrate -path migrations -database "postgres://postgres:password@localhost:5432/todo_app?sslmode=disable" down 1
```

## ORM Support

Primary runtime path is still `database/sql` repositories.

Optional GORM bootstrap is available:

- File: `internal/database/gorm.go`
- Config:

```yaml
orm:
  enabled: false
  auto_migrate: false
```

When enabled, app creates a GORM connection and can auto-migrate models.

## Local Run

```bash
go mod tidy
go test ./...
go run ./cmd/api
```

API base URL: `http://localhost:8080/api/v1`

## Docker Run

`docker-compose.yml` includes 3 services:

- `db` (PostgreSQL)
- `migrate` (runs schema migrations)
- `api` (Echo server)

Run:

```bash
docker compose up --build
```

## Notes

- Graceful shutdown handles `SIGINT` and `SIGTERM`.
- Server timeout is configurable in `config`.
- Input validation is enforced with `go-playground/validator` tags.

## Simple Structure View

If you want a simpler mental model, focus on these folders only:

- `internal/app` : startup and dependency wiring
- `internal/handlers` : HTTP layer
- `internal/service` : business logic
- `internal/repository` : database queries
- `internal/models` : domain models
- `internal/middleware` : auth/logging/error handling
- `internal/config` : runtime configuration
