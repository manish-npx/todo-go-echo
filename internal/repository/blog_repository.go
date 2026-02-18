package repository

import (
	"database/sql"
	"time"
	"todo-go-echo/internal/models"
)

// BlogRepository defines the interface for blog operations
type BlogRepository interface {
	GetAll() ([]models.Blog, error)
	GetByID(id int) (*models.Blog, error)
	GetByCategory(categoryID int) ([]models.Blog, error)
	GetPublished() ([]models.Blog, error)
	GetByAuthor(author string) ([]models.Blog, error)
	Create(blog *models.Blog) error
	Update(blog *models.Blog) error
	Delete(id int) error
	IncrementViews(id int) error
	Search(query string) ([]models.Blog, error)
}

type blogRepository struct {
	db *sql.DB
}

// NewBlogRepository creates a new blog repository
func NewBlogRepository(db *sql.DB) BlogRepository {
	return &blogRepository{db: db}
}

// GetAll retrieves all blogs with their categories
func (r *blogRepository) GetAll() ([]models.Blog, error) {
	query := `
        SELECT
            b.id, b.title, b.content, b.author, b.category_id,
            b.status, b.views, b.created_at, b.updated_at, b.published_at,
            c.id, c.name, c.description, c.created_at
        FROM blogs b
        LEFT JOIN categories c ON b.category_id = c.id
        ORDER BY b.created_at DESC
    `

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blogs []models.Blog
	for rows.Next() {
		blog, err := r.scanBlogWithCategory(rows)
		if err != nil {
			return nil, err
		}
		blogs = append(blogs, *blog)
	}
	return blogs, nil
}

// GetByID retrieves a blog by ID with its category
func (r *blogRepository) GetByID(id int) (*models.Blog, error) {
	query := `
        SELECT
            b.id, b.title, b.content, b.author, b.category_id,
            b.status, b.views, b.created_at, b.updated_at, b.published_at,
            c.id, c.name, c.description, c.created_at
        FROM blogs b
        LEFT JOIN categories c ON b.category_id = c.id
        WHERE b.id = $1
    `

	row := r.db.QueryRow(query, id)
	return r.scanBlogWithCategory(row)
}

