package main

import (
	"log"
	"net/http"
)

const (
	AUTHORIZATION_CODE_DURATION = 600
	ACCESS_TOKEN_DURATION       = 86400
)

var (
	SUPPORTED_SCOPES = []string{"read", "write"}
)

func authorization(w http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()

	session := Session{
		Client:      query.Get("client_id"),
		State:       query.Get("state"),
		Scope:       query.Get("scope"),
		RedirectUri: query.Get("redirect_uri"),
	}
}

func main() {
	log.Println("start oauth server on localhost:8081...")
	http.HandleFunc("/authorization", authorization)
	err := http.ListenAndServe("localhost:8081", nil)
	if err != nil {
		log.Fatal(err)
	}
}
