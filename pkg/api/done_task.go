package api

import (
	"net/http"
	"time"

	"main.go/pkg/db"
)

func taskDoneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")

	if id == "" {
		writeJson(w, http.StatusBadRequest, taskError{Error: "Empty task id"})
		return
	}
	task, err := db.GetTask(id)
	if err != nil {
		writeJson(w, http.StatusBadRequest, taskError{Error: err.Error()})
		return
	}

	if task.Repeat == "" {
		err := db.DeleteTask(id)
		if err != nil {
			writeJson(w, http.StatusInternalServerError, taskError{Error: err.Error()})
			return
		}
	} else {
		now := time.Now()
		next, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeJson(w, http.StatusBadRequest, taskError{Error: err.Error()})
			return
		}
		if err := db.UpdateDate(next, id); err != nil {
			writeJson(w, http.StatusInternalServerError, taskError{Error: err.Error()})
			return
		}
	}
	writeJson(w, http.StatusOK, struct{}{})
}
