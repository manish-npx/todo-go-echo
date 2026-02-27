package utils

import "github.com/labstack/echo/v4"

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func Success(c echo.Context, message string, data interface{}) error {
	return c.JSON(200, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Error(c echo.Context, status int, message string, err interface{}) error {
	return c.JSON(status, APIResponse{
		Success: false,
		Message: message,
		Error:   err,
	})
}
