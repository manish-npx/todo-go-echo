package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/manish-npx/todo-go-echo/internal/models"
	"github.com/manish-npx/todo-go-echo/internal/repository"
)

// CategoryHandler handles HTTP requests for categories
type CategoryHandler struct {
	repo repository.CategoryRepository
}

// NewCategoryHandler creates a new category handler
func NewCategoryHandler(repo repository.CategoryRepository) *CategoryHandler {
	return &CategoryHandler{repo: repo}
}

// GetCategories handles GET /categories
func (h *CategoryHandler) GetCategories(c echo.Context) error {
	ctx := c.Request().Context()

	categories, err := h.repo.GetAll(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch categories",
		})
	}

	return c.JSON(http.StatusOK, categories)
}

// GetCategory handles GET /categories/:id
func (h *CategoryHandler) GetCategory(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid category ID",
		})
	}

	category, err := h.repo.GetByID(ctx, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch category",
		})
	}

	if category == nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Category not found",
		})
	}

	return c.JSON(http.StatusOK, category)
}

// CreateCategory handles POST /categories
func (h *CategoryHandler) CreateCategory(c echo.Context) error {
	ctx := c.Request().Context()

	var req models.CreateCategoryRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// Simple validation
	if req.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Category name is required",
		})
	}

	category := &models.Category{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := h.repo.Create(ctx, category); err != nil {
		// Check for duplicate name (PostgreSQL unique violation)
		// You can add more specific error checking here
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create category. Name might already exist.",
		})
	}

	return c.JSON(http.StatusCreated, category)
}

// UpdateCategory handles PUT /categories/:id
func (h *CategoryHandler) UpdateCategory(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid category ID",
		})
	}

	// Get existing category
	category, err := h.repo.GetByID(ctx, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch category",
		})
	}

	if category == nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Category not found",
		})
	}

	// Bind updates
	var req models.UpdateCategoryRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// Update fields if provided
	if req.Name != nil {
		category.Name = *req.Name
	}
	if req.Description != nil {
		category.Description = *req.Description
	}

	if err := h.repo.Update(ctx, category); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to update category",
		})
	}

	return c.JSON(http.StatusOK, category)
}

// DeleteCategory handles DELETE /categories/:id
func (h *CategoryHandler) DeleteCategory(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid category ID",
		})
	}

	if err := h.repo.Delete(ctx, id); err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "Category not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to delete category",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Category deleted successfully",
	})
}
