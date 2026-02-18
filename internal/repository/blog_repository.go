package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/manish-npx/todo-go-echo/internal/models"
)

type blogRepository struct {
	db *sql.DB
}

type BlogRepository interface {
	GetAll(ctx context.Context) ([]models.Blog, error)
	GetByID(ctx context.Context, id int) (*models.Blog, error)
	GetByCategory(ctx context.Context, categoryID int) ([]models.Blog, error)
	GetPublished(ctx context.Context) ([]models.Blog, error) // ADD CONTEXT
	GetByAuthor(ctx context.Context, author string) ([]models.Blog, error)
	Create(ctx context.Context, blog *models.Blog) error
	Update(ctx context.Context, blog *models.Blog) error
	Delete(ctx context.Context, id int) error
	IncrementViews(ctx context.Context, id int) error
	Search(ctx context.Context, query string) ([]models.Blog, error)
}

// NewBlogRepository creates a new blog repository
func NewBlogRepository(db *sql.DB) BlogRepository {
	return &blogRepository{db: db}
}

// GetAll retrieves all blogs
func (r *blogRepository) GetAll(ctx context.Context) ([]models.Blog, error) {
	query := `SELECT id, title, content, author, category_id, status, views, created_at, updated_at, published_at
			  FROM blogs ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blogs []models.Blog
	for rows.Next() {
		var blog models.Blog
		var categoryID sql.NullInt64

		err := rows.Scan(
			&blog.ID,
			&blog.Title,
			&blog.Content,
			&blog.Author,
			&categoryID,
			&blog.Status,
			&blog.Views,
			&blog.CreatedAt,
			&blog.UpdatedAt,
			&blog.PublishedAt,
		)
		if err != nil {
			return nil, err
		}

		if categoryID.Valid {
			id := int(categoryID.Int64)
			blog.CategoryID = &id
		}

		blogs = append(blogs, blog)
	}

	return blogs, nil
}

// GetByID retrieves a blog by ID
func (r *blogRepository) GetByID(ctx context.Context, id int) (*models.Blog, error) {
	query := `SELECT id, title, content, author, category_id, status, views, created_at, updated_at, published_at
			  FROM blogs WHERE id = $1`

	var blog models.Blog
	var categoryID sql.NullInt64

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&blog.ID,
		&blog.Title,
		&blog.Content,
		&blog.Author,
		&categoryID,
		&blog.Status,
		&blog.Views,
		&blog.CreatedAt,
		&blog.UpdatedAt,
		&blog.PublishedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if categoryID.Valid {
		id := int(categoryID.Int64)
		blog.CategoryID = &id
	}

	return &blog, nil
}

// GetByCategory retrieves blogs by category
func (r *blogRepository) GetByCategory(ctx context.Context, categoryID int) ([]models.Blog, error) {
	query := `SELECT id, title, content, author, category_id, status, views, created_at, updated_at, published_at
			  FROM blogs WHERE category_id = $1 ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blogs []models.Blog
	for rows.Next() {
		var blog models.Blog
		var catID sql.NullInt64

		err := rows.Scan(
			&blog.ID,
			&blog.Title,
			&blog.Content,
			&blog.Author,
			&catID,
			&blog.Status,
			&blog.Views,
			&blog.CreatedAt,
			&blog.UpdatedAt,
			&blog.PublishedAt,
		)
		if err != nil {
			return nil, err
		}

		if catID.Valid {
			id := int(catID.Int64)
			blog.CategoryID = &id
		}

		blogs = append(blogs, blog)
	}

	return blogs, nil
}

// GetPublished retrieves published blogs - FIXED: Added ctx parameter
func (r *blogRepository) GetPublished(ctx context.Context) ([]models.Blog, error) {
	query := `SELECT id, title, content, author, category_id, status, views, created_at, updated_at, published_at
			  FROM blogs WHERE status = 'published' ORDER BY published_at DESC`

	rows, err := r.db.QueryContext(ctx, query) // Use QueryContext with ctx
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blogs []models.Blog
	for rows.Next() {
		var blog models.Blog
		var categoryID sql.NullInt64

		err := rows.Scan(
			&blog.ID,
			&blog.Title,
			&blog.Content,
			&blog.Author,
			&categoryID,
			&blog.Status,
			&blog.Views,
			&blog.CreatedAt,
			&blog.UpdatedAt,
			&blog.PublishedAt,
		)
		if err != nil {
			return nil, err
		}

		if categoryID.Valid {
			id := int(categoryID.Int64)
			blog.CategoryID = &id
		}

		blogs = append(blogs, blog)
	}

	return blogs, nil
}

