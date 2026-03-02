package database

import (
	"fmt"

	"github.com/manish-npx/todo-go-echo/internal/config"
	"github.com/manish-npx/todo-go-echo/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewGormConnection creates an optional GORM connection.
// The current app still uses sql repositories for runtime CRUD paths.
func NewGormConnection(cfg config.DatabaseConfig, ormCfg config.ORMConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open gorm connection: %w", err)
	}

	if ormCfg.AutoMigrate {
		if err := db.AutoMigrate(&models.User{}, &models.Category{}, &models.Blog{}, &models.Todo{}); err != nil {
			return nil, fmt.Errorf("failed gorm automigrate: %w", err)
		}
	}

	return db, nil
}
