package server

import (
	"fmt"
	"log"
	"net/http"

	"main.go/pkg/api"
)

var port = 7540

const webDir = "./web"

func RunServer() error {
	log.Printf("Running server on %d", port)
	api.Init()
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
