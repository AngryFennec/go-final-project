package api

import (
	"net/http"

	"main.go/pkg/db"
)

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, http.StatusBadRequest, taskError{Error: "Empty task id"})
		return
	}
	if err := db.DeleteTask(id); err != nil {
		writeJson(w, http.StatusInternalServerError, taskError{Error: err.Error()})

		return
	}
	writeJson(w, http.StatusOK, struct{}{})
}
