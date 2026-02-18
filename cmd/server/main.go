package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/manish-npx/todo-go-echo/internal/config"
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

    // Initialize repositories
    todoRepo := repository.NewTodoRepository(db)
    categoryRepo := repository.NewCategoryRepository(db)
    blogRepo := repository.NewBlogRepository(db)

    // Initialize handlers
    todoHandler := handlers.NewTodoHandler(todoRepo)
    categoryHandler := handlers.NewCategoryHandler(categoryRepo)
    blogHandler := handlers.NewBlogHandler(blogRepo, categoryRepo)

    // Create Echo instance
    e := echo.New()

    // Middleware
	// Global middleware
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())
    e.Use(middleware.CORS())
    e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
        Timeout: 30 * time.Second,
    }))
    e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))


    // API Routes
    api := e.Group("/api/v1")

    // Todo routes (existing)
    todos := api.Group("/todos")
    todos.GET("", todoHandler.GetTodos)
    todos.POST("", todoHandler.CreateTodo)
    todos.GET("/:id", todoHandler.GetTodo)
    todos.PUT("/:id", todoHandler.UpdateTodo)
    todos.DELETE("/:id", todoHandler.DeleteTodo)

    // 🔴 NEW: Category routes
    categories := api.Group("/categories")
    categories.GET("", categoryHandler.GetCategories)
    categories.POST("", categoryHandler.CreateCategory)
    categories.GET("/:id", categoryHandler.GetCategory)
    categories.PUT("/:id", categoryHandler.UpdateCategory)
    categories.DELETE("/:id", categoryHandler.DeleteCategory)

    // 🔴 NEW: Blog routes
    blogs := api.Group("/blogs")
    blogs.GET("", blogHandler.GetBlogs)
    blogs.POST("", blogHandler.CreateBlog)
    blogs.GET("/search", blogHandler.SearchBlogs)
    blogs.GET("/:id", blogHandler.GetBlog)
    blogs.PUT("/:id", blogHandler.UpdateBlog)
    blogs.DELETE("/:id", blogHandler.DeleteBlog)
    blogs.PATCH("/:id/publish", blogHandler.PublishBlog)

    // Health check
