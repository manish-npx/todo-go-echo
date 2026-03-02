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

type BlogHandler struct {
	service service.BlogService
}

func NewBlogHandler(service service.BlogService) *BlogHandler {
	return &BlogHandler{service: service}
}

// GetBlogs handles GET /api/v1/blogs.
func (h *BlogHandler) GetBlogs(c echo.Context) error {
	ctx := c.Request().Context()

	categoryID := c.QueryParam("category")
	author := c.QueryParam("author")
	status := c.QueryParam("status")

	blogs, err := h.service.GetBlogs(ctx, categoryID, author, status)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse(constants.ErrInternal, err.Error()))
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse(constants.MsgBlogsFetched, blogs))
}

// GetBlog handles GET /api/v1/blogs/:id.
func (h *BlogHandler) GetBlog(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse(constants.ErrInvalidID, err.Error()))
	}

	blog, err := h.service.GetByID(ctx, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse(constants.ErrInternal, err.Error()))
	}
	if blog == nil {
		return c.JSON(http.StatusNotFound, dto.ErrorResponse("Blog not found", nil))
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse(constants.MsgBlogFetched, blog))
}

// CreateBlog handles POST /api/v1/blogs.
func (h *BlogHandler) CreateBlog(c echo.Context) error {
	ctx := c.Request().Context()

	var req models.CreateBlogRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse(constants.ErrValidation, err.Error()))
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse(constants.ErrValidation, err.Error()))
	}

	blog, err := h.service.Create(ctx, req)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid category ID", nil))
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse(constants.ErrInternal, err.Error()))
	}

	return c.JSON(http.StatusCreated, dto.SuccessResponse(constants.MsgBlogCreated, blog))
}

// UpdateBlog handles PUT /api/v1/blogs/:id.
func (h *BlogHandler) UpdateBlog(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse(constants.ErrInvalidID, err.Error()))
	}

	var req models.UpdateBlogRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse(constants.ErrValidation, err.Error()))
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse(constants.ErrValidation, err.Error()))
	}

	blog, err := h.service.Update(ctx, id, req)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse("Blog not found", nil))
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse(constants.ErrInternal, err.Error()))
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse(constants.MsgBlogUpdated, blog))
}

// DeleteBlog handles DELETE /api/v1/blogs/:id.
func (h *BlogHandler) DeleteBlog(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse(constants.ErrInvalidID, err.Error()))
	}

	if err := h.service.Delete(ctx, id); err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse("Blog not found", nil))
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse(constants.ErrInternal, err.Error()))
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse(constants.MsgBlogDeleted, nil))
}

// SearchBlogs handles GET /api/v1/blogs/search?q=term.
func (h *BlogHandler) SearchBlogs(c echo.Context) error {
	ctx := c.Request().Context()

	query := c.QueryParam("q")
	if query == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse(constants.ErrValidation, "search query is required"))
	}

	blogs, err := h.service.Search(ctx, query)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse(constants.ErrInternal, err.Error()))
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse(constants.MsgBlogsFetched, blogs))
}

// PublishBlog handles PATCH /api/v1/blogs/:id/publish.
func (h *BlogHandler) PublishBlog(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse(constants.ErrInvalidID, err.Error()))
	}

	blog, err := h.service.Publish(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse("Blog not found", nil))
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse(constants.ErrInternal, err.Error()))
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse(constants.MsgBlogPublished, blog))
}
