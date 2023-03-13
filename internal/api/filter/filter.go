package filter

import (
	"fmt"
	"go_api/internal/api/error"
	"net/http"
	"regexp"
	"golang.org/x/exp/slices"
)
type HttpMethod string
const (
	GET HttpMethod = http.MethodGet
	POST HttpMethod = http.MethodPost
	PATCH HttpMethod = http.MethodPatch
	PUT HttpMethod = http.MethodPut
	DELETE HttpMethod = http.MethodDelete
)

type Header struct {
	Name string
	Required bool
	Validation string
}

type CustomFilter interface {
	Filter(req *http.Request) RequestFilterError
}

type RequestConfig struct {
	ValidMethods []HttpMethod
	ValidHeaders []Header
	CustomFilter []CustomFilter
}

type RequestFilterError struct {
	HttpStatusCode int
	Err error
}

func RequestFilter(req *http.Request, config RequestConfig) RequestFilterError {

	var currentMethod HttpMethod = HttpMethod(req.Method)
	if !slices.Contains(config.ValidMethods, currentMethod){
		return RequestFilterError{
			HttpStatusCode: http.StatusBadRequest,
			Err: &api.ApiError {
				Code: api.InvalidHttpMethod,
				Message: fmt.Sprintf("Invalid Http Method (%s)", currentMethod),
			},
		}
	}
	headers := req.Header
	for _, validHeader := range config.ValidHeaders {
		if validHeader.Required {
			header := headers.Get(validHeader.Name)
			if header == "" {
				return RequestFilterError{
					HttpStatusCode: http.StatusBadRequest,
					Err: &api.ApiError{
						Code: api.InvalidHeader,
						Message: fmt.Sprintf("Missing Header (%s)", validHeader.Name),
					},
				}
			}
			regex, err := regexp.Compile(validHeader.Validation)
			if err != nil {
				return RequestFilterError{
					HttpStatusCode: http.StatusBadRequest,
					Err: &api.ApiError{
						Code: api.ProcessingError,
						Message: fmt.Sprintf("Invalid Header Validation (%s)", validHeader.Name),
					},
				}
			}
			if !regex.MatchString(header) {
				return RequestFilterError{
					HttpStatusCode: http.StatusBadRequest,
					Err: &api.ApiError{
						Code: api.InvalidHeader,
						Message: fmt.Sprintf("Invalid Header Value (%s)", validHeader.Name),
					},
				}
			}
		}
	}

	for _, customFilter := range config.CustomFilter {
		response := customFilter.Filter(req)
		if response.Err != nil {
			return response
		}
	}

	return RequestFilterError{}
}