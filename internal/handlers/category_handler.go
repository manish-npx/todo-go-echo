package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/manish-npx/todo-go-echo/internal/constants"
	"github.com/manish-npx/todo-go-echo/internal/dto"
	"github.com/manish-npx/todo-go-echo/internal/models"
	"github.com/manish-npx/todo-go-echo/internal/repository"
)

type CategoryHandler struct {
	repo repository.CategoryRepository
}

func NewCategoryHandler(repo repository.CategoryRepository) *CategoryHandler {
	return &CategoryHandler{repo: repo}
}

// GetCategories returns all categories
func (h *CategoryHandler) GetCategories(c echo.Context) error {
	ctx := c.Request().Context()

	categories, err := h.repo.GetAll(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			dto.ErrorResponse(constants.ErrInternal, err.Error()))
	}

	return c.JSON(http.StatusOK,
		dto.SuccessResponse("Categories fetched successfully", categories))
}

// CreateCategory creates a new category
func (h *CategoryHandler) CreateCategory(c echo.Context) error {
	ctx := c.Request().Context()

	var req models.CreateCategoryRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest,
			dto.ErrorResponse(constants.ErrValidation, err.Error()))
	}

	// Validate using tags
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest,
			dto.ErrorResponse(constants.ErrValidation, err.Error()))
	}

	category := &models.Category{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := h.repo.Create(ctx, category); err != nil {
		return c.JSON(http.StatusInternalServerError,
			dto.ErrorResponse(constants.ErrInternal, err.Error()))
	}

	return c.JSON(http.StatusCreated,
		dto.SuccessResponse("Category created successfully", category))
}
