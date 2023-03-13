package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	api "go_api/internal/api/error"
	"go_api/internal/api/filter"
	"go_api/internal/api/filter/custom"
	"io"
	"log"
	"net/http"
)

const CONTENT_TYPE_APPLICATION_JSON = "application/json"

func HandleRequest() {
	http.HandleFunc("/query", Query)
	http.HandleFunc("/login", Login)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
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
		AddContentType(res, CONTENT_TYPE_APPLICATION_JSON)
		res.WriteHeader(requestFilterError.HttpStatusCode)
		res.Write(Serialize(requestFilterError.Err))
		return
	}

	body, err := io.ReadAll(req.Body)

	if err != nil {
		var apiError api.ApiError
		apiError.Code = 1
		apiError.Message = "Error trying to parse Request."
		res.WriteHeader(http.StatusBadRequest)

		res.Write(Serialize(apiError))
		return
	}
	
	DoStuff(body, res)
}
func HandleInvalidMethod(method string, res http.ResponseWriter, req *http.Request) bool {
	if method == req.Method {
		return true
	}

	AddContentType(res, CONTENT_TYPE_APPLICATION_JSON)

	var apiError api.ApiError
	apiError.Code = 0
	apiError.Message = fmt.Sprintf("Error trying to parse Request. Invalid Method(%s)", req.Method)
	res.WriteHeader(http.StatusMethodNotAllowed)

	res.Write(Serialize(apiError))
	return false
}
func Serialize(strct interface{}) []byte {
	byteArray := new(bytes.Buffer)
	json.NewEncoder(byteArray).Encode(strct)
	return byteArray.Bytes()
}

func AddContentType(res http.ResponseWriter, contentType string) {
	res.Header().Add("Content-Type", contentType)
}
func Query(res http.ResponseWriter, req *http.Request) {
	querys := make(map[string]string)
	for parameter := range req.URL.Query() {
		querys[parameter] = req.URL.Query().Get(parameter)
	}
	AddContentType(res, CONTENT_TYPE_APPLICATION_JSON)
	b, err := json.Marshal(querys)
	if err != nil {
		log.Println("Error trying to parse Response.")
		var apiError api.ApiError
		apiError.Code = 1
		apiError.Message = "Error trying to parse Response."
		res.WriteHeader(http.StatusInternalServerError)
		b = Serialize(apiError)
	}
	res.Write(b)
	fmt.Println(res, "Querys")
}


type Response struct {
	Message string
}
func DoStuff(body []byte, res http.ResponseWriter) {
	AddContentType(res, CONTENT_TYPE_APPLICATION_JSON)
	res.WriteHeader(http.StatusOK)
	var response Response
	response.Message = "Finished doing stuff."

	res.Write(Serialize(response))
}
