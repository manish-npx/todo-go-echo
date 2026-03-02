package main

import (
	"log"

	"github.com/labstack/echo/v4"

	_ "github.com/lib/pq"
	"github.com/manish-npx/todo-go-echo/internal/config"
	"github.com/manish-npx/todo-go-echo/internal/database"
	"github.com/manish-npx/todo-go-echo/internal/handlers"
	middleware "github.com/manish-npx/todo-go-echo/internal/middleware"
	"github.com/manish-npx/todo-go-echo/internal/repository"
	"github.com/manish-npx/todo-go-echo/internal/routes"
	"github.com/manish-npx/todo-go-echo/internal/server"
	"github.com/manish-npx/todo-go-echo/internal/service"
	"github.com/manish-npx/todo-go-echo/internal/validator"
)

func main() {

	// Load config
	cfg, err := config.LoadConfig("config/config.yaml")

	if err != nil {
		log.Fatal(err)
	}

	// DB
	db, err := database.NewPostgresConnection(cfg.Database)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Dependency Injection todos
	todoRepo := repository.NewTodoRepository(db)
	todoService := service.NewTodoService(todoRepo)
	todoHandler := handlers.NewTodoHandler(todoService)

	//  Dependency Injection  category
	categoryRepo := repository.NewCategoryRepository(db)
	categoryService := service.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	//  Dependency Injection  Blogs
	blogRepo := repository.NewBlogRepository(db)
	blogService := service.NewBlogService(blogRepo, categoryRepo)
	blogHandler := handlers.NewBlogHandler(blogService)

	// USER Dependency Injection
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, cfg.JWT.Secret)
	userHandler := handlers.NewUserHandler(userService)

	// Echo Start
	e := echo.New()
	e.Validator = validator.New()

	// Middleware
	middleware.Setup(e)

	// Routes
	routes.RegisterRoutes(e, routes.RouteHandlers{
		TodoHandler:     todoHandler,
		CategoryHandler: categoryHandler,
		BlogHandler:     blogHandler,
		UserHandler:     userHandler,
		JWTSecret:       cfg.JWT.Secret,
	})

	// Static React build
	e.Static("/", "dist")

	// Start server
	server.Start(e, cfg.Server.Port, cfg.Server.Timeout)

}
