package api

func NewErrorResponse(err error) *ErrorResponse {
	return &ErrorResponse{Errors: []string{err.Error()}}
}

type ErrorResponse struct {
	Errors []string `json:"errors"`
}
