package service

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lib/pq"
	"github.com/manish-npx/todo-go-echo/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type userRepoMock struct {
	users     map[string]*models.User
	createErr error
}

func (m *userRepoMock) Create(user *models.User) error {
	if m.createErr != nil {
		return m.createErr
	}
	if m.users == nil {
		m.users = map[string]*models.User{}
	}
	user.ID = len(m.users) + 1
	m.users[user.Email] = user
	return nil
}

func (m *userRepoMock) GetByEmail(email string) (*models.User, error) {
	user, ok := m.users[email]
	if !ok {
		return nil, sql.ErrNoRows
	}
	return user, nil
}

func (m *userRepoMock) GetByID(id int) (*models.User, error) {
	for _, user := range m.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (m *userRepoMock) GetAll() ([]models.User, error) {
	users := make([]models.User, 0, len(m.users))
	for _, user := range m.users {
		users = append(users, *user)
	}
	return users, nil
}

func TestUserServiceRegisterHashesPassword(t *testing.T) {
	repo := &userRepoMock{users: map[string]*models.User{}}
	svc := NewUserService(repo, "test-secret")

	req := models.RegisterRequest{
		Name:     "Manish",
		Email:    "manish@example.com",
		Mobile:   "9999999999",
		Password: "plain-password",
	}

	if err := svc.Register(req); err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	created := repo.users[req.Email]
	if created == nil {
		t.Fatal("Register() did not persist user")
	}
	if created.Password == req.Password {
		t.Fatal("Register() stored plain text password")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(created.Password), []byte(req.Password)); err != nil {
		t.Fatalf("stored hash does not match password: %v", err)
	}
}

func TestUserServiceRegisterDuplicateEmail(t *testing.T) {
	repo := &userRepoMock{
		createErr: &pq.Error{Code: "23505", Constraint: "uni_users_email"},
	}
	svc := NewUserService(repo, "test-secret")

	err := svc.Register(models.RegisterRequest{
		Name:     "Manish",
		Email:    "manish@example.com",
		Mobile:   "9999999999",
		Password: "plain-password",
	})

	if err == nil {
		t.Fatal("expected duplicate email error, got nil")
	}
	if !errors.Is(err, ErrEmailAlreadyExists) {
		t.Fatalf("expected ErrEmailAlreadyExists, got: %v", err)
	}
}

func TestUserServiceLoginReturnsJWT(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("secret-pass"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to build test hash: %v", err)
	}

	repo := &userRepoMock{
		users: map[string]*models.User{
			"manish@example.com": {
				ID:       7,
				Email:    "manish@example.com",
				Password: string(hash),
			},
		},
	}
	const signingSecret = "test-secret"
	svc := NewUserService(repo, signingSecret)

	tokenString, err := svc.Login(models.LoginRequest{
		Email:    "manish@example.com",
		Password: "secret-pass",
	})
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}
	if tokenString == "" {
		t.Fatal("Login() returned empty token")
	}

	parsed, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(signingSecret), nil
	})
	if err != nil {
		t.Fatalf("failed parsing JWT: %v", err)
	}
	if !parsed.Valid {
		t.Fatal("JWT token is not valid")
	}
}

func TestUserServiceLoginInvalidPassword(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("right-pass"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to build test hash: %v", err)
	}

	repo := &userRepoMock{
		users: map[string]*models.User{
			"manish@example.com": {
				ID:       10,
				Email:    "manish@example.com",
				Password: string(hash),
			},
		},
	}
	svc := NewUserService(repo, "test-secret")

	_, err = svc.Login(models.LoginRequest{
		Email:    "manish@example.com",
		Password: "wrong-pass",
	})
	if err == nil || err.Error() != "invalid credentials" {
		t.Fatalf("expected invalid credentials error, got: %v", err)
	}
}
