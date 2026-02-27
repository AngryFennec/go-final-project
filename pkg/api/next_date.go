package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const DateFormat = "20060102"
const MaxDays = 400

func afterNow(date, now time.Time) bool {
	dateWithoutTime := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	nowWithoutTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return dateWithoutTime.After(nowWithoutTime)
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	date, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", err
	}

	repeatRule := strings.Split(repeat, " ")

	switch repeatRule[0] {
	case "y":
		if len(repeatRule) > 1 {
			return "", errors.New("wrong year format")
		}

		newDate := date
		for {
			newDate = newDate.AddDate(1, 0, 0)
			if afterNow(newDate, now) {
				break
			}
		}
		return newDate.Format(DateFormat), nil

	case "d":
		if len(repeatRule) != 2 {
			return "", errors.New("wrong day format")
		}
		daysCount, err := strconv.Atoi(repeatRule[1])
		if err != nil {
			return "", err
		}
		if daysCount > MaxDays {
			return "", errors.New("wrong day format")
		}

		newDate := date
		for {
			newDate = newDate.AddDate(0, 0, daysCount)
			if afterNow(newDate, now) {
				break
			}
		}
		return newDate.Format(DateFormat), nil
	}

	return "", errors.New("wrong date format")
}

func parseNow(rawNow string) (time.Time, error) {
	if rawNow == "" {
		return time.Now(), nil
	}
	now, err := time.Parse(DateFormat, rawNow)
	return now, err
}

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rawNow := r.FormValue("now")
	rawDate := r.FormValue("date")
	rawRepeat := r.FormValue("repeat")

	if rawRepeat == "" {
		http.Error(w, "date parameter is required", http.StatusBadRequest)
	}

	if rawDate == "" {
		http.Error(w, "date parameter is required", http.StatusBadRequest)
		return
	}

	parsedNow, err := parseNow(rawNow)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	nextDate, err := NextDate(parsedNow, rawDate, rawRepeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(nextDate))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
