package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/manish-npx/todo-go-echo/internal/constants"
	"github.com/manish-npx/todo-go-echo/internal/dto"
	"github.com/manish-npx/todo-go-echo/internal/services"
	"github.com/manish-npx/todo-go-echo/internal/utils"
)

type TodoHandler struct {
	service services.TodoService
}

func NewTodoHandler(s services.TodoService) *TodoHandler {
	return &TodoHandler{s}
}

func (h *TodoHandler) Create(c echo.Context) error {
	var req dto.CreateTodoRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, constants.ErrInvalidPayload, err)
	}

	if err := c.Validate(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "Validation failed", err)
	}

	todo, err := h.service.Create(req.Title)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, "Failed to create", err)
	}

	return utils.Success(c, constants.Success, todo)
}

func (h *TodoHandler) GetAll(c echo.Context) error {
	todos, err := h.service.GetAll()
	if err != nil {
		return utils.Error(c, 500, "Failed", err)
	}
	return utils.Success(c, constants.Success, todos)
}

func (h *TodoHandler) Update(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return utils.Error(c, 400, constants.ErrInvalidID, err)
	}

	var req dto.UpdateTodoRequest
	c.Bind(&req)

	todo, err := h.service.Update(id, req.Title, req.Completed)
	if err != nil {
		return utils.Error(c, 500, "Update failed", err)
	}

	return utils.Success(c, constants.Success, todo)
}

func (h *TodoHandler) Delete(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return utils.Error(c, 400, constants.ErrInvalidID, err)
	}

	if err := h.service.Delete(id); err != nil {
		return utils.Error(c, 500, "Delete failed", err)
	}

	return utils.Success(c, "Deleted successfully", nil)
}
