package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/manish-npx/todo-go-echo/internal/constants"
	"github.com/manish-npx/todo-go-echo/internal/dto"
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
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse(constants.ErrValidation, err.Error()))
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse(constants.ErrValidation, err.Error()))
	}

	if err := h.service.Register(req); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse(constants.ErrInternal, err.Error()))
	}

	return c.JSON(http.StatusCreated, dto.SuccessResponse(constants.MsgUserRegistered, nil))
}

// POST /login
func (h *UserHandler) Login(c echo.Context) error {
	var req models.LoginRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse(constants.ErrValidation, err.Error()))
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse(constants.ErrValidation, err.Error()))
	}

	token, err := h.service.Login(req)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse("Invalid credentials", err.Error()))
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse(constants.MsgLoginSuccess, map[string]string{"token": token}))
}

// GET /users
func (h *UserHandler) GetUsers(c echo.Context) error {
	users, err := h.service.GetUsers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse(constants.ErrInternal, err.Error()))
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse(constants.MsgUsersFetched, users))
}
