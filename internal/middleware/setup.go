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

	// ✅ Proper Request Logger (New API)
	e.Use(echoMiddleware.RequestLoggerWithConfig(
		echoMiddleware.RequestLoggerConfig{
			LogURI:     true,
			LogMethod:  true,
			LogStatus:  true,
			LogLatency: true,

			LogValuesFunc: func(c echo.Context, v echoMiddleware.RequestLoggerValues) error {
				log.Printf(
					"METHOD=%s URI=%s STATUS=%d LATENCY=%s",
					v.Method,
					v.URI,
					v.Status,
					v.Latency,
				)
				return nil
			},
		},
	))

	// ✅ Context Timeout (race-free)
	e.Use(echoMiddleware.ContextTimeoutWithConfig(
		echoMiddleware.ContextTimeoutConfig{
			Timeout: 30 * time.Second,
		},
	))
}
