package utils

import "errors"

// Custom error types
var (
	ErrNotFound     = errors.New("resource not found")
	ErrInvalidInput = errors.New("invalid input")
	ErrDuplicate    = errors.New("resource already exists")
	ErrUnauthorized = errors.New("unauthorized")
	ErrDatabase     = errors.New("database error")
)

// AppError represents a structured application error
type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"-"`
}

func (e AppError) Error() string {
	return e.Message
}

// Common application errors
var (
	ErrBlogNotFound = AppError{
		Code:    "BLOG_NOT_FOUND",
		Message: "Blog post not found",
		Status:  404,
	}

	ErrCategoryNotFound = AppError{
		Code:    "CATEGORY_NOT_FOUND",
		Message: "Category not found",
		Status:  404,
	}

	ErrInvalidBlogData = AppError{
		Code:    "INVALID_BLOG_DATA",
		Message: "Invalid blog data provided",
		Status:  400,
	}

	ErrCategoryInUse = AppError{
		Code:    "CATEGORY_IN_USE",
		Message: "Category has blogs and cannot be deleted",
		Status:  400,
	}
)
