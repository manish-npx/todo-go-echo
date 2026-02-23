package appmiddleware

import (
	"time"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func Setup(e *echo.Echo) {

	e.Use(echoMiddleware.Recover())

	e.Use(echoMiddleware.RequestLoggerWithConfig(
		echoMiddleware.RequestLoggerConfig{
			LogURI:    true,
			LogMethod: true,
			LogStatus: true,
		},
	))

	e.Use(echoMiddleware.CORS())

	e.Use(echoMiddleware.ContextTimeoutWithConfig(
		echoMiddleware.ContextTimeoutConfig{
			Timeout: 30 * time.Second,
		},
	))
}
