package api

import "net/http"

func Init() {
	http.HandleFunc("/api/nextdate", nextDateHandler)
	http.Handle("/", http.FileServer(http.Dir("./web")))
}
