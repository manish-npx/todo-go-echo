package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

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

// GetCategory returns a category by id.
func (h *CategoryHandler) GetCategory(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			dto.ErrorResponse(constants.ErrInvalidID, err.Error()))
	}

	category, err := h.repo.GetByID(ctx, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			dto.ErrorResponse(constants.ErrInternal, err.Error()))
	}
	if category == nil {
		return c.JSON(http.StatusNotFound,
			dto.ErrorResponse("Category not found", nil))
	}

	return c.JSON(http.StatusOK,
		dto.SuccessResponse("Category fetched successfully", category))
}

// UpdateCategory updates an existing category.
func (h *CategoryHandler) UpdateCategory(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			dto.ErrorResponse(constants.ErrInvalidID, err.Error()))
	}

	existing, err := h.repo.GetByID(ctx, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			dto.ErrorResponse(constants.ErrInternal, err.Error()))
	}
	if existing == nil {
		return c.JSON(http.StatusNotFound,
			dto.ErrorResponse("Category not found", nil))
	}

	var req models.UpdateCategoryRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest,
			dto.ErrorResponse(constants.ErrValidation, err.Error()))
	}

	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Description != nil {
		existing.Description = *req.Description
	}

	if err := h.repo.Update(ctx, existing); err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound,
				dto.ErrorResponse("Category not found", nil))
		}
		return c.JSON(http.StatusInternalServerError,
			dto.ErrorResponse(constants.ErrInternal, err.Error()))
	}

	return c.JSON(http.StatusOK,
		dto.SuccessResponse("Category updated successfully", existing))
}

// DeleteCategory deletes an existing category.
func (h *CategoryHandler) DeleteCategory(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			dto.ErrorResponse(constants.ErrInvalidID, err.Error()))
	}

	if err := h.repo.Delete(ctx, id); err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound,
				dto.ErrorResponse("Category not found", nil))
		}
		return c.JSON(http.StatusInternalServerError,
			dto.ErrorResponse(constants.ErrInternal, err.Error()))
	}

	return c.JSON(http.StatusOK,
		dto.SuccessResponse("Category deleted successfully", nil))
}
