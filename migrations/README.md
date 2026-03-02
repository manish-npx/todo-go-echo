# Migrations

This project keeps SQL migrations in this folder for production-safe schema changes.

## golang-migrate

Create:

```bash
migrate create -ext sql -dir migrations add_new_table
```

Run:

```bash
migrate -path migrations -database "postgres://postgres:password@localhost:5432/todo_app?sslmode=disable" up
```

## goose

Create:

```bash
goose -dir migrations create add_new_table sql
```

Run:

```bash
goose -dir migrations postgres "host=localhost port=5432 user=postgres password=password dbname=todo_app sslmode=disable" up
```
