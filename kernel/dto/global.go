package dto

// ErrorResponse is the standard error envelope returned by all handlers.
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

type SuccessResponse struct {
	Message string `json:"message,omitempty"`
}
