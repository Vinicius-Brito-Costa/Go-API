package query

import (
	"log"
	"net/http"
	"encoding/json"
	"fmt"
	requestUtil "go_api/internal/api/utils"
	api "go_api/internal/api/error"
)

func Query(res http.ResponseWriter, req *http.Request) {
	querys := make(map[string]string)
	for parameter := range req.URL.Query() {
		querys[parameter] = req.URL.Query().Get(parameter)
	}
	requestUtil.AddContentType(res, requestUtil.CONTENT_TYPE_APPLICATION_JSON)
	b, err := json.Marshal(querys)
	if err != nil {
		log.Println("Error trying to parse Response.")
		var apiError api.ApiError
		apiError.Code = 1
		apiError.Message = "Error trying to parse Response."
		res.WriteHeader(http.StatusInternalServerError)
		b = requestUtil.Serialize(apiError)
	}
	res.Write(b)
	fmt.Println(res, "Querys")
}