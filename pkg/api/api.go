package api

import "net/http"

func Init() {
	http.HandleFunc("/api/nextdate", nextDateHandler)
	http.HandleFunc("/api/task/done", taskDoneHandler)
	http.HandleFunc("/api/task", taskHandler)
	http.HandleFunc("/api/tasks", tasksHandler)
	http.Handle("/", http.FileServer(http.Dir("./web")))
}