// GetByCategory retrieves blogs in a specific category
func (r *blogRepository) GetByCategory(categoryID int) ([]models.Blog, error) {
	query := `
        SELECT
            b.id, b.title, b.content, b.author, b.category_id,
            b.status, b.views, b.created_at, b.updated_at, b.published_at,
            c.id, c.name, c.description, c.created_at
        FROM blogs b
        LEFT JOIN categories c ON b.category_id = c.id
        WHERE b.category_id = $1
        ORDER BY b.created_at DESC
    `

	rows, err := r.db.Query(query, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blogs []models.Blog
	for rows.Next() {
		blog, err := r.scanBlogWithCategory(rows)
		if err != nil {
			return nil, err
		}
		blogs = append(blogs, *blog)
	}
	return blogs, nil
}

// GetPublished retrieves only published blogs
func (r *blogRepository) GetPublished() ([]models.Blog, error) {
	query := `
        SELECT
            b.id, b.title, b.content, b.author, b.category_id,
            b.status, b.views, b.created_at, b.updated_at, b.published_at,
            c.id, c.name, c.description, c.created_at
        FROM blogs b
        LEFT JOIN categories c ON b.category_id = c.id
        WHERE b.status = 'published'
        ORDER BY b.published_at DESC
    `

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blogs []models.Blog
	for rows.Next() {
		blog, err := r.scanBlogWithCategory(rows)
		if err != nil {
			return nil, err
		}
		blogs = append(blogs, *blog)
	}
	return blogs, nil
}

// GetByAuthor retrieves blogs by author
func (r *blogRepository) GetByAuthor(author string) ([]models.Blog, error) {
	query := `
        SELECT
            b.id, b.title, b.content, b.author, b.category_id,
            b.status, b.views, b.created_at, b.updated_at, b.published_at,
            c.id, c.name, c.description, c.created_at
        FROM blogs b
        LEFT JOIN categories c ON b.category_id = c.id
        WHERE b.author ILIKE $1
        ORDER BY b.created_at DESC
    `

	rows, err := r.db.Query(query, "%"+author+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blogs []models.Blog
	for rows.Next() {
		blog, err := r.scanBlogWithCategory(rows)
		if err != nil {
			return nil, err
		}
		blogs = append(blogs, *blog)
	}
	return blogs, nil
}

// Create inserts a new blog
func (r *blogRepository) Create(blog *models.Blog) error {
	query := `
        INSERT INTO blogs (title, content, author, category_id, status, created_at, updated_at, published_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING id, views
    `

	now := time.Now()
	blog.CreatedAt = now
	blog.UpdatedAt = now
	blog.Views = 0

	// Set published_at if status is published
	if blog.Status == models.StatusPublished {
		blog.PublishedAt = &now
	}

	return r.db.QueryRow(query,
		blog.Title,
		blog.Content,
		blog.Author,
		blog.CategoryID,
		blog.Status,
		blog.CreatedAt,
		blog.UpdatedAt,
		blog.PublishedAt,
	).Scan(&blog.ID, &blog.Views)
}

// Update modifies an existing blog
func (r *blogRepository) Update(blog *models.Blog) error {
	// Always update updated_at
	blog.UpdatedAt = time.Now()

	// If status changed to published and wasn't published before, set published_at
	if blog.Status == models.StatusPublished && blog.PublishedAt == nil {
		now := time.Now()
		blog.PublishedAt = &now
	}

	query := `
        UPDATE blogs
        SET title = $1, content = $2, author = $3, category_id = $4,
            status = $5, updated_at = $6, published_at = $7
        WHERE id = $8
    `

	result, err := r.db.Exec(query,
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

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// Delete removes a blog
func (r *blogRepository) Delete(id int) error {
	query := `DELETE FROM blogs WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// IncrementViews increases the view count
func (r *blogRepository) IncrementViews(id int) error {
	query := `UPDATE blogs SET views = views + 1 WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// Search searches blogs by title or content
func (r *blogRepository) Search(searchTerm string) ([]models.Blog, error) {
	query := `
        SELECT
            b.id, b.title, b.content, b.author, b.category_id,
            b.status, b.views, b.created_at, b.updated_at, b.published_at,
            c.id, c.name, c.description, c.created_at
        FROM blogs b
        LEFT JOIN categories c ON b.category_id = c.id
        WHERE b.title ILIKE $1 OR b.content ILIKE $1
        ORDER BY b.created_at DESC
    `

	pattern := "%" + searchTerm + "%"
	rows, err := r.db.Query(query, pattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blogs []models.Blog
	for rows.Next() {
		blog, err := r.scanBlogWithCategory(rows)
		if err != nil {
			return nil, err
		}
		blogs = append(blogs, *blog)
	}
	return blogs, nil
}

// Helper function to scan a blog row with its category
func (r *blogRepository) scanBlogWithCategory(row interface {
	Scan(dest ...interface{}) error
}) (*models.Blog, error) {
	var blog models.Blog
	var category models.Category

	// Variables for category fields (might be NULL)
	var categoryID sql.NullInt64
	var categoryName, categoryDescription sql.NullString
	var categoryCreatedAt sql.NullTime

	err := row.Scan(
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
		&category.ID,
		&categoryName,
		&categoryDescription,
		&categoryCreatedAt,
	)

	if err != nil {
		return nil, err
	}

	// If category exists, populate it
	if categoryID.Valid {
		blog.CategoryID = new(int)
		*blog.CategoryID = int(categoryID.Int64)

		category.ID = int(categoryID.Int64)
		category.Name = categoryName.String
		category.Description = categoryDescription.String
		category.CreatedAt = categoryCreatedAt.Time
		blog.Category = &category
	}

	return &blog, nil
}
