package handlers

import (
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
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
		if errors.Is(err, service.ErrEmailAlreadyExists) {
			return c.JSON(http.StatusConflict, dto.ErrorResponse("Email already exists", err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse(constants.ErrInternal, err.Error()))
	}

	return c.JSON(http.StatusCreated, dto.SuccessResponse(constants.MsgUserRegistered, nil))
}

// POST /users (protected)
func (h *UserHandler) CreateUser(c echo.Context) error {
	var req models.CreateUserRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse(constants.ErrValidation, err.Error()))
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse(constants.ErrValidation, err.Error()))
	}

	if err := h.service.CreateUser(req); err != nil {
		if errors.Is(err, service.ErrEmailAlreadyExists) {
			return c.JSON(http.StatusConflict, dto.ErrorResponse("Email already exists", err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse(constants.ErrInternal, err.Error()))
	}

	return c.JSON(http.StatusCreated, dto.SuccessResponse(constants.MsgUserCreated, nil))
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

// GET /users/profile (protected)
func (h *UserHandler) Profile(c echo.Context) error {
	userID, err := getUserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse("Unauthorized", err.Error()))
	}

	user, err := h.service.GetProfile(userID)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse("User not found", err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse(constants.ErrInternal, err.Error()))
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse(constants.MsgUserProfileFetched, user))
}

// GET /users
func (h *UserHandler) GetUsers(c echo.Context) error {
	users, err := h.service.GetUsers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse(constants.ErrInternal, err.Error()))
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse(constants.MsgUsersFetched, users))
}

func getUserIDFromToken(c echo.Context) (int, error) {
	rawToken := c.Get("user")
	token, ok := rawToken.(*jwt.Token)
	if !ok || token == nil {
		return 0, errors.New("invalid auth token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid auth claims")
	}

	rawUserID, ok := claims["user_id"]
	if !ok {
		return 0, errors.New("missing user_id claim")
	}

	userIDFloat, ok := rawUserID.(float64)
	if !ok {
		return 0, errors.New("invalid user_id claim")
	}

	return int(userIDFloat), nil
}
