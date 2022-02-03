package server

// APIError - wrapper struct for returning errors in JSON form.
type APIError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// NewApiError - returns new APIError instance.
func NewApiError(status int, err error) *APIError {
	return &APIError{
		Status:  status,
		Message: err.Error(),
	}
}
