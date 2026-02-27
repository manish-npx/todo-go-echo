package repository

import (
	"database/sql"

	"github.com/manish-npx/todo-go-echo/internal/models"
)

// UserRepository defines database operations
// This is just a CONTRACT (interface)
type UserRepository interface {
	Create(user *models.User) error
	GetByEmail(email string) (*models.User, error)
	GetById(id int) (*models.User, error)
	GetAll() ([]models.User, error)
}

// actual struct that talks to database
type userRepository struct {
	db *sql.DB // database connection
}

// constructor function (Dependency Injection)
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

// Create inserts user into database
func (r *userRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (name, email, mobile, password)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.Exec(query,
		user.Name,
		user.Email,
		user.Mobile,
		user.Password,
	)
	return err
}

// GetByEmail fetches user by email
func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, name, email, mobile, password, created_at
		FROM users WHERE email=$1
	`

	var user models.User
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Mobile,
		&user.Password,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetByID(id int64) (*models.User, error) {
	var u models.User

	err := r.db.QueryRow(
		`SELECT id,name,email,created_at FROM users WHERE id=$1`,
		id,
	).Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt)

	if err != nil {
		return nil, err
	}
	return &u, nil
}

// GetAll returns all users
func (r *userRepository) GetAll() ([]models.User, error) {
	rows, err := r.db.Query("SELECT id, name, email, mobile, created_at FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Mobile,
			&user.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
