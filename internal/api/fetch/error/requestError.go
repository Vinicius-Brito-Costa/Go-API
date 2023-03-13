package requests

import (
	"fmt"
)
type RequestError struct {
	StatusCode int
	Body []byte
	Err error
}

func (re *RequestError) Error() string{
	return fmt.Sprintf("There's a request error. Status Code: %d. Body: %s. Error: %s", re.StatusCode, re.Body, re.Err.Error())
}

