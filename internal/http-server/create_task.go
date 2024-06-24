package http_server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type TaskCreator interface {
	CreateTask(*Task) (string, error)
}

func (c *Controller) AddTask(log *slog.Logger, taskCreator TaskCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log = log.With(slog.String("request_id", middleware.GetReqID(r.Context())))
		// TODO: log something?

		task := Task{}
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			c.sendError(w, r, fmt.Errorf("invalid request body"))
			return
		}

		if task.Title == "" {
			c.sendError(w, r, fmt.Errorf("'title' is empty"))
			return
		}

		if task.Date != "" {
			if _, err := time.Parse(TimeLayout, task.Date); err != nil {
				c.sendError(w, r, fmt.Errorf("invalid 'date' format"))
				return
			}
		}
		if task.Date == "" || task.Date < time.Now().Format(TimeLayout) {
			task.Date = time.Now().Format(TimeLayout)
		}

		if task.Repeat != "" {
			nextDate, err := NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				c.sendError(w, r, fmt.Errorf("invalid 'repeat' format"))
				return
			}
			if task.Date < time.Now().Format(TimeLayout) {
				task.Date = nextDate
			}
		}

		id, err := taskCreator.CreateTask(&task)
		if err != nil {
			c.sendError(w, r, fmt.Errorf("invalid request"))
			return
		}

		render.JSON(w, r, map[string]string{"id": id})
	}
}
