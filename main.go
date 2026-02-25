package main

import (
	"log"
	"os"

	"main.go/pkg/db"
	"main.go/pkg/server"
)

func main() {
	logger := log.New(os.Stderr, "todo list", log.LstdFlags)

	err := db.Init("scheduler.db")
	if err != nil {
		logger.Fatalf("%s", err.Error())
	}

	err = server.RunServer()
	if err != nil {
		logger.Fatalf("%s", err.Error())
	}

}
