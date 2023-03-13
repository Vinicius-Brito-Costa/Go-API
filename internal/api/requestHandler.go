package api

import (
	"go_api/internal/api/service/login"
	"go_api/internal/api/service/query"
	"log"
	"net/http"
)

func HandleRequest() {
	http.HandleFunc("/query", query.Query)
	http.HandleFunc("/login", login.Login)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
