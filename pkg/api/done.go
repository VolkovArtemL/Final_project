package api

import (
	"net/http"
	"time"

	"github.com/VolkovArtemL/Final_project/pkg/db"
)

// doneTaskHandler обрабатывает POST /api/task/done?id=<число>
func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeError(w, "не указан идентификатор")
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeError(w, "задача не найдена")
		return
	}

	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			writeError(w, err.Error())
			return
		}
	} else {
		next, err := NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			writeError(w, err.Error())
			return
		}

		if err := db.UpdateDate(next, id); err != nil {
			writeError(w, err.Error())
			return
		}
	}

	writeJSON(w, map[string]any{})
}
