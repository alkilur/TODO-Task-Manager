package handlers

import (
	"log/slog"
	"net/http"
	"time"

	srv "yet-another-todo-list/internal/http-server"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

func GetNextDate(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log = log.With(slog.String("request_id", middleware.GetReqID(r.Context())))
		// TODO: log something?

		now, err := time.Parse(srv.TimeLayout, r.FormValue("now"))
		if err != nil {
			srv.SendError(w, r, srv.ErrInvalidNow)
			return
		}

		nextDate, err := srv.NextDate(now, r.FormValue("date"), r.FormValue("repeat"))
		if err != nil {
			srv.SendError(w, r, err)
			return
		}

		render.PlainText(w, r, nextDate)
	}
}
