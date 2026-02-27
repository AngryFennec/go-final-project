package api

import (
	"net/http"

	"main.go/pkg/db"
)

const TasksLimit = 50

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.Tasks(TasksLimit)
	if err != nil {
		writeJson(w, http.StatusInternalServerError, taskError{err.Error()})
		return
	}
	writeJson(w, http.StatusOK, TasksResp{
		Tasks: tasks,
	})
}
