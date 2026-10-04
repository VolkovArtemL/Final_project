package api

import (
	"encoding/json"
	"net/http"

	"github.com/VolkovArtemL/Final_project/pkg/db"
)

// getTaskHandler обрабатывает GET /api/task?id=<число>
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
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

	writeJSON(w, http.StatusOK, task)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, http.StatusBadRequest, "ошибка десериализации JSON")
		return
	}

	if task.ID == "" {
		writeError(w, http.StatusBadRequest, "не указан идентификатор")
		return
	}

	if task.Title == "" {
		writeError(w, http.StatusBadRequest, "не указан заголовок задачи")
		return
	}

	if err := checkDate(&task); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := db.UpdateTask(&task); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{})
}
