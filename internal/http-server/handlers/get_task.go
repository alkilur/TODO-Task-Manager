package handlers

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	srv "yet-another-todo-list/internal/http-server"
)

type TaskGetter interface {
	GetTask(id string) (*srv.Task, error)
}

func GetTask(log *slog.Logger, taskGetter TaskGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log = log.With(slog.String("request_id", middleware.GetReqID(r.Context())))

		id := r.FormValue("id")
		if id == "" {
			srv.SendError(w, r, srv.ErrInvalidID)
			return
		}

		task, err := taskGetter.GetTask(id)
		if err != nil {
			srv.SendError(w, r, err)
			return
		}

		render.JSON(w, r, task)
	}
}
