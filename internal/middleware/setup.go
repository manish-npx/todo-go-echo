package appmiddleware

import (
	"log"
	"time"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
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
				log.Printf(
					"REQUEST_ID=%s METHOD=%s URI=%s STATUS=%d LATENCY=%s",
					v.RequestID,
					v.Method,
					v.URI,
					v.Status,
					v.Latency,
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
