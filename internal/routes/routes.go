package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/manish-npx/todo-go-echo/internal/handlers"
	"github.com/manish-npx/todo-go-echo/internal/middleware"
)

type RouteHandlers struct {
	TodoHandler     *handlers.TodoHandler
	CategoryHandler *handlers.CategoryHandler
	BlogHandler     *handlers.BlogHandler
	UserHandler     *handlers.UserHandler
	JWTSecret       string // Secret injected once and used only for protected route middleware.
}

func RegisterRoutes(router *echo.Echo, routeHandlers RouteHandlers) {
	api := router.Group("/api/v1")

	// Health
	api.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "OK",
		})
	})

	// Todos
	todos := api.Group("/todos")
	todos.GET("", routeHandlers.TodoHandler.GetTodos)
	todos.POST("", routeHandlers.TodoHandler.CreateTodo)
	todos.GET("/:id", routeHandlers.TodoHandler.GetTodo)
	todos.PUT("/:id", routeHandlers.TodoHandler.UpdateTodo)
	todos.DELETE("/:id", routeHandlers.TodoHandler.DeleteTodo)

	// Categories
	categories := api.Group("/categories")
	categories.GET("", routeHandlers.CategoryHandler.GetCategories)
	categories.POST("", routeHandlers.CategoryHandler.CreateCategory)
	categories.GET("/:id", routeHandlers.CategoryHandler.GetCategory)
	categories.PUT("/:id", routeHandlers.CategoryHandler.UpdateCategory)
	categories.DELETE("/:id", routeHandlers.CategoryHandler.DeleteCategory)

	// Blogs
	blogs := api.Group("/blogs")
	blogs.GET("", routeHandlers.BlogHandler.GetBlogs)
	blogs.POST("", routeHandlers.BlogHandler.CreateBlog)
	blogs.GET("/search", routeHandlers.BlogHandler.SearchBlogs)
	blogs.GET("/:id", routeHandlers.BlogHandler.GetBlog)
	blogs.PUT("/:id", routeHandlers.BlogHandler.UpdateBlog)
	blogs.DELETE("/:id", routeHandlers.BlogHandler.DeleteBlog)
	blogs.PATCH("/:id/publish", routeHandlers.BlogHandler.PublishBlog)

	// Auth
	auth := api.Group("/auth")
	auth.POST("/register", routeHandlers.UserHandler.Register)
	auth.POST("/login", routeHandlers.UserHandler.Login)

	// Backward-compatible aliases without /auth prefix.
	api.POST("/register", routeHandlers.UserHandler.Register)
	api.POST("/login", routeHandlers.UserHandler.Login)

	// Users (protected)
	users := api.Group("/users")
	users.Use(middleware.JWTMiddleware(routeHandlers.JWTSecret))
	users.POST("", routeHandlers.UserHandler.CreateUser)
	users.GET("/profile", routeHandlers.UserHandler.Profile)
	users.GET("", routeHandlers.UserHandler.GetUsers)
}
