package login

import (
	"net/http"
	"go_api/internal/api/filter"
	"go_api/internal/api/filter/custom"
	requestUtil "go_api/internal/api/utils"
	api "go_api/internal/api/error"
	"io"
)


func Login(res http.ResponseWriter, req *http.Request) {

	var validMethods []filter.HttpMethod
	validMethods = append(validMethods, filter.POST)
	validMethods = append(validMethods, filter.PUT)

	var validHeaders []filter.Header
	validHeaders = append(validHeaders, filter.Header{
		Name:       "Content-Type",
		Required:   true,
		Validation: "application\\/json",
	})

	var customFilters []filter.CustomFilter
	loginFilter := custom.LoginFilter{}
	customFilters = append(customFilters, &loginFilter)

	config := filter.RequestConfig{
		ValidMethods: validMethods,
		ValidHeaders: validHeaders,
		CustomFilter: customFilters,
	}
	requestFilterError := filter.RequestFilter(req, config)
	if requestFilterError.Err != nil {
		requestUtil.AddContentType(res, requestUtil.CONTENT_TYPE_APPLICATION_JSON)
		res.WriteHeader(requestFilterError.HttpStatusCode)
		res.Write(requestUtil.Serialize(requestFilterError.Err))
		return
	}

	body, err := io.ReadAll(req.Body)

	if err != nil {
		var apiError api.ApiError
		apiError.Code = 1
		apiError.Message = "Error trying to parse Request."
		res.WriteHeader(http.StatusBadRequest)

		res.Write(requestUtil.Serialize(apiError))
		return
	}
	
	DoStuff(body, res)
}

type Response struct {
	Message string
}
func DoStuff(body []byte, res http.ResponseWriter) {
	requestUtil.AddContentType(res, requestUtil.CONTENT_TYPE_APPLICATION_JSON)
	res.WriteHeader(http.StatusOK)
	var response Response
	response.Message = "Finished doing stuff."

	res.Write(requestUtil.Serialize(response))
}