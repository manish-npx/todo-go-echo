package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/manish-npx/todo-go-echo/internal/models"
	"gorm.io/gorm"
)

type BlogRepository interface {
	GetAll(ctx context.Context) ([]models.Blog, error)
	GetByID(ctx context.Context, id int) (*models.Blog, error)
	GetByCategory(ctx context.Context, categoryID int) ([]models.Blog, error)
	GetPublished(ctx context.Context) ([]models.Blog, error)
	GetByAuthor(ctx context.Context, author string) ([]models.Blog, error)
	Create(ctx context.Context, blog *models.Blog) error
	Update(ctx context.Context, blog *models.Blog) error
	Delete(ctx context.Context, id int) error
	IncrementViews(ctx context.Context, id int) error
	Search(ctx context.Context, query string) ([]models.Blog, error)
}

type blogRepository struct {
	db *gorm.DB
}

func NewBlogRepository(db *gorm.DB) BlogRepository {
	return &blogRepository{db: db}
}

func (r *blogRepository) GetAll(ctx context.Context) ([]models.Blog, error) {
	var blogs []models.Blog
	err := r.db.WithContext(ctx).Order("created_at DESC").Find(&blogs).Error
	if err != nil {
		return nil, err
	}
	return blogs, nil
}

func (r *blogRepository) GetByID(ctx context.Context, id int) (*models.Blog, error) {
	var blog models.Blog
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&blog).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &blog, nil
}

func (r *blogRepository) GetByCategory(ctx context.Context, categoryID int) ([]models.Blog, error) {
	var blogs []models.Blog
	err := r.db.WithContext(ctx).
		Where("category_id = ?", categoryID).
		Order("created_at DESC").
		Find(&blogs).Error
	if err != nil {
		return nil, err
	}
	return blogs, nil
}

func (r *blogRepository) GetPublished(ctx context.Context) ([]models.Blog, error) {
	var blogs []models.Blog
	err := r.db.WithContext(ctx).
		Where("status = ?", models.StatusPublished).
		Order("published_at DESC").
		Find(&blogs).Error
	if err != nil {
		return nil, err
	}
	return blogs, nil
}

func (r *blogRepository) GetByAuthor(ctx context.Context, author string) ([]models.Blog, error) {
	var blogs []models.Blog
	err := r.db.WithContext(ctx).
		Where("author ILIKE ?", "%"+author+"%").
		Order("created_at DESC").
		Find(&blogs).Error
	if err != nil {
		return nil, err
	}
	return blogs, nil
}

func (r *blogRepository) Create(ctx context.Context, blog *models.Blog) error {
	now := time.Now()
	blog.CreatedAt = now
	blog.UpdatedAt = now
	blog.Views = 0
	if blog.Status == models.StatusPublished && blog.PublishedAt == nil {
		blog.PublishedAt = &now
	}
	return r.db.WithContext(ctx).Create(blog).Error
}

func (r *blogRepository) Update(ctx context.Context, blog *models.Blog) error {
	blog.UpdatedAt = time.Now()
	if blog.Status == models.StatusPublished && blog.PublishedAt == nil {
		now := time.Now()
		blog.PublishedAt = &now
	}

	result := r.db.WithContext(ctx).Model(&models.Blog{}).
		Where("id = ?", blog.ID).
		Updates(map[string]any{
			"title":        blog.Title,
			"content":      blog.Content,
			"author":       blog.Author,
			"category_id":  blog.CategoryID,
			"status":       blog.Status,
			"updated_at":   blog.UpdatedAt,
			"published_at": blog.PublishedAt,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *blogRepository) Delete(ctx context.Context, id int) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Blog{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *blogRepository) IncrementViews(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Model(&models.Blog{}).
		Where("id = ?", id).
		Update("views", gorm.Expr("views + 1")).Error
}

func (r *blogRepository) Search(ctx context.Context, searchTerm string) ([]models.Blog, error) {
	var blogs []models.Blog
	pattern := "%" + searchTerm + "%"
	err := r.db.WithContext(ctx).
		Where("title ILIKE ? OR content ILIKE ?", pattern, pattern).
		Order("created_at DESC").
		Find(&blogs).Error
	if err != nil {
		return nil, err
	}
	return blogs, nil
}
