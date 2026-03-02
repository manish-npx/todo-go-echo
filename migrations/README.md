# Migrations

Production apps should version schema changes through migration files.

This folder uses `golang-migrate` file naming:

- `NNN_name.up.sql`
- `NNN_name.down.sql`

## Create a new migration

```bash
migrate create -ext sql -dir migrations add_indexes
```

## Apply migrations (local)

```bash
migrate -path migrations -database "postgres://postgres:password@localhost:5432/todo_app?sslmode=disable" up
```

## Rollback one migration

```bash
migrate -path migrations -database "postgres://postgres:password@localhost:5432/todo_app?sslmode=disable" down 1
```

## Docker compose flow

`docker-compose.yml` contains a dedicated `migrate` service that runs `up` before API startup.
