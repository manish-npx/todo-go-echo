package models

import (
	"time"
)

// models Todo represents a task in our todo list
type Todo struct {
	ID          int       `json:"id" db:"id"`
	Title       string    `json:"title" validate:"required,min=3"`
	Description string    `json:"description" validate:"required,min=5"`
	Completed   bool      `json:"completed" db:"completed"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// CreateTodoRequest is used when creating a new todo
type CreateTodoRequest struct {
	Title       string `json:"title" validate:"required,min=3"`
	Description string `json:"description" validate:"required,min=5"`
}

// UpdateTodoRequest is used when updating an existing todo
type UpdateTodoRequest struct {
	Title       *string `json:"title,omitempty" validate:"omitempty,min=3"`
	Description *string `json:"description,omitempty" validate:"omitempty,min=5"`
	Completed   *bool   `json:"completed,omitempty"`
}
