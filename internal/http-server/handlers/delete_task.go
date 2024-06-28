package handlers

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	srv "yet-another-todo-list/internal/http-server"
	"yet-another-todo-list/internal/lib/slwrap"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type TaskDeleter interface {
	DeleteTask(id string) error
}

func DeleteTask(log *slog.Logger, taskDeleter TaskDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log = log.With(slog.String("request_id", middleware.GetReqID(r.Context())))

		id := r.FormValue("id")
		if id == "" {
			srv.SendError(w, r, srv.ErrInvalidID)
			return
		}

		err := taskDeleter.DeleteTask(id)
		if err != nil {
			if errors.Is(err, srv.ErrInvalidID) {
				srv.SendError(w, r, err)
				return
			}
			log.Error("taskDeleter.DeleteTask", slwrap.Wrap(err))
			srv.SendError(w, r, fmt.Errorf("error to delete task"))
			return
		}

		render.JSON(w, r, map[string]string{})
	}
}
