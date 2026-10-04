package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const timeFormat = "20060102"

// afterNow сравнивает только календарные даты (без учёта времени),
// возвращает true, если date больше now
func afterNow(date, now time.Time) bool {
	y1, m1, d1 := date.Date()
	y2, m2, d2 := now.Date()

	if y1 != y2 {
		return y1 > y2
	}
	if m1 != m2 {
		return m1 > m2
	}
	return d1 > d2
}

// NextDate вычисляет следующую дату выполнения задачи по правилу repeat,
// начиная с даты dstart, так чтобы результат был строго больше now
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("правило повторения не может быть пустым")
	}

	date, err := time.Parse(timeFormat, dstart)
	if err != nil {
		return "", errors.New("некорректная дата dstart")
	}

	parts := strings.Split(repeat, " ")

	switch parts[0] {
	case "d":
		date, err = nextDateDay(date, now, parts)
	case "y":
		date = nextDateYear(date, now)
	case "w":
		date, err = nextDateWeek(date, now, parts)
	case "m":
		date, err = nextDateMonth(date, now, parts)
	default:
		return "", errors.New("неподдерживаемый формат правила повторения")
	}

	if err != nil {
		return "", err
	}

	return date.Format(timeFormat), nil
}

// nextDateDay обрабатывает правило "d <число>"
func nextDateDay(date, now time.Time, parts []string) (time.Time, error) {
	if len(parts) != 2 {
		return date, errors.New("не указан интервал в днях")
	}

	interval, err := strconv.Atoi(parts[1])
	if err != nil || interval < 1 || interval > 400 {
		return date, errors.New("недопустимый интервал в днях")
	}

	for {
		date = date.AddDate(0, 0, interval)
		if afterNow(date, now) {
			break
		}
	}

	return date, nil
}

// nextDateYear обрабатывает правило "y"
func nextDateYear(date, now time.Time) time.Time {
	for {
		date = date.AddDate(1, 0, 0)
		if afterNow(date, now) {
			break
		}
	}

	return date
}

// nextDateWeek обрабатывает правило "w <дни недели через запятую>"
func nextDateWeek(date, now time.Time, parts []string) (time.Time, error) {
	if len(parts) != 2 {
		return date, errors.New("не указаны дни недели")
	}

	weekdays, err := parseIntList(parts[1], 1, 7, "недопустимый день недели")
	if err != nil {
		return date, err
	}

	for {
		date = date.AddDate(0, 0, 1)

		wd := int(date.Weekday())
		if wd == 0 {
			wd = 7 // воскресенье
		}

		if contains(weekdays, wd) && afterNow(date, now) {
			break
		}
	}

	return date, nil
}

// nextDateMonth обрабатывает правило "m <дни месяца> [месяцы]"
func nextDateMonth(date, now time.Time, parts []string) (time.Time, error) {
	if len(parts) < 2 {
		return date, errors.New("не указаны дни месяца")
	}

	days, err := parseMonthDays(parts[1])
	if err != nil {
		return date, err
	}

	var months []int
	if len(parts) == 3 {
		months, err = parseIntList(parts[2], 1, 12, "недопустимый месяц")
		if err != nil {
			return date, err
		}
	}

	for {
		date = date.AddDate(0, 0, 1)

		if len(months) > 0 && !contains(months, int(date.Month())) {
			continue
		}

		if dayMatches(date, days) && afterNow(date, now) {
			break
		}
	}

	return date, nil
}

// dayMatches проверяет, подходит ли день месяца date под список days,
// с учётом специальных значений -1 (последний день) и -2 (предпоследний)
func dayMatches(date time.Time, days []int) bool {
	lastDay := lastDayOfMonth(date)

	for _, d := range days {
		switch {
		case d > 0 && date.Day() == d:
			return true
		case d == -1 && date.Day() == lastDay:
			return true
		case d == -2 && date.Day() == lastDay-1:
			return true
		}
	}

	return false
}

// lastDayOfMonth возвращает номер последнего дня месяца для даты date
func lastDayOfMonth(date time.Time) int {
	firstOfNextMonth := time.Date(date.Year(), date.Month()+1, 1, 0, 0, 0, 0, date.Location())
	lastDay := firstOfNextMonth.AddDate(0, 0, -1)
	return lastDay.Day()
}

// parseIntList парсит строку вида "1,4,5" в список чисел,
// проверяя, что каждое значение в диапазоне [min, max]
func parseIntList(s string, min, max int, errMsg string) ([]int, error) {
	parts := strings.Split(s, ",")
	result := make([]int, 0, len(parts))

	for _, p := range parts {
		n, err := strconv.Atoi(p)
		if err != nil || n < min || n > max {
			return nil, errors.New(errMsg)
		}
		result = append(result, n)
	}

	return result, nil
}

// parseMonthDays парсит список дней месяца, разрешая 1..31, -1, -2
func parseMonthDays(s string) ([]int, error) {
	parts := strings.Split(s, ",")
	result := make([]int, 0, len(parts))

	for _, p := range parts {
		n, err := strconv.Atoi(p)
		if err != nil || n == 0 || n < -2 || n > 31 {
			return nil, errors.New("недопустимый день месяца")
		}
		result = append(result, n)
	}

	return result, nil
}

// contains проверяет наличие значения v в слайсе list
func contains(list []int, v int) bool {
	for _, x := range list {
		if x == v {
			return true
		}
	}
	return false
}

// nextDateHandler — обработчик GET /api/nextdate
func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	now := time.Now()
	if nowStr != "" {
		parsedNow, err := time.Parse(timeFormat, nowStr)
		if err != nil {
			http.Error(w, "некорректный параметр now", http.StatusBadRequest)
			return
		}
		now = parsedNow
	}

	next, err := NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(next))
}
