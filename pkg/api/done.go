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
		writeError(w, http.StatusBadRequest, "не указан идентификатор")
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "задача не найдена")
		return
	}

	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		next, err := NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		if err := db.UpdateDate(next, id); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	writeJSON(w, http.StatusOK, map[string]any{})
}
