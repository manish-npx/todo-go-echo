package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/manish-npx/todo-go-echo/internal/config"
	"github.com/manish-npx/todo-go-echo/internal/database"
	"github.com/manish-npx/todo-go-echo/internal/handlers"
	"github.com/manish-npx/todo-go-echo/internal/repository"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Connect to database
	db, err := database.NewPostgresConnection(cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize repository and handler
	todoRepo := repository.NewTodoRepository(db)
	todoHandler := handlers.NewTodoHandler(todoRepo)

	// Create Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Routes
	api := e.Group("/api/v1")
	{
		todos := api.Group("/todos")
		todos.GET("", todoHandler.GetTodos)
		todos.POST("", todoHandler.CreateTodo)
		todos.GET("/:id", todoHandler.GetTodo)
		todos.PUT("/:id", todoHandler.UpdateTodo)
		todos.DELETE("/:id", todoHandler.DeleteTodo)
	}

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.String(200, "OK")
	})

	// Start server
	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := e.Start(":" + cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
