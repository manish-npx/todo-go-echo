package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/manish-npx/todo-go-echo/internal/logger"
	"go.uber.org/zap"
)

func Setup(e *echo.Echo) {
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.CORS())
	e.Use(echoMiddleware.RequestID())

	// ✅ Structured request logging for observability.
	e.Use(echoMiddleware.RequestLoggerWithConfig(
		echoMiddleware.RequestLoggerConfig{
			LogRequestID: true,
			LogURI:       true,
			LogMethod:    true,
			LogStatus:    true,
			LogLatency:   true,
			LogValuesFunc: func(c echo.Context, v echoMiddleware.RequestLoggerValues) error {
				logger.L().Info("http_request",
					zap.String("request_id", v.RequestID),
					zap.String("method", v.Method),
					zap.String("uri", v.URI),
					zap.Int("status", v.Status),
					zap.Duration("latency", v.Latency),
				)
				return nil
			},
		},
	))

	// ✅ Request timeout to protect server resources.
	e.Use(echoMiddleware.ContextTimeoutWithConfig(
		echoMiddleware.ContextTimeoutConfig{
			Timeout: 30 * time.Second,
		},
	))
}
