package app

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/manish-npx/todo-go-echo/internal/config"
	"github.com/manish-npx/todo-go-echo/internal/database"
	"github.com/manish-npx/todo-go-echo/internal/handlers"
	"github.com/manish-npx/todo-go-echo/internal/logger"
	"github.com/manish-npx/todo-go-echo/internal/middleware"
	"github.com/manish-npx/todo-go-echo/internal/repository"
	"github.com/manish-npx/todo-go-echo/internal/routes"
	"github.com/manish-npx/todo-go-echo/internal/service"
	"github.com/manish-npx/todo-go-echo/internal/validator"
)

// App keeps runtime dependencies in one place.
type App struct {
	Config *config.Config
	Echo   *echo.Echo
	DB     *sql.DB
}

// New builds the full application graph (config, db, handlers, routes, middleware).
func New(configPath string) (*App, error) {
	if err := logger.Init(); err != nil {
		return nil, fmt.Errorf("logger init failed: %w", err)
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("config load failed: %w", err)
	}

	db, err := database.NewPostgresConnection(cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("database connection failed: %w", err)
	}

	// Optional ORM bootstrap. CRUD paths still use sql repositories.
	if cfg.ORM.Enabled {
		if _, err := database.NewGormConnection(cfg.Database, cfg.ORM); err != nil {
			_ = db.Close()
			return nil, fmt.Errorf("gorm bootstrap failed: %w", err)
		}
	}

	// Dependency injection chain
	todoRepo := repository.NewTodoRepository(db)
	todoService := service.NewTodoService(todoRepo)
	todoHandler := handlers.NewTodoHandler(todoService)

	categoryRepo := repository.NewCategoryRepository(db)
	categoryService := service.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	blogRepo := repository.NewBlogRepository(db)
	blogService := service.NewBlogService(blogRepo, categoryRepo)
	blogHandler := handlers.NewBlogHandler(blogService)

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, cfg.JWT.Secret)
	userHandler := handlers.NewUserHandler(userService)

	e := echo.New()
	e.Validator = validator.New()
	e.HTTPErrorHandler = middleware.ErrorHandler
	middleware.Setup(e)

	routes.RegisterRoutes(e, routes.RouteHandlers{
		TodoHandler:     todoHandler,
		CategoryHandler: categoryHandler,
		BlogHandler:     blogHandler,
		UserHandler:     userHandler,
		JWTSecret:       cfg.JWT.Secret,
	})

	e.Static("/", "dist")

	return &App{
		Config: cfg,
		Echo:   e,
		DB:     db,
	}, nil
}

// Run starts server with graceful shutdown.
func (a *App) Run() {
	go func() {
		if err := a.Echo.Start(":" + a.Config.Server.Port); err != nil && err != http.ErrServerClosed {
			a.Echo.Logger.Fatal("shutting down")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	shutdownTimeout := 10 * time.Second
	if a.Config.Server.Timeout > 0 {
		shutdownTimeout = time.Duration(a.Config.Server.Timeout) * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := a.Echo.Shutdown(ctx); err != nil {
		a.Echo.Logger.Fatal(err)
	}
}

// Close releases app resources.
func (a *App) Close() {
	if a.DB != nil {
		_ = a.DB.Close()
	}
	logger.Sync()
}
