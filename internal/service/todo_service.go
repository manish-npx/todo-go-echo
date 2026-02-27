package service

import (
	"github.com/manish-npx/todo-go-echo/internal/models"
	"github.com/manish-npx/todo-go-echo/internal/repository"
)

type TodoService interface {
	GetAll() ([]models.Todo, error)
	GetByID(id int) (*models.Todo, error)
	Create(req models.CreateTodoRequest) (*models.Todo, error)
	Update(id int, req models.UpdateTodoRequest) (*models.Todo, error)
	Delete(id int) error
}

type todoService struct {
	repo repository.TodoRepository
}

func NewTodoService(repo repository.TodoRepository) TodoService {
	return &todoService{repo: repo}
}

func (s *todoService) GetAll() ([]models.Todo, error) {
	return s.repo.GetAll()
}

func (s *todoService) GetByID(id int) (*models.Todo, error) {
	return s.repo.GetByID(id)
}

func (s *todoService) Create(req models.CreateTodoRequest) (*models.Todo, error) {

	todo := &models.Todo{
		Title:       req.Title,
		Description: req.Description,
	}

	err := s.repo.Create(todo)
	return todo, err
}

func (s *todoService) Update(id int, req models.UpdateTodoRequest) (*models.Todo, error) {

	todo, err := s.repo.GetByID(id)
	if err != nil || todo == nil {
		return nil, err
	}

	if req.Title != nil {
		todo.Title = *req.Title
	}
	if req.Description != nil {
		todo.Description = *req.Description
	}
	if req.Completed != nil {
		todo.Completed = *req.Completed
	}

	err = s.repo.Update(todo)
	return todo, err
}

func (s *todoService) Delete(id int) error {
	return s.repo.Delete(id)
}
