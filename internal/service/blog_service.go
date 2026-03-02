package service

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"github.com/manish-npx/todo-go-echo/internal/models"
	"github.com/manish-npx/todo-go-echo/internal/repository"
)

type BlogService interface {
	GetBlogs(ctx context.Context, categoryID, author, status string) ([]models.Blog, error)
	GetByID(ctx context.Context, id int) (*models.Blog, error)
	Create(ctx context.Context, req models.CreateBlogRequest) (*models.Blog, error)
	Update(ctx context.Context, id int, req models.UpdateBlogRequest) (*models.Blog, error)
	Delete(ctx context.Context, id int) error
	Search(ctx context.Context, query string) ([]models.Blog, error)
	Publish(ctx context.Context, id int) (*models.Blog, error)
}

type blogService struct {
	blogRepo     repository.BlogRepository
	categoryRepo repository.CategoryRepository
}

func NewBlogService(blogRepo repository.BlogRepository, categoryRepo repository.CategoryRepository) BlogService {
	return &blogService{
		blogRepo:     blogRepo,
		categoryRepo: categoryRepo,
	}
}

func (s *blogService) GetBlogs(ctx context.Context, categoryID, author, status string) ([]models.Blog, error) {
	switch {
	case categoryID != "":
		id, err := strconv.Atoi(categoryID)
		if err != nil {
			return nil, errors.New("invalid category id")
		}
		return s.blogRepo.GetByCategory(ctx, id)
	case author != "":
		return s.blogRepo.GetByAuthor(ctx, author)
	case status == "published":
		return s.blogRepo.GetPublished(ctx)
	default:
		return s.blogRepo.GetAll(ctx)
	}
}

func (s *blogService) GetByID(ctx context.Context, id int) (*models.Blog, error) {
	return s.blogRepo.GetByID(ctx, id)
}

func (s *blogService) Create(ctx context.Context, req models.CreateBlogRequest) (*models.Blog, error) {
	if req.CategoryID != nil {
		category, err := s.categoryRepo.GetByID(ctx, *req.CategoryID)
		if err != nil {
			return nil, err
		}
		if category == nil {
			return nil, sql.ErrNoRows
		}
	}

	status := models.StatusDraft
	if req.Status == string(models.StatusPublished) {
		status = models.StatusPublished
	}

	blog := &models.Blog{
		Title:      req.Title,
		Content:    req.Content,
		Author:     req.Author,
		CategoryID: req.CategoryID,
		Status:     status,
	}
	if err := s.blogRepo.Create(ctx, blog); err != nil {
		return nil, err
	}

	return blog, nil
}

func (s *blogService) Update(ctx context.Context, id int, req models.UpdateBlogRequest) (*models.Blog, error) {
	blog, err := s.blogRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if blog == nil {
		return nil, sql.ErrNoRows
	}

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
		if *req.CategoryID != 0 {
			category, err := s.categoryRepo.GetByID(ctx, *req.CategoryID)
			if err != nil {
				return nil, err
			}
			if category == nil {
				return nil, sql.ErrNoRows
			}
		}
		blog.CategoryID = req.CategoryID
	}
	if req.Status != nil {
		blog.Status = models.BlogStatus(*req.Status)
	}

	if err := s.blogRepo.Update(ctx, blog); err != nil {
		return nil, err
	}

	return blog, nil
}

func (s *blogService) Delete(ctx context.Context, id int) error {
	return s.blogRepo.Delete(ctx, id)
}

func (s *blogService) Search(ctx context.Context, query string) ([]models.Blog, error) {
	return s.blogRepo.Search(ctx, query)
}

func (s *blogService) Publish(ctx context.Context, id int) (*models.Blog, error) {
	blog, err := s.blogRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if blog == nil {
		return nil, sql.ErrNoRows
	}

	blog.Status = models.StatusPublished
	if err := s.blogRepo.Update(ctx, blog); err != nil {
		return nil, err
	}

	return blog, nil
}
