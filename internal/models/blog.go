package models

import "time"

// BlogStatus represents the possible states of a blog post
type BlogStatus string

const (
	StatusDraft     BlogStatus = "draft"
	StatusPublished BlogStatus = "published"
)

// Blog represents a blog post
type Blog struct {
	ID          int        `json:"id" db:"id"`
	Title       string     `json:"title" db:"title"`
	Content     string     `json:"content" db:"content"`
	Author      string     `json:"author" db:"author"`
	CategoryID  *int       `json:"category_id,omitempty" db:"category_id"`
	Category    *Category  `json:"category,omitempty"` // This will be populated when joining
	Status      BlogStatus `json:"status" db:"status"`
	Views       int        `json:"views" db:"views"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	PublishedAt *time.Time `json:"published_at,omitempty" db:"published_at"`
}

// CreateBlogRequest is used when creating a blog
type CreateBlogRequest struct {
	Title      string `json:"title" validate:"required,min=3,max=255"`
	Content    string `json:"content" validate:"required,min=10"`
	Author     string `json:"author" validate:"required,min=2"`
	CategoryID *int   `json:"category_id"`
	Status     string `json:"status"` // "draft" or "published"
}

// UpdateBlogRequest is used when updating a blog
type UpdateBlogRequest struct {
	Title      *string `json:"title" validate:"omitempty,min=3,max=255"`
	Content    *string `json:"content" validate:"omitempty,min=10"`
	Author     *string `json:"author" validate:"omitempty,min=2"`
	CategoryID *int    `json:"category_id"`
	Status     *string `json:"status"`
}
