package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/manish-npx/todo-go-echo/internal/models"
	"github.com/manish-npx/todo-go-echo/internal/repository"

	"github.com/labstack/echo/v4"
)

// TodoHandler handles HTTP requests for todos
type TodoHandler struct {
	repo repository.TodoRepository
}

// NewTodoHandler creates a new todo handler
func NewTodoHandler(repo repository.TodoRepository) *TodoHandler {
	return &TodoHandler{repo: repo}
}

// GetTodos handles GET /todos
func (h *TodoHandler) GetTodos(c echo.Context) error {
	todos, err := h.repo.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch todos",
		})
	}

	return c.JSON(http.StatusOK, todos)
}

// GetTodo handles GET /todos/:id
func (h *TodoHandler) GetTodo(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid todo ID",
		})
	}

	todo, err := h.repo.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch todo",
		})
	}

	if todo == nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Todo not found",
		})
	}

	return c.JSON(http.StatusOK, todo)
}

// CreateTodo handles POST /todos
func (h *TodoHandler) CreateTodo(c echo.Context) error {
	var req models.CreateTodoRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// Basic validation
	if req.Title == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Title is required",
		})
	}

	todo := &models.Todo{
		Title:       req.Title,
		Description: req.Description,
	}

	if err := h.repo.Create(todo); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create todo",
		})
	}

	return c.JSON(http.StatusCreated, todo)
}

// UpdateTodo handles PUT /todos/:id
func (h *TodoHandler) UpdateTodo(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid todo ID",
		})
	}

	// Get existing todo
	existingTodo, err := h.repo.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch todo",
		})
	}

	if existingTodo == nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Todo not found",
		})
	}

	// Bind update request
	var req models.UpdateTodoRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// Update fields if provided
	if req.Title != nil {
		existingTodo.Title = *req.Title
	}
	if req.Description != nil {
		existingTodo.Description = *req.Description
	}
	if req.Completed != nil {
		existingTodo.Completed = *req.Completed
	}

	// Save updates
	if err := h.repo.Update(existingTodo); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to update todo",
		})
	}

	return c.JSON(http.StatusOK, existingTodo)
}

// DeleteTodo handles DELETE /todos/:id
func (h *TodoHandler) DeleteTodo(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid todo ID",
		})
	}

	if err := h.repo.Delete(id); err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "Todo not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to delete todo",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Todo deleted successfully",
	})
}
