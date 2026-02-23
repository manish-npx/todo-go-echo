package main

import (
	"log"

	"github.com/labstack/echo/v4"

	"github.com/manish-npx/todo-go-echo/internal/config"
	"github.com/manish-npx/todo-go-echo/internal/database"
	"github.com/manish-npx/todo-go-echo/internal/handlers"
	appmiddleware "github.com/manish-npx/todo-go-echo/internal/middleware"
	"github.com/manish-npx/todo-go-echo/internal/repository"
	"github.com/manish-npx/todo-go-echo/internal/routes"
	"github.com/manish-npx/todo-go-echo/internal/server"
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

	// Repositories
	todoRepo := repository.NewTodoRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	blogRepo := repository.NewBlogRepository(db)

	// Handlers
	todoHandler := handlers.NewTodoHandler(todoRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryRepo)
	blogHandler := handlers.NewBlogHandler(blogRepo, categoryRepo)

	// Echo
	e := echo.New()

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

	// Start server
	server.Start(e, cfg.Server.Port)

}
