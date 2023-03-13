package serialize

import (
	"bytes"
	"encoding/json"
	"net/http"
	"fmt"
	api "go_api/internal/api/error"
)

const CONTENT_TYPE_APPLICATION_JSON = "application/json"

func Serialize(strct interface{}) []byte {
	byteArray := new(bytes.Buffer)
	json.NewEncoder(byteArray).Encode(strct)
	return byteArray.Bytes()
}

func AddContentType(res http.ResponseWriter, contentType string) {
	res.Header().Add("Content-Type", contentType)
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