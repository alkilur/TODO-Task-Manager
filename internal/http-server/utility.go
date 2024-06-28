package http_server

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/render"
)

const TimeLayout string = "20060102"

var (
	ErrInvalidNow       = errors.New("invalid 'now' format")
	ErrUnmarshal        = errors.New("error unmarshalling request body")
	ErrEmptyTitle       = errors.New("'title' cannot be empty")
	ErrInvalidDate      = errors.New("invalid 'date' format")
	ErrInvalidRepeat    = errors.New("invalid 'repeat' format")
	ErrInvalidID        = errors.New("invalid 'id' format")
	ErrTaskNotFound     = errors.New("task not found")
	ErrMethodNotAllowed = errors.New("invalid http method")
)

type Task struct {
	ID      string `json:"id" db:"id"`
	Date    string `json:"date" db:"date"`
	Title   string `json:"title" db:"title"`
	Comment string `json:"comment" db:"comment"`
	Repeat  string `json:"repeat" db:"repeat"`
}

func SendError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, ErrInvalidNow) ||
		errors.Is(err, ErrInvalidRepeat) ||
		errors.Is(err, ErrUnmarshal) ||
		errors.Is(err, ErrEmptyTitle) ||
		errors.Is(err, ErrInvalidID) ||
		errors.Is(err, ErrInvalidDate):
		w.WriteHeader(http.StatusBadRequest)
	case errors.Is(err, ErrMethodNotAllowed):
		w.WriteHeader(http.StatusMethodNotAllowed)
	case errors.Is(err, ErrTaskNotFound):
		w.WriteHeader(http.StatusNotFound)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	render.JSON(w, r, map[string]string{"error": err.Error()})
}

func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", ErrInvalidRepeat
	}

	nextDate, err := time.Parse(TimeLayout, date)
	if err != nil {
		return "", ErrInvalidDate
	}

	repeatParts := strings.Split(repeat, " ")
	switch repeatParts[0] {
	case "d":
		if len(repeatParts) != 2 {
			return "", ErrInvalidRepeat
		}
		days, err := strconv.Atoi(repeatParts[1])
		if err != nil || days >= 400 {
			return "", ErrInvalidRepeat
		}
		nextDate = nextDate.AddDate(0, 0, days)
		for nextDate.Before(now) || nextDate.Equal(now) {
			nextDate = nextDate.AddDate(0, 0, days)
		}
	case "y":
		nextDate = nextDate.AddDate(1, 0, 0)
		for nextDate.Before(now) || nextDate.Equal(now) {
			nextDate = nextDate.AddDate(1, 0, 0)
		}
	default:
		return "", ErrInvalidRepeat
	}

	return nextDate.Format(TimeLayout), nil
}
