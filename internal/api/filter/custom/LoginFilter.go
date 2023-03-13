package custom

import (
	"encoding/json"
	api "go_api/internal/api/error"
	"go_api/internal/api/filter"
	"go_api/internal/api/model/login"
	"io"
	"net/http"
)

type LoginFilter struct {
}

func (f *LoginFilter) Filter(req *http.Request) filter.RequestFilterError {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return filter.RequestFilterError{
			HttpStatusCode: http.StatusBadRequest,
			Err: &api.ApiError{
				Code: api.InvalidRequestMessage,
				Message: "Invalid request.",
			},
		}
	}
	var user login.User
	json.Unmarshal(body, &user)
	if user.Username == "Vinicius" && user.Password == "vinicius123" {
		return filter.RequestFilterError{}
	}
	return filter.RequestFilterError {
		HttpStatusCode: http.StatusForbidden,
		Err: &api.ApiError{
			Code: api.AuthenticationError,
			Message: "Invalid user or password.",
		},
	}
}