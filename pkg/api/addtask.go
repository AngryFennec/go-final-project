package api

import (
	"encoding/json"
	"fmt"
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

func checkDate(req taskRequest) (string, error) {
	now := time.Now()
	if req.Date == "" {
		return now.Format(DateFormat), nil
	}
	date, err := time.Parse(DateFormat, req.Date)
	if err != nil {
		return "", err
	}

	var next string

	if req.Repeat != "" {
		next, err = NextDate(now, req.Date, req.Repeat)
		if err != nil {
			return "", err
		}
	}

	if afterNow(now, date) {
		if len(req.Repeat) == 0 {
			return now.Format(DateFormat), nil
		} else {
			return next, nil
		}
	}

	return "", fmt.Errorf("wrong date parsing")

}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var req taskRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Decoding Error: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		http.Error(w, "Empty title", http.StatusBadRequest)
		return
	}

	date, err := checkDate(req)
	if err != nil {
		http.Error(w, "Date error: "+err.Error(), http.StatusBadRequest)
		return
	}

	task := &db.Task{
		Date:    date,
		Title:   req.Title,
		Comment: req.Comment,
		Repeat:  req.Repeat,
	}

	id, err := db.AddTask(task)

	if err != nil {
		http.Error(w, "Add task error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJson(w, http.StatusOK, taskID{id})

}
