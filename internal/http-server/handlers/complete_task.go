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

type TaskCompleter interface {
	CompleteTask(id string) error
}

func CompleteTask(log *slog.Logger, taskCompleter TaskCompleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log = log.With(slog.String("request_id", middleware.GetReqID(r.Context())))

		if r.Method != http.MethodPost {
			srv.SendError(w, r, srv.ErrMethodNotAllowed)
			return
		}

		id := r.FormValue("id")
		if id == "" {
			srv.SendError(w, r, srv.ErrInvalidID)
			return
		}

		err := taskCompleter.CompleteTask(id)
		if err != nil {
			if errors.Is(err, srv.ErrInvalidID) {
				srv.SendError(w, r, err)
				return
			}
			log.Error("taskCompleter.CompleteTask", slwrap.Wrap(err))
			srv.SendError(w, r, fmt.Errorf("error to complete task"))
			return
		}

		render.JSON(w, r, map[string]string{})
	}
}
