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
	categories, err := h.repo.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch categories",
		})
	}
	return c.JSON(http.StatusOK, categories)
}

// GetCategory handles GET /categories/:id
func (h *CategoryHandler) GetCategory(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid category ID",
		})
	}

	category, err := h.repo.GetByID(id)
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

	// Get blog count for this category
	blogCount, _ := h.repo.GetBlogCount(id)

	// Create response with additional info
	response := struct {
		models.Category
		BlogCount int `json:"blog_count"`
	}{
		Category:  *category,
		BlogCount: blogCount,
	}

	return c.JSON(http.StatusOK, response)
}

// CreateCategory handles POST /categories
func (h *CategoryHandler) CreateCategory(c echo.Context) error {
	var req models.CreateCategoryRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// Validate
	if req.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Category name is required",
		})
	}

	category := &models.Category{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := h.repo.Create(category); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create category",
		})
	}

	return c.JSON(http.StatusCreated, category)
}

// UpdateCategory handles PUT /categories/:id
func (h *CategoryHandler) UpdateCategory(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid category ID",
		})
	}

	// Get existing category
	category, err := h.repo.GetByID(id)
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

	if err := h.repo.Update(category); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to update category",
		})
	}

	return c.JSON(http.StatusOK, category)
}

// DeleteCategory handles DELETE /categories/:id
func (h *CategoryHandler) DeleteCategory(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid category ID",
		})
	}

	// Check if category has blogs
	count, err := h.repo.GetBlogCount(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to check category usage",
		})
	}

	if count > 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Cannot delete category with existing blogs. Move or delete the blogs first.",
		})
	}

	if err := h.repo.Delete(id); err != nil {
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
