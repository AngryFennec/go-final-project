package api

import (
	"encoding/json"
	"net/http"
	"time"

	"main.go/pkg/db"
)

type taskRequest struct {
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type taskID struct {
	ID int64 `json:"id"`
}

type taskError struct {
	Error string `json:"error,omitempty"`
}

func writeJson(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	enc := json.NewEncoder(w)
	if err := enc.Encode(data); err != nil {
		http.Error(w, "Failed to encode json: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func checkDate(task *db.Task) error {
	now := time.Now()
	if len(task.Date) == 0 {
		task.Date = now.Format(DateFormat)
		return nil

	}
	date, err := time.Parse(DateFormat, task.Date)
	if err != nil {
		return err
	}

	if afterNow(now, date) {
		if len(task.Repeat) == 0 {
			task.Date = now.Format(DateFormat)
			return nil
		}
		next, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
		task.Date = next
	}

	return nil
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var req taskRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJson(w, http.StatusBadRequest, taskError{err.Error()})
		return
	}

	if req.Title == "" {
		writeJson(w, http.StatusBadRequest, taskError{Error: "Empty title"})
		return
	}

	task := &db.Task{
		Date:    req.Date,
		Title:   req.Title,
		Comment: req.Comment,
		Repeat:  req.Repeat,
	}

	err := checkDate(task)
	if err != nil {
		writeJson(w, http.StatusBadRequest, taskError{err.Error()})
		return
	}

	id, err := db.AddTask(task)

	if err != nil {
		writeJson(w, http.StatusInternalServerError, taskError{err.Error()})
		return
	}

	writeJson(w, http.StatusOK, taskID{id})
}
