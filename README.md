# 📝 Todo Application with Go, Echo & PostgreSQL

A complete beginner-friendly CRUD Todo API built using:

- Go (Golang)
- Echo Web Framework
- PostgreSQL
- YAML Configuration
- Clean Architecture (cmd + internal structure)

This project helps you understand real-world backend structure while learning important Go concepts like structs, interfaces, dependency injection, and database handling.

---

# 🚀 Tech Stack

- Go 1.21+
- Echo Framework
- PostgreSQL
- pgx (Postgres driver)
- YAML Configuration

---

# 🎯 Important Go Topics Covered

## Core Go Concepts

- Variables & Data Types
- Functions
- Packages & Modules
- Error Handling
- Pointers
- Structs
- Methods
- Interfaces
- Composition
- Dependency Injection
- Context package
- JSON encoding/decoding
- Struct Tags (`json`, `yaml`)

## Backend & Architecture Concepts

- Clean Architecture
- Repository Pattern
- Layered Architecture
- Configuration Management (YAML)
- Connection Pooling
- REST API Design
- HTTP Status Codes
- Middleware basics
- Environment configuration

---

# 📂 Project Structure

```
todo-go-echo/
│
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
│
├── config/
│   └── config.yaml              # Application configuration
│
├── internal/
│   ├── config/
│   │   └── config.go            # Config loader logic
│   │
│   ├── database/
│   │   └── postgres.go          # PostgreSQL connection setup
│   │
│   ├── handlers/
│   │   ├── blog_handler.go
│   │   ├── category_handler.go
│   │   └── todo_handler.go      # HTTP handlers (controllers)
│   │
│   ├── middleware/
│   │   └── appmiddleware.go     # App-level middleware setup
│   │
│   ├── models/
│   │   ├── blog.go
│   │   ├── category.go
│   │   └── todo.go              # Domain models
│   │
│   ├── repository/
│   │   ├── blog_repository.go
│   │   ├── category_repository.go
│   │   └── todo_repository.go   # Database access layer
│   │
│   ├── routes/
│   │   └── routes.go            # Route registration
│   │
│   ├── server/
│   │   └── server.go            # Graceful server startup & shutdown
│   │
│   └── utils/
│       ├── errors.go
│       └── validator.go
│
├── go.mod
├── go.sum
└── .gitignore

```

---

# 🧠 Application Flow

Client

HTTP Request
↓
Routes
↓
Handlers (Controllers)
↓
Repositories
↓
Database

Each layer has a single responsibility.

---

# 📋 Prerequisites

- Go 1.21+
- PostgreSQL 14+
- Basic understanding of REST APIs

---

# 🐘 Database Setup

Create database:

```sql
CREATE DATABASE todo_db;
```

Create table:

```sql
CREATE TABLE todos (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    completed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

---

# ⚙️ Configuration (config/config.yaml)

```yaml
server:
  port: "8080"

database:
  host: "localhost"
  port: "5432"
  user: "postgres"
  password: "postgres"
  dbname: "todo_db"
  sslmode: "disable"
```

---

# 🚀 Running the Project

```bash
git clone https://github.com/manish-npx/todo-go-echo.git
cd todo-go-echo

go mod download
go run cmd/server/main.go
```

Server runs on:

```
http://localhost:8080
```

---

# 📌 API Endpoints

| Method | Endpoint   | Description     |
| ------ | ---------- | --------------- |
| GET    | /todos     | Get all todos   |
| GET    | /todos/:id | Get single todo |
| POST   | /todos     | Create new todo |
| PUT    | /todos/:id | Update todo     |
| DELETE | /todos/:id | Delete todo     |

---

# 🧪 API Testing Examples

## Get all todos

```bash
curl http://localhost:8080/todos
```

## Create a todo

```bash
curl -X POST http://localhost:8080/todos \
  -H "Content-Type: application/json" \
  -d '{"title": "Learn Go", "description": "Study structs and interfaces"}'
```

## Update a todo

```bash
curl -X PUT http://localhost:8080/todos/1 \
  -H "Content-Type: application/json" \
  -d '{"completed": true}'
```

## Delete a todo

```bash
curl -X DELETE http://localhost:8080/todos/1
```

---

# 🧩 Struct Example

```go
type Todo struct {
    ID        int       `json:"id"`
    Title     string    `json:"title"`
    Completed bool      `json:"completed"`
}
```

---

# 🧩 Interface Example

```go
type TodoRepository interface {
    Create(todo *Todo) error
    GetAll() ([]Todo, error)
}
```

Interfaces allow flexible and testable design.

---

# 🏗 Why Use internal/ Folder?

Go convention:
Packages inside `internal/` cannot be imported outside this project.
It protects your private application logic.

---

# 🔮 Future Improvements

- JWT Authentication
- Middleware (Logger, Recovery)
- Docker support
- Unit testing
- Service layer
- Pagination
- Swagger documentation

---

# 👨‍💻 Author

Manish

---

# ⭐ If this project helped you

Give it a star on GitHub.
