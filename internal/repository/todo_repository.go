package repository

import (
	"database/sql"
	"time"

	"github.com/manish-npx/todo-go-echo/internal/models"
)

// _ TodoRepository defines the interface for todo database operations
type TodoRepository interface {
	GetAll() ([]models.Todo, error)
	GetByID(id int) (*models.Todo, error)
	Create(todo *models.Todo) error
	Update(todo *models.Todo) error
	Delete(id int) error
}

type todoRepository struct {
	db *sql.DB
}

// NewTodoRepository creates a new todo repository
func NewTodoRepository(db *sql.DB) TodoRepository {
	return &todoRepository{db: db}
}

// GetAll retrieves all todos from the database
func (r *todoRepository) GetAll() ([]models.Todo, error) {
	query := `SELECT id, title, description, completed, created_at, updated_at FROM todos ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []models.Todo
	for rows.Next() {
		var todo models.Todo
		err := rows.Scan(
			&todo.ID,
			&todo.Title,
			&todo.Description,
			&todo.Completed,
			&todo.CreatedAt,
			&todo.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	return todos, nil
}

// GetByID retrieves a single todo by its ID
func (r *todoRepository) GetByID(id int) (*models.Todo, error) {
	query := `SELECT id, title, description, completed, created_at, updated_at FROM todos WHERE id = $1`

	var todo models.Todo
	err := r.db.QueryRow(query, id).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Description,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // _Todo not found
	}
	if err != nil {
		return nil, err
	}

	return &todo, nil
}

// Create inserts a new todo into the database
func (r *todoRepository) Create(todo *models.Todo) error {
	query := `INSERT INTO todos (title, description, completed, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5) RETURNING id`

	now := time.Now()
	todo.CreatedAt = now
	todo.UpdatedAt = now
	todo.Completed = false // New todos are not completed by default

	err := r.db.QueryRow(query, todo.Title, todo.Description, todo.Completed, todo.CreatedAt, todo.UpdatedAt).Scan(&todo.ID)
	return err
}

// Update modifies an existing todo
func (r *todoRepository) Update(todo *models.Todo) error {
	query := `UPDATE todos SET title = $1, description = $2, completed = $3, updated_at = $4 WHERE id = $5`

	todo.UpdatedAt = time.Now()

	result, err := r.db.Exec(query, todo.Title, todo.Description, todo.Completed, todo.UpdatedAt, todo.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// Delete removes a todo from the database
func (r *todoRepository) Delete(id int) error {
	query := `DELETE FROM todos WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}