// GetByAuthor retrieves blogs by author
func (r *blogRepository) GetByAuthor(ctx context.Context, author string) ([]models.Blog, error) {
	query := `SELECT id, title, content, author, category_id, status, views, created_at, updated_at, published_at
			  FROM blogs WHERE author ILIKE $1 ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, "%"+author+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blogs []models.Blog
	for rows.Next() {
		var blog models.Blog
		var categoryID sql.NullInt64

		err := rows.Scan(
			&blog.ID,
			&blog.Title,
			&blog.Content,
			&blog.Author,
			&categoryID,
			&blog.Status,
			&blog.Views,
			&blog.CreatedAt,
			&blog.UpdatedAt,
			&blog.PublishedAt,
		)
		if err != nil {
			return nil, err
		}

		if categoryID.Valid {
			id := int(categoryID.Int64)
			blog.CategoryID = &id
		}

		blogs = append(blogs, blog)
	}

	return blogs, nil
}

// Create inserts a new blog
func (r *blogRepository) Create(ctx context.Context, blog *models.Blog) error {
	query := `INSERT INTO blogs (title, content, author, category_id, status, created_at, updated_at, published_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, views`

	now := time.Now()
	blog.CreatedAt = now
	blog.UpdatedAt = now
	blog.Views = 0

	if blog.Status == models.StatusPublished && blog.PublishedAt == nil {
		blog.PublishedAt = &now
	}

	err := r.db.QueryRowContext(ctx, query,
		blog.Title,
		blog.Content,
		blog.Author,
		blog.CategoryID,
		blog.Status,
		blog.CreatedAt,
		blog.UpdatedAt,
		blog.PublishedAt,
	).Scan(&blog.ID, &blog.Views)

	return err
}

// Update modifies an existing blog
func (r *blogRepository) Update(ctx context.Context, blog *models.Blog) error {
	blog.UpdatedAt = time.Now()

	if blog.Status == models.StatusPublished && blog.PublishedAt == nil {
		now := time.Now()
		blog.PublishedAt = &now
	}

	query := `UPDATE blogs
			  SET title = $1, content = $2, author = $3, category_id = $4,
				  status = $5, updated_at = $6, published_at = $7
			  WHERE id = $8`

	result, err := r.db.ExecContext(ctx, query,
		blog.Title,
		blog.Content,
		blog.Author,
		blog.CategoryID,
		blog.Status,
		blog.UpdatedAt,
		blog.PublishedAt,
		blog.ID,
	)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// Delete removes a blog
func (r *blogRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM blogs WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// IncrementViews increases the view count
func (r *blogRepository) IncrementViews(ctx context.Context, id int) error {
	query := `UPDATE blogs SET views = views + 1 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// Search searches blogs by title or content
func (r *blogRepository) Search(ctx context.Context, searchTerm string) ([]models.Blog, error) {
	query := `SELECT id, title, content, author, category_id, status, views, created_at, updated_at, published_at
			  FROM blogs
			  WHERE title ILIKE $1 OR content ILIKE $1
			  ORDER BY created_at DESC`

	pattern := "%" + searchTerm + "%"
	rows, err := r.db.QueryContext(ctx, query, pattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blogs []models.Blog
	for rows.Next() {
		var blog models.Blog
		var categoryID sql.NullInt64

		err := rows.Scan(
			&blog.ID,
			&blog.Title,
			&blog.Content,
			&blog.Author,
			&categoryID,
			&blog.Status,
			&blog.Views,
			&blog.CreatedAt,
			&blog.UpdatedAt,
			&blog.PublishedAt,
		)
		if err != nil {
			return nil, err
		}

		if categoryID.Valid {
			id := int(categoryID.Int64)
			blog.CategoryID = &id
		}

		blogs = append(blogs, blog)
	}

	return blogs, nil
}
