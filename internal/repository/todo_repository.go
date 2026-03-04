package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/manish-npx/todo-go-echo/internal/models"
	"gorm.io/gorm"
)

type TodoRepository interface {
	GetAll() ([]models.Todo, error)
	GetByID(id int) (*models.Todo, error)
	Create(todo *models.Todo) error
	Update(todo *models.Todo) error
	Delete(id int) error
}

type todoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) TodoRepository {
	return &todoRepository{db: db}
}

func (r *todoRepository) GetAll() ([]models.Todo, error) {
	var todos []models.Todo
	err := r.db.Order("created_at DESC").Find(&todos).Error
	if err != nil {
		return nil, err
	}
	return todos, nil
}

func (r *todoRepository) GetByID(id int) (*models.Todo, error) {
	var todo models.Todo
	err := r.db.Where("id = ?", id).First(&todo).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

func (r *todoRepository) Create(todo *models.Todo) error {
	now := time.Now()
	todo.CreatedAt = now
	todo.UpdatedAt = now
	todo.Completed = false
	return r.db.Create(todo).Error
}

func (r *todoRepository) Update(todo *models.Todo) error {
	todo.UpdatedAt = time.Now()
	result := r.db.Model(&models.Todo{}).
		Where("id = ?", todo.ID).
		Updates(map[string]any{
			"title":       todo.Title,
			"description": todo.Description,
			"completed":   todo.Completed,
			"updated_at":  todo.UpdatedAt,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *todoRepository) Delete(id int) error {
	result := r.db.Where("id = ?", id).Delete(&models.Todo{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
