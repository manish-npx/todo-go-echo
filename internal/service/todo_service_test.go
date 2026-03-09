package service

import (
	"database/sql"
	"testing"

	"github.com/manish-npx/todo-go-echo/internal/models"
)

type todoRepoMock struct {
	todos map[int]*models.Todo
}

func (m *todoRepoMock) GetAll() ([]models.Todo, error) {
	result := make([]models.Todo, 0, len(m.todos))
	for _, todo := range m.todos {
		result = append(result, *todo)
	}
	return result, nil
}

func (m *todoRepoMock) GetByID(id int) (*models.Todo, error) {
	todo, ok := m.todos[id]
	if !ok {
		return nil, nil
	}
	return todo, nil
}

func (m *todoRepoMock) Create(todo *models.Todo) error {
	nextID := len(m.todos) + 1
	todo.ID = nextID
	m.todos[nextID] = todo
	return nil
}

func (m *todoRepoMock) Update(todo *models.Todo) error {
	if _, ok := m.todos[todo.ID]; !ok {
		return sql.ErrNoRows
	}
	m.todos[todo.ID] = todo
	return nil
}

func (m *todoRepoMock) Delete(id int) error {
	if _, ok := m.todos[id]; !ok {
		return sql.ErrNoRows
	}
	delete(m.todos, id)
	return nil
}

func TestTodoServiceCreate(t *testing.T) {
	repo := &todoRepoMock{todos: map[int]*models.Todo{}}
	svc := NewTodoService(repo)

	req := models.CreateTodoRequest{
		Title:       "Write tests",
		Description: "Add service unit tests",
	}

	todo, err := svc.Create(req)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if todo == nil {
		t.Fatal("Create() returned nil todo")
	}
	if todo.ID == 0 {
		t.Fatalf("Create() expected non-zero ID, got %d", todo.ID)
	}
	if todo.Title != req.Title || todo.Description != req.Description {
		t.Fatalf("Create() returned unexpected todo: %+v", todo)
	}
}

func TestTodoServiceUpdate(t *testing.T) {
	repo := &todoRepoMock{
		todos: map[int]*models.Todo{
			1: {ID: 1, Title: "Old", Description: "Old desc", Completed: false},
		},
	}
	svc := NewTodoService(repo)

	newTitle := "New title"
	completed := true
	req := models.UpdateTodoRequest{
		Title:     &newTitle,
		Completed: &completed,
	}

	todo, err := svc.Update(1, req)
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if todo == nil {
		t.Fatal("Update() returned nil todo")
	}
	if todo.Title != newTitle {
		t.Fatalf("Update() expected title %q, got %q", newTitle, todo.Title)
	}
	if !todo.Completed {
		t.Fatal("Update() expected completed=true")
	}
}
