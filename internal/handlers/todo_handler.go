package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/manish-npx/todo-go-echo/internal/constants"
	"github.com/manish-npx/todo-go-echo/internal/dto"
	"github.com/manish-npx/todo-go-echo/internal/models"
	"github.com/manish-npx/todo-go-echo/internal/service"
)

type TodoHandler struct {
	service service.TodoService
}

func NewTodoHandler(s service.TodoService) *TodoHandler {
	return &TodoHandler{s}
}

// CreateTodo handles POST /api/v1/todos.
func (h *TodoHandler) CreateTodo(c echo.Context) error {
	var req models.CreateTodoRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse(constants.ErrValidation, err.Error()))
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse(constants.ErrValidation, err.Error()))
	}

	todo, err := h.service.Create(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse(constants.ErrInternal, err.Error()))
	}

	return c.JSON(http.StatusCreated, dto.SuccessResponse(constants.MsgTodoCreated, todo))
}

// GetTodos handles GET /api/v1/todos.
func (h *TodoHandler) GetTodos(c echo.Context) error {
	todos, err := h.service.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse(constants.ErrInternal, err.Error()))
	}
	return c.JSON(http.StatusOK, dto.SuccessResponse(constants.MsgTodosFetched, todos))
}

// GetTodo handles GET /api/v1/todos/:id.
func (h *TodoHandler) GetTodo(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse(constants.ErrInvalidID, err.Error()))
	}

	todo, err := h.service.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse(constants.ErrInternal, err.Error()))
	}
	if todo == nil {
		return c.JSON(http.StatusNotFound, dto.ErrorResponse("Todo not found", nil))
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse(constants.MsgTodoFetched, todo))
}

// UpdateTodo handles PUT /api/v1/todos/:id.
func (h *TodoHandler) UpdateTodo(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse(constants.ErrInvalidID, err.Error()))
	}

	var req models.UpdateTodoRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse(constants.ErrValidation, err.Error()))
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse(constants.ErrValidation, err.Error()))
	}

	todo, err := h.service.Update(id, req)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse("Todo not found", nil))
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse(constants.ErrInternal, err.Error()))
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse(constants.MsgTodoUpdated, todo))
}

// DeleteTodo handles DELETE /api/v1/todos/:id.
func (h *TodoHandler) DeleteTodo(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse(constants.ErrInvalidID, err.Error()))
	}

	if err := h.service.Delete(id); err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse("Todo not found", nil))
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse(constants.ErrInternal, err.Error()))
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse(constants.MsgTodoDeleted, nil))
}
