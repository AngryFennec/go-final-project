package api

import (
	"net/http"

	"main.go/pkg/db"
)

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, http.StatusBadRequest, taskError{Error: "Empty task id"})
		return
	}

	task, err := db.GetTask(id)

	if err != nil {
		writeJson(w, http.StatusInternalServerError, taskError{err.Error()})
		return
	}

	writeJson(w, http.StatusOK, task)

}
