package repository

import (
	"context"
	"database/sql"
	"todo-go-echo/internal/models"
)

// CategoryRepository defines the interface for category operations
type CategoryRepository interface {
	GetAll() ([]models.Category, error)
	GetByID(id int) (*models.Category, error)
	Create(category *models.Category) error
	Update(category *models.Category) error
	Delete(id int) error
	GetBlogCount(categoryID int) (int, error)
}

type categoryRepository struct {
	db *sql.DB
}

// NewCategoryRepository creates a new category repository
func NewCategoryRepository(db *sql.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

// GetAll retrieves all categories
func (r *categoryRepository) GetAll(ctx context.Context) ([]models.Category, error) {
	query := `SELECT id, name, description, created_at FROM categories ORDER BY name`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var cat models.Category
		err := rows.Scan(&cat.ID, &cat.Name, &cat.Description, &cat.CreatedAt)
		if err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}
	return categories, nil
}

// GetByID retrieves a category by ID
func (r *categoryRepository) GetByID(id int) (*models.Category, error) {
	query := `SELECT id, name, description, created_at FROM categories WHERE id = $1`

	var cat models.Category
	err := r.db.QueryRow(query, id).Scan(&cat.ID, &cat.Name, &cat.Description, &cat.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &cat, nil
}

// Create inserts a new category
func (r *categoryRepository) Create(category *models.Category) error {
	query := `INSERT INTO categories (name, description) VALUES ($1, $2) RETURNING id, created_at`
	return r.db.QueryRow(query, category.Name, category.Description).Scan(&category.ID, &category.CreatedAt)
}

// Update modifies an existing category
func (r *categoryRepository) Update(category *models.Category) error {
	query := `UPDATE categories SET name = $1, description = $2 WHERE id = $3`
	result, err := r.db.Exec(query, category.Name, category.Description, category.ID)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// Delete removes a category
func (r *categoryRepository) Delete(id int) error {
	query := `DELETE FROM categories WHERE id = $1`
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

// GetBlogCount returns number of blogs in a category
func (r *categoryRepository) GetBlogCount(categoryID int) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM blogs WHERE category_id = $1`
	err := r.db.QueryRow(query, categoryID).Scan(&count)
	return count, err
}
