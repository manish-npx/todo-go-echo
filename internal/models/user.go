package models

import "time"

// User represents database table structure
type User struct {
	ID        int       `json:"id"`         // primary key
	Name      string    `json:"name"`       // user name
	Email     string    `json:"email"`      // unique email
	Mobile    string    `json:"mobile"`     // phone number
	Password  string    `json:"-"`          // hide password in API response
	CreatedAt time.Time `json:"created_at"` // auto timestamp
}

// Request struct for registration
type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Mobile   string `json:"mobile"`
	Password string `json:"password"`
}

// Request struct for login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
