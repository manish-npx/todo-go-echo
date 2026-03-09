package models

import "time"

// User represents database table structure
type User struct {
	ID        int       `json:"id"`                                       // primary key
	Name      string    `json:"name"`                                     // user name
	Email     string    `json:"email" gorm:"uniqueIndex:uni_users_email"` // unique email
	Mobile    string    `json:"mobile"`                                   // phone number
	Password  string    `json:"-" gorm:"column:password_hash"`            // hide password and map to DB hash column
	CreatedAt time.Time `json:"created_at"`                               // auto timestamp
	UpdatedAt time.Time `json:"updated_at"`                               // auto timestamp
}

// Request struct for registration
type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Mobile   string `json:"mobile" validate:"required,min=8,max=20"`
	Password string `json:"password" validate:"required,min=6,max=100"`
}

// Request struct for protected user creation
type CreateUserRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Mobile   string `json:"mobile" validate:"required,min=8,max=20"`
	Password string `json:"password" validate:"required,min=6,max=100"`
}

// Request struct for login
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=100"`
}
