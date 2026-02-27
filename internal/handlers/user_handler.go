package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/manish-npx/todo-go-echo/internal/models"
	"github.com/manish-npx/todo-go-echo/internal/service"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(s service.UserService) *UserHandler {
	return &UserHandler{service: s}
}

// POST /register
func (h *UserHandler) Register(c echo.Context) error {
	var req models.RegisterRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request",
		})
	}

	if err := h.service.Register(req); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"message": "User registered successfully",
	})
}

// POST /login
func (h *UserHandler) Login(c echo.Context) error {
	var req models.LoginRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request",
		})
	}

	token, err := h.service.Login(req)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}

// GET /users
func (h *UserHandler) GetUsers(c echo.Context) error {
	users, err := h.service.GetUsers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch users",
		})
	}

	return c.JSON(http.StatusOK, users)
}
