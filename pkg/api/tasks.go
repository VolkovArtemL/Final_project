package api

import (
	"net/http"

	"github.com/VolkovArtemL/Final_project/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

const tasksLimit = 50

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.FormValue("search")

	var tasks []*db.Task
	var err error

	if search != "" {
		tasks, err = db.TasksSearch(search, tasksLimit)
	} else {
		tasks, err = db.Tasks(tasksLimit)
	}

	if err != nil {
		writeError(w, err.Error())
		return
	}

	writeJSON(w, TasksResp{
		Tasks: tasks,
	})
}
