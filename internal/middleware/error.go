package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/manish-npx/todo-go-echo/internal/dto"
	"github.com/manish-npx/todo-go-echo/internal/logger"
	"go.uber.org/zap"
)

// ErrorHandler centralizes HTTP error responses.
func ErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	code := http.StatusInternalServerError
	msg := "internal server error"

	if httpErr, ok := err.(*echo.HTTPError); ok {
		code = httpErr.Code
		switch m := httpErr.Message.(type) {
		case string:
			msg = m
		default:
			msg = http.StatusText(httpErr.Code)
			if msg == "" {
				msg = "request failed"
			}
		}
	} else if err != nil {
		msg = err.Error()
	}

	logger.L().Error("request_failed",
		zap.Int("status", code),
		zap.String("path", c.Path()),
		zap.String("method", c.Request().Method),
		zap.Error(err),
	)

	_ = c.JSON(code, dto.ErrorResponse(msg, nil))
}
