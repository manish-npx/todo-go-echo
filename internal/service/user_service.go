package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/manish-npx/todo-go-echo/internal/models"
	"github.com/manish-npx/todo-go-echo/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// Service CONTRACT
type UserService interface {
	Register(req models.RegisterRequest) error
	Login(req models.LoginRequest) (string, error)
	GetUsers() ([]models.User, error)
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

	// hash password before storing
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Mobile:   req.Mobile,
		Password: string(hashedPassword),
	}

	return s.repo.Create(&user)
}

// Login verifies password and generates JWT
func (s *userService) Login(req models.LoginRequest) (string, error) {

	user, err := s.repo.GetByEmail(req.Email)
	if err != nil {
		return "", errors.New("user not found")
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
