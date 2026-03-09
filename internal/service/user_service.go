package service

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lib/pq"
	"github.com/manish-npx/todo-go-echo/internal/models"
	"github.com/manish-npx/todo-go-echo/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var ErrEmailAlreadyExists = errors.New("email already exists")
var ErrUserNotFound = errors.New("user not found")

// Service CONTRACT
type UserService interface {
	Register(req models.RegisterRequest) error
	CreateUser(req models.CreateUserRequest) error
	Login(req models.LoginRequest) (string, error)
	GetUsers() ([]models.User, error)
	GetProfile(userID int) (*models.User, error)
}

// actual service struct
type userService struct {
	repo      repository.UserRepository // dependency
	jwtSecret string                    // secret key for JWT
}

// constructor
func NewUserService(repo repository.UserRepository, secret string) UserService {
	return &userService{
		repo:      repo,
		jwtSecret: secret,
	}
}

// Register handles password hashing + save to DB
func (s *userService) Register(req models.RegisterRequest) error {
	return s.createUser(req.Name, req.Email, req.Mobile, req.Password)
}

// CreateUser is a protected API path to create users.
func (s *userService) CreateUser(req models.CreateUserRequest) error {
	return s.createUser(req.Name, req.Email, req.Mobile, req.Password)
}

func (s *userService) createUser(name, email, mobile, plainPassword string) error {
	// hash password before storing
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(plainPassword),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	user := models.User{
		Name:     name,
		Email:    email,
		Mobile:   mobile,
		Password: string(hashedPassword),
	}

	if err := s.repo.Create(&user); err != nil {
		if isDuplicateEmailError(err) {
			return ErrEmailAlreadyExists
		}
		return err
	}

	return nil
}

// Login verifies password and generates JWT
func (s *userService) Login(req models.LoginRequest) (string, error) {

	user, err := s.repo.GetByEmail(req.Email)
	if err != nil {
		return "", ErrUserNotFound
	}

	// compare hashed password
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.Password),
	)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	return token.SignedString([]byte(s.jwtSecret))
}

// GetUsers returns all users
func (s *userService) GetUsers() ([]models.User, error) {
	return s.repo.GetAll()
}

// GetProfile fetches currently authenticated user by ID.
func (s *userService) GetProfile(userID int) (*models.User, error) {
	user, err := s.repo.GetByID(userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func isDuplicateEmailError(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return string(pqErr.Code) == "23505" &&
			(pqErr.Constraint == "uni_users_email" || pqErr.Constraint == "users_email_key")
	}

	var pgxErr *pgconn.PgError
	if errors.As(err, &pgxErr) {
		return pgxErr.Code == "23505" &&
			(pgxErr.ConstraintName == "uni_users_email" || pgxErr.ConstraintName == "users_email_key")
	}

	// Fallback for wrapped/driver-agnostic errors.
	return strings.Contains(err.Error(), "SQLSTATE 23505") &&
		strings.Contains(err.Error(), "users_email")
}
