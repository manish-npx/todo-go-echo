package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/manish-npx/todo-go-echo/internal/config"
	"github.com/manish-npx/todo-go-echo/internal/database"
	"github.com/manish-npx/todo-go-echo/internal/handlers"
	"github.com/manish-npx/todo-go-echo/internal/logger"
	"github.com/manish-npx/todo-go-echo/internal/middleware"
	"github.com/manish-npx/todo-go-echo/internal/repository"
	"github.com/manish-npx/todo-go-echo/internal/routes"
	"github.com/manish-npx/todo-go-echo/internal/service"
	"github.com/manish-npx/todo-go-echo/internal/validator"
	"gorm.io/gorm"
)

// App keeps runtime dependencies in one place.
type App struct {
	Config *config.Config
	Echo   *echo.Echo
	GormDB *gorm.DB
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

	gormDB, err := database.NewGormConnection(cfg.Database, cfg.ORM)
	if err != nil {
		return nil, fmt.Errorf("gorm bootstrap failed: %w", err)
	}
	todoRepo := repository.NewTodoRepository(gormDB)
	categoryRepo := repository.NewCategoryRepository(gormDB)
	blogRepo := repository.NewBlogRepository(gormDB)
	userRepo := repository.NewUserRepository(gormDB)

	todoService := service.NewTodoService(todoRepo)
	todoHandler := handlers.NewTodoHandler(todoService)

	categoryService := service.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	blogService := service.NewBlogService(blogRepo, categoryRepo)
	blogHandler := handlers.NewBlogHandler(blogService)

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
		GormDB: gormDB,
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
	if a.GormDB != nil {
		sqlDB, err := a.GormDB.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
	}

	logger.Sync()
}
