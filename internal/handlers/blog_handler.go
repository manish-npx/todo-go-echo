package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/manish-npx/todo-go-echo/internal/models"
	"github.com/manish-npx/todo-go-echo/internal/repository"
)

// BlogHandler handles HTTP requests for blogs
type BlogHandler struct {
	blogRepo     repository.BlogRepository
	categoryRepo repository.CategoryRepository
}

// NewBlogHandler creates a new blog handler
func NewBlogHandler(blogRepo repository.BlogRepository, categoryRepo repository.CategoryRepository) *BlogHandler {
	return &BlogHandler{
		blogRepo:     blogRepo,
		categoryRepo: categoryRepo,
	}
}

// GetBlogs handles GET /blogs
// Supports query parameters: ?category=1&author=john&status=published
func (h *BlogHandler) GetBlogs(c echo.Context) error {
	// Check for query parameters
	categoryID := c.QueryParam("category")
	author := c.QueryParam("author")
	status := c.QueryParam("status")

	var blogs []models.Blog
	var err error

	// Apply filters based on query parameters
	switch {
	case categoryID != "":
		id, _ := strconv.Atoi(categoryID)
		blogs, err = h.blogRepo.GetByCategory(id)
	case author != "":
		blogs, err = h.blogRepo.GetByAuthor(author)
	case status == "published":
		blogs, err = h.blogRepo.GetPublished()
	default:
		blogs, err = h.blogRepo.GetAll()
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch blogs",
		})
	}

	return c.JSON(http.StatusOK, blogs)
}

// GetBlog handles GET /blogs/:id
func (h *BlogHandler) GetBlog(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid blog ID",
		})
	}

	// Increment view count (don't return error to client if this fails)
	h.blogRepo.IncrementViews(id)

	// Get blog details
	blog, err := h.blogRepo.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch blog",
		})
	}

	if blog == nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Blog not found",
		})
	}

	return c.JSON(http.StatusOK, blog)
}

// CreateBlog handles POST /blogs
func (h *BlogHandler) CreateBlog(c echo.Context) error {
	var req models.CreateBlogRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// Validate required fields
	if req.Title == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Title is required",
		})
	}
	if req.Content == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Content is required",
		})
	}
	if req.Author == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Author is required",
		})
	}

	// Validate category if provided
	if req.CategoryID != nil {
		category, err := h.categoryRepo.GetByID(*req.CategoryID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to validate category",
			})
		}
		if category == nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid category ID",
			})
		}
	}

	// Set status (default to draft)
	status := models.StatusDraft
	if req.Status == "published" {
		status = models.StatusPublished
	}

	blog := &models.Blog{
		Title:      req.Title,
		Content:    req.Content,
		Author:     req.Author,
		CategoryID: req.CategoryID,
		Status:     status,
	}

	if err := h.blogRepo.Create(blog); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create blog",
		})
	}

	return c.JSON(http.StatusCreated, blog)
}

// UpdateBlog handles PUT /blogs/:id
func (h *BlogHandler) UpdateBlog(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid blog ID",
		})
	}

	// Get existing blog
	blog, err := h.blogRepo.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch blog",
		})
	}

	if blog == nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Blog not found",
		})
	}

	// Bind updates
	var req models.UpdateBlogRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// Update fields if provided
	if req.Title != nil {
		blog.Title = *req.Title
	}
	if req.Content != nil {
		blog.Content = *req.Content
	}
	if req.Author != nil {
		blog.Author = *req.Author
	}
	if req.CategoryID != nil {
		// Validate category if provided
		if *req.CategoryID != 0 {
			category, err := h.categoryRepo.GetByID(*req.CategoryID)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{
					"error": "Failed to validate category",
				})
			}
			if category == nil {
				return c.JSON(http.StatusBadRequest, map[string]string{
					"error": "Invalid category ID",
				})
			}
		}
		blog.CategoryID = req.CategoryID
	}
	if req.Status != nil {
		blog.Status = models.BlogStatus(*req.Status)
	}

	if err := h.blogRepo.Update(blog); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to update blog",
		})
	}

	return c.JSON(http.StatusOK, blog)
}

// DeleteBlog handles DELETE /blogs/:id
func (h *BlogHandler) DeleteBlog(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid blog ID",
		})
	}

	if err := h.blogRepo.Delete(id); err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "Blog not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to delete blog",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Blog deleted successfully",
	})
}

// SearchBlogs handles GET /blogs/search?q=term
func (h *BlogHandler) SearchBlogs(c echo.Context) error {
	query := c.QueryParam("q")
	if query == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Search query is required",
		})
	}

	blogs, err := h.blogRepo.Search(query)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to search blogs",
		})
	}

	return c.JSON(http.StatusOK, blogs)
}

// PublishBlog handles PATCH /blogs/:id/publish
func (h *BlogHandler) PublishBlog(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid blog ID",
		})
	}

	// Get existing blog
	blog, err := h.blogRepo.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch blog",
		})
	}

	if blog == nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Blog not found",
		})
	}

	// Update status to published
	blog.Status = models.StatusPublished

	if err := h.blogRepo.Update(blog); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to publish blog",
		})
	}

	return c.JSON(http.StatusOK, blog)
}
