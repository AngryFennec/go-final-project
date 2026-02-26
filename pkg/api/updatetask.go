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

	date, err := checkDate(req)
	if err != nil {
		writeJson(w, http.StatusBadRequest, taskError{err.Error()})
		return
	}

	task := &db.Task{
		ID:      req.ID,
		Date:    date,
		Title:   req.Title,
		Comment: req.Comment,
		Repeat:  req.Repeat,
	}

	err = db.UpdateTask(task)

	if err != nil {
		writeJson(w, http.StatusInternalServerError, taskError{err.Error()})
		return
	}

	writeJson(w, http.StatusOK, struct{}{})
}
