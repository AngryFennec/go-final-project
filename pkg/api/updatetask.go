package api

import (
	"encoding/json"
	"net/http"

	"main.go/pkg/db"
)

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var req db.Task

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJson(w, http.StatusBadRequest, taskError{err.Error()})
		return
	}

	if req.Title == "" {
		writeJson(w, http.StatusBadRequest, taskError{Error: "Empty title"})
		return
	}

	err := checkDate(&req)
	if err != nil {
		writeJson(w, http.StatusBadRequest, taskError{err.Error()})
		return
	}

	err = db.UpdateTask(&req)

	if err != nil {
		writeJson(w, http.StatusInternalServerError, taskError{err.Error()})
		return
	}

	writeJson(w, http.StatusOK, struct{}{})
}
