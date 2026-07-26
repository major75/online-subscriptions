package models

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func SuccessResponse(data any, message string) *Response {
	return &Response{
		Success: true,
		Message: message,
		Data:    data,
	}
}

func ErrorResponse(data any, message string) *Response {
	return &Response{
		Success: false,
		Message: message,
		Data:    data,
	}
}
