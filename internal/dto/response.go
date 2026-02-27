package dto

// APIResponse is a standard API response structure
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// SuccessResponse helper
func SuccessResponse(message string, data interface{}) APIResponse {
	return APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}

// ErrorResponse helper
func ErrorResponse(message string, err interface{}) APIResponse {
	return APIResponse{
		Success: false,
		Message: message,
		Error:   err,
	}
}
