package errortools

import "net/http"

// ErrorDto is used to return the error back to the client.
type ErrorDto struct {
	Status  int    `json:"status"`
	Error   string `json:"error"`
	Message any    `json:"message"`
} //	@name	ErrorDto

// NewErrorDto creates a new [ErrorDto].
func NewErrorDto(status int, message any) ErrorDto {
	return ErrorDto{
		Status:  status,
		Error:   http.StatusText(status),
		Message: message,
	}
}
