package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"yet-another-todo-list/internal/lib/slwrap"

	srv "yet-another-todo-list/internal/http-server"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type TasksGetter interface {
	GetTasks(searchQuery string) ([]srv.Task, error)
}

func GetTasks(log *slog.Logger, tasksGetter TasksGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log = log.With(slog.String("request_id", middleware.GetReqID(r.Context())))

		if r.Method != http.MethodGet {
			srv.SendError(w, r, srv.ErrMethodNotAllowed)
			return
		}

		tasks, err := tasksGetter.GetTasks(r.FormValue("search"))
		if err != nil {
			log.Error("tasksGetter.GetTasks", slwrap.Wrap(err))
			srv.SendError(w, r, fmt.Errorf("error to get tasks"))
			return
		}

		render.JSON(w, r, map[string]any{"tasks": tasks})
	}
}
