package database

import (
	"database/sql"
	"fmt"

	"github.com/manish-npx/todo-go-echo/internal/config"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// NewPostgresConnection establishes a connection to PostgreSQL
func NewPostgresConnection(cfg config.DatabaseConfig) (*sql.DB, error) {
	// Create connection string
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	// Open database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	return db, nil
}
