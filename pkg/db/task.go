package db

import (
	"database/sql"
	"fmt"
	"time"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// AddTask добавляет задачу в таблицу scheduler и возвращает идентификатор новой записи
func AddTask(task *Task) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`

	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// Tasks возвращает список задач, отсортированных по дате, максимум limit штук
func Tasks(limit int) ([]*Task, error) {
	rows, err := db.Query(
		`SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?`,
		limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]*Task, 0)

	for rows.Next() {
		task := &Task{}
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// TasksSearch возвращает задачи, отфильтрованные по подстроке search
// (в заголовке или комментарии) либо по конкретной дате, максимум limit штук
func TasksSearch(search string, limit int) ([]*Task, error) {
	// проверяем, не является ли search датой в формате 02.01.2006
	if date, err := time.Parse("02.01.2006", search); err == nil {
		return tasksByDate(date.Format("20060102"), limit)
	}

	return tasksByText(search, limit)
}

func tasksByText(search string, limit int) ([]*Task, error) {
	like := "%" + search + "%"

	rows, err := db.Query(
		`SELECT id, date, title, comment, repeat FROM scheduler
		 WHERE title LIKE ? OR comment LIKE ?
		 ORDER BY date LIMIT ?`,
		like, like, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanTasks(rows)
}

func tasksByDate(date string, limit int) ([]*Task, error) {
	rows, err := db.Query(
		`SELECT id, date, title, comment, repeat FROM scheduler
		 WHERE date = ?
		 ORDER BY date LIMIT ?`,
		date, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanTasks(rows)
}

func scanTasks(rows *sql.Rows) ([]*Task, error) {
	tasks := make([]*Task, 0)

	for rows.Next() {
		task := &Task{}
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// GetTask возвращает задачу по её идентификатору
func GetTask(id string) (*Task, error) {
	task := &Task{}

	row := db.QueryRow(
		`SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`,
		id)

	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return nil, err
	}

	return task, nil
}

// UpdateTask обновляет параметры задачи по её идентификатору
func UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`

	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}

	return nil
}
