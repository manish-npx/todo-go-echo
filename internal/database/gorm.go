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
		// Keep ORM migration flow simple and predictable.
		if err := normalizeUsersEmailUniqueConstraint(db); err != nil {
			return nil, fmt.Errorf("failed normalizing users email constraint: %w", err)
		}
		if err := normalizeCategoriesNameUniqueConstraint(db); err != nil {
			return nil, fmt.Errorf("failed normalizing categories name constraint: %w", err)
		}

		if err := db.AutoMigrate(&models.User{}, &models.Category{}, &models.Blog{}, &models.Todo{}); err != nil {
			return nil, fmt.Errorf("failed gorm automigrate: %w", err)
		}
	}

	return db, nil
}

func normalizeUsersEmailUniqueConstraint(db *gorm.DB) error {
	const query = `
DO $$
BEGIN
	IF EXISTS (
		SELECT 1
		FROM pg_constraint
		WHERE conname = 'users_email_key'
		  AND conrelid = 'users'::regclass
	)
	AND NOT EXISTS (
		SELECT 1
		FROM pg_constraint
		WHERE conname = 'uni_users_email'
		  AND conrelid = 'users'::regclass
	) THEN
		ALTER TABLE users RENAME CONSTRAINT users_email_key TO uni_users_email;
	END IF;
END $$;
`

	return db.Exec(query).Error
}

func normalizeCategoriesNameUniqueConstraint(db *gorm.DB) error {
	const query = `
DO $$
BEGIN
	IF EXISTS (
		SELECT 1
		FROM pg_constraint
		WHERE conname = 'categories_name_key'
		  AND conrelid = 'categories'::regclass
	)
	AND NOT EXISTS (
		SELECT 1
		FROM pg_constraint
		WHERE conname = 'uni_categories_name'
		  AND conrelid = 'categories'::regclass
	) THEN
		ALTER TABLE categories RENAME CONSTRAINT categories_name_key TO uni_categories_name;
	END IF;
END $$;
`

	return db.Exec(query).Error
}
