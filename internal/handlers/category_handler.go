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

type CategoryHandler struct {
	service service.CategoryService
}

func NewCategoryHandler(service service.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

// GetCategories returns all categories
func (h *CategoryHandler) GetCategories(c echo.Context) error {
	ctx := c.Request().Context()

	categories, err := h.service.GetAll(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			dto.ErrorResponse(constants.ErrInternal, err.Error()))
	}

	return c.JSON(http.StatusOK,
		dto.SuccessResponse(constants.MsgCategoriesFetched, categories))
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

	category, err := h.service.Create(ctx, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			dto.ErrorResponse(constants.ErrInternal, err.Error()))
	}

	return c.JSON(http.StatusCreated,
		dto.SuccessResponse(constants.MsgCategoryCreated, category))
}

// GetCategory returns a category by id.
func (h *CategoryHandler) GetCategory(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			dto.ErrorResponse(constants.ErrInvalidID, err.Error()))
	}

	category, err := h.service.GetByID(ctx, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			dto.ErrorResponse(constants.ErrInternal, err.Error()))
	}
	if category == nil {
		return c.JSON(http.StatusNotFound,
			dto.ErrorResponse("Category not found", nil))
	}

	return c.JSON(http.StatusOK,
		dto.SuccessResponse(constants.MsgCategoryFetched, category))
}

// UpdateCategory updates an existing category.
func (h *CategoryHandler) UpdateCategory(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			dto.ErrorResponse(constants.ErrInvalidID, err.Error()))
	}

	var req models.UpdateCategoryRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest,
			dto.ErrorResponse(constants.ErrValidation, err.Error()))
	}

	category, err := h.service.Update(ctx, id, req)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound,
				dto.ErrorResponse("Category not found", nil))
		}
		return c.JSON(http.StatusInternalServerError,
			dto.ErrorResponse(constants.ErrInternal, err.Error()))
	}

	return c.JSON(http.StatusOK,
		dto.SuccessResponse(constants.MsgCategoryUpdated, category))
}

// DeleteCategory deletes an existing category.
func (h *CategoryHandler) DeleteCategory(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			dto.ErrorResponse(constants.ErrInvalidID, err.Error()))
	}

	if err := h.service.Delete(ctx, id); err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound,
				dto.ErrorResponse("Category not found", nil))
		}
		return c.JSON(http.StatusInternalServerError,
			dto.ErrorResponse(constants.ErrInternal, err.Error()))
	}

	return c.JSON(http.StatusOK,
		dto.SuccessResponse(constants.MsgCategoryDeleted, nil))
}
