package main

import (
	"log"

	"github.com/labstack/echo/v4"

	_ "github.com/lib/pq"
	"github.com/manish-npx/todo-go-echo/internal/config"
	"github.com/manish-npx/todo-go-echo/internal/database"
	"github.com/manish-npx/todo-go-echo/internal/handlers"
	appmiddleware "github.com/manish-npx/todo-go-echo/internal/middleware"
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

	// todos
	todoRepo := repository.NewTodoRepository(db)
	todoService := service.NewTodoService(todoRepo)
	todoHandler := handlers.NewTodoHandler(todoService)

	//category
	categoryRepo := repository.NewCategoryRepository(db)
	categoryHandler := handlers.NewCategoryHandler(categoryRepo)

	// Blogs
	blogRepo := repository.NewBlogRepository(db)
	blogHandler := handlers.NewBlogHandler(blogRepo, categoryRepo)

	// Dependency Injection
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, cfg.JWT.Secret)
	userHandler := handlers.NewUserHandler(userService)

	// Echo
	e := echo.New()
	e.Validator = validator.New()

	// Middleware
	appmiddleware.Setup(e)

	// Routes
	routes.RegisterRoutes(e, routes.RouteHandlers{
		TodoHandler:     todoHandler,
		CategoryHandler: categoryHandler,
		BlogHandler:     blogHandler,
	})

	// Static React build
	e.Static("/", "dist")

	// Routes
	e.POST("/register", userHandler.Register)
	e.POST("/login", userHandler.Login)
	authGroup := e.Group("")
	authGroup.Use(appmiddleware.JWTMiddleware(cfg.JWT.Secret))
	authGroup.GET("/users", userHandler.GetUsers)

	// Start server
	server.Start(e, cfg.Server.Port)

}
