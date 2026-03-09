package models

import "time"

// Category represents a blog category
type Category struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name" gorm:"uniqueIndex:uni_categories_name"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// CreateCategoryRequest is used when creating a category
type CreateCategoryRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=100"`
	Description string `json:"description" validate:"required,min=5,max=500"`
}

// UpdateCategoryRequest is used when updating a category
type UpdateCategoryRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}
