package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	srv "yet-another-todo-list/internal/http-server"
	"yet-another-todo-list/internal/lib/slwrap"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type TaskCreator interface {
	CreateTask(*srv.Task) (string, error)
}

func AddTask(log *slog.Logger, taskCreator TaskCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log = log.With(slog.String("request_id", middleware.GetReqID(r.Context())))

		task := srv.Task{}
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			srv.SendError(w, r, srv.ErrUnmarshal)
			return
		}

		if task.Title == "" {
			srv.SendError(w, r, srv.ErrEmptyTitle)
			return
		}

		if task.Date != "" {
			if _, err := time.Parse(srv.TimeLayout, task.Date); err != nil {
				srv.SendError(w, r, srv.ErrInvalidDate)
				return
			}
		}
		if task.Date == "" || task.Date < time.Now().Format(srv.TimeLayout) {
			task.Date = time.Now().Format(srv.TimeLayout)
		}

		if task.Repeat != "" {
			nextDate, err := srv.NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				srv.SendError(w, r, err)
				return
			}
			if task.Date < time.Now().Format(srv.TimeLayout) {
				task.Date = nextDate
			}
		}

		id, err := taskCreator.CreateTask(&task)
		if err != nil {
			log.Error("taskCreator.CreateTask", slwrap.Wrap(err))
			srv.SendError(w, r, fmt.Errorf("error to create task"))
			return
		}

		render.JSON(w, r, map[string]string{"id": id})
	}
}
