package service

import (
	"context"
	"database/sql"

	"github.com/manish-npx/todo-go-echo/internal/models"
	"github.com/manish-npx/todo-go-echo/internal/repository"
)

type CategoryService interface {
	GetAll(ctx context.Context) ([]models.Category, error)
	GetByID(ctx context.Context, id int) (*models.Category, error)
	Create(ctx context.Context, req models.CreateCategoryRequest) (*models.Category, error)
	Update(ctx context.Context, id int, req models.UpdateCategoryRequest) (*models.Category, error)
	Delete(ctx context.Context, id int) error
}

type categoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) GetAll(ctx context.Context) ([]models.Category, error) {
	return s.repo.GetAll(ctx)
}

func (s *categoryService) GetByID(ctx context.Context, id int) (*models.Category, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *categoryService) Create(ctx context.Context, req models.CreateCategoryRequest) (*models.Category, error) {
	category := &models.Category{
		Name:        req.Name,
		Description: req.Description,
	}
	if err := s.repo.Create(ctx, category); err != nil {
		return nil, err
	}
	return category, nil
}

func (s *categoryService) Update(ctx context.Context, id int, req models.UpdateCategoryRequest) (*models.Category, error) {
	category, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, sql.ErrNoRows
	}

	if req.Name != nil {
		category.Name = *req.Name
	}
	if req.Description != nil {
		category.Description = *req.Description
	}

	if err := s.repo.Update(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *categoryService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
