package api

import (
	"net/http"

	"github.com/VolkovArtemL/Final_project/pkg/db"
)

// deleteTaskHandler обрабатывает DELETE /api/task?id=<число>
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeError(w, "не указан идентификатор")
		return
	}

	if err := db.DeleteTask(id); err != nil {
		writeError(w, err.Error())
		return
	}

	writeJSON(w, map[string]any{})
}
