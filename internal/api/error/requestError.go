package api

import (
	"fmt"
)

const (
	InvalidRequestMessage = 0
	InvalidHttpMethod = 1
	InvalidHeader = 2
	ProcessingError = 3
	AuthenticationError = 4
)

type ApiError struct {
	Code int `json:"code"`
	Message string `json:"message"`
}
func (r *ApiError) Error() string {
	return fmt.Sprintf("code: %d, message: %s", r.Code, r.Message)
}