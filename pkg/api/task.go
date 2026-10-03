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
		writeError(w, "не указан идентификатор")
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeError(w, "задача не найдена")
		return
	}

	writeJSON(w, task)
}

// updateTaskHandler обрабатывает PUT /api/task — изменение существующей задачи
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, "ошибка десериализации JSON")
		return
	}

	if task.ID == "" {
		writeError(w, "не указан идентификатор")
		return
	}

	if task.Title == "" {
		writeError(w, "не указан заголовок задачи")
		return
	}

	if err := checkDate(&task); err != nil {
		writeError(w, err.Error())
		return
	}

	if err := db.UpdateTask(&task); err != nil {
		writeError(w, err.Error())
		return
	}

	writeJSON(w, map[string]any{})
}
