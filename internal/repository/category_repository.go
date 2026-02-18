package repository

import (
	"context"
	"database/sql"

	"github.com/manish-npx/todo-go-echo/internal/models"
)

// CategoryRepository defines the interface for category operations
type CategoryRepository interface {
	GetAll(ctx context.Context) ([]models.Category, error)
	GetByID(ctx context.Context, id int) (*models.Category, error)
	Create(ctx context.Context, category *models.Category) error
	Update(ctx context.Context, category *models.Category) error
	Delete(ctx context.Context, id int) error
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

	rows, err := r.db.QueryContext(ctx, query)
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

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

// GetByID retrieves a category by ID
func (r *categoryRepository) GetByID(ctx context.Context, id int) (*models.Category, error) {
	query := `SELECT id, name, description, created_at FROM categories WHERE id = $1`

	var cat models.Category
	err := r.db.QueryRowContext(ctx, query, id).Scan(&cat.ID, &cat.Name, &cat.Description, &cat.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &cat, nil
}

// Create inserts a new category
func (r *categoryRepository) Create(ctx context.Context, category *models.Category) error {
	query := `INSERT INTO categories (name, description) VALUES ($1, $2) RETURNING id, created_at`

	err := r.db.QueryRowContext(ctx, query, category.Name, category.Description).Scan(&category.ID, &category.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

// Update modifies an existing category
func (r *categoryRepository) Update(ctx context.Context, category *models.Category) error {
	query := `UPDATE categories SET name = $1, description = $2 WHERE id = $3`

	result, err := r.db.ExecContext(ctx, query, category.Name, category.Description, category.ID)
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

// Delete removes a category
func (r *categoryRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM categories WHERE id = $1`

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
