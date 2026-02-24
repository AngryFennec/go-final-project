package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

var port = 7540

const webDir = "./web"

func main() {
	logger := log.New(os.Stderr, "todo list", log.LstdFlags)

	http.Handle("/", http.FileServer(http.Dir(webDir)))
	log.Printf("Running server on %d", port)

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		logger.Fatalf("%s", err.Error())
	}
}
