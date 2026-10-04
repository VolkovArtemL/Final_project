package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/VolkovArtemL/Final_project/pkg/db"
)

// checkDate проверяет и при необходимости корректирует дату задачи
func checkDate(task *db.Task) error {
	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format(timeFormat)
	}

	t, err := time.Parse(timeFormat, task.Date)
	if err != nil {
		return errors.New("дата представлена в неверном формате")
	}

	var next string
	if task.Repeat != "" {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}

	if afterNow(now, t) {
		if task.Repeat == "" {
			task.Date = now.Format(timeFormat)
		} else {
			task.Date = next
		}
	}

	return nil
}

// writeJSON сериализует data в JSON и записывает в ответ
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)

	resp, err := json.Marshal(data)
	if err != nil {
		w.Write([]byte(`{"error":"ошибка сериализации ответа"}`))
		return
	}

	w.Write(resp)
}

// writeError записывает в ответ JSON с текстом ошибки
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

// addTaskHandler обрабатывает POST /api/task — добавление новой задачи
func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, http.StatusBadRequest, "ошибка десериализации JSON")
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

	id, err := db.AddTask(&task)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"id": strconv.FormatInt(id, 10)})
}
