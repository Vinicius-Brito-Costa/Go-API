package fetch

import (
	"encoding/json"
	"errors"
	"fmt"
	fetchError "go_api/internal/api/fetch/error"
	"io"
	"net/http"
	"time"
)

const TIMEOUT_SECONDS = 30

func Get(url string, headers map[string] string, body io.Reader, response interface{}) (err error) {
	return BaseRequest(http.MethodGet, url, headers, body, response)
}

func Post(url string, headers map[string] string, body io.Reader, response interface{}) (err error){
	return BaseRequest(http.MethodPost, url, headers, body, response)
}

func BaseRequest(method string ,url string, headers map[string] string, body io.Reader, response interface{}) (err error){
	req, err := http.NewRequest(method, url, body);
	if err != nil {
		return err
	}
	if headers != nil {
		for header := range headers {
			fmt.Println(header, headers[header])
			req.Header.Set(header, headers[header])
		}
	}
	
	client := http.Client{
		Timeout: TIMEOUT_SECONDS * time.Second,
	}

	resp, err := client.Do(req)
	if err !=  nil {
		return err
	}
	
	defer resp.Body.Close()
	res, err := io.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return &fetchError.RequestError{
			StatusCode: resp.StatusCode,
			Body: res,
			Err: errors.New("Invalid response code."),
		}
	}
	err = json.Unmarshal(res, &response)
	return err
}