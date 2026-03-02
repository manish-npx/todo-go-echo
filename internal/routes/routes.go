package routes

import (
	"github.com/labstack/echo/v4"

	"github.com/manish-npx/todo-go-echo/internal/handlers"
	appmiddleware "github.com/manish-npx/todo-go-echo/internal/middleware"
)

type RouteHandlers struct {
	TodoHandler     *handlers.TodoHandler
	CategoryHandler *handlers.CategoryHandler
	BlogHandler     *handlers.BlogHandler
	UserHandler     *handlers.UserHandler
	JWTSecret       string
}

func RegisterRoutes(e *echo.Echo, h RouteHandlers) {
	api := e.Group("/api/v1")

	// Health
	api.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"status": "OK",
		})
	})

	// Todos
	todos := api.Group("/todos")
	todos.GET("", h.TodoHandler.GetTodos)
	todos.POST("", h.TodoHandler.CreateTodo)
	todos.GET("/:id", h.TodoHandler.GetTodo)
	todos.PUT("/:id", h.TodoHandler.UpdateTodo)
	todos.DELETE("/:id", h.TodoHandler.DeleteTodo)

	// Categories
	categories := api.Group("/categories")
	categories.GET("", h.CategoryHandler.GetCategories)
	categories.POST("", h.CategoryHandler.CreateCategory)
	categories.GET("/:id", h.CategoryHandler.GetCategory)
	categories.PUT("/:id", h.CategoryHandler.UpdateCategory)
	categories.DELETE("/:id", h.CategoryHandler.DeleteCategory)

	// Blogs
	blogs := api.Group("/blogs")
	blogs.GET("", h.BlogHandler.GetBlogs)
	blogs.POST("", h.BlogHandler.CreateBlog)
	blogs.GET("/search", h.BlogHandler.SearchBlogs)
	blogs.GET("/:id", h.BlogHandler.GetBlog)
	blogs.PUT("/:id", h.BlogHandler.UpdateBlog)
	blogs.DELETE("/:id", h.BlogHandler.DeleteBlog)
	blogs.PATCH("/:id/publish", h.BlogHandler.PublishBlog)

	// Auth
	auth := api.Group("/auth")
	auth.POST("/register", h.UserHandler.Register)
	auth.POST("/login", h.UserHandler.Login)

	// Users (protected)
	users := api.Group("/users")
	users.Use(appmiddleware.JWTMiddleware(h.JWTSecret))
	users.GET("", h.UserHandler.GetUsers)
}
